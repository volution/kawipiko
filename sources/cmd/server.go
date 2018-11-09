

package main


import "bytes"
import "flag"
import "fmt"
import "io"
import "io/ioutil"
import "log"
import "net/http"
import "net/url"
import "os"
import "time"

// import "github.com/colinmarc/cdb"
import cdb "github.com/cipriancraciun/go-cdb-lib"

import . "../lib"




type server struct {
	cdbReader *cdb.CDB
	debug bool
}


func (_server *server) ServeHTTP (_response http.ResponseWriter, _request *http.Request) () {
	
	_timestamp := time.Now ()
	_timestampHttp := _timestamp.Format (http.TimeFormat)
	
	_responseHeaders := _response.Header ()
	
	// _responseHeaders.Set ("Content-Security-Policy", "upgrade-insecure-requests")
	_responseHeaders.Set ("Referrer-Policy", "strict-origin-when-cross-origin")
	_responseHeaders.Set ("X-Frame-Options", "SAMEORIGIN")
	_responseHeaders.Set ("X-content-type-Options", "nosniff")
	_responseHeaders.Set ("X-XSS-Protection", "1; mode=block")
	
	_responseHeaders.Set ("date", _timestampHttp)
	_responseHeaders.Set ("last-modified", _timestampHttp)
	_responseHeaders.Set ("age", "0")
	
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
		_found : for _, _namespace := range []string {NamespaceFilesContent, NamespaceFoldersContent, NamespaceFoldersEntries} {
			_key := fmt.Sprintf ("%s:%s", _namespace, _path_0)
			if _value, _error := _server.cdbReader.Get ([]byte (_key)); _error == nil {
				if _value != nil {
					_fingerprint = string (_value)
					if (_namespace == NamespaceFoldersContent) || (_namespace == NamespaceFoldersEntries) {
						if (_path == _path_0) && (_path != "/") {
							_server.ServeRedirect (_response, http.StatusTemporaryRedirect, _path + "/")
							return
						}
					}
					if _namespace == NamespaceFoldersEntries {
						for _, _index := range []string {
								"index.html", "index.htm",
								"index.xhtml", "index.xht",
								"index.txt",
								"index.json",
								"index.xml",
						} {
							_pathIndex := _path_0 + "/" + _index
							if _path_0 == "/" {
								_pathIndex = "/" + _index
							}
							_key := fmt.Sprintf ("%s:%s", NamespaceFilesContent, _pathIndex)
							if _value, _error := _server.cdbReader.Get ([]byte (_key)); _error == nil {
								_fingerprint = string (_value)
								break _found
							} else {
								_server.ServeError (_response, http.StatusInternalServerError, _error)
								return
							}
						}
					}
					break _found
				}
			} else {
				_server.ServeError (_response, http.StatusInternalServerError, _error)
				return
			}
		}
	}
	if _fingerprint == "" {
		if _path != "/favicon.ico" {
			log.Printf ("[ww] [7416f61d]  not found for `%s`!", _path)
			_server.ServeError (_response, http.StatusNotFound, nil)
		} else {
			_data, _dataContentType := FaviconData ()
			_responseHeaders.Set ("content-type", _dataContentType)
			_responseHeaders.Set ("content-encoding", "identity")
			_responseHeaders.Set ("etag", "f00f5f99bb3d45ef9806547fe5fe031a")
			_responseHeaders.Set ("cache-control", "public, immutable, max-age=3600")
			_response.WriteHeader (http.StatusOK)
			_response.Write (_data)
		}
		return
	}
	
	var _data []byte
	{
		_key := fmt.Sprintf ("%s:%s", NamespaceDataContent, _fingerprint)
		if _value, _error := _server.cdbReader.Get ([]byte (_key)); _error == nil {
			if _value != nil {
				_data = _value
			} else {
				_server.ServeError (_response, http.StatusInternalServerError, fmt.Errorf ("[0165c193]  missing data content:  `%s`", _fingerprint))
				return
			}
		} else {
			_server.ServeError (_response, http.StatusInternalServerError, _error)
			return
		}
	}
	
	var _metadata [][2]string
	{
		_key := fmt.Sprintf ("%s:%s", NamespaceDataMetadata, _fingerprint)
		if _value, _error := _server.cdbReader.Get ([]byte (_key)); _error == nil {
			if _value != nil {
				if _metadata_0, _error := MetadataDecode (_value); _error == nil {
					_metadata = _metadata_0
				} else {
					_server.ServeError (_response, http.StatusInternalServerError, _error)
					return
				}
			} else {
				_server.ServeError (_response, http.StatusInternalServerError, fmt.Errorf ("[e8702411]  missing data metadata:  `%s`", _fingerprint))
				return
			}
		} else {
			_server.ServeError (_response, http.StatusInternalServerError, _error)
			return
		}
	}
	
	if _server.debug {
		log.Printf ("[dd] [b15f3cad]  serving for `%s`...\n", _path)
	}
	
	for _, _metadata := range _metadata {
		_responseHeaders.Set (_metadata[0], _metadata[1])
	}
	_responseHeaders.Set ("cache-control", "public, immutable, max-age=3600")
	
	_response.WriteHeader (http.StatusOK)
	_response.Write (_data)
}




func (_server *server) ServeRedirect (_response http.ResponseWriter, _status uint, _urlRaw string) () {
	_responseHeaders := _response.Header ()
	
	var _url string
	if _url_0, _error := url.Parse (_urlRaw); _error == nil {
		_url = _url_0.String ()
	} else {
		_server.ServeError (_response, http.StatusInternalServerError, _error)
		return
	}
	
	_responseHeaders.Set ("content-type", MimeTypeText)
	_responseHeaders.Set ("content-encoding", "identity")
	_responseHeaders.Set ("etag", "7aa652d8d607b85808c87c1c2105fbb5")
	_responseHeaders.Set ("cache-control", "public, immutable, max-age=3600")
	_responseHeaders.Set ("location", _url)
	
	_response.WriteHeader (int (_status))
	_response.Write ([]byte (fmt.Sprintf ("[%d] %s", _status, _url)))
}


func (_server *server) ServeError (_response http.ResponseWriter, _status uint, _error error) () {
	_responseHeaders := _response.Header ()
	
	_responseHeaders.Set ("content-type", MimeTypeText)
	_responseHeaders.Set ("content-encoding", "identity")
	_responseHeaders.Set ("cache-control", "no-cache")
	
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
	var _preload bool
	var _debug bool
	
	{
		_flags := flag.NewFlagSet ("cdb-server", flag.ContinueOnError)
		
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
	
	if _debug {
		log.Printf ("[ii] [f11e4e37]  listening on `http://%s/`", _bind)
	}
	
	if _error := http.ListenAndServe (_bind, _server); _error != nil {
		AbortError (_error, "[44f45c67]  failed starting server!")
	}
	
	return nil
}

