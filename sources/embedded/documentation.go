

package embedded


import _ "embed"




//go:embed documentation/sbom.txt
var SbomTxt string

//go:embed documentation/sbom.html
var SbomHtml string

//go:embed documentation/sbom.json
var SbomJson string

