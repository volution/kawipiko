

package archiver


import "bytes"
import "compress/gzip"
import "fmt"

import "github.com/google/brotli/go/cbrotli"




func Compress (_data []byte, _algorithm string) ([]byte, string, error) {
	switch _algorithm {
		case "gz", "gzip" :
			return CompressGzip (_data)
		case "br", "brotli" :
			return CompressBrotli (_data)
		case "", "none", "identity" :
			return _data, "identity", nil
		default :
			return nil, "", fmt.Errorf ("[ea23f966]  invalid compression algorithm `%s`", _algorithm)
	}
}




func CompressGzip (_data []byte) ([]byte, string, error) {
	
	_buffer := & bytes.Buffer {}
	
	var _encoder *gzip.Writer
	if _encoder_0, _error := gzip.NewWriterLevel (_buffer, gzip.BestCompression); _error == nil {
		_encoder = _encoder_0
	} else {
		return nil, "", _error
	}
	
	if _, _error := _encoder.Write (_data); _error != nil {
		return nil, "", _error
	}
	if _error := _encoder.Close (); _error != nil {
		return nil, "", _error
	}
	
	_data = _buffer.Bytes ()
	return _data, "gzip", nil
}


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

