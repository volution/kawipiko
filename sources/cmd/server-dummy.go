
package main


import "runtime"
import "runtime/debug"
import "time"

import "github.com/valyala/fasthttp"
import "github.com/valyala/fasthttp/reuseport"




func main () () {
	
	
	runtime.GOMAXPROCS (2)
	
	debug.SetGCPercent (50)
	debug.SetMaxThreads (128)
	debug.SetMaxStack (16 * 1024)
	
	
	_listener, _error := reuseport.Listen ("tcp4", "127.9.185.194:8080")
	if _error != nil { panic (_error) }
	
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
			
		//	ReadTimeout : 30 * time.Second,
		//	WriteTimeout : 30 * time.Second,
		//	IdleTimeout : 360 * time.Second,
			
			TCPKeepalive : true,
			TCPKeepalivePeriod : 60 * time.Second,
			
			ReduceMemoryUsage : false,
			KeepHijackedConns : true,
			
			ErrorHandler : nil,
			ConnState : nil,
			HeaderReceived : nil,
			
			Logger : nil,
			LogAllErrors : true,
			
		}
	
	_error = _server.Serve (_listener)
	if _error != nil { panic (_error) }
}




func serve (_context *fasthttp.RequestCtx) () {
	_context.Response.Header.SetRaw (serveMeta)
	_context.Response.SetBodyRaw (serveData)
}

var serveMeta = []byte ("HTTP/1.1 200 OK\r\nContent-Length: 13\r\n\r\n")
var serveData = []byte ("hello world!\n")

