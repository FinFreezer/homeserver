package functionality

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

type FileNode struct {
	Name     string     `json:"name"`
	IsDir    bool       `json:"isDir"`
	Children []FileNode `json:"children,omitempty"`
}

func buildFileTree(fromPath string) (FileNode, error) {
	itemList, err := os.ReadDir(fromPath)
	info, err := os.Stat(fromPath)
	if err != nil || !info.IsDir() {
		return FileNode{Name: info.Name(), IsDir: false}, err
	}
	currentDir := FileNode{Name: info.Name(), IsDir: true, Children: []FileNode{}}

	for _, item := range itemList {
		if !item.IsDir() {
			currentDir.Children = append(currentDir.Children,
				FileNode{
					Name:  item.Name(),
					IsDir: item.IsDir(),
				})

		} else {
			childNode, err := buildFileTree(path.Join(fromPath, item.Name()))
			if err != nil {
				return childNode, err
			}
			currentDir.Children = append(currentDir.Children, childNode)
		}
	}
	return currentDir, nil
}

func readByteRange(file *os.File, byteRange string, w http.ResponseWriter) ([]byte, error, int) {
	fileInfo, err := file.Stat()
	fileSize := fileInfo.Size()
	if err != nil {
		return []byte{}, errors.New("Error getting fileinfo."), 500
	}

	if !strings.HasPrefix(byteRange, "bytes=") {
		return []byte{}, errors.New("No range found or corrupted header."), 416
	}
	trimmedByteRange := strings.TrimPrefix(byteRange, "bytes=")
	if start, end, found := strings.Cut(trimmedByteRange, "-"); !found {
		return []byte{}, errors.New("No range found or corrupted header."), 416

	} else {

		maxBytesPerWrite := 5 << 20
		writeData := make([]byte, maxBytesPerWrite)
		startInt, err := strconv.Atoi(start)
		if err != nil {
			log.Println(start)
			return []byte{}, errors.New("Error converting range."), 416
		}

		if end == "" && startInt < int(fileSize) {
			log.Println("Streaming until the end.")
			w.Header().Set("Content-Length", strconv.Itoa(int(fileSize-int64(startInt)+1)))
			for n := startInt; n < int(fileSize); {
				nRead, err := file.ReadAt(writeData, int64(startInt))
				n += nRead
				if err != nil && err != io.EOF {
					return writeData[:nRead], err, 500
				}
				w.Write(writeData[:nRead])
			}
			return writeData, err, 200
		}

		endInt, err := strconv.Atoi(end)
		if err != nil {
			log.Println(end)
			return []byte{}, errors.New("Error converting range."), 416
		}
		if startInt < endInt && endInt < int(fileSize) {
			log.Println("Streaming range.")
			w.Header().Set("Content-Length", strconv.Itoa(int(endInt-startInt+1)))
			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", startInt, endInt, fileSize))
			for n := startInt; n < endInt; {
				//file.Seek(int64(startInt), 0)
				nRead, err := file.ReadAt(writeData, int64(n))
				n += nRead
				if err != nil && err != io.EOF {
					return writeData, err, 500
				}
				nWrote, err := w.Write(writeData[:nRead])
				if err != nil {
					return writeData[nWrote:], err, 500
				}
			}
			return writeData, err, 200
		} else {
			return []byte{}, errors.New("Invalid byte range."), 416
		}
	}
}
