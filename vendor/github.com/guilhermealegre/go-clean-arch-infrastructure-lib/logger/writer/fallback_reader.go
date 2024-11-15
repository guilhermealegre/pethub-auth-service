package writer

// FallbackReader
type FallbackReader struct {
	file *File
}

// NewFallbackReader
func NewFallbackReader(path, fileName string) *FallbackReader {
	return &FallbackReader{
		file: NewFile(path, fileName),
	}
}

// ReadLines read file lines
func (f *FallbackReader) ReadLines() (lines []string, err error) {
	return f.file.ReadLines()
}
