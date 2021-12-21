

package server


import _ "embed"




//go:embed files/tls/server-rsa-public.pem
var DefaultTlsRsaCertificatePublic []byte

//go:embed files/tls/server-rsa-private.pem
var DefaultTlsRsaCertificatePrivate []byte


//go:embed files/tls/server-ed25519-public.pem
var DefaultTlsEd25519CertificatePublic []byte

//go:embed files/tls/server-ed25519-private.pem
var DefaultTlsEd25519CertificatePrivate []byte

