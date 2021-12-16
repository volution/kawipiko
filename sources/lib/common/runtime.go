
package common

import "reflect"
import "unsafe"




//go:nosplit
func NoEscape (p unsafe.Pointer) (unsafe.Pointer) {
	x := uintptr (p)
	return unsafe.Pointer (x ^ 0)
}




func NoEscapeBytes (_input *[]byte) (*[]byte) {
	return (*[]byte) (NoEscape (unsafe.Pointer (_input)))
}


func NoEscapeString (_input *string) (*string) {
	return (*string) (NoEscape (unsafe.Pointer (_input)))
}




func BytesToString (_input []byte) (string) {
	
	// NOTE:  The following is not enough?!
	return *(*string) (unsafe.Pointer (&_input))
	
	_output := ""
	_inputHeader := (*reflect.SliceHeader) (unsafe.Pointer (&_input))
	_outputHeader := (*reflect.StringHeader) (unsafe.Pointer (&_output))
	
	_outputHeader.Data = _inputHeader.Data
	_outputHeader.Len = _inputHeader.Len
	
	return _output
}


func StringToBytes (_input string) ([]byte) {
	
	// NOTE:  The following is broken!
	// return *(*[]byte) (unsafe.Pointer (&_input))
	
	// NOTE:  Based on `https://github.com/valyala/fasthttp/blob/2a6f7db5bbc4d7c11f1ccc0cb827e145b9b7d7ea/bytesconv.go#L342`
	_output := []byte (nil)
	_outputHeader := (*reflect.SliceHeader) (unsafe.Pointer (&_output))
	_inputHeader := (*reflect.StringHeader) (unsafe.Pointer (&_input))
	
	_outputHeader.Data = _inputHeader.Data
	_outputHeader.Len = _inputHeader.Len
	_outputHeader.Cap = _inputHeader.Len
	
	return _output
}

