

package server


import _ "embed"

import . "github.com/volution/kawipiko/lib/common"




var HeartbeatContentType = MimeTypeText
var HeartbeatContentEncoding = "identity"
var HeartbeatDataOk = []byte ("OK\n")
var HeartbeatDataNok = []byte ("NOK\n")

var DummyContentType = MimeTypeText
var DummyContentEncoding = "identity"
var DummyData = []byte ("hello world!\n")


var FaviconContentType = "image/vnd.microsoft.icon"
var FaviconContentEncoding = "identity"

//go:embed files/favicon.ico
var FaviconData []byte

