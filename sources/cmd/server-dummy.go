
package main


import "fmt"
import "os"
import "runtime"
import "runtime/debug"
import "strconv"
import "time"

import "github.com/valyala/fasthttp"
import "github.com/valyala/fasthttp/reuseport"




func main () () {
	
	
	runtime.GOMAXPROCS (1)
	
	debug.SetGCPercent (-1)
	debug.SetMaxThreads (128)
	debug.SetMaxStack (16 * 1024)
	
	
	_endpoint := "127.0.0.1:8080"
	_threads := 1
	_timeouts := false
	
	switch len (os.Args) {
		case 1 :
			// NOP
		case 2 :
			_endpoint = os.Args[1]
		case 3 :
			_endpoint = os.Args[1]
			if _threads_0, _error := strconv.Atoi (os.Args[2]); (_error == nil) && (_threads > 0) {
				_threads = _threads_0
			} else {
				panic ("[40396d14]  invalid arguments!")
			}
		default :
			panic ("[60023f00]  invalid arguments!")
	}
	
	if _threads > 1 {
		runtime.GOMAXPROCS (_threads)
		debug.SetMaxThreads (int (128 * (_threads / 64 + 1)))
	}
	
	_listener, _error := reuseport.Listen ("tcp4", _endpoint)
	if _error != nil {
		panic (fmt.Sprintf ("[8c30a625]  failed to listen:  %s", _error))
	}
	
	fmt.Fprintf (os.Stderr, "[ii] [04fa2421]  listening on `http://%s/` (using %d threads)...\n", _endpoint, _threads)
	
	_server := & fasthttp.Server {
			
			Name : "kawipiko",
			Handler : serve,
			GetOnly : true,
			
			NoDefaultServerHeader : true,
			NoDefaultContentType : true,
			NoDefaultDate : true,
			DisableHeaderNamesNormalizing : true,
			
			Concurrency : 16 * 1024 + 128,
			MaxRequestsPerConn : 256 * 1024,
			
			ReadBufferSize : 16 * 1024,
			WriteBufferSize : 16 * 1024,
			MaxRequestBodySize : 16 * 1024,
			
			ReadTimeout : 30 * time.Second,
			WriteTimeout : 30 * time.Second,
			IdleTimeout : 360 * time.Second,
			
			TCPKeepalive : true,
			TCPKeepalivePeriod : 60 * time.Second,
			
			ReduceMemoryUsage : false,
			
			CloseOnShutdown : true,
			DisableKeepalive : false,
			
			ErrorHandler : nil,
			ConnState : nil,
			HeaderReceived : nil,
			
			Logger : nil,
			LogAllErrors : true,
			
		}
	
	if !_timeouts {
		_server.ReadTimeout = 0
		_server.WriteTimeout = 0
		_server.IdleTimeout = 0
	}
	
	_error = _server.Serve (_listener)
	if _error != nil {
		panic (fmt.Sprintf ("[ee9bc0a5]  failed to serve:  %s", _error))
	}
}




func serve (_context *fasthttp.RequestCtx) () {
	_context.Response.SetBodyRaw (serveData)
}

var serveData = []byte ("hello world!\n")

