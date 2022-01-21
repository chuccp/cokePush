package km

import (
	"bytes"
	"testing"
	"time"
)

func TestChunk2(t *testing.T) {

	bytesArray := make([]byte, 0)

	bytesArray = append(bytesArray, 1, 2, 3, 4, 5, 5, 6, 7, 8, 8, 9, 9, 0, 0)

}
func TestAppend(b *testing.T) {


	t := time.Now().UnixNano()

	bytesArray := make([]byte, 0)
	for i := 0; i < 1000_0; i++ {
		bytesArray = append(bytesArray, 1, 2, 3, 4, 5, 5, 6, 7, 8, 8, 9, 9, 0, 0)
	}

	b.Log(time.Now().UnixNano() - t)

}
func TestAppend2(b *testing.T) {
	t := time.Now().UnixNano()
	bytesArray := new(bytes.Buffer)
	for i := 0; i < 1000_0; i++ {
		bytesArray.Write([]byte{1, 2, 3, 4, 5, 5, 6, 7, 8, 8, 9, 9, 0, 0})
	}
	bytesArray.Bytes()
	b.Log(time.Now().UnixNano() - t)
}
