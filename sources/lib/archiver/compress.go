

package archiver


import "bytes"
import "compress/gzip"
import "fmt"


import "github.com/foobaz/go-zopfli/zopfli"
import brotli "github.com/itchio/go-brotli/enc"




func Compress (_data []byte, _algorithm string) ([]byte, string, error) {
	switch _algorithm {
		case "gz", "gzip" :
			return CompressGzip (_data)
		case "zopfli" :
			return CompressZopfli (_data)
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




func CompressZopfli (_data []byte) ([]byte, string, error) {
	
	_buffer := & bytes.Buffer {}
	
	_options := zopfli.DefaultOptions ()
	_options.NumIterations = 15
	_options.BlockSplitting = true
	_options.BlockSplittingLast = true
	_options.BlockSplittingMax = 0
	_options.BlockType = zopfli.DYNAMIC_BLOCK
	
	if _error := zopfli.GzipCompress (&_options, _data, _buffer); _error != nil {
		return nil, "", _error
	}
	
	_data = _buffer.Bytes ()
	return _data, "gzip", nil
}




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

