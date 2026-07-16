package functionality

import (
	"os"
	"path"
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
