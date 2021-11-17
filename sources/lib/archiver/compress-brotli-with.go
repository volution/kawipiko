
// +build !nobrotli


package archiver


import "bytes"

import brotli "github.com/itchio/go-brotli/enc"




func CompressBrotli (_data []byte) ([]byte, string, error) {
	
	_buffer := & bytes.Buffer {}
	
	_options := brotli.BrotliWriterOptions { Quality : 11, LGWin : 24}
	
	_encoder := brotli.NewBrotliWriter (_buffer, &_options)
	
	if _, _error := _encoder.Write (_data); _error != nil {
		return nil, "", _error
	}
	if _error := _encoder.Close (); _error != nil {
		return nil, "", _error
	}
	
	_data = _buffer.Bytes ()
	return _data, "br", nil
}

