package writer

// FallbackWriter
type FallbackWriter struct {
	file *File
}

// NewFallbackWriter
func NewFallbackWriter(path, fileName string) *FallbackWriter {
	return &FallbackWriter{
		file: NewFile(path, fileName),
	}
}

// Write write bytes to file
func (f *FallbackWriter) Write(bytes []byte) (n int, err error) {
	return f.file.Read(bytes)
}

// Remove remove file
func (f *FallbackWriter) Remove() (err error) {
	return f.file.Remove()
}
