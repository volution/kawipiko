

package main


import "bytes"
import "flag"
import "fmt"
import "io"
import "log"
import "net"
import "net/http"
import "os"
import "os/signal"
import "runtime"
import "runtime/debug"
import "runtime/pprof"
import "strconv"
import "sync"
import "syscall"
import "time"
import "unsafe"

import "github.com/colinmarc/cdb"

import "github.com/valyala/fasthttp"
import "github.com/valyala/fasthttp/reuseport"

import . "github.com/volution/kawipiko/lib/common"
import . "github.com/volution/kawipiko/lib/server"




type server struct {
	httpServer *fasthttp.Server
	cdbReader *cdb.CDB
	cachedFileFingerprints map[string][]byte
	cachedDataMeta map[string][]byte
	cachedDataContent map[string][]byte
	securityHeadersEnabled bool
	securityHeadersTls bool
	debug bool
	dummy bool
}




func (_server *server) Serve (_context *fasthttp.RequestCtx) () {
	
	// _request := (*fasthttp.Request) (NoEscape (unsafe.Pointer (&_context.Request)))
	_requestHeaders := (*fasthttp.RequestHeader) (NoEscape (unsafe.Pointer (&_context.Request.Header)))
	_response := (*fasthttp.Response) (NoEscape (unsafe.Pointer (&_context.Response)))
	_responseHeaders := (*fasthttp.ResponseHeader) (NoEscape (unsafe.Pointer (&_context.Response.Header)))
	
	_keyBuffer := [1024]byte {}
	_pathBuffer := [1024]byte {}
	
	_method := _requestHeaders.Method ()
	
	_path := _pathBuffer[:0]
	_path = append (_path, _requestHeaders.RequestURI () ...)
	if _pathLimit := bytes.IndexByte (_path, '?'); _pathLimit > 0 {
		_path = _path[: _pathLimit]
	}
	// FIXME:  Decode path according to `decodeArgAppendNoPlus`!
	
	_pathLen := len (_path)
	_pathIsRoot := _pathLen == 1
	_pathHasSlash := !_pathIsRoot && (_path[_pathLen - 1] == '/')
	
	if ! bytes.Equal (StringToBytes (http.MethodGet), _method) {
		log.Printf ("[ww] [bce7a75b]  invalid method `%s` for `%s`!\n", _requestHeaders.Method (), _requestHeaders.RequestURI ())
		_server.ServeError (_context, http.StatusMethodNotAllowed, nil, true)
		return
	}
	if (_pathLen == 0) || (_path[0] != '/') {
		log.Printf ("[ww] [fa6b1923]  invalid path `%s`!\n", _requestHeaders.RequestURI ())
		_server.ServeError (_context, http.StatusBadRequest, nil, true)
		return
	}
	
	if bytes.HasPrefix (_path, StringToBytes ("/__/")) {
		if bytes.Equal (_path, StringToBytes ("/__/heartbeat")) || bytes.HasPrefix (_path, StringToBytes ("/__/heartbeat/")) {
			_server.ServeStatic (_context, http.StatusOK, HeartbeatDataOk, HeartbeatContentType, HeartbeatContentEncoding, false)
			return
		} else if bytes.Equal (_path, StringToBytes ("/__/about")) {
			_server.ServeStatic (_context, http.StatusOK, AboutBannerData, AboutBannerContentType, AboutBannerContentEncoding, true)
			return
		} else if bytes.HasPrefix (_path, StringToBytes ("/__/errors/banners/")) {
			_code := _path[len ("/__/errors/banners/") :]
			if _code, _error := strconv.Atoi (BytesToString (*NoEscapeBytes (&_code))); _error == nil {
				_banner, _bannerFound := ErrorBannersData[uint (_code)]
				if (_code > 0) && _bannerFound {
					_server.ServeStatic (_context, http.StatusOK, _banner, ErrorBannerContentType, ErrorBannerContentEncoding, true)
					return
				}
			}
			_server.ServeError (_context, http.StatusNotFound, nil, true)
			return
		} else {
			_server.ServeError (_context, http.StatusNotFound, nil, true)
			return
		}
	}
	
	if _server.securityHeadersEnabled {
		if _server.securityHeadersTls {
			_responseHeaders.AddRaw (StringToBytes ("Strict-Transport-Security"), StringToBytes ("max-age=31536000"))
			_responseHeaders.AddRaw (StringToBytes ("Content-Security-Policy"), StringToBytes ("upgrade-insecure-requests"))
		}
		{
			_responseHeaders.AddRaw (StringToBytes ("Referrer-Policy"), StringToBytes ("strict-origin-when-cross-origin"))
			_responseHeaders.AddRaw (StringToBytes ("X-Content-Type-Options"), StringToBytes ("nosniff"))
			_responseHeaders.AddRaw (StringToBytes ("X-XSS-Protection"), StringToBytes ("1; mode=block"))
			_responseHeaders.AddRaw (StringToBytes ("X-Frame-Options"), StringToBytes ("sameorigin"))
		}
	}
	
	var _fingerprints []byte
	
	var _namespaceAndPathSuffixes = [][2]string {
			{NamespaceFilesContent, ""},
			{NamespaceFilesContent, "/"},
			{NamespaceFoldersContent, ""},
		}
	
	if _fingerprints == nil {
		_loop_1 : for _namespaceAndPathSuffixIndex := range _namespaceAndPathSuffixes {
			_namespaceAndPathSuffix := _namespaceAndPathSuffixes[_namespaceAndPathSuffixIndex]
			_namespace := _namespaceAndPathSuffix[0]
			_pathSuffix := _namespaceAndPathSuffix[1]
			
			switch {
				case !_pathIsRoot && !_pathHasSlash :
					break
				case _pathSuffix == "/" :
					continue _loop_1
				case _pathSuffix == "" :
					break
				case _pathSuffix[0] == '/' :
					_pathSuffix = _pathSuffix[1:]
			}
			_pathSuffixHasSlash := (len (_pathSuffix) != 0) && (_pathSuffix[0] == '/')
			
			if _server.cachedFileFingerprints != nil {
				_key := _keyBuffer[:0]
				_key = append (_key, _path ...)
				_key = append (_key, _pathSuffix ...)
				_fingerprints, _ = _server.cachedFileFingerprints[BytesToString (_key)]
			} else {
				_key := _keyBuffer[:0]
				_key = append (_key, _namespace ...)
				_key = append (_key, ':')
				_key = append (_key, _path ...)
				_key = append (_key, _pathSuffix ...)
				if _value, _error := _server.cdbReader.GetWithCdbHash (_key); _error == nil {
					_fingerprints = _value
				} else {
					_server.ServeError (_context, http.StatusInternalServerError, _error, false)
					return
				}
			}
			
			if _fingerprints != nil {
				if ((_namespace == NamespaceFoldersContent) || _pathSuffixHasSlash) && (!_pathIsRoot && !_pathHasSlash) {
					_path = append (_path, '/')
					_server.ServeRedirect (_context, http.StatusTemporaryRedirect, _path, true)
					return
				}
				break _loop_1
			}
		}
	}
	
	if _fingerprints == nil {
		if bytes.Equal (StringToBytes ("/favicon.ico"), _path) {
			_server.ServeStatic (_context, http.StatusOK, FaviconData, FaviconContentType, FaviconContentEncoding, true)
			return
		}
	}
	
	if _fingerprints == nil {
		_loop_2 : for
				_pathLimit := bytes.LastIndexByte (_path, '/');
				_pathLimit >= 0;
				_pathLimit = bytes.LastIndexByte (_path[: _pathLimit], '/') {
			
			if _server.cachedFileFingerprints != nil {
				_key := _keyBuffer[:0]
				_key = append (_key, _path[: _pathLimit] ...)
				_key = append (_key, "/*" ...)
				_fingerprints, _ = _server.cachedFileFingerprints[BytesToString (_key)]
			} else {
				_key := _keyBuffer[:0]
				_key = append (_key, NamespaceFilesContent ...)
				_key = append (_key, ':')
				_key = append (_key, _path[: _pathLimit] ...)
				_key = append (_key, "/*" ...)
				if _value, _error := _server.cdbReader.GetWithCdbHash (_key); _error == nil {
					_fingerprints = _value
				} else {
					_server.ServeError (_context, http.StatusInternalServerError, _error, false)
					return
				}
			}
			
			if _fingerprints != nil {
				break _loop_2
			}
		}
	}
	
	if _fingerprints == nil {
		log.Printf ("[ww] [7416f61d]  not found `%s`!\n", _requestHeaders.RequestURI ())
		_server.ServeError (_context, http.StatusNotFound, nil, true)
		return
	}
	
	if len (_fingerprints) != 129 {
		log.Printf ("[ee] [7ee6c981]  invalid data fingerprints for `%s`!\n", _requestHeaders.RequestURI ())
		_server.ServeError (_context, http.StatusInternalServerError, nil, false)
		return
	}
	_fingerprintContent := _fingerprints[0:64]
	_fingerprintMeta := _fingerprints[65:129]
	
	_responseHeaders.AddRaw (StringToBytes ("Cache-Control"), StringToBytes ("public, immutable, max-age=3600"))
	
	var _data []byte
	if _server.cachedDataContent != nil {
		_data, _ = _server.cachedDataContent[BytesToString (_fingerprintContent)]
	} else {
		_key := _keyBuffer[:0]
		_key = append (_key, NamespaceDataContent ...)
		_key = append (_key, ':')
		_key = append (_key, _fingerprintContent ...)
		if _value, _error := _server.cdbReader.GetWithCdbHash (_key); _error == nil {
			_data = _value
		} else {
			_server.ServeError (_context, http.StatusInternalServerError, _error, false)
			return
		}
	}
	if _data == nil {
		log.Printf ("[ee] [0165c193]  missing data content for `%s`!\n", _requestHeaders.RequestURI ())
		_server.ServeError (_context, http.StatusInternalServerError, nil, false)
		return
	}
	
	var _dataMetaRaw []byte
	if _server.cachedDataMeta != nil {
		_dataMetaRaw, _ = _server.cachedDataMeta[BytesToString (_fingerprintMeta)]
	} else {
		_key := _keyBuffer[:0]
		_key = append (_key, NamespaceDataMetadata ...)
		_key = append (_key, ':')
		_key = append (_key, _fingerprintMeta ...)
		if _value, _error := _server.cdbReader.GetWithCdbHash (_key); _error == nil {
			_dataMetaRaw = _value
		} else {
			_server.ServeError (_context, http.StatusInternalServerError, _error, false)
			return
		}
	}
	if _dataMetaRaw == nil {
		log.Printf ("[ee] [e8702411]  missing data metadata for `%s`!\n", _requestHeaders.RequestURI ())
		_server.ServeError (_context, http.StatusInternalServerError, nil, false)
		return
	}
	
	_responseStatus := http.StatusOK
	_handleHeader := func (_name []byte, _value []byte) {
			if _name[0] != '_' {
				_responseHeaders.AddRaw (_name, _value)
			} else {
				switch BytesToString (_name) {
					case "_Status" :
						if _value, _error := strconv.Atoi (BytesToString (_value)); _error == nil {
							if (_value >= 200) && (_value <= 599) {
								_responseStatus = _value
							} else {
								log.Printf ("[c2f7ec36]  invalid data metadata for `%s`!\n", _requestHeaders.RequestURI ())
								_responseStatus = http.StatusInternalServerError
								}
							} else {
								log.Printf ("[beedae55]  invalid data metadata for `%s`!\n", _requestHeaders.RequestURI ())
								_responseStatus = http.StatusInternalServerError
							}
					default :
						log.Printf ("[7acc7d90]  invalid data metadata for `%s`!\n", _requestHeaders.RequestURI ())
				}
			}
		}
	if _error := MetadataDecodeIterate (_dataMetaRaw, _handleHeader); _error != nil {
		_server.ServeError (_context, http.StatusInternalServerError, _error, false)
		return
	}
	
	if _server.debug {
		log.Printf ("[dd] [b15f3cad]  serving for `%s`...\n", _requestHeaders.RequestURI ())
	}
	
	_response.SetStatusCode (_responseStatus)
	_response.SetBodyRaw (_data)
}




