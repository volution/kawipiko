

package main


import "bytes"
import "flag"
import "io"
import "io/ioutil"
import "log"
import "net/http"
import "os"
import "time"

// import "github.com/colinmarc/cdb"
import cdb "github.com/cipriancraciun/go-cdb-lib"

import "github.com/valyala/fasthttp"

import . "github.com/cipriancraciun/go-cdb-http/lib/common"
import . "github.com/cipriancraciun/go-cdb-http/lib/server"




type server struct {
	cdbReader *cdb.CDB
	debug bool
}


func (_server *server) HandleHTTP (_context *fasthttp.RequestCtx) () {
	_uri := _context.URI ()
	_request := &_context.Request
	_requestHeaders := &_request.Header
	_response := &_context.Response
	_responseHeaders := &_response.Header
	
	_keyBuffer := [1024]byte {}
	_pathNewBuffer := [1024]byte {}
	_timestampBuffer := [128]byte {}
	
	_timestamp := time.Now ()
	_timestampHttp := _timestamp.AppendFormat (_timestampBuffer[:0], http.TimeFormat)
	
	// _responseHeaders.Set ("Content-Security-Policy", "upgrade-insecure-requests")
	_responseHeaders.Set ("Referrer-Policy", "strict-origin-when-cross-origin")
	_responseHeaders.Set ("X-Frame-Options", "SAMEORIGIN")
	_responseHeaders.Set ("X-content-type-Options", "nosniff")
	_responseHeaders.Set ("X-XSS-Protection", "1; mode=block")
	
	_responseHeaders.SetBytesV ("date", _timestampHttp)
	_responseHeaders.SetBytesV ("last-modified", _timestampHttp)
	_responseHeaders.Set ("age", "0")
	
	_method := _requestHeaders.Method ()
	_path := _uri.Path ()
	_pathLen := len (_path)
	_pathIsRoot := _pathLen == 1
	_pathHasSlash := !_pathIsRoot && (_path[_pathLen - 1] == '/')
	
	if ! bytes.Equal ([]byte (http.MethodGet), _method) {
		// log.Printf ("[ww] [bce7a75b] invalid method `%s` for `%s`!\n", _method, _path)
		_server.ServeError (_context, http.StatusMethodNotAllowed, nil)
		return
	}
	if (_pathLen == 0) || (_path[0] != '/') {
		// log.Printf ("[ww] [fa6b1923] invalid path `%s`!\n", _path)
		_server.ServeError (_context, http.StatusBadRequest, nil)
		return
	}
	
	var _fingerprint []byte
	{
		_path_0 := _path
		if _pathHasSlash {
			_path_0 = _path[: _pathLen - 1]
		}
		_found : for _, _namespace := range []string {NamespaceFilesContent, NamespaceFoldersContent, NamespaceFoldersEntries} {
			_key := _keyBuffer[:0]
			_key = append (_key, _namespace ...)
			_key = append (_key, ':')
			_key = append (_key, _path_0 ...)
			if _value, _error := _server.cdbReader.Get (_key); _error == nil {
				if _value != nil {
					_fingerprint = _value
					if (_namespace == NamespaceFoldersContent) || (_namespace == NamespaceFoldersEntries) {
						if !_pathHasSlash {
							_pathNew := _pathNewBuffer[:0]
							_pathNew = append (_pathNew, _path_0 ...)
							_pathNew = append (_pathNew, '/')
							_server.ServeRedirect (_context, http.StatusTemporaryRedirect, _pathNew)
							return
						}
					}
					if _namespace == NamespaceFoldersEntries {
						for _, _indexName := range []string {
								"index.html", "index.htm",
								"index.xhtml", "index.xht",
								"index.txt",
								"index.json",
								"index.xml",
						} {
							_key := _keyBuffer[:0]
							_key = append (_key, NamespaceFilesContent ...)
							_key = append (_key, '/')
							if !_pathIsRoot {
								_key = append (_key, _path_0 ...)
							}
							_key = append (_key, '/')
							_key = append (_key, _indexName ...)
							if _value, _error := _server.cdbReader.Get (_key); _error == nil {
								_fingerprint = _value
								break _found
							} else {
								_server.ServeError (_context, http.StatusInternalServerError, _error)
								return
							}
						}
					}
					break _found
				}
			} else {
				_server.ServeError (_context, http.StatusInternalServerError, _error)
				return
			}
		}
	}
	
	if _fingerprint == nil {
		if ! bytes.Equal ([]byte ("/favicon.ico"), _path) {
			// log.Printf ("[ww] [7416f61d]  not found `%s`!\n", _path)
			_server.ServeError (_context, http.StatusNotFound, nil)
		} else {
			_data, _dataContentType := FaviconData ()
			_responseHeaders.Set ("content-type", _dataContentType)
			_responseHeaders.Set ("content-encoding", "identity")
			_responseHeaders.Set ("etag", "f00f5f99bb3d45ef9806547fe5fe031a")
			_responseHeaders.Set ("cache-control", "public, immutable, max-age=3600")
			_response.SetStatusCode (http.StatusOK)
			_response.SetBody (_data)
		}
		return
	}
	
	var _data []byte
	{
		_key := _keyBuffer[:0]
		_key = append (_key, NamespaceDataContent ...)
		_key = append (_key, ':')
		_key = append (_key, _fingerprint ...)
		if _value, _error := _server.cdbReader.Get (_key); _error == nil {
			if _value != nil {
				_data = _value
			} else {
				// log.Printf ("[ee] [0165c193]  missing data content for `%s`!\n", _fingerprint)
				_server.ServeError (_context, http.StatusInternalServerError, nil)
				return
			}
		} else {
			_server.ServeError (_context, http.StatusInternalServerError, _error)
			return
		}
	}
	
	var _metadata [][2]string
	{
		_key := _keyBuffer[:0]
		_key = append (_key, NamespaceDataMetadata ...)
		_key = append (_key, ':')
		_key = append (_key, _fingerprint ...)
		if _value, _error := _server.cdbReader.Get (_key); _error == nil {
			if _value != nil {
				if _metadata_0, _error := MetadataDecode (_value); _error == nil {
					_metadata = _metadata_0
				} else {
					_server.ServeError (_context, http.StatusInternalServerError, _error)
					return
				}
			} else {
				// log.Printf ("[ee] [e8702411]  missing data metadata for `%s`!\n", _fingerprint)
				_server.ServeError (_context, http.StatusInternalServerError, nil)
				return
			}
		} else {
			_server.ServeError (_context, http.StatusInternalServerError, _error)
			return
		}
	}
	
	if _server.debug {
		// log.Printf ("[dd] [b15f3cad]  serving for `%s`...\n", _path)
	}
	
	for _, _metadata := range _metadata {
		_responseHeaders.Set (_metadata[0], _metadata[1])
	}
	_responseHeaders.Set ("cache-control", "public, immutable, max-age=3600")
	
	_response.SetStatusCode (http.StatusOK)
	_response.SetBody (_data)
}




