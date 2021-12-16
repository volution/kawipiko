

package server


import _ "embed"




//go:embed files/tls/server-public.pem
var DefaultTlsCertificatePublic []byte

//go:embed files/tls/server-private.pem
var DefaultTlsCertificatePrivate []byte

