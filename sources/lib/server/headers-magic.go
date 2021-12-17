

package server


import "fmt"
import "net/http"
import "log"
import "reflect"
import "sync/atomic"
import "unsafe"




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

