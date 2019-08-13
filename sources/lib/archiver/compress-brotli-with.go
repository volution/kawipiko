
// +build !nobrotli


package archiver


import "bytes"

import "github.com/google/brotli/go/cbrotli"




func CompressBrotli (_data []byte) ([]byte, string, error) {
	
	_buffer := & bytes.Buffer {}
	
	_encoder := cbrotli.NewWriter (_buffer, cbrotli.WriterOptions { Quality : 11, LGWin : 24})
	
	if _, _error := _encoder.Write (_data); _error != nil {
		return nil, "", _error
	}
	if _error := _encoder.Close (); _error != nil {
		return nil, "", _error
	}
	
	_data = _buffer.Bytes ()
	return _data, "br", nil
}

