

package server


import "encoding/base64"




func FaviconData () ([]byte, string) {
	_data, _ := base64.StdEncoding.DecodeString (FaviconDataBase64)
	return _data, "image/vnd.microsoft.icon"
}




var FaviconDataBase64 = `
AAABAAMAMDAAAAEAIACoJQAANgAAACAgAAABACAAqBAAAN4lAAAQEAAAAQAgAGgEAACGNgAAKAAA
ADAAAABgAAAAAQAgAAAAAAAAJAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAICAgASAgIAdgICAO4CAgEuAgIBOgICA
ToCAgE6AgIBOgICAToCAgE6AgIBOgICAToCAgE6AgIBOgICAToCAgE6AgIBOgICAToCAgE6AgIBO
gICAToCAgE6AgIBOgICAToCAgE6AgIBOgICAToCAgEyAgIA8gICAHICAgAMAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAgIAFgICAQYCA
gKGAgIDcgICA9ICAgPqAgID7gICA+4CAgPuAgID7gICA+4CAgPuAgID7gICA+4CAgPuAgID7gICA
+4CAgPuAgID7gICA+4CAgPuAgID7gICA+4CAgPuAgID7gICA+4CAgPuAgID7gICA+4CAgPqAgID0
gICA2oCAgJeAgIAzgICAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAICAgBOAgICKgICA74CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CA
gP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA
/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgIDkgICAaICAgAUAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgICAFYCAgKeAgID+gICA/4CAgP+AgID/gICA/4CAgP+A
gID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgPyAgID9gICA/4CAgP+AgID/gICA/4CA
gP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA
9ICAgGyAgIABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAgIAGgICAkYCAgP6AgID/
gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA54CAgG2A
gIB1gICA7ICAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CA
gP+AgID/gICA/4CAgP+AgID/gICA/4CAgOeAgIA7AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAACAgIBMgICA8oCAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/
gICA/4CAgP+AgIDogICAVgAAAACAgIACgICAYYCAgO2AgID/gICA/4CAgP+AgID/gICA/4CAgP+A
gID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgICkgICABQAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAICAgAeAgICugICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA
/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgOiAgIBWgYGBAQAAAAAAAAAAgICAAoCAgGGAgIDt
gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+A
gID/gICA/4CAgP+AgIDkgICAJwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAICAgCqAgIDngICA/4CA
gP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA6ICAgFWBgYEBAAAA
AAAAAAAAAAAAAAAAAICAgAKAgIBhgICA7YCAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/
gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID6gICATwAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAICAgE+AgID6gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CA
gP+AgIDogICAVYGBgQEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAgIACgICAYYCAgO2AgID/gICA
/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/
gICAZQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAICAgGOAgID/gICA/4CAgP+AgID/gICA/4CAgP+A
gID/gICA/4CAgP+AgID/gICA/4CAgOiAgIBVgICAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAgICAAoCAgGGAgIDtgICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA
/4CAgP+AgID/gICA/4CAgP+AgID/gICAYQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAICAgGGAgID/
gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA6ICAgFWAgIABAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAICAgAKAgIBhgICA7YCAgP+AgID/gICA/4CA
gP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID3gICARwAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAICAgEmAgID5gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID9
gICAbAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAgICAe4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CA
gP+AgIDagICAHQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAICAgCOAgIDhgICA/4CAgP+AgID/gICA
/4CAgP+AgID/gICA/4CAgP+AgID+gICAgoCAgCOAgIAkgICAJICAgCWAgIAPAAAAAAAAAAAAAAAA
AAAAAICAgBKAgIAlgICAJICAgCSAgIAkgICAjoCAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+A
gID/gICA/4CAgP+AgID/gICA/4CAgP+AgICKgYGBAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAICA
gASAgIChgICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA+ICAgOeAgIDlgICA
5YCAgOeAgIBiAAAAAAAAAAAAAAAAAAAAAICAgHCAgIDogICA5YCAgOWAgIDngICA+YCAgP+AgID/
gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgM+AgIAjAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAgIA9gICA6oCAgP+AgID/gICA/4CAgP+AgID/gICA/4CA
gP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgIBsAAAAAAAAAAAAAAAAAAAAAICAgHyAgID/gICA
/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/
gICA3YCAgEEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAgIACgICAeoCAgPqA
gID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgIBsAAAAAAAA
AAAAAAAAAAAAAICAgHyAgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA
/4CAgP+AgID/gICA/4CAgPuAgIC+gICAOgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAgICAC4CAgI2AgID5gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+A
gID/gICA/4CAgP+AgIBsAAAAAAAAAAAAAAAAAAAAAICAgHyAgID/gICA/4CAgP+AgID/gICA/4CA
gP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID2gICAroCAgF2AgIATAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAICAgAqAgIBsgICA3oCAgP+AgID/
gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgIBsAAAAAAAAAAAAAAAAAAAAAICAgHyA
gID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgIDqgICALwAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAACAgIABgICAJYCAgJeAgID+gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgIBs
AAAAAAAAAAAAAAAAAAAAAICAgHyAgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+A
gID/gICA/4CAgP+AgID4gICARwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAICAgEqAgID6gICA/4CAgP+AgID/gICA
/4CAgP+AgID/gICA/4CAgP+AgIBtAAAAAAAAAAAAAAAAAAAAAICAgHyAgID/gICA/4CAgP+AgID/
gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID5gICASgAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAICA
gDmAgIDygICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgICRgICAEoCAgBGAgIARgICA
FICAgJ+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgIDs
gICAMQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAICAgB2AgIDbgICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CA
gP+AgID1gICA04CAgM6AgIDOgICA1ICAgPeAgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA
/4CAgP+AgID/gICA/4CAgP+AgIC6gICADAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAICAgASAgICkgICA/4CAgP+A
gID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CA
gP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgO6AgIBNAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAACAgIBNgICA9ICAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+A
gID/gICA/4CAgP+AgID/gICA/4CAgP+AgID+gICA/oCAgP+AgID/gICA/4CAgP+AgID/gICA7YCA
gGuAgIADAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAgIAJgICApoCAgP+AgID/gICA/4CAgP+AgID/
gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP6AgICggICAjYCAgNmA
gIDygICA9ICAgOWAgICxgICARoCAgAMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgICA
KYCAgNCAgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/
gICA/4CAgL6AgIAbgICAAYCAgBuAgIA3gICAO4CAgCaAgIAJAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAICAgDqAgIDRgICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA
/4CAgP+AgID/gICA/4CAgP+AgID/gICAv4CAgCkAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAgIApgICAqICA
gPWAgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgPCAgICVgICAHAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAgICACoCAgEyAgIClgICA3ICAgPOAgID6gICA+YCAgPGAgIDYgICAm4CA
gECAgIAGAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAgIAEgICAHYCAgDqA
gIBKgICASYCAgDeAgIAZgICAAwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAP///////wAA////////AAD///////8AAP///////wAA////////
AAD///////8AAP///////wAA////////AAD///////8AAP///////wAA/4AAAAD/AAD+AAAAAH8A
APwAAAAAPwAA+AAGAAAfAAD4AA8AAA8AAPAAH4AADwAA8AA/wAAPAADwAH/gAA8AAPAA//AADwAA
8AH/+AAPAADwA//8AA8AAPAB//gADwAA8AAfgAAfAAD4AB+AAD8AAPwAH4AAfwAA/AAfgAH/AAD/
AB+AA/8AAP+AH4AD/wAA/8AfgAP/AAD/wA8AA/8AAP/AAAAD/wAA/8AAAAf/AAD/4AAAD/8AAP/g
AAAf/wAA//AAD///AAD/+AAf//8AAP/8AD///wAA//8A////AAD///////8AAP///////wAA////
////AAD///////8AAP///////wAA////////AAD///////8AAP///////wAA////////AAD/////
//8AACgAAAAgAAAAQAAAAAEAIAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAICAgAmAgIA6gICAboCAgIKAgICDgICA
g4CAgIOAgICDgICAg4CAgIOAgICDgICAg4CAgIOAgICDgICAg4CAgIOAgICDgICAg4CAgIOAgIB/
gICAXoCAgCGAgIABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAgIAogICApYCAgPCA
gID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CA
gP+AgID/gICA/4CAgP+AgID9gICA24CAgGeAgIAGAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgICA
KoCAgMuAgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgPuAgIDXgICA8oCAgP+AgID/
gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA9YCAgGsAAAAAAAAAAAAAAAAA
AAAAAAAAAICAgAqAgICrgICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID9gICAnICA
gBqAgIBrgICA8YCAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA
3oCAgCgAAAAAAAAAAAAAAAAAAAAAgICAQ4CAgPOAgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/
gICA/YCAgJyAgIARAAAAAICAgAOAgIBrgICA8YCAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+A
gID/gICA/4CAgP+AgID/gICAagAAAAAAAAAAAAAAAAAAAACAgIB7gICA/4CAgP+AgID/gICA/4CA
gP+AgID/gICA/4CAgP2AgICcgICAEQAAAAAAAAAAAAAAAICAgAOAgIBrgICA8YCAgP+AgID/gICA
/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgICPAAAAAAAAAAAAAAAAAAAAAICAgJGAgID/
gICA/4CAgP+AgID/gICA/4CAgP+AgID9gICAnICAgBEAAAAAAAAAAAAAAAAAAAAAAAAAAICAgAOA
gIBrgICA8YCAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgIwAAAAAAAAAAAAA
AAAAAAAAgICAhICAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgK+AgIAPAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAICAgAGAgIB3gICA/oCAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID9
gICAYQAAAAAAAAAAAAAAAAAAAACAgIBSgICA+YCAgP+AgID/gICA/4CAgP+AgID/gICAv4CAgGCA
gIBhgICAUYCAgAkAAAAAAAAAAICAgD2AgIBigICAXoCAgJiAgID9gICA/4CAgP+AgID/gICA/4CA
gP+AgID/gICA/4CAgNCAgIAdAAAAAAAAAAAAAAAAAAAAAICAgBSAgIDDgICA/4CAgP+AgID/gICA
/4CAgP+AgID/gICA/4CAgP+AgIDYgICAGAAAAACBgYEBgICAooCAgP+AgID/gICA/4CAgP+AgID/
gICA/4CAgP+AgID/gICA/4CAgP+AgIDogICATwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAICAgEWA
gIDlgICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgNiAgIAYAAAAAIGBgQGAgICigICA/4CA
gP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgIDzgICAvoCAgEaAgIABAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAgICAAYCAgEmAgIDNgICA/YCAgP+AgID/gICA/4CAgP+AgID/gICA2ICAgBgAAAAA
gYGBAYCAgKKAgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgIKAgIAMAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAICAgBmAgICngICA/4CAgP+AgID/gICA/4CA
gP+AgIDYgICAGAAAAACBgYEBgICAooCAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgID/gICA
fgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAICAgHaAgID/
gICA/4CAgP+AgID/gICA/4CAgOKAgIAzgICAE4CAgBqAgIC1gICA/4CAgP+AgID/gICA/4CAgP+A
gID/gICA/4CAgP+AgIBwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAgICATYCAgPmAgID/gICA/4CAgP+AgID/gICA/YCAgOCAgIDUgICA2YCAgPeAgID/gICA
/4CAgP+AgID/gICA/4CAgP+AgID/gICA54CAgDMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAACAgIAXgICAzYCAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+A
gID/gICA/4CAgP+AgID+gICA/ICAgP+AgID/gICA/4CAgPGAgIBugICAAgAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAgIBhgICA94CAgP+AgID/gICA
/4CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgM2AgIBogICApICAgL2AgICggICASYCAgAUAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAICAgAiA
gICEgICA94CAgP+AgID/gICA/4CAgP+AgID/gICA/4CAgP+AgIDYgICAPwAAAACAgIADgICACICA
gAIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAICAgAiAgIBigICAz4CAgPmAgID/gICA/4CAgP+AgIDugICAp4CAgC8AAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAgIAXgICAToCAgHeAgICAgICAaoCA
gDaAgIAIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA////////////////////////////
/////wAA//gAAB/wAAAP4AMAB+AHgAfgD8ADwB/gA8A/8AfgP+AH4AcAD/AHAB/4BwA//AcAf/4H
AH/+AAB//gAA//8AEf//AD///8B////7//////////////////////////////////8oAAAAEAAA
ACAAAAABACAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AACAgIACgICABoCAgAaAgIAGgICABoCAgAaAgIAGgICABoCAgAaAgIAEAAAAAAAAAAAAAAAAAAAA
AICAgAGAgIBAgICAnoCAgLeAgIC3gICAuYCAgLqAgIC4gICAt4CAgLeAgIC4gICArICAgF2AgIAH
AAAAAAAAAACAgIBBgICA5YCAgP+AgID/gICA/4CAgNmAgICfgICA9ICAgP+AgID/gICA/4CAgP+A
gID1gICAYAAAAACAgIADgICAooCAgP+AgID/gICA/4CAgNeAgIA8gICABICAgHWAgID0gICA/4CA
gP+AgID/gICA/4CAgLGAgIAGgICACICAgLeAgID/gICA/4CAgOOAgIBAAAAAAAAAAACAgIAGgICA
fYCAgPyAgID/gICA/4CAgP+AgICugICABQAAAACAgICEgICA/4CAgP+AgIDpgICAq4CAgE+AgIAF
gICAf4CAgLyAgID7gICA/4CAgP+AgIDvgICAVQAAAAAAAAAAgICAHICAgLGAgID7gICA/4CAgP+A
gIB4gICACICAgMWAgID/gICA/4CAgP+AgIDZgICAToCAgAQAAAAAAAAAAAAAAACAgIAUgICAwICA
gP+AgID/gICAiYCAgCKAgIDNgICA/4CAgP+AgID/gICAr4CAgAQAAAAAAAAAAAAAAAAAAAAAAAAA
AICAgIeAgID/gICA/4CAgO6AgIDdgICA+YCAgPqAgID+gICA84CAgGQAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAACAgIAmgICAzICAgP+AgID/gICA/4CAgOmAgIBrgICAXICAgEOAgIAGAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAICAgCeAgICHgICAtICAgJyAgIBCgICAAgAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgICAAYCAgAWAgIACAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAP//AAD//wAA//8AAOAHAADAAwAAg4EAAIfB
AACDgwAAwwcAAOEHAADgDwAA8H8AAPj/AAD//wAA//8AAP//AAA=
`