func (_server *server) ServeStatic (_context *fasthttp.RequestCtx, _status uint, _data []byte, _contentType string, _contentEncoding string, _cache bool) () {
	
	_response := (*fasthttp.Response) (NoEscape (unsafe.Pointer (&_context.Response)))
	_responseHeaders := (*fasthttp.ResponseHeader) (NoEscape (unsafe.Pointer (&_context.Response.Header)))
	
	_responseHeaders.AddRaw (StringToBytes ("Content-Type"), StringToBytes (_contentType))
	_responseHeaders.AddRaw (StringToBytes ("Content-Encoding"), StringToBytes (_contentEncoding))
	
	if _cache {
		_responseHeaders.AddRaw (StringToBytes ("Cache-Control"), StringToBytes ("public, immutable, max-age=3600"))
	} else {
		_responseHeaders.AddRaw (StringToBytes ("Cache-Control"), StringToBytes ("no-cache"))
	}
	
	_response.SetStatusCode (int (_status))
	_response.SetBodyRaw (_data)
}


func (_server *server) ServeRedirect (_context *fasthttp.RequestCtx, _status uint, _path []byte, _cache bool) () {
	
	_response := (*fasthttp.Response) (NoEscape (unsafe.Pointer (&_context.Response)))
	_responseHeaders := (*fasthttp.ResponseHeader) (NoEscape (unsafe.Pointer (&_context.Response.Header)))
	
	_responseHeaders.AddRaw (StringToBytes ("Content-Encoding"), StringToBytes ("identity"))
	_responseHeaders.AddRaw (StringToBytes ("Location"), _path)
	
	if _cache {
		_responseHeaders.AddRaw (StringToBytes ("Cache-Control"), StringToBytes ("public, immutable, max-age=3600"))
	} else {
		_responseHeaders.AddRaw (StringToBytes ("Cache-Control"), StringToBytes ("no-cache"))
	}
	
	_responseHeaders.AddRaw (StringToBytes ("Content-Type"), StringToBytes (MimeTypeText))
	_responseHeaders.AddRaw (StringToBytes ("Content-Encoding"), StringToBytes ("identity"))
	
	_response.SetStatusCode (int (_status))
}


