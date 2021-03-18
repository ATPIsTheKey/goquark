package utils

import "fmt"

type SourceIndex struct {
	FileName string
	Line     int
	Column   int
}

func (sourceIndex *SourceIndex) Fmt() string {
	return fmt.Sprintf("(file: %s, line: %d, col: %d)", sourceIndex.FileName, sourceIndex.Line, sourceIndex.Column)
}
