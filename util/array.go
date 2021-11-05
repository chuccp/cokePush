package util

import "bytes"

func Equal(a []byte, b []byte) bool {


	return len(a)==len(b)&&bytes.Equal(a,b)
}
