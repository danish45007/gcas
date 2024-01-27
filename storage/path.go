package storage

import "fmt"

type Path struct {
	PathName string
	FileName string
}

func (p Path) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.FileName)
}
