package functionality

import (
	"os"
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
		return FileNode{}, err
	}
	for _, item := range itemList {
		if item.IsDir() {
			newNode := FileNode{
				Name:     item.Name(),
				IsDir:    item.IsDir(),
				Children: []FileNode{},
			}
		}
	}
}
