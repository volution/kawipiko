

package archiver


import "bytes"
import "compress/gzip"
import "fmt"


import "github.com/foobaz/go-zopfli/zopfli"
import "github.com/andybalholm/brotli"




func Compress (_data []byte, _algorithm string, _level int) ([]byte, error) {
	switch _algorithm {
		case "gz", "gzip" :
			return CompressGzip (_data, _level)
		case "zopfli" :
			return CompressZopfli (_data, _level)
		case "br", "brotli" :
			return CompressBrotli (_data, _level)
		case "", "none", "identity" :
			return _data, nil
		default :
			return nil, fmt.Errorf ("[ea23f966]  invalid compression algorithm `%s`", _algorithm)
	}
}


func CompressEncoding (_algorithm string) (string, string, error) {
	switch _algorithm {
		case "gz", "gzip" :
			return "gzip", "gzip", nil
		case "zopfli" :
			return "zopfli", "gzip", nil
		case "br", "brotli" :
			return "brotli", "br", nil
		case "", "none", "identity" :
			return "identity", "identity", nil
		default :
			return "", "", fmt.Errorf ("[6026a403]  invalid compression algorithm `%s`", _algorithm)
	}
}




func CompressGzip (_data []byte, _level int) ([]byte, error) {
	
	_buffer := & bytes.Buffer {}
	
	if (_level < -2) || (_level > 9) {
		return nil, fmt.Errorf ("[a6c02c58]  invalid compression level `%d` (-1 for default, -2 for Huffman only, 0 to 9 for fast to slow)", _level)
	}
	
	var _encoder *gzip.Writer
	if _encoder_0, _error := gzip.NewWriterLevel (_buffer, _level); _error == nil {
		_encoder = _encoder_0
	} else {
		return nil, _error
	}
	
	if _, _error := _encoder.Write (_data); _error != nil {
		return nil, _error
	}
	if _error := _encoder.Close (); _error != nil {
		return nil, _error
	}
	
	_data = _buffer.Bytes ()
	return _data, nil
}




func CompressZopfli (_data []byte, _level int) ([]byte, error) {
	
	if (_level < -1) || (_level > 30) {
		return nil, fmt.Errorf ("[fe30ea07]  invalid compression level `%d` (-1 for default, 0 to 30 iterations for fast to slow)", _level)
	}
	
	_buffer := & bytes.Buffer {}
	
	_options := zopfli.DefaultOptions ()
	if _level != -1 {
		_options.NumIterations = _level
		_options.BlockSplitting = true
		_options.BlockSplittingLast = false
		_options.BlockSplittingMax = 0
		_options.BlockType = zopfli.DYNAMIC_BLOCK
	}
	
	if _error := zopfli.GzipCompress (&_options, _data, _buffer); _error != nil {
		return nil, _error
	}
	
	_data = _buffer.Bytes ()
	return _data, nil
}




func CompressBrotli (_data []byte, _level int) ([]byte, error) {
	
	if (_level < -2) || (_level > 9) {
		return nil, fmt.Errorf ("[4aa20d1b]  invalid compression level `%d` (-1 for default, 0 to 9 for fast to slow, -2 for extreme)", _level)
	}
	
	_buffer := & bytes.Buffer {}
	
	_options := brotli.WriterOptions {}
	if _level == -2 {
		_options.Quality = 11
		_options.LGWin = 24
	} else if _level == -1 {
		_options.Quality = 6
	} else {
		_options.Quality = _level
	}
	
	_encoder := brotli.NewWriterOptions (_buffer, _options)
	
	if _, _error := _encoder.Write (_data); _error != nil {
		return nil, _error
	}
	if _error := _encoder.Close (); _error != nil {
		return nil, _error
	}
	
	_data = _buffer.Bytes ()
	return _data, nil
}

