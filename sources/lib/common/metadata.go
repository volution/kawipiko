

package common


import "bytes"
import "fmt"
import "regexp"
import "sort"




func MetadataEncode (_metadata map[string]string) ([]byte, error) {
	
	_metadataArray := make ([][2]string, 0, len (_metadata))
	for _key, _value := range _metadata {
		_metadataArray = append (_metadataArray, [2]string {_key, _value})
	}
	sort.Slice (_metadataArray,
			func (i int, j int) (bool) {
				return _metadataArray[i][0] < _metadataArray[j][0]
			})
	
	_buffer := & bytes.Buffer {}
	
	for _, _metadata := range _metadataArray {
		_key := _metadata[0]
		_value := _metadata[1]
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
		_buffer.Write ([]byte (": "))
		_buffer.Write ([]byte (_value))
		_buffer.Write ([]byte ("\r\n"))
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
		_headerLimit := bytes.IndexByte (_data, '\r')
		if (_headerLimit == -1) {
			return fmt.Errorf ("[2d0d442a]  invalid metadata encoding")
		}
		if ((_headerOffset + _headerLimit) == (_dataSize - 1)) || (_data[_headerLimit + 1] != '\n') {
			return fmt.Errorf ("[0e319685]  invalid metadata encoding")
		}
		_headerOffset += _headerLimit + 2
		
		_data = _data[: _headerLimit]
		_separator := bytes.Index (_data, []byte (": "))
		if _separator == -1 {
			return fmt.Errorf ("[41f3756c]  invalid metadata encoding")
		}
		_key := _data[: _separator]
		_value := _data[_separator + 3 :]
		if len (_key) == 0 {
			return fmt.Errorf ("[c3f5e8f3]  invalid metadata encoding (empty key)")
		}
		if len (_value) == 0 {
			return fmt.Errorf ("[d6a923b6]  invalid metadata encoding (empty value)")
		}
		
		_callback (_key, _value)
	}
}




var metadataKeyRegex = regexp.MustCompile (`\A(?:[A-Z0-9](?:[a-z0-9]?[a-z]+)(?:-[A-Z0-9](?:[a-z0-9]?[a-z]+))*)|ETag\z`)
var metadataValueRegex = regexp.MustCompile (`\A[[:graph:]](?: ?[[:graph:]]+)*\z`)

