

package server


import "runtime"
import "net"


import "github.com/valyala/tcplisten"
import "github.com/valyala/fasthttp/reuseport"




func listenTcp (_endpoint string) (net.Listener, error) {
	
	if runtime.GOOS == "android" {
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

