
package common

import "unsafe"




//go:nosplit
func NoEscape (p unsafe.Pointer) (unsafe.Pointer) {
	x := uintptr (p)
	return unsafe.Pointer (x ^ 0)
}

func NoEscapeBytes (p *[]byte) (*[]byte) {
	return (*[]byte) (NoEscape (unsafe.Pointer (&p)))
}

func NoEscapeString (p *string) (*string) {
	return (*string) (NoEscape (unsafe.Pointer (&p)))
}