func (_server *server) ServeError (_context *fasthttp.RequestCtx, _status uint, _error error, _cache bool) () {
	
	_response := (*fasthttp.Response) (NoEscape (unsafe.Pointer (&_context.Response)))
	_responseHeaders := (*fasthttp.ResponseHeader) (NoEscape (unsafe.Pointer (&_context.Response.Header)))
	
	if _cache {
		_responseHeaders.AddRaw (StringToBytes ("Cache-Control"), StringToBytes ("public, immutable, max-age=3600"))
	} else {
		_responseHeaders.AddRaw (StringToBytes ("Cache-Control"), StringToBytes ("no-cache"))
	}
	
	_responseHeaders.AddRaw (StringToBytes ("Content-Type"), StringToBytes (ErrorBannerContentType))
	_responseHeaders.AddRaw (StringToBytes ("Content-Encoding"), StringToBytes (ErrorBannerContentEncoding))
	
	if _banner, _bannerFound := ErrorBannersData[_status]; _bannerFound {
		_response.SetBodyRaw (_banner)
	}
	
	_response.SetStatusCode (int (_status))
	
	LogError (_error, "")
}


func (_server *server) ServeDummy (_context *fasthttp.RequestCtx) () {
	if true {
		_server.ServeStatic (_context, http.StatusOK, DummyData, DummyContentType, DummyContentEncoding, false)
	} else {
		ServeDummyRaw (_context)
	}
}

