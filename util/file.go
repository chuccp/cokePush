package util

import "os"

type FileStream struct {
	file *os.File
}

func NewFileStream(filePath string)(*FileStream,error)  {
	file,err:=os.Create(filePath)
	return &FileStream{file:file},err
}
func( file *FileStream)Read(p []byte) (n int, err error){
	return file.file.Read(p)
}

func( file *FileStream)Write(p []byte) (n int, err error){
	return file.file.Write(p)
}
