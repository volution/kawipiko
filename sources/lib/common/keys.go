

package common


import "encoding/binary"
import "fmt"




func PrepareKeyToString (_namespace string, _key uint64) (string, error) {
	if _key_0, _error := PrepareKey (_namespace, _key); _error == nil {
		return EncodeKeyToString (_namespace, _key_0)
	} else {
		return "", _error
	}
}




func EncodeKeyToString (_namespace string, _key uint64) (string, error) {
	var _buffer [8]byte
	if _error := EncodeKeyToBytes_0 (_namespace, _key, _buffer[:]); _error == nil {
		return string (_buffer[:]), nil
	} else {
		return "", _error
	}
}


func EncodeKeyToBytes (_namespace string, _key uint64) ([]byte, error) {
	_buffer := make ([]byte, 8)
	_error := EncodeKeyToBytes_0 (_namespace, _key, _buffer)
	return _buffer, _error
}


func EncodeKeyToBytes_0 (_namespace string, _key uint64, _buffer []byte) (error) {
	if len (_buffer) != 8 {
		return fmt.Errorf ("[890eef13]  invalid key buffer length!")
	}
	_prefix := keyNamespacePrefix (_namespace)
	if _prefix == 0 {
		return fmt.Errorf ("[feece73b]  invalid key namespace `%s`!", _namespace)
	}
	if (_key << 8) == 0 {
		return fmt.Errorf ("[c8fa7817]  invalid key zero!")
	}
	_keyPrefix := byte (_key >> 56)
	if _keyPrefix != _prefix {
		return fmt.Errorf ("[85a5c362]  invalid key prefix `%0x` for `%d`!", _key, _keyPrefix)
	}
	binary.BigEndian.PutUint64 (_buffer, _key)
	return nil
}




func PrepareKey (_namespace string, _key uint64) (uint64, error) {
	if _key == 0 {
		return 0, fmt.Errorf ("[355f64d0]  invalid key zero!")
	}
	if _key >= (1 << 24) {
		return 0, fmt.Errorf ("[aba09b4d]  invalid key value `%d`!", _key)
	}
	_prefix := keyNamespacePrefix (_namespace)
	if _prefix == 0 {
		return 0, fmt.Errorf ("[feece73b]  invalid key namespace `%s`!", _namespace)
	}
	_keyString := fmt.Sprintf ("%08x", _key)
	_keyBytes := []byte (_keyString)
	_keyBytes[0] = _prefix
	_keyValue := binary.BigEndian.Uint64 (_keyBytes)
	return _keyValue, nil
}




func keyNamespacePrefix (_namespace string) (byte) {
	switch _namespace {
		case NamespaceFilesContent : return 'f'
		case NamespaceFilesIndex : return 'F'
		case NamespaceFoldersContent : return 'l'
		case NamespaceFoldersIndex : return 'L'
		case NamespaceDataContent : return 'd'
		case NamespaceDataMetadata : return 'm'
		default : return '0'
	}
}




func EncodeKeysPairToString (_namespace1 string, _key1 uint64, _namespace2 string, _key2 uint64) (string, error) {
	var _buffer [16]byte
	if _error := EncodeKeyToBytes_0 (_namespace1, _key1, _buffer[0:8]); _error != nil {
		return "", _error
	}
	if _error := EncodeKeyToBytes_0 (_namespace2, _key2, _buffer[8:16]); _error != nil {
		return "", _error
	}
	return string (_buffer[:]), nil
}




func DecodeKeysPair (_buffer []byte) (uint64, uint64, error) {
	if len (_buffer) != 16 {
		return 0, 0, fmt.Errorf ("[dd6ba461]  invalid key buffer length!")
	}
	_key1 := binary.BigEndian.Uint64 (_buffer[0:8])
	_key2 := binary.BigEndian.Uint64 (_buffer[8:16])
	return _key1, _key2, nil
}

