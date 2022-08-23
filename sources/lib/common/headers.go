

package common


import "fmt"
import "net/http"
import "strings"




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

func CanonicalHeaderValueArrayFromBytes (_value []byte) ([]string) {
	if _canonical, _found := CanonicalHeaderValuesArraysMap[BytesToString (*NoEscapeBytes (&_value))]; _found {
		return _canonical
	} else {
		return []string { string (_value) }
	}
}


func CanonicalHeaderValueRegister (_value string) () {
	CanonicalHeaderValues = append (CanonicalHeaderValues, _value)
	CanonicalHeaderValuesMap[_value] = _value
}




var CanonicalHeaderNamesMap map[string]string
var CanonicalHeaderNamesToKey map[string]uint64
var CanonicalHeaderNamesFromKey map[uint64]string

var CanonicalHeaderValuesMap map[string]string
var CanonicalHeaderValuesToKey map[string]uint64
var CanonicalHeaderValuesFromKey map[uint64]string
var CanonicalHeaderValuesArraysMap map[string][]string




var CanonicalHeaderNames = []string {
		
		// FIXME:  Move this somewhere else!
		"!Status",
		
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
		
		// FIXME:  Move this somewhere else!
		"200", "404",
		
	}




func init () {
	
	
	CanonicalHeaderNamesMap = make (map[string]string, len (CanonicalHeaderNames) * 4)
	CanonicalHeaderNamesToKey = make (map[string]uint64, len (CanonicalHeaderNames) * 4)
	CanonicalHeaderNamesFromKey = make (map[uint64]string, len (CanonicalHeaderNames))
	
	for _index, _header := range CanonicalHeaderNames {
		
		_http := http.CanonicalHeaderKey (_header)
		_toLower := strings.ToLower (_header)
		_toUpper := strings.ToUpper (_header)
		
		if _, _duplicate := CanonicalHeaderNamesMap[strings.ToLower (_header)]; _duplicate {
			panic (fmt.Sprintf ("[f0dffe23]  invalid duplicate header `%s`", _header))
		}
		
		_key, _error := PrepareKey (NamespaceHeaderName, uint64 (_index + 1))
		if _error != nil {
			panic (_error)
		}
		
		CanonicalHeaderNamesMap[_header] = _header
		CanonicalHeaderNamesMap[_toLower] = _header
		CanonicalHeaderNamesMap[_toUpper] = _header
		CanonicalHeaderNamesMap[_http] = _header
		
		CanonicalHeaderNamesToKey[_header] = _key
		CanonicalHeaderNamesToKey[_toLower] = _key
		CanonicalHeaderNamesToKey[_toUpper] = _key
		CanonicalHeaderNamesToKey[_http] = _key
		CanonicalHeaderNamesFromKey[_key] = _header
	}
	
	
	CanonicalHeaderValues = append (CanonicalHeaderValues, MimeTypes ...)
	CanonicalHeaderValues = append (CanonicalHeaderValues, MimeTypesExtras ...)
	
	CanonicalHeaderValuesMap = make (map[string]string, len (CanonicalHeaderValues))
	CanonicalHeaderValuesArraysMap = make (map[string][]string, len (CanonicalHeaderValues))
	CanonicalHeaderValuesToKey = make (map[string]uint64, len (CanonicalHeaderValues))
	CanonicalHeaderValuesFromKey = make (map[uint64]string, len (CanonicalHeaderValues))
	
	for _index, _value := range CanonicalHeaderValues {
		
		_key, _error := PrepareKey (NamespaceHeaderValue, uint64 (_index + 1))
		if _error != nil {
			panic (_error)
		}
		
		CanonicalHeaderValuesMap[_value] = _value
		CanonicalHeaderValuesArraysMap[_value] = []string { _value }
		
		CanonicalHeaderValuesToKey[_value] = _key
		CanonicalHeaderValuesFromKey[_key] = _value
	}
}

