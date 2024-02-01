package split_writer

import (
	"fmt"
	"log"
	"os"
)

type SplitWriter struct {
	currentIndex     uint
	currentSize      uint
	currentFile      *os.File
	fileNameTemplate string
	splitSize        uint
}

func min(a uint, b uint) uint {
	if a < b {
		return a
	}
	return b
}

func New(fileNameTemplate string, splitSize uint) (SplitWriter, error) {
	if splitSize == 0 {
		return SplitWriter{}, fmt.Errorf("split file size must be > 0")
	}
	firstname := fmt.Sprintf("%s.zip", fileNameTemplate)
	f, err := os.Create(firstname)
	if err != nil {
		return SplitWriter{}, fmt.Errorf("failed to create file %s: %s", firstname, err.Error())
	}
	return SplitWriter{
		currentIndex:     0,
		currentSize:      0,
		currentFile:      f,
		fileNameTemplate: fileNameTemplate,
		splitSize:        splitSize,
	}, nil
}

func (w *SplitWriter) Close() error {
	return w.currentFile.Close()
}

func (w *SplitWriter) formatFileName() string {
	return fmt.Sprintf("%s.zip%03d", w.fileNameTemplate, w.currentIndex)
}

func (w *SplitWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	bytestotake := w.splitSize - w.currentSize
	bytescurr := min(uint(len(p)), bytestotake)
	bufcurr, bufnext := p[0:bytescurr], p[bytescurr:]
	_, err = w.currentFile.Write(bufcurr)
	if err != nil {
		return 0, err
	}
	w.currentSize += bytescurr
	if w.currentSize == w.splitSize {
		w.currentFile.Close()
		w.currentIndex += 1
		w.currentSize = 0
		newfile, err := os.Create(w.formatFileName())
		if err != nil {
			err = fmt.Errorf("failed to rotate file: %s", err.Error())
			return 0, err
		}
		_, err = newfile.Write(bufnext)
		w.currentFile = newfile
		if err != nil {
			return 0, err
		}
	} else if w.currentSize > w.splitSize {
		log.Fatal("should not happen")
	}
	return
}
