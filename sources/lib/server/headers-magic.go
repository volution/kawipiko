

package server


import "fmt"
import "net/http"
import "log"
import "reflect"
import "sync/atomic"
import "unsafe"


import . "github.com/volution/kawipiko/lib/common"




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


func (_buffer *HttpResponseWriterHeadersBuffer) IncludeBytes (_name []byte, _value []byte) () {
	
	if _buffer.headersCount == 128 {
		panic ("[ca806ede]")
	}
	
	_buffer.headers[_buffer.headersCount] = [2][]byte {_name, _value}
	_buffer.headersCount += 1
}


func (_buffer *HttpResponseWriterHeadersBuffer) IncludeString (_name string, _value string) () {
	_buffer.IncludeBytes (StringToBytes (_name), StringToBytes (_value))
}




func (_buffer *HttpResponseWriterHeadersBuffer) WriteToGenericResponse (_response http.ResponseWriter) () {
	
	_headers := _response.Header ()
	
	_buffer.WriteToGenericHeaders (_headers)
	
	_response.WriteHeader (_buffer.status)
}




func (_buffer *HttpResponseWriterHeadersBuffer) WriteToGenericHeaders (_headers http.Header) () {
	
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
}




func (_buffer *HttpResponseWriterHeadersBuffer) WriteTo (_response http.ResponseWriter) () {
	
	
	if !_httpResponseWriterHeadersMagic_enabled {
		_buffer.WriteToGenericResponse (_response)
		return
	}
	
	
	_responseRaw := *(*[2]uintptr) (unsafe.Pointer (&_response))
	_responseType := _responseRaw[0]
	
	
	_redo :
	
	_netHttp1Type := atomic.LoadUintptr (&_httpResponseWriterHeadersMagic_netHttp1_type)
	if _responseType == _netHttp1Type {
		_buffer.WriteToNetHttp1 (_response)
		return
	}
	
	_netHttp2Type := atomic.LoadUintptr (&_httpResponseWriterHeadersMagic_netHttp2_type)
	if _responseType == _netHttp2Type {
		_buffer.WriteToNetHttp2 (_response)
		return
	}
	
	_quicHttp3Type := atomic.LoadUintptr (&_httpResponseWriterHeadersMagic_quicHttp3_type)
	if _responseType == _quicHttp3Type {
		_buffer.WriteToQuicHttp3 (_response)
		return
	}
	
	_httpResponseWriterHeadersMagic_detect (_response)
	
	goto _redo
}

var _httpResponseWriterHeadersMagic_netHttp1_type uintptr
var _httpResponseWriterHeadersMagic_netHttp2_type uintptr
var _httpResponseWriterHeadersMagic_quicHttp3_type uintptr




func (_buffer *HttpResponseWriterHeadersBuffer) WriteToNetHttp1 (_response http.ResponseWriter) () {
	
	_responseRaw := *(*[2]uintptr) (unsafe.Pointer (&_response))
	_responseAddress := unsafe.Pointer (_responseRaw[1])
	
	_headers := make (map[string][]string, 16)
	
	{
		_handlerHeaderOffset := atomic.LoadInt32 (&_httpResponseWriterHeadersMagic_netHttp1_handlerHeaderOffset)
		_cwHeaderOffset := atomic.LoadInt32 (&_httpResponseWriterHeadersMagic_netHttp1_cwHeaderOffset)
		
		_handlerHeaderValue := (*http.Header) (unsafe.Add (_responseAddress, _handlerHeaderOffset))
		_cwHeaderValue := (*http.Header) (unsafe.Add (_responseAddress, _cwHeaderOffset))
		
		*_handlerHeaderValue = _headers
		*_cwHeaderValue = _headers
	}
	
	_buffer.WriteToGenericHeaders (_headers)
}




func (_buffer *HttpResponseWriterHeadersBuffer) WriteToNetHttp2 (_response http.ResponseWriter) () {
	
	_buffer.WriteToGenericResponse (_response)
}




