
// +build nobrotli


package archiver


import "fmt"




func CompressBrotli (_data []byte) ([]byte, string, error) {
	
	return nil, "", fmt.Errorf ("[9252cf70]  unsupported compression algorithm `%s`", "brotli")
}

