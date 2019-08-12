
package main

import "runtime"
import "runtime/debug"
import "github.com/valyala/fasthttp"
import "github.com/valyala/fasthttp/reuseport"


func main () () {
	
	runtime.GOMAXPROCS (2)
	
	debug.SetGCPercent (50)
	debug.SetMaxThreads (128)
	debug.SetMaxStack (16 * 1024)
	
	_data := []byte ("hello world!\n")
	
	_listener, _error := reuseport.Listen ("tcp4", "127.9.185.194:8080")
	
	if _error != nil { panic (_error) }
	
	_error = fasthttp.Serve (
			_listener,
			func (_context *fasthttp.RequestCtx) () {
				_context.Response.SetBodyRaw (_data)
			})
	
	if _error != nil { panic (_error) }
}

