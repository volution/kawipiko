
package common

import "reflect"
import "unsafe"




//go:nosplit
func NoEscape (p unsafe.Pointer) (unsafe.Pointer) {
	x := uintptr (p)
	return unsafe.Pointer (x ^ 0)
}

func NoEscapeBytes (p *[]byte) (*[]byte) {
	return (*[]byte) (NoEscape (unsafe.Pointer (p)))
}

func NoEscapeString (p *string) (*string) {
	return (*string) (NoEscape (unsafe.Pointer (p)))
}


func BytesToString (b []byte) (string) {
	return *(*string) (unsafe.Pointer (&b))
}

func StringToBytes (_string string) ([]byte) {
	
	// NOTE:  The following is broken!
	// return *(*[]byte) (unsafe.Pointer (&_string))
	
	// NOTE:  Based on `https://github.com/valyala/fasthttp/blob/2a6f7db5bbc4d7c11f1ccc0cb827e145b9b7d7ea/bytesconv.go#L342`
	_bytes := []byte (nil)
	_bytesHeader := (*reflect.SliceHeader) (unsafe.Pointer (&_bytes))
	_stringHeader := (*reflect.StringHeader) (unsafe.Pointer (&_string))
	_bytesHeader.Data = _stringHeader.Data
	_bytesHeader.Len = _stringHeader.Len
	_bytesHeader.Cap = _stringHeader.Len
	return _bytes
}

