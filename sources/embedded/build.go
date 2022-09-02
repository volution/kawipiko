

package embedded


import _ "embed"




//go:embed build/version.txt
var buildVersion string

//go:embed build/number.txt
var buildNumber string

//go:embed build/timestamp.txt
var buildTimestamp string


//go:embed build/sources.hash
var buildSourcesHash string

//go:embed build/sources.md5
var BuildSourcesMd5 string

//go:embed build/sources.cpio.gz
var BuildSourcesCpioGz []byte

