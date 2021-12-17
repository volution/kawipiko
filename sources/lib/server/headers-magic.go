

package server


import "fmt"
import "net/http"
import "log"
import "reflect"
import "sync/atomic"
import "unsafe"




type HttpResponseWriterHeadersBuffer struct {
	status int
	headers [128][2][]byte
	headersCount int
}


func NewHttpResponseWriterHeadersBuffer (_status int) (HttpResponseWriterHeadersBuffer) {
	
	return HttpResponseWriterHeadersBuffer {
			status : _status,
			headersCount : 0,
		}
}


func (_buffer *HttpResponseWriterHeadersBuffer) Include (_name []byte, _value []byte) () {
	
	if _buffer.headersCount == 128 {
		panic ("[ca806ede]")
	}
	
	_buffer.headers[_buffer.headersCount] = [2][]byte {_name, _value}
	_buffer.headersCount += 1
}


func (_buffer *HttpResponseWriterHeadersBuffer) WriteTo (_response http.ResponseWriter) () {
	
	_headers := HttpResponseWriterHeaderDoMagic (_response)
	
	for _index := 0; _index < _buffer.headersCount; _index += 1 {
		
		_nameAndValue_0 := _buffer.headers[_index]
		_name_0 := _nameAndValue_0[0]
		_value_0 := _nameAndValue_0[1]
		
		_name := CanonicalHeaderNameFromBytes (_name_0)
		if _values, _ := _headers[_name]; _values == nil {
			_values = CanonicalHeaderValueArrayFromBytes (_value_0)
			_headers[_name] = _values
		} else {
			_value := CanonicalHeaderValueFromBytes (_value_0)
			_values = append (_values, _value)
			_headers[_name] = _values
		}
	}
	
	_headers["Date"] = nil
	
	_response.WriteHeader (_buffer.status)
}




func HttpResponseWriterHeaderDoMagic (_response http.ResponseWriter) (http.Header) {
	
	
	if !_httpResponseWriteHeaderMagic_enabled {
		return _response.Header ()
	}
	
	
	// NOTE:  Because we don't modify the headers after calling `WriteHeader`,
	//        the following code tricks `net/http` into believing it didn't gave us the headers.
	//        This eliminates the `http.Header.Clone()` call on `WriteHeader`.
	
	_responseReflect := reflect.ValueOf (_response) .Elem ()
	_responseAddress := unsafe.Pointer (_responseReflect.UnsafeAddr ())
	_responseType := _responseReflect.Type ()
	_responsePackage := _responseType.PkgPath ()
	_responseTypeName := _responseType.Name ()
	
	var _header http.Header
	
	switch {
		
		case (_responsePackage == "net/http") && (_responseTypeName == "response") : {
			
			_handlerHeaderOffset := atomic.LoadInt32 (&_httpResponseWriteHeaderMagic_netHttp_handlerHeaderOffset)
			_cwHeaderOffset := atomic.LoadInt32 (&_httpResponseWriteHeaderMagic_netHttp_cwHeaderOffset)
			
			if (_handlerHeaderOffset == 0) || (_cwHeaderOffset == 0) {
				
				_handlerHeaderReflect := _responseReflect.FieldByName ("handlerHeader")
				_handlerHeaderAddress := unsafe.Pointer (_handlerHeaderReflect.UnsafeAddr ())
				
				_cwHeaderReflect := _responseReflect.FieldByName ("cw") .FieldByName ("header")
				_cwHeaderAddress := unsafe.Pointer (_cwHeaderReflect.UnsafeAddr ())
				
				_handlerHeaderOffset = int32 (int64 (uintptr (_handlerHeaderAddress)) - int64 (uintptr (_responseAddress)))
				_cwHeaderOffset = int32 (int64 (uintptr (_cwHeaderAddress)) - int64 (uintptr (_responseAddress)))
				
				atomic.StoreInt32 (&_httpResponseWriteHeaderMagic_netHttp_handlerHeaderOffset, _handlerHeaderOffset)
				atomic.StoreInt32 (&_httpResponseWriteHeaderMagic_netHttp_cwHeaderOffset, _cwHeaderOffset)
			}
			
			_handlerHeaderValue := (*http.Header) (unsafe.Add (_responseAddress, _handlerHeaderOffset))
			_cwHeaderValue := (*http.Header) (unsafe.Add (_responseAddress, _cwHeaderOffset))
			
			_header = make (map[string][]string, 16)
			
			*_handlerHeaderValue = _header
			*_cwHeaderValue = _header
		}
		
		case (_responsePackage == "net/http") && (_responseTypeName == "http2responseWriter") : {
			
			_header = _response.Header ()
		}
		
		case (_responsePackage == "github.com/lucas-clemente/quic-go/http3") && (_responseTypeName == "responseWriter") : {
			
			_header = _response.Header ()
		}
		
		default : {
			
			log.Printf ("[ee] [64583df9]  unsupported HTTP ResponseWriter `%s.%s`!\n", _responsePackage, _responseTypeName)
			
			if _httpResponseWriteHeaderMagic_panic {
				panic (fmt.Sprintf ("[09274c17]  unsupported HTTP ResponseWriter `%s.%s`!", _responsePackage, _responseTypeName))
			}
			
			_header = _response.Header ()
		}
	}
	
	return _header
}


var _httpResponseWriteHeaderMagic_enabled = true
var _httpResponseWriteHeaderMagic_panic = true

var _httpResponseWriteHeaderMagic_netHttp_handlerHeaderOffset int32
var _httpResponseWriteHeaderMagic_netHttp_cwHeaderOffset int32