func (_buffer *HttpResponseWriterHeadersBuffer) WriteToQuicHttp3 (_response http.ResponseWriter) () {
	
	_buffer.WriteToGenericResponse (_response)
}




func _httpResponseWriterHeadersMagic_detect (_response http.ResponseWriter) {
	
	
	// NOTE:  Because we don't modify the headers after calling `WriteHeader`,
	//        the following code tricks `net/http` into believing it didn't gave us the headers.
	//        This eliminates the `http.Header.Clone()` call on `WriteHeader`.
	
	_responseReflect := reflect.ValueOf (_response) .Elem ()
	_responseAddress := unsafe.Pointer (_responseReflect.UnsafeAddr ())
	_responseType := _responseReflect.Type ()
	_responsePackage := _responseType.PkgPath ()
	_responseTypeName := _responseType.Name ()
	
	_responseRaw := *(*[2]uintptr) (unsafe.Pointer (&_response))
	_responseRawType := _responseRaw[0]
	
	switch {
		
		case (_responsePackage == "net/http") && (_responseTypeName == "response") : {
			
			log.Printf ("[dd] [73f6be95]  [magic...]  detected NetHttp1 (`%s.%s`) with type `%08x`;", _responsePackage, _responseTypeName, _responseRawType)
			
			_handlerHeaderReflect := _responseReflect.FieldByName ("handlerHeader")
			_handlerHeaderAddress := unsafe.Pointer (_handlerHeaderReflect.UnsafeAddr ())
			
			_cwHeaderReflect := _responseReflect.FieldByName ("cw") .FieldByName ("header")
			_cwHeaderAddress := unsafe.Pointer (_cwHeaderReflect.UnsafeAddr ())
			
			_handlerHeaderOffset := int32 (int64 (uintptr (_handlerHeaderAddress)) - int64 (uintptr (_responseAddress)))
			_cwHeaderOffset := int32 (int64 (uintptr (_cwHeaderAddress)) - int64 (uintptr (_responseAddress)))
			
			atomic.StoreInt32 (&_httpResponseWriterHeadersMagic_netHttp1_handlerHeaderOffset, _handlerHeaderOffset)
			atomic.StoreInt32 (&_httpResponseWriterHeadersMagic_netHttp1_cwHeaderOffset, _cwHeaderOffset)
			
			atomic.StoreUintptr (&_httpResponseWriterHeadersMagic_netHttp1_type, _responseRawType)
		}
		
		case (_responsePackage == "net/http") && (_responseTypeName == "http2responseWriter") : {
			
			log.Printf ("[dd] [cfb457eb]  [magic...]  detected NetHttp2 (`%s.%s`) with type `%08x`;", _responsePackage, _responseTypeName, _responseRawType)
			
			atomic.StoreUintptr (&_httpResponseWriterHeadersMagic_netHttp2_type, _responseRawType)
		}
		
		case (_responsePackage == "github.com/lucas-clemente/quic-go/http3") && (_responseTypeName == "responseWriter") : {
			
			log.Printf ("[dd] [90b8f7c6]  [magic...]  detected QuicHttp3 (`%s.%s`) with type `%08x`;", _responsePackage, _responseTypeName, _responseRawType)
			
			atomic.StoreUintptr (&_httpResponseWriterHeadersMagic_quicHttp3_type, _responseRawType)
		}
		
		default : {
			
			log.Printf ("[ee] [64583df9]  unsupported HTTP ResponseWriter `%s.%s`!\n", _responsePackage, _responseTypeName)
			
			if _httpResponseWriterHeadersMagic_panic {
				panic (fmt.Sprintf ("[09274c17]  unsupported HTTP ResponseWriter `%s.%s`!", _responsePackage, _responseTypeName))
			}
		}
	}
}


var _httpResponseWriterHeadersMagic_enabled = true
var _httpResponseWriterHeadersMagic_panic = true

var _httpResponseWriterHeadersMagic_netHttp1_handlerHeaderOffset int32
var _httpResponseWriterHeadersMagic_netHttp1_cwHeaderOffset int32

