

package main


import "flag"
import "fmt"
import "log"
import "net/http"
import "net/url"

import "github.com/colinmarc/cdb"

import . "../lib"



type server struct {
	cdbReader *cdb.CDB
	debug bool
}


func (_server *server) ServeHTTP (_response http.ResponseWriter, _request *http.Request) () {
	
	_responseHeaders := _response.Header ()
	
	_responseHeaders.Set ("Content-Security-Policy", "upgrade-insecure-requests")
	_responseHeaders.Set ("Referrer-Policy", "strict-origin-when-cross-origin")
	_responseHeaders.Set ("X-Frame-Options", "SAMEORIGIN")
	_responseHeaders.Set ("X-Content-Type-Options", "nosniff")
	_responseHeaders.Set ("X-XSS-Protection", "1; mode=block")
	
	_method := _request.Method
	_path := _request.URL.Path
	
	if _method != http.MethodGet {
		log.Printf ("[ww] [bce7a75b] invalid method `%s` for `%s`!", _method, _path)
		_server.ServeError (_response, http.StatusMethodNotAllowed, nil)
		return
	}
	
	var _fingerprint string
	{
		_path_0 := _path
		if (_path != "/") && (_path[len (_path) - 1] == '/') {
			_path_0 = _path[: len (_path) - 1]
		}
		for _, _namespace := range []string {NamespaceFilesContent, NamespaceFoldersContent, NamespaceFoldersMetadata} {
			_key := fmt.Sprintf ("%s:%s", _namespace, _path_0)
			if _value, _error := _server.cdbReader.Get ([]byte (_key)); _error == nil {
				if _value != nil {
					_fingerprint = string (_value)
					if ((_namespace == NamespaceFoldersContent) || (_namespace == NamespaceFoldersMetadata)) && (_path == _path_0) {
						_server.ServeRedirect (_response, http.StatusTemporaryRedirect, _path + "/")
						return
					}
					break
				}
			} else {
				_server.ServeError (_response, http.StatusInternalServerError, _error)
				return
			}
		}
	}
	if _fingerprint == "" {
		log.Printf ("[ww] [7416f61d]  not found for `%s`!", _path)
		_server.ServeError (_response, http.StatusNotFound, nil)
		return
	}
	
	var _data []byte
	{
		_key := fmt.Sprintf ("%s:%s", NamespaceDataContent, _fingerprint)
		if _value, _error := _server.cdbReader.Get ([]byte (_key)); _error == nil {
			if _value != nil {
				_data = _value
			}
		} else {
			_server.ServeError (_response, http.StatusInternalServerError, _error)
			return
		}
	}
	if _data == nil {
		_server.ServeError (_response, http.StatusNotFound, nil)
		return
	}
	
	{
		_key := fmt.Sprintf ("%s:%s", NamespaceDataContentType, _fingerprint)
		if _value, _error := _server.cdbReader.Get ([]byte (_key)); _error == nil {
			if _value != nil {
				_responseHeaders.Set ("Content-Type", string (_value))
			}
		} else {
			_server.ServeError (_response, http.StatusInternalServerError, _error)
			return
		}
	}
	
	if _server.debug {
		log.Printf ("[dd] [b15f3cad]  serving for `%s`...\n", _path)
	}
	
	_response.WriteHeader (http.StatusOK)
	_response.Write (_data)
}




func (_server *server) ServeRedirect (_response http.ResponseWriter, _status uint, _urlRaw string) () {
	var _url string
	if _url_0, _error := url.Parse (_urlRaw); _error == nil {
		_url = _url_0.String ()
	} else {
		_server.ServeError (_response, http.StatusInternalServerError, _error)
		return
	}
	_response.Header () .Set ("Content-Type", "text/plain; charset=utf-8")
	_response.Header () .Set ("Location", _url)
	_response.WriteHeader (int (_status))
	_response.Write ([]byte (fmt.Sprintf ("[%d] %s", _status, _url)))
}


func (_server *server) ServeError (_response http.ResponseWriter, _status uint, _error error) () {
	_response.Header () .Set ("Content-Type", "text/plain; charset=utf-8")
	_response.WriteHeader (int (_status))
	_response.Write ([]byte (fmt.Sprintf ("[%d]", _status)))
	LogError (_error, "")
}




func main () () {
	Main (main_0)
}


func main_0 () (error) {
	
	
	var _bind string
	var _archive string
	var _debug bool
	
	{
		_flags := flag.NewFlagSet ("cdb-server", flag.ContinueOnError)
		
		_bind_0 := _flags.String ("bind", "", "<ip>:<port>")
		_archive_0 := _flags.String ("archive", "", "<path>")
		_debug_0 := _flags.Bool ("debug", false, "")
		
		FlagsParse (_flags, 0, 0)
		
		_bind = *_bind_0
		_archive = *_archive_0
		_debug = *_debug_0
		
		if _bind == "" {
			AbortError (nil, "[eefe1a38]  expected bind address argument!")
		}
		if _archive == "" {
			AbortError (nil, "[eefe1a38]  expected archive file argument!")
		}
	}
	
	
	var _cdbReader *cdb.CDB
	if _cdbReader_0, _error := cdb.Open (_archive); _error == nil {
		_cdbReader = _cdbReader_0
	} else {
		AbortError (_error, "[85234ba0]  failed opening archive!")
	}
	
	
	_server := & server {
			cdbReader : _cdbReader,
			debug : _debug,
		}
	
	if _error := http.ListenAndServe (_bind, _server); _error != nil {
		AbortError (_error, "[44f45c67]  failed starting server!")
	}
	
	return nil
}

