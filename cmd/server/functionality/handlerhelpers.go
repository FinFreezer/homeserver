package functionality

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type FileNode struct {
	Name     string     `json:"name"`
	IsDir    bool       `json:"isDir"`
	Children []FileNode `json:"children,omitempty"`
}

func buildFileTree(fromPath string, dirOnly bool, depth int) (FileNode, error) {
	itemList, err := os.ReadDir(fromPath)
	info, err := os.Stat(fromPath)
	if err != nil || !info.IsDir() {
		return FileNode{Name: info.Name(), IsDir: false}, err
	}
	currentDir := FileNode{Name: info.Name(), IsDir: true, Children: []FileNode{}}
	if dirOnly {
		for _, item := range itemList {
			if item.IsDir() {
				childNode := FileNode{}
				err = nil
				if depth > 0 {
					childNode, err = buildFileTree(path.Join(fromPath, item.Name()), dirOnly, depth-1)
				} else {
					childNode.Name = item.Name()
					childNode.IsDir = item.IsDir()
					childNode.Children = nil
				}
				currentDir.Children = append(currentDir.Children, childNode)
				if err != nil {
					return childNode, err
				}
			}
		}
		return currentDir, nil

	} else {
		for _, item := range itemList {
			if !item.IsDir() {
				currentDir.Children = append(currentDir.Children,
					FileNode{
						Name:  item.Name(),
						IsDir: item.IsDir(),
					})

			} else {
				childNode := FileNode{}
				err = nil
				if depth > 0 {
					childNode, err = buildFileTree(path.Join(fromPath, item.Name()), dirOnly, depth-1)
				} else {
					childNode.Name = item.Name()
					childNode.IsDir = item.IsDir()
					childNode.Children = nil
				}
				if err != nil {
					return childNode, err
				}
				currentDir.Children = append(currentDir.Children, childNode)
			}
		}
		return currentDir, nil
	}

}

func readByteRange(file *os.File, byteRange string, w http.ResponseWriter) (error, int) {
	fileInfo, err := file.Stat()
	fileSize := fileInfo.Size()
	if err != nil {
		return errors.New("Error getting fileinfo."), 500
	}

	if !strings.HasPrefix(byteRange, "bytes=") {
		return errors.New("No range found or corrupted header."), 416
	}
	trimmedByteRange := strings.TrimPrefix(byteRange, "bytes=")
	if start, end, found := strings.Cut(trimmedByteRange, "-"); !found {
		return errors.New("No range found or corrupted header."), 416

	} else {

		maxBytesPerWrite := 5 << 20
		writeData := make([]byte, maxBytesPerWrite)
		startInt, err := strconv.Atoi(start)
		if err != nil {
			log.Println(start)
			return errors.New("Error converting range."), 416
		}

		if end == "" && startInt < int(fileSize) {
			w.Header().Set("Content-Length", strconv.Itoa(int(fileSize-int64(startInt)+1)))
			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", startInt, fileSize-1, fileSize))
			log.Println("Streaming until the end.")
			w.WriteHeader(206)
			for n := startInt; n < int(fileSize); {
				nRead, err := file.ReadAt(writeData, int64(startInt))
				if nRead == 0 {
					if flusher, ok := w.(http.Flusher); ok {
						flusher.Flush()
					}
					return nil, 0
				}
				startInt += nRead
				if err != nil && err != io.EOF {
					return err, 0
				}
				_, err = w.Write(writeData[:nRead])
				if err != nil {
					if flusher, ok := w.(http.Flusher); ok {
						flusher.Flush()
					}
					return nil, 0
				}
			}
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
			return err, 0
		} else {
			endInt, err := strconv.Atoi(end)
			if err != nil {
				log.Println(end)
				return errors.New("Error converting range."), 416
			}
			if startInt < endInt && endInt < int(fileSize) {
				w.Header().Set("Content-Length", strconv.Itoa(int(endInt-startInt+1)))
				w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", startInt, endInt, fileSize))
				log.Println("Streaming range.")
				w.WriteHeader(206)
				for n := startInt; n < endInt; {
					if endInt-n < len(writeData) {
						writeData = make([]byte, endInt-n-1)
					}
					//file.Seek(int64(startInt), 0)
					nRead, err := file.ReadAt(writeData, int64(n))
					if nRead == 0 {
						if flusher, ok := w.(http.Flusher); ok {
							flusher.Flush()
						}
						return nil, 0
					}
					n += nRead
					if err != nil && err != io.EOF {
						return err, 500
					}
					_, err = w.Write(writeData[:nRead])
					if err != nil {
						if flusher, ok := w.(http.Flusher); ok {
							flusher.Flush()
						}
						return nil, 0
					}
				}
				return err, 200
			} else {
				return errors.New("Invalid byte range."), 416
			}
		}
	}
}

