

package server


import "runtime"
import "net"
import "os"


import "github.com/valyala/tcplisten"
import "github.com/valyala/fasthttp/reuseport"




func listenTcp (_endpoint string) (net.Listener, error) {
	
	_useTcpListen := false
	
	if ! _useTcpListen {
		if seccompApplied {
			_useTcpListen = true
		}
	}
	
	if ! _useTcpListen {
		if (runtime.GOOS == "android") {
			_useTcpListen = true
		}
	}
	if ! _useTcpListen {
		if _, _error := os.Stat ("/proc/sys/net/core/somaxconn"); _error != nil {
			_useTcpListen = true
		}
	}
	
	if _useTcpListen {
		var _config = & tcplisten.Config {
				ReusePort : true,
				DeferAccept : true,
				FastOpen : true,
				Backlog : 1024,
			}
		return _config.NewListener ("tcp4", _endpoint)
	}
	
	// FIXME:  Perhaps always use `tcplisten`...
	
	return reuseport.Listen ("tcp4", _endpoint)
}




func listenUdp (_endpoint string) (net.PacketConn, error) {
	return net.ListenPacket ("udp4", _endpoint)
}

