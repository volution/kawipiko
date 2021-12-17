

package server


import "fmt"
import "net/http"
import "strings"


import . "github.com/volution/kawipiko/lib/common"




func CanonicalHeaderNameFromBytes (_header []byte) (string) {
	if _canonical, _found := CanonicalHeaderNamesMap[BytesToString (*NoEscapeBytes (&_header))]; _found {
		return _canonical
	} else {
		return string (_header)
	}
}


func CanonicalHeaderValueFromBytes (_value []byte) (string) {
	if _canonical, _found := CanonicalHeaderValuesMap[BytesToString (*NoEscapeBytes (&_value))]; _found {
		return _canonical
	} else {
		return string (_value)
	}
}


func CanonicalHeaderValueRegister (_value string) () {
	CanonicalHeaderValues = append (CanonicalHeaderValues, _value)
	CanonicalHeaderValuesMap[_value] = _value
}




var CanonicalHeaderNamesMap map[string]string
var CanonicalHeaderValuesMap map[string]string


var CanonicalHeaderNames = []string {
		"Accept",
		"Accept-CH",
		"Accept-CH-Lifetime",
		"Accept-Charset",
		"Accept-Encoding",
		"Accept-Language",
		"Accept-Patch",
		"Accept-Push-Policy",
		"Accept-Ranges",
		"Accept-Signature",
		"Access-Control-Allow-Credentials",
		"Access-Control-Allow-Headers",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Origin",
		"Access-Control-Expose-Headers",
		"Access-Control-Max-Age",
		"Access-Control-Request-Headers",
		"Access-Control-Request-Method",
		"Age",
		"Allow",
		"Alt-Svc",
		"Authorization",
		"Cache-Control",
		"Clear-Site-Data",
		"Connection",
		"Content-DPR",
		"Content-Disposition",
		"Content-Encoding",
		"Content-Id",
		"Content-Language",
		"Content-Length",
		"Content-Location",
		"Content-Range",
		"Content-Security-Policy",
		"Content-Security-Policy-Report-Only",
		"Content-Transfer-Encoding",
		"Content-Type",
		"Cookie",
		"Cross-Origin-Resource-Policy",
		"DNT",
		"DPR",
		"Date",
		"ETag",
		"Early-Data",
		"Expect",
		"Expect-CT",
		"Expires",
		"Feature-Policy",
		"Forwarded",
		"From",
		"Host",
		"If-Match",
		"If-Modified-Since",
		"If-None-Match",
		"If-Range",
		"If-Unmodified-Since",
		"Index",
		"Keep-Alive",
		"Large-Allocation",
		"Last-Event-ID",
		"Last-Modified",
		"Link",
		"Location",
		"Max-Forwards",
		"NEL",
		"Origin",
		"Ping-From",
		"Ping-To",
		"Pragma",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Public-Key-Pins",
		"Public-Key-Pins-Report-Only",
		"Push-Policy",
		"Range",
		"Referer",
		"Referrer-Policy",
		"Report-To",
		"Retry-After",
		"Save-Data",
		"Sec-WebSocket-Accept",
		"Sec-WebSocket-Extensions",
		"Sec-WebSocket-Key",
		"Sec-WebSocket-Protocol",
		"Sec-WebSocket-Version",
		"Server",
		"Server-Timing",
		"Set-Cookie",
		"Signature",
		"Signed-Headers",
		"SourceMap",
		"Strict-Transport-Security",
		"TE",
		"Timing-Allow-Origin",
		"Tk",
		"Trailer",
		"Transfer-Encoding",
		"Upgrade",
		"Upgrade-Insecure-Requests",
		"User-Agent",
		"Vary",
		"Via",
		"Viewport-Width",
		"WWW-Authenticate",
		"Warning",
		"Width",
		"X-Content-Type-Options",
		"X-DNS-Prefetch-Control",
		"X-Download-Options",
		"X-Forwarded-By",
		"X-Forwarded-For",
		"X-Forwarded-Host",
		"X-Forwarded-Proto",
		"X-Frame-Options",
		"X-Imforwards",
		"X-Permitted-Cross-Domain-Policies",
		"X-Pingback",
		"X-Powered-By",
		"X-Requested-With",
		"X-Robots-Tag",
		"X-UA-Compatible",
		"X-XSS-Protection",
	}




var CanonicalHeaderValues = []string {
		
		"public, immutable, max-age=3600",
		"no-store, max-age=0",
		
		"identity",
		"deflate",
		"gzip",
		"br",
		
		"max-age=31536000",
		"upgrade-insecure-requests",
		
		"strict-origin-when-cross-origin",
		"nosniff",
		"1; mode=block",
		"sameorigin",
		
	}




func init () {
	
	
	CanonicalHeaderNamesMap = make (map[string]string, len (CanonicalHeaderNames) * 4)
	for _, _header := range CanonicalHeaderNames {
		
		_http := http.CanonicalHeaderKey (_header)
		_toLower := strings.ToLower (_header)
		_toUpper := strings.ToUpper (_header)
		
		if _, _duplicate := CanonicalHeaderNamesMap[strings.ToLower (_header)]; _duplicate {
			panic (fmt.Sprintf ("[f0dffe23]  invalid duplicate header `%s`", _header))
		}
		
		CanonicalHeaderNamesMap[_header] = _header
		CanonicalHeaderNamesMap[_toLower] = _header
		CanonicalHeaderNamesMap[_toUpper] = _header
		CanonicalHeaderNamesMap[_http] = _header
	}
	
	
	CanonicalHeaderValues = append (CanonicalHeaderValues, MimeTypes ...)
	
	CanonicalHeaderValuesMap = make (map[string]string, len (CanonicalHeaderValues))
	for _, _value := range CanonicalHeaderValues {
		CanonicalHeaderValuesMap[_value] = _value
	}
}