func ServeDummyRaw (_context *fasthttp.RequestCtx) () {
	_context.Response.SetBodyRaw (DummyData)
	_context.Response.SetStatusCode (200)
}





func main () () {
	Main (main_0)
}


func main_0 () (error) {
	
	
	var _bind string
	var _archivePath string
	var _archiveInmem bool
	var _archiveMmap bool
	var _archivePreload bool
	var _indexPaths bool
	var _indexDataMeta bool
	var _indexDataContent bool
	var _securityHeadersEnabled bool
	var _securityHeadersTls bool
	var _timeoutDisabled bool
	var _processes uint
	var _threads uint
	var _slave uint
	var _debug bool
	var _dummy bool
	var _isFirst bool
	var _isMaster bool
	
	var _profileCpu string
	var _profileMem string
	
	{
		_flags := flag.NewFlagSet ("kawipiko-server", flag.ContinueOnError)
		
		_flags.Usage = func () () {
			fmt.Fprintf (os.Stderr, "%s",
`
  ====  kawipiko -- blazingly fast static HTTP server  ====

  |  Documentation, issues and sources:
  |      * https://github.com/volution/kawipiko
  |  Authors:
  |      * Ciprian Dorin Craciun
  |          ciprian@volution.ro
  |          ciprian.craciun@gmail.com
  |          https://volution.ro/ciprian
  -----------------------------------------------------------

  kawipiko-server

    --archive <path>
    --archive-inmem      (memory-loaded archive file)
    --archive-mmap       (memory-mapped archive file)
    --archive-preload    (preload archive file)

    --index-all
    --index-paths
    --index-data-meta
    --index-data-content

    --bind <ip>:<port>

    --processes <count>  (of slave processes)
    --threads <count>    (of threads per process)

    --security-headers-tls
    --security-headers-disable
    --timeout-disable

    --profile-cpu <path>
    --profile-mem <path>

    --debug
    --dummy

  ** for details see:
     https://github.com/volution/kawipiko#kawipiko-server

`)
		}
		
		_bind_0 := _flags.String ("bind", "", "")
		_archivePath_0 := _flags.String ("archive", "", "")
		_archiveInmem_0 := _flags.Bool ("archive-inmem", false, "")
		_archiveMmap_0 := _flags.Bool ("archive-mmap", false, "")
		_archivePreload_0 := _flags.Bool ("archive-preload", false, "")
		_indexAll_0 := _flags.Bool ("index-all", false, "")
		_indexPaths_0 := _flags.Bool ("index-paths", false, "")
		_indexDataMeta_0 := _flags.Bool ("index-data-meta", false, "")
		_indexDataContent_0 := _flags.Bool ("index-data-content", false, "")
		_timeoutDisabled_0 := _flags.Bool ("timeout-disable", false, "")
		_securityHeadersTls_0 := _flags.Bool ("security-headers-tls", false, "")
		_securityHeadersDisabled_0 := _flags.Bool ("security-headers-disable", false, "")
		_processes_0 := _flags.Uint ("processes", 0, "")
		_threads_0 := _flags.Uint ("threads", 0, "")
		_slave_0 := _flags.Uint ("slave", 0, "")
		_profileCpu_0 := _flags.String ("profile-cpu", "", "")
		_profileMem_0 := _flags.String ("profile-mem", "", "")
		_debug_0 := _flags.Bool ("debug", false, "")
		_dummy_0 := _flags.Bool ("dummy", false, "")
		
		FlagsParse (_flags, 0, 0)
		
		_bind = *_bind_0
		_archivePath = *_archivePath_0
		_archiveInmem = *_archiveInmem_0
		_archiveMmap = *_archiveMmap_0
		_archivePreload = *_archivePreload_0
		_indexAll := *_indexAll_0
		_indexPaths = _indexAll || *_indexPaths_0
		_indexDataMeta = _indexAll || *_indexDataMeta_0
		_indexDataContent = _indexAll || *_indexDataContent_0
		_securityHeadersTls = *_securityHeadersTls_0
		_securityHeadersEnabled = ! *_securityHeadersDisabled_0
		_timeoutDisabled = *_timeoutDisabled_0
		_processes = *_processes_0
		_threads = *_threads_0
		_slave = *_slave_0
		_debug = *_debug_0
		_dummy = *_dummy_0
		
		_profileCpu = *_profileCpu_0
		_profileMem = *_profileMem_0
		
		if _slave == 0 {
			_isMaster = true
		}
		if _slave <= 1 {
			_isFirst = true
		}
		
		if _bind == "" {
			AbortError (nil, "[6edd9512]  expected bind address argument!")
		}
		
		if !_dummy {
			if _archivePath == "" {
				AbortError (nil, "[eefe1a38]  expected archive file argument!")
			}
			if _archiveInmem && _archiveMmap {
				AbortError (nil, "[a2101041]  archive 'memory-loaded' and 'memory-mapped' are mutually exclusive!")
			}
			if _archiveInmem && _archivePreload {
				log.Printf ("[ww] [3e8a40e4]  archive 'memory-loaded' implies preloading!\n")
				_archivePreload = false
			}
		} else {
			if _isMaster {
				log.Printf ("[ww] [8e014192]  running in dummy mode;  all archive related arguments are ignored!\n")
			}
			_archivePath = ""
			_archiveInmem = false
			_archiveMmap = false
			_archivePreload = false
			_indexAll = false
			_indexPaths = false
			_indexDataMeta = false
			_indexDataContent = false
		}
		
		if (_processes > 1) && ((_profileCpu != "") || (_profileMem != "")) {
			AbortError (nil, "[cd18d250]  multi-process and profiling are mutually exclusive!")
		}
		
		if _processes < 1 {
			_processes = 1
		}
		if _threads < 1 {
			_threads = 1
		}
		
		if _processes > 1024 {
			AbortError (nil, "[45736c1d]  maximum number of allowed processes is 1024!")
		}
		if _threads > 1024 {
			AbortError (nil, "[c5df3c8d]  maximum number of allowed threads is 1024!")
		}
		if (_processes * _threads) > 1024 {
			AbortError (nil, "[b0177488]  maximum number of allowed threads in total is 1024!")
		}
	}
	
	
	runtime.GOMAXPROCS (int (_threads))
	
	debug.SetGCPercent (50)
	debug.SetMaxThreads (int (128 * (_threads / 64 + 1)))
	debug.SetMaxStack (16 * 1024)
	
	
	_httpServerReduceMemory := false
	
	if false {
		if _error := syscall.Setrlimit (syscall.RLIMIT_DATA, & syscall.Rlimit { Max : 4 * 1024 * 1024 * 1024 }); _error != nil {
			AbortError (_error, "[f661b4fe]  failed to configure limits!")
		}
	}
	
	
	if _processes > 1 {
		
		log.Printf ("[ii] [06f8c944]  sub-processes starting (`%d` processes with `%d` threads each)...\n", _processes, _threads)
		
		_processesJoin := & sync.WaitGroup {}
		
		_processesPid := make ([]*os.Process, _processes)
		
		_processName := os.Args[0]
		_processArguments := make ([]string, 0, len (os.Args))
		_processArguments = append (_processArguments,
				"--bind", _bind,
			)
		if _archivePath != "" {
			_processArguments = append (_processArguments, "--archive", _archivePath)
		}
		if _archiveInmem {
			_processArguments = append (_processArguments, "--archive-inmem")
		}
		if _archiveMmap {
			_processArguments = append (_processArguments, "--archive-mmap")
		}
		if _archivePreload {
			_processArguments = append (_processArguments, "--archive-preload")
		}
		if _indexPaths {
			_processArguments = append (_processArguments, "--index-paths")
		}
		if _indexDataMeta {
			_processArguments = append (_processArguments, "--index-data-meta")
		}
		if _indexDataContent {
			_processArguments = append (_processArguments, "--index-data-content")
		}
		if _securityHeadersTls {
			_processArguments = append (_processArguments, "--security-headers-tls")
		}
		if !_securityHeadersEnabled {
			_processArguments = append (_processArguments, "--security-headers-disable")
		}
		if _timeoutDisabled {
			_processArguments = append (_processArguments, "--timeout-disable")
		}
		if _debug {
			_processArguments = append (_processArguments, "--debug")
		}
		if _dummy {
			_processArguments = append (_processArguments, "--dummy")
		}
		_processArguments = append (_processArguments, "--threads", fmt.Sprintf ("%d", _threads))
		
		_processAttributes := & os.ProcAttr {
				Env : []string {},
				Files : []*os.File {
						os.Stdin,
						os.Stdout,
						os.Stderr,
					},
				Sys : nil,
			}
		
		for _processIndex, _ := range _processesPid {
			_processArguments := append ([]string { _processName, "--slave", fmt.Sprintf ("%d", _processIndex + 1) }, _processArguments ...)
			if _processPid, _error := os.StartProcess (_processName, _processArguments, _processAttributes); _error == nil {
				_processesJoin.Add (1)
				_processesPid[_processIndex] = _processPid
				if _debug {
					log.Printf ("[ii] [63cb22f8]  sub-process `%d` started (with `%d` threads);\n", _processPid.Pid, _threads)
				}
				go func (_index int, _processPid *os.Process) () {
					if _processStatus, _error := _processPid.Wait (); _error == nil {
						if _processStatus.Success () {
							if _debug {
								log.Printf ("[ii] [66b60b81]  sub-process `%d` succeeded;\n", _processPid.Pid)
							}
						} else {
							log.Printf ("[ww] [5d25046b]  sub-process `%d` failed:  `%s`;  ignoring!\n", _processPid.Pid, _processStatus)
						}
					} else {
						LogError (_error, fmt.Sprintf ("[f1bfc927]  failed waiting for sub-process `%d`;  ignoring!", _processPid.Pid))
					}
					_processesPid[_processIndex] = nil
					_processesJoin.Done ()
				} (_processIndex, _processPid)
			} else {
				LogError (_error, "[8892b34d]  failed starting sub-process;  ignoring!")
			}
		}
		
		{
			_signals := make (chan os.Signal, 32)
			signal.Notify (_signals, syscall.SIGINT, syscall.SIGTERM)
			go func () () {
				for {
					_signal := <- _signals
					if _debug {
						log.Printf ("[ii] [a9243ecb]  signaling sub-processes...\n")
					}
					for _, _processPid := range _processesPid {
						if _processPid != nil {
							if _error := _processPid.Signal (_signal); _error != nil {
								LogError (_error, fmt.Sprintf ("[ab681164]  failed signaling sub-process `%d`;  ignoring!", _processPid.Pid))
							}
						}
					}
				}
			} ()
		}
		
		_processesJoin.Wait ()
		
		if _debug {
			log.Printf ("[ii] [b949bafc]  sub-processes terminated;\n")
		}
		
		return nil
	}
	
	
	if _isMaster {
		log.Printf ("[ii] [6602a54a]  starting (with `%d` threads)...\n", _threads)
	}
	
	
	var _cdbReader *cdb.CDB
	if _archivePath != "" {
		
		if _debug || _isFirst {
			log.Printf ("[ii] [3b788396]  opening archive file `%s`...\n", _archivePath)
		}
		
		var _cdbFile *os.File
		if _cdbFile_0, _error := os.Open (_archivePath); _error == nil {
			_cdbFile = _cdbFile_0
		} else {
			AbortError (_error, "[9e0b5ed3]  failed opening archive file!")
		}
		
		var _cdbFileSize int
		{
			var _cdbFileSize_0 int64
			if _cdbFileStat, _error := _cdbFile.Stat (); _error == nil {
				_cdbFileSize_0 = _cdbFileStat.Size ()
			} else {
				AbortError (_error, "[0ccf0a3b]  failed opening archive file!")
			}
			if _cdbFileSize_0 < 1024 {
				AbortError (nil, "[6635a2a8]  failed opening archive:  file is too small (or empty)!")
			}
			if _cdbFileSize_0 >= (2 * 1024 * 1024 * 1024) {
				AbortError (nil, "[545bf6ce]  failed opening archive:  file is too large!")
			}
			_cdbFileSize = int (_cdbFileSize_0)
		}
		
		if _archivePreload {
			if _debug {
				log.Printf ("[ii] [13f4ebf7]  preloading archive file...\n")
			}
			_buffer := [16 * 1024]byte {}
			_loop : for {
				switch _, _error := _cdbFile.Read (_buffer[:]); _error {
					case io.EOF :
						break _loop
					case nil :
						continue _loop
					default :
						AbortError (_error, "[a1c3b922]  failed preloading archive file...\n")
				}
			}
		}
		
		if _archiveInmem || _archiveMmap {
			
			var _cdbData []byte
			
			if _archiveInmem {
				
				if _debug {
					log.Printf ("[ii] [216e584b]  opening memory-loaded archive...\n")
				}
				
				_cdbData = make ([]byte, _cdbFileSize)
				if _, _error := io.ReadFull (_cdbFile, _cdbData); _error != nil {
					AbortError (_error, "[73039784]  failed loading archive file!")
				}
				
			} else if _archiveMmap {
				
				if _debug {
					log.Printf ("[ii] [f47fae8a]  opening memory-mapped archive...\n")
				}
				
				if _cdbData_0, _error := syscall.Mmap (int (_cdbFile.Fd ()), 0, int (_cdbFileSize), syscall.PROT_READ, syscall.MAP_SHARED); _error == nil {
					_cdbData = _cdbData_0
				} else {
					AbortError (_error, "[c0e2632c]  failed mapping archive file!")
				}
				
				if _archivePreload {
					if _debug {
						log.Printf ("[ii] [d96b06c9]  preloading memory-loaded archive...\n")
					}
					_buffer := [16 * 1024]byte {}
					_bufferOffset := 0
					for {
						if _bufferOffset == _cdbFileSize {
							break
						}
						_bufferOffset += copy (_buffer[:], _cdbData[_bufferOffset:])
					}
				}
				
			} else {
				panic ("e4fffcd8")
			}
			
			if _error := _cdbFile.Close (); _error != nil {
				AbortError (_error, "[5e0449c2]  failed closing archive file!")
			}
			
			if _cdbReader_0, _error := cdb.NewFromBufferWithHasher (_cdbData, nil); _error == nil {
				_cdbReader = _cdbReader_0
			} else {
				AbortError (_error, "[27e4813e]  failed opening archive!")
			}
			
		} else {
			
			if _debug || _isFirst {
				log.Printf ("[ww] [dd697a66]  using `read`-based archive (with significant performance impact)!\n")
			}
			
			if _cdbReader_0, _error := cdb.NewFromReaderWithHasher (_cdbFile, nil); _error == nil {
				_cdbReader = _cdbReader_0
			} else {
				AbortError (_error, "[35832022]  failed opening archive!")
			}
			
		}
		
		if _schemaVersion, _error := _cdbReader.GetWithCdbHash ([]byte (NamespaceSchemaVersion)); _error == nil {
			if _schemaVersion == nil {
				AbortError (nil, "[09316866]  missing archive schema version!")
			} else if string (_schemaVersion) != CurrentSchemaVersion {
				AbortError (nil, "[e6482cf7]  invalid archive schema version!")
			}
		} else {
			AbortError (_error, "[87cae197]  failed opening archive!")
		}
	}
	
	
	var _cachedFileFingerprints map[string][]byte
	if _indexPaths {
		_cachedFileFingerprints = make (map[string][]byte, 128 * 1024)
	}
	var _cachedDataMeta map[string][]byte
	if _indexDataMeta {
		_cachedDataMeta = make (map[string][]byte, 128 * 1024)
	}
	var _cachedDataContent map[string][]byte
	if _indexDataContent {
		_cachedDataContent = make (map[string][]byte, 128 * 1024)
	}
	
	if _indexPaths || _indexDataMeta || _indexDataContent {
		if _debug {
			log.Printf ("[ii] [fa5338fd]  indexing archive...\n")
		}
		if _filesIndex, _error := _cdbReader.GetWithCdbHash ([]byte (NamespaceFilesIndex)); _error == nil {
			if _filesIndex != nil {
				_keyBuffer := [1024]byte {}
				for {
					_offset := bytes.IndexByte (_filesIndex, '\n')
					if _offset == 0 {
						continue
					}
					if _offset == -1 {
						break
					}
					_filePath := _filesIndex[: _offset]
					_filesIndex = _filesIndex[_offset + 1 :]
					var _fingerprints []byte
					var _fingerprintContent []byte
					var _fingerprintMeta []byte
					{
						_key := _keyBuffer[:0]
						_key = append (_key, NamespaceFilesContent ...)
						_key = append (_key, ':')
						_key = append (_key, _filePath ...)
						if _fingerprints_0, _error := _cdbReader.GetWithCdbHash (_key); _error == nil {
							if _fingerprints_0 != nil {
								_fingerprints = _fingerprints_0
								_fingerprintContent = _fingerprints[0:64]
								_fingerprintMeta = _fingerprints[65:129]
							} else {
								AbortError (_error, "[460b3cf1]  failed indexing archive!")
							}
						} else {
							AbortError (_error, "[216f2075]  failed indexing archive!")
						}
					}
					if _indexPaths {
						_cachedFileFingerprints[BytesToString (_filePath)] = _fingerprints
					}
					if _indexDataMeta {
						if _, _wasCached := _cachedDataMeta[BytesToString (_fingerprintMeta)]; !_wasCached {
							_key := _keyBuffer[:0]
							_key = append (_key, NamespaceDataMetadata ...)
							_key = append (_key, ':')
							_key = append (_key, _fingerprintMeta ...)
							if _dataMeta, _error := _cdbReader.GetWithCdbHash (_key); _error == nil {
								if _dataMeta != nil {
									_cachedDataMeta[BytesToString (_fingerprintMeta)] = _dataMeta
								} else {
									AbortError (_error, "[6df556bf]  failed indexing archive!")
								}
							} else {
								AbortError (_error, "[0d730134]  failed indexing archive!")
							}
						}
					}
					if _indexDataContent {
						if _, _wasCached := _cachedDataContent[BytesToString (_fingerprintContent)]; !_wasCached {
							_key := _keyBuffer[:0]
							_key = append (_key, NamespaceDataContent ...)
							_key = append (_key, ':')
							_key = append (_key, _fingerprintContent ...)
							if _dataContent, _error := _cdbReader.GetWithCdbHash (_key); _error == nil {
								if _dataContent != nil {
									_cachedDataContent[BytesToString (_fingerprintContent)] = _dataContent
								} else {
									AbortError (_error, "[4e27fe46]  failed indexing archive!")
								}
							} else {
								AbortError (_error, "[532845ad]  failed indexing archive!")
							}
						}
					}
				}
			} else {
				log.Printf ("[ww] [30314f31]  missing archive files index;  ignoring!\n")
			}
		} else {
			AbortError (_error, "[82299b3d]  failed indexing arcdive!")
		}
	}
	
	if _indexPaths && _indexDataMeta && _indexDataContent {
		if _error := _cdbReader.Close (); _error == nil {
			_cdbReader = nil
		} else {
			AbortError (_error, "[d7aa79e1]  failed closing archive!")
		}
	}
	
	
	_server := & server {
			httpServer : nil,
			cdbReader : _cdbReader,
			cachedFileFingerprints : _cachedFileFingerprints,
			cachedDataMeta : _cachedDataMeta,
			cachedDataContent : _cachedDataContent,
			securityHeadersTls : _securityHeadersTls,
			securityHeadersEnabled : _securityHeadersEnabled,
			debug : _debug,
			dummy : _dummy,
		}
	
	
	if _profileCpu != "" {
		log.Printf ("[ii] [70c210f3]  profiling CPU to `%s`...\n", _profileCpu)
		_stream, _error := os.Create (_profileCpu)
		if _error != nil {
			AbortError (_error, "[fd4e0009]  failed opening CPU profile!")
		}
		_error = pprof.StartCPUProfile (_stream)
		if _error != nil {
			AbortError (_error, "[ac721629]  failed starting CPU profile!")
		}
		defer pprof.StopCPUProfile ()
	}
	if _profileMem != "" {
		log.Printf ("[ii] [9196ee90]  profiling MEM to `%s`...\n", _profileMem)
		_stream, _error := os.Create (_profileMem)
		if _error != nil {
			AbortError (_error, "[907d08b5]  failed opening MEM profile!")
		}
		_profile := pprof.Lookup ("heap")
		defer func () () {
			runtime.GC ()
			if _profile != nil {
				if _error := _profile.WriteTo (_stream, 0); _error != nil {
					AbortError (_error, "[4b1e5112]  failed writing MEM profile!")
				}
			} else {
				AbortError (nil, "[385dc8f0]  failed loading MEM profile!")
			}
			_stream.Close ()
		} ()
	}
	
	
	_httpServer := & fasthttp.Server {
			
			Name : "kawipiko",
			Handler : _server.Serve,
			GetOnly : true,
			
			NoDefaultServerHeader : true,
			NoDefaultContentType : true,
			DisableHeaderNamesNormalizing : true,
			
			Concurrency : 4 * 1024 + 128,
			MaxRequestsPerConn : 256 * 1024,
			
			ReadBufferSize : 16 * 1024,
			WriteBufferSize : 16 * 1024,
			MaxRequestBodySize : 16 * 1024,
			
			ReadTimeout : 30 * time.Second,
			WriteTimeout : 30 * time.Second,
			IdleTimeout : 360 * time.Second,
			
			TCPKeepalive : true,
			TCPKeepalivePeriod : 60 * time.Second,
			
			ReduceMemoryUsage : _httpServerReduceMemory,
			
		}
	
	if _timeoutDisabled {
		_httpServer.ReadTimeout = 0
		_httpServer.WriteTimeout = 0
		_httpServer.IdleTimeout = 0
	}
	
	if _dummy {
		_httpServer.Handler = _server.ServeDummy
	}
	
	_server.httpServer = _httpServer
	
	
	{
		_signals := make (chan os.Signal, 32)
		signal.Notify (_signals, syscall.SIGINT, syscall.SIGTERM)
		go func () () {
			<- _signals
			if _debug {
				log.Printf ("[ii] [691cb695]  shutingdown...\n")
			}
			_server.httpServer.Shutdown ()
		} ()
	}
	
	
	if _debug || _isFirst {
		log.Printf ("[ii] [f11e4e37]  listening on `http://%s/`;\n", _bind)
	}
	
	var _httpListener net.Listener
	if _httpListener_0, _error := reuseport.Listen ("tcp4", _bind); _error == nil {
		_httpListener = _httpListener_0
	} else {
		AbortError (_error, "[d5f51e9f]  failed starting listener!")
	}
	
	if _error := _httpServer.Serve (_httpListener); _error != nil {
		AbortError (_error, "[44f45c67]  failed executing server!")
	}
	
	
	if _debug {
		defer log.Printf ("[ii] [a49175db]  done!\n")
	}
	return nil
}