func readContentType(f *os.File) (string, error) {
	readBuffer := make([]byte, 512)
	n, err := f.Read(readBuffer)
	f.Seek(0, 0)
	if err != nil && err != io.EOF && n != 0 {
		return "", err
	}
	return http.DetectContentType(readBuffer), nil
}

func createDefaultPlaylist(path string, a *ApiConfig) *os.File {
	log.Printf("Creating a playlist with current directory for %s", path)
	streamingPath := os.Getenv("DST_SERVER") + os.Getenv("DFLT_PORT") + "/stream/"
	startFrom := filepath.Base(path)
	log.Printf("Start from episode '%s'\n", startFrom)
	dir := filepath.Dir(path)
	dirStat, err := os.Stat(dir)
	if err != nil {
		log.Println(err)
		return nil
	}
	if _, err := os.Stat(dir + "/playlist.m3u"); !errors.Is(err, os.ErrNotExist) {
		os.Remove(dir + "/playlist.m3u")
	}

	playlist, err := os.Create(dir + "/playlist.m3u")
	if err != nil {
		log.Println(err)
		return nil
	}
	playlist.Write([]byte("#EXTM3U\n"))
	playlist.Write([]byte("#PLAYLIST: Streaming\n"))
	if dirStat.IsDir() {
		entryList, err := os.ReadDir(dir)
		entryList, episode := moveEpisodeToStart(entryList, startFrom)
		if err != nil {
			log.Println(err)
			return nil
		}
		for _, entry := range entryList {
			infoStr := fmt.Sprintf("Episode %d", episode)
			urlPath, err := url.Parse(streamingPath + (cleanPath(dir+"/"+entry.Name(), a.CurrentRoot)))
			if err != nil {
				log.Println(err)
				return nil
			}
			if !entry.IsDir() {
				if isValidType(entry.Name()) {

					toWrite := fmt.Sprintf("#EXTINF:-1,%s\n%s\n",
						infoStr, urlPath.String(),
					)
					playlist.Write([]byte(toWrite))
					episode += 1
				}
			}
		}
		playlist.Seek(0, 0)
		return playlist
	} else {
		return nil
	}
}

func isValidType(name string) bool {
	validTypes := []string{".mp3", ".mp4", ".mkv", ".wav", ".avi", ".webm"}
	for _, typeName := range validTypes {
		if strings.Contains(name, typeName) {
			return true
		}
	}
	return false
}

func cleanPath(internal string, currentRoot string) string {
	pathList := strings.Split(internal, "/")
	rootList := strings.Split(currentRoot, "/")
	rootList[0] = "stream"
	pathList[0] = "stream"
	relativeRoot := path.Join(rootList...) + "/"
	newPath := path.Join(pathList...)
	log.Printf("Moving %s to match %s", newPath, relativeRoot)
	pathToReturn := strings.ReplaceAll(newPath, relativeRoot, "")
	return pathToReturn
}

func moveEpisodeToStart(files []os.DirEntry, firstEp string) ([]os.DirEntry, int) {
	firstIdx := 0
	newFiles := []os.DirEntry{}
	for i, file := range files {
		if strings.Contains(file.Name(), firstEp) {
			firstIdx = i
		}
	}
	newFiles = append(newFiles, files[firstIdx:]...)
	newFiles = append(newFiles, files[:firstIdx]...)
	return newFiles, firstIdx + 1
}
