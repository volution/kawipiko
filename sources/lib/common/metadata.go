

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
	for _, _data := range bytes.Split (_data, []byte ("\n")) {
		if len (_data) == 0 {
			continue
		}
		_data := bytes.SplitN (_data, []byte (" : "), 2)
		if len (_data) == 2 {
			_metadata = append (_metadata, [2]string { string (_data[0]), string (_data[1]) })
		} else {
			return nil, fmt.Errorf ("[7cb30bf7]  invalid metadata encoding")
		}
	}
	return _metadata, nil
}




var metadataKeyRegex *regexp.Regexp = regexp.MustCompile (`\A[a-z0-9](?:[a-z0-9-]?[a-z]+)*\z`)
var metadataValueRegex *regexp.Regexp = regexp.MustCompile (`\A[[:graph:]](?: ?[[:graph:]]+)*\z`)

