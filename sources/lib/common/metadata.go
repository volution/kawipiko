

package common


import "bytes"
import "fmt"
import "regexp"




func MetadataEncode (_metadata map[string]string) ([]byte, error) {
	
	_buffer := & bytes.Buffer {}
	
	for _key, _value := range _metadata {
		if ! metadataKeyRegex.MatchString (_key) {
			return nil, fmt.Errorf ("[2f761e02]  invalid metadata key:  `%s`", _key)
		}
		if _value == "" {
			continue
		}
		if ! metadataValueRegex.MatchString (_value) {
			return nil, fmt.Errorf ("[e8faf5bd]  invalid metadata value:  `%s`", _value)
		}
		_buffer.Write ([]byte (_key))
		_buffer.Write ([]byte (" : "))
		_buffer.Write ([]byte (_value))
		_buffer.Write ([]byte ("\n"))
	}
	
	_data := _buffer.Bytes ()
	return _data, nil
}




func MetadataDecode (_data []byte) ([][2]string, error) {
	_metadata := make ([][2]string, 0, 16)
	_metadataAppend := func (_key []byte, _value []byte) () {
		_metadata = append (_metadata, [2]string { string (_key), string (_value) })
	}
	if _error := MetadataDecodeIterate (_data, _metadataAppend); _error != nil {
		return nil, _error
	} else {
		return _metadata, nil
	}
}

func MetadataDecodeIterate (_data []byte, _callback func ([]byte, []byte) ()) (error) {
	
	_dataSize := len (_data)
	_headerOffset := 0
	
	for {
		
		if _headerOffset == _dataSize {
			return nil
		}
		
		_data := _data[_headerOffset :]
		_headerLimit := bytes.IndexByte (_data, '\n')
		if (_headerLimit == -1) {
			return fmt.Errorf ("[2d0d442a]  invalid metadata encoding")
		}
		_headerOffset += _headerLimit + 1
		
		_data = _data[: _headerLimit]
		_separator := bytes.Index (_data, []byte (" : "))
		if _separator == -1 {
			return fmt.Errorf ("[41f3756c]  invalid metadata encoding")
		}
		_key := _data[: _separator]
		_value := _data[_separator + 3 :]
		
		_callback (_key, _value)
	}
}




var metadataKeyRegex *regexp.Regexp = regexp.MustCompile (`\A(?:[A-Z0-9](?:[a-z0-9]?[a-z]+)(?:-[A-Z0-9](?:[a-z0-9]?[a-z]+))*)|ETag\z`)
var metadataValueRegex *regexp.Regexp = regexp.MustCompile (`\A[[:graph:]](?: ?[[:graph:]]+)*\z`)

