package writer

import (
	"bufio"
	"fmt"
	"os"
	"path"

	errorCodes "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors"
)

// File
type File struct {
	// path
	path string
	// file Name
	fileName string
	// file
	file *os.File
}

// NewFile creates a new writer to file
func NewFile(path, fileName string) *File {
	return &File{
		path:     path,
		fileName: fileName,
	}
}

// Write write bytes to file
func (f *File) Write(bytes []byte) (n int, err error) {
	if _, err = os.Stat(f.path); os.IsNotExist(err) {
		err = os.MkdirAll(f.path, 0777)
	}

	if err != nil {
		err = errorCodes.ErrorCreatingDirectory().Formats(f.path)
		fmt.Println(err)
		return 0, err
	}

	if f.file, err = os.OpenFile(f.filePath(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777); err != nil {
		err = errorCodes.ErrorOpeningFile().Formats(f.filePath())
		fmt.Println(err)
		return 0, err
	}

	defer f.file.Close()

	if n, err = f.file.Write(bytes); err != nil {
		err = errorCodes.ErrorWritingFile().Formats(f.filePath())
		fmt.Println(err)
		return 0, err
	}

	return n, nil
}

// open open file
func (f *File) open() (err error) {
	if f.file, err = os.Open(f.filePath()); err != nil {
		return err
	}
	return nil
}

// Read read file bytes
func (f *File) Read(p []byte) (n int, err error) {
	if f.file == nil {
		if err = f.open(); err != nil {
			return 0, err
		}
	}

	return f.file.Read(p)
}

// ReadLines read file lines
func (f *File) ReadLines() (lines []string, err error) {
	var file *os.File
	if file, err = os.Open(f.filePath()); err != nil {
		// ignore non existing file
		return nil, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

// Remove remove file
func (f *File) Remove() error {
	return os.Remove(f.filePath())
}

// filePath returns file path
func (f *File) filePath() string {
	return path.Join(f.path, f.fileName)
}