func (_server *server) ServeRedirect (_context *fasthttp.RequestCtx, _status uint, _path []byte) () {
	_response := &_context.Response
	_responseHeaders := &_response.Header
	
	_responseHeaders.Set ("content-type", MimeTypeText)
	_responseHeaders.Set ("content-encoding", "identity")
	_responseHeaders.Set ("etag", "7aa652d8d607b85808c87c1c2105fbb5")
	_responseHeaders.Set ("cache-control", "public, immutable, max-age=3600")
	_responseHeaders.SetBytesV ("location", _path)
	
	_response.SetStatusCode (int (_status))
	// _response.SetBody ([]byte (fmt.Sprintf ("[%d] %s", _status, _path)))
}


func (_server *server) ServeError (_context *fasthttp.RequestCtx, _status uint, _error error) () {
	_response := &_context.Response
	_responseHeaders := &_response.Header
	
	_responseHeaders.Set ("content-type", MimeTypeText)
	_responseHeaders.Set ("content-encoding", "identity")
	_responseHeaders.Set ("cache-control", "no-cache")
	
	_response.SetStatusCode (int (_status))
	// _response.SetBody ([]byte (fmt.Sprintf ("[%d]", _status)))
	
	LogError (_error, "")
}




func main () () {
	Main (main_0)
}


func main_0 () (error) {
	
	
	var _bind string
	var _archive string
	var _preload bool
	var _debug bool
	
	{
		_flags := flag.NewFlagSet ("cdb-http-server", flag.ContinueOnError)
		
		_bind_0 := _flags.String ("bind", "", "<ip>:<port>")
		_archive_0 := _flags.String ("archive", "", "<path>")
		_preload_0 := _flags.Bool ("preload", false, "")
		_debug_0 := _flags.Bool ("debug", false, "")
		
		FlagsParse (_flags, 0, 0)
		
		_bind = *_bind_0
		_archive = *_archive_0
		_preload = *_preload_0
		_debug = *_debug_0
		
		if _bind == "" {
			AbortError (nil, "[eefe1a38]  expected bind address argument!")
		}
		if _archive == "" {
			AbortError (nil, "[eefe1a38]  expected archive file argument!")
		}
	}
	
	
	var _cdbFile io.ReaderAt
	{
		var _file *os.File
		if _file_0, _error := os.Open (_archive); _error == nil {
			_file = _file_0
		} else {
			AbortError (_error, "[9e0b5ed3]  failed opening archive!")
		}
		
		if _preload {
			log.Printf ("[ii] [922ca187]  preloading archive...\n")
			if _data, _error := ioutil.ReadAll (_file); _error == nil {
				_cdbFile = bytes.NewReader (_data)
			} else {
				AbortError (_error, "[73039784]  failed preloading archive!")
			}
		} else {
			_cdbFile = _file
		}
	}
	
	var _cdbReader *cdb.CDB
	if _cdbReader_0, _error := cdb.NewWithHasher (_cdbFile, nil); _error == nil {
		_cdbReader = _cdbReader_0
	} else {
		AbortError (_error, "[85234ba0]  failed opening archive!")
	}
	
	
	_server := & server {
			cdbReader : _cdbReader,
			debug : _debug,
		}
	
	log.Printf ("[ii] [f11e4e37]  listening on `http://%s/`;\n", _bind)
	
	if _error := fasthttp.ListenAndServe (_bind, _server.HandleHTTP); _error != nil {
		AbortError (_error, "[44f45c67]  failed starting server!")
	}
	
	return nil
}

