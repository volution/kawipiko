

package common


import "bytes"
import "fmt"
import "regexp"
import "sort"
import "strconv"




func MetadataEncodeHttp (_metadata map[string]string) ([]byte, error) {
	
	if len (_metadata) > 128 {
		return nil, fmt.Errorf ("[c56fce8f]  invalid metadata:  too large!")
	}
	
	_metadataArray := make ([][2]string, 0, len (_metadata))
	for _key, _value := range _metadata {
		if _value == "" {
			continue
		}
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




func MetadataDecodeHttp (_data []byte) ([][2]string, error) {
	_metadata := make ([][2]string, 0, 16)
	_metadataAppend := func (_key []byte, _value []byte) () {
		_metadata = append (_metadata, [2]string { string (_key), string (_value) })
	}
	if _error := MetadataDecodeHttpIterate (_data, _metadataAppend); _error != nil {
		return nil, _error
	} else {
		return _metadata, nil
	}
}




func MetadataDecodeHttpIterate (_data []byte, _callback func ([]byte, []byte) ()) (error) {
	
	_dataSize := len (_data)
	_headerOffset := 0
	
	for {
		
		if _headerOffset == _dataSize {
			return nil
		}
		
		_data := _data[_headerOffset :]
		_headerLimit := bytes.Index (_data, []byte ("\r\n"))
		if (_headerLimit == -1) {
			return fmt.Errorf ("[2d0d442a]  invalid metadata encoding")
		}
		
		_data = _data[: _headerLimit]
		_separator := bytes.Index (_data, []byte (": "))
		if _separator == -1 {
			return fmt.Errorf ("[41f3756c]  invalid metadata encoding")
		}
		
		_key := _data[: _separator]
		_value := _data[_separator + 2 :]
		if _separator == 0 {
			return fmt.Errorf ("[c3f5e8f3]  invalid metadata encoding (empty key)")
		}
		if _separator == (_headerLimit - 2) {
			return fmt.Errorf ("[d6a923b6]  invalid metadata encoding (empty value)")
		}
		
		_callback (_key, _value)
		
		_headerOffset += _headerLimit + 2
	}
}




func MetadataEncodeBinary (_metadata map[string]string) ([]byte, error) {
	
	if len (_metadata) > 128 {
		return nil, fmt.Errorf ("[2249daa0]  invalid metadata:  too large!")
	}
	
	_metadataArray := make ([][2]string, 0, len (_metadata))
	for _key, _value := range _metadata {
		if _value == "" {
			continue
		}
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
		if (_key != "") && (_key[0] == '!') {
			if _key == "!Status" {
				if _value, _error := strconv.Atoi (_value); _error == nil {
					if (_value >= 200) && (_value <= 599) {
						// NOP
					} else {
						return nil, fmt.Errorf ("[08d97429]  invalid metadata value:  `%d`", _value)
					}
				} else {
					return nil, fmt.Errorf ("[7a36c814]  invalid metadata value:  `%s`", _value)
				}
			} else {
				return nil, fmt.Errorf ("[777a334d]  invalid metadata key:  `%s`", _key)
			}
		} else {
			if ! metadataKeyRegex.MatchString (_key) {
				return nil, fmt.Errorf ("[9c53ceb6]  invalid metadata key:  `%s`", _key)
			}
			if ! metadataValueRegex.MatchString (_value) {
				return nil, fmt.Errorf ("[f932f38f]  invalid metadata value:  `%s`", _value)
			}
		}
		
		_keyId, _keyFound := CanonicalHeaderNamesToKey[_key]
		_valueId, _valueFound := CanonicalHeaderValuesToKey[_value]
		if !_keyFound {
			return nil, fmt.Errorf ("[a2a62863]  invalid metadata key:  `%s` (not canonical)", _key)
		}
		if _valueFound {
			var _pairBuffer [16]byte
			if _error := EncodeKeysPairToBytes_0 (NamespaceHeaderName, _keyId, NamespaceHeaderValue, _valueId, _pairBuffer[:]); _error != nil {
				return nil, _error
			}
			_buffer.Write (_pairBuffer[:])
		} else {
			_valueSize := len (_value)
			if _valueSize > 4096 {
				return nil, fmt.Errorf ("[85c25588]  invalid metadata value:  `%s`", _value)
			}
			var _pairBuffer [12]byte
			if _error := EncodeKeyToBytes_0 (NamespaceHeaderName, _keyId, _pairBuffer[:8]); _error != nil {
				return nil, _error
			}
			_pairBuffer[8] = 'Z'
			_pairBuffer[9] = byte (_valueSize / 25 / 25) + 'a'
			_pairBuffer[10] = byte (_valueSize / 25 % 25) + 'a'
			_pairBuffer[11] = byte (_valueSize % 25) + 'a'
			_buffer.Write (_pairBuffer[:])
			_buffer.Write ([]byte (_value))
		}
	}
	
	_data := _buffer.Bytes ()
	return _data, nil
}




func MetadataDecodeBinary (_data []byte) ([][2]string, error) {
	_metadata := make ([][2]string, 0, 16)
	_metadataAppend := func (_key []byte, _value []byte) () {
		_metadata = append (_metadata, [2]string { string (_key), string (_value) })
	}
	if _error := MetadataDecodeBinaryIterate (_data, _metadataAppend); _error != nil {
		return nil, _error
	} else {
		return _metadata, nil
	}
}




func MetadataDecodeBinaryIterate (_data []byte, _callback func ([]byte, []byte) ()) (error) {
	
	_dataLimit := len (_data)
	_dataOffset := 0
	
	for {
		
		_sliceSize := _dataLimit - _dataOffset
		if _sliceSize == 0 {
			return nil
		}
		_slice := _data[_dataOffset:]
		
		if _slice[0] != NamespaceHeaderNamePrefix {
			return fmt.Errorf ("[f49c93cb]  invalid metadata encoding")
		}
		if _sliceSize < 8 {
			return fmt.Errorf ("[e8d008dc]  invalid metadata encoding")
		}
		var _key []byte
		if _key_0, _found := CanonicalHeaderNamesFromKey[DecodeKey_9 (_slice[0:8])]; _found {
			_key = StringToBytes (_key_0)
		} else {
			return fmt.Errorf ("[7aa09c0f]  invalid metadata encoding")
		}
		
		_dataOffset += 8
		
		_sliceSize = _dataLimit - _dataOffset
		if _sliceSize == 0 {
			return fmt.Errorf ("[77c8bef7]  invalid metadata encoding")
		}
		_slice = _data[_dataOffset:]
		
		var _value []byte
		if _slice[0] == NamespaceHeaderValuePrefix {
			
			if _sliceSize < 8 {
				return fmt.Errorf ("[7cd40b03]  invalid metadata encoding")
			}
			if _value_0, _found := CanonicalHeaderValuesFromKey[DecodeKey_9 (_slice[0:8])]; _found {
				_value = StringToBytes (_value_0)
			} else {
				return fmt.Errorf ("[334e65ef]  invalid metadata encoding")
			}
			
			_dataOffset += 8
			
		} else if _slice[0] == 'Z' {
			
			if _sliceSize < 4 {
				return fmt.Errorf ("[e52b70b0]  invalid metadata encoding")
			}
			
			_valueSize := 0
			_valueSize += int (_slice[1] - 'a') * 25 * 25
			_valueSize += int (_slice[2] - 'a') * 25
			_valueSize += int (_slice[3] - 'a')
			
			if _sliceSize < (4 + _valueSize) {
				return fmt.Errorf ("[3c4a6b51]  invalid metadata encoding")
			}
			
			_value = _slice[4 : 4 + _valueSize]
			
			_dataOffset += 4 + _valueSize
			
		} else {
			return fmt.Errorf ("[2b43651c]  invalid metadata encoding")
		}
		
		_callback (_key, _value)
	}
}




var metadataKeyRegex = regexp.MustCompile (`\A(?:[A-Z0-9](?:[a-z0-9]?[a-z]+)(?:-[A-Z0-9](?:[a-z0-9]?[a-z]+))*)|ETag\z`)
var metadataValueRegex = regexp.MustCompile (`\A[[:graph:]](?: ?[[:graph:]]+)*\z`)

