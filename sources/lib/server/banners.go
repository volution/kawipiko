

package server


import _ "embed"

import . "github.com/volution/kawipiko/lib/common"




var AboutBannerContentType = MimeTypeText
var AboutBannerContentEncoding = "identity"

//go:embed files/about.txt
var AboutBannerData []byte




var ErrorBannerContentType = MimeTypeText
var ErrorBannerContentEncoding = "identity"

var ErrorBannersData = map[uint][]byte {
		
		100 : ErrorBanner100Data,
		101 : ErrorBanner101Data,
		102 : ErrorBanner102Data,
		103 : ErrorBanner103Data,
		200 : ErrorBanner200Data,
		201 : ErrorBanner201Data,
		202 : ErrorBanner202Data,
		203 : ErrorBanner203Data,
		204 : ErrorBanner204Data,
		205 : ErrorBanner205Data,
		206 : ErrorBanner206Data,
		207 : ErrorBanner207Data,
		208 : ErrorBanner208Data,
		226 : ErrorBanner226Data,
		300 : ErrorBanner300Data,
		301 : ErrorBanner301Data,
		302 : ErrorBanner302Data,
		303 : ErrorBanner303Data,
		304 : ErrorBanner304Data,
		305 : ErrorBanner305Data,
		306 : ErrorBanner306Data,
		307 : ErrorBanner307Data,
		308 : ErrorBanner308Data,
		400 : ErrorBanner400Data,
		401 : ErrorBanner401Data,
		402 : ErrorBanner402Data,
		403 : ErrorBanner403Data,
		404 : ErrorBanner404Data,
		405 : ErrorBanner405Data,
		406 : ErrorBanner406Data,
		407 : ErrorBanner407Data,
		408 : ErrorBanner408Data,
		409 : ErrorBanner409Data,
		410 : ErrorBanner410Data,
		411 : ErrorBanner411Data,
		412 : ErrorBanner412Data,
		413 : ErrorBanner413Data,
		414 : ErrorBanner414Data,
		415 : ErrorBanner415Data,
		416 : ErrorBanner416Data,
		417 : ErrorBanner417Data,
		421 : ErrorBanner421Data,
		422 : ErrorBanner422Data,
		423 : ErrorBanner423Data,
		424 : ErrorBanner424Data,
		425 : ErrorBanner425Data,
		426 : ErrorBanner426Data,
		428 : ErrorBanner428Data,
		429 : ErrorBanner429Data,
		431 : ErrorBanner431Data,
		451 : ErrorBanner451Data,
		500 : ErrorBanner500Data,
		501 : ErrorBanner501Data,
		502 : ErrorBanner502Data,
		503 : ErrorBanner503Data,
		504 : ErrorBanner504Data,
		505 : ErrorBanner505Data,
		506 : ErrorBanner506Data,
		507 : ErrorBanner507Data,
		508 : ErrorBanner508Data,
		510 : ErrorBanner510Data,
		511 : ErrorBanner511Data,
		
	}


//go:embed files/errors/100.txt
var ErrorBanner100Data []byte

//go:embed files/errors/101.txt
var ErrorBanner101Data []byte

//go:embed files/errors/102.txt
var ErrorBanner102Data []byte

//go:embed files/errors/103.txt
var ErrorBanner103Data []byte

//go:embed files/errors/200.txt
var ErrorBanner200Data []byte

//go:embed files/errors/201.txt
var ErrorBanner201Data []byte

//go:embed files/errors/202.txt
var ErrorBanner202Data []byte

//go:embed files/errors/203.txt
var ErrorBanner203Data []byte

//go:embed files/errors/204.txt
var ErrorBanner204Data []byte

//go:embed files/errors/205.txt
var ErrorBanner205Data []byte

//go:embed files/errors/206.txt
var ErrorBanner206Data []byte

//go:embed files/errors/207.txt
var ErrorBanner207Data []byte

//go:embed files/errors/208.txt
var ErrorBanner208Data []byte

//go:embed files/errors/226.txt
var ErrorBanner226Data []byte

//go:embed files/errors/300.txt
var ErrorBanner300Data []byte

//go:embed files/errors/301.txt
var ErrorBanner301Data []byte

//go:embed files/errors/302.txt
var ErrorBanner302Data []byte

//go:embed files/errors/303.txt
var ErrorBanner303Data []byte

//go:embed files/errors/304.txt
var ErrorBanner304Data []byte

//go:embed files/errors/305.txt
var ErrorBanner305Data []byte

//go:embed files/errors/306.txt
var ErrorBanner306Data []byte

//go:embed files/errors/307.txt
var ErrorBanner307Data []byte

//go:embed files/errors/308.txt
var ErrorBanner308Data []byte

//go:embed files/errors/400.txt
var ErrorBanner400Data []byte

//go:embed files/errors/401.txt
var ErrorBanner401Data []byte

//go:embed files/errors/402.txt
var ErrorBanner402Data []byte

//go:embed files/errors/403.txt
var ErrorBanner403Data []byte

//go:embed files/errors/404.txt
var ErrorBanner404Data []byte

//go:embed files/errors/405.txt
var ErrorBanner405Data []byte

//go:embed files/errors/406.txt
var ErrorBanner406Data []byte

//go:embed files/errors/407.txt
var ErrorBanner407Data []byte

//go:embed files/errors/408.txt
var ErrorBanner408Data []byte

//go:embed files/errors/409.txt
var ErrorBanner409Data []byte

//go:embed files/errors/410.txt
var ErrorBanner410Data []byte

//go:embed files/errors/411.txt
var ErrorBanner411Data []byte

//go:embed files/errors/412.txt
var ErrorBanner412Data []byte

//go:embed files/errors/413.txt
var ErrorBanner413Data []byte

//go:embed files/errors/414.txt
var ErrorBanner414Data []byte

//go:embed files/errors/415.txt
var ErrorBanner415Data []byte

//go:embed files/errors/416.txt
var ErrorBanner416Data []byte

//go:embed files/errors/417.txt
var ErrorBanner417Data []byte

//go:embed files/errors/421.txt
var ErrorBanner421Data []byte

//go:embed files/errors/422.txt
var ErrorBanner422Data []byte

//go:embed files/errors/423.txt
var ErrorBanner423Data []byte

//go:embed files/errors/424.txt
var ErrorBanner424Data []byte

//go:embed files/errors/425.txt
var ErrorBanner425Data []byte

//go:embed files/errors/426.txt
var ErrorBanner426Data []byte

//go:embed files/errors/428.txt
var ErrorBanner428Data []byte

//go:embed files/errors/429.txt
var ErrorBanner429Data []byte

//go:embed files/errors/431.txt
var ErrorBanner431Data []byte

//go:embed files/errors/451.txt
var ErrorBanner451Data []byte

//go:embed files/errors/500.txt
var ErrorBanner500Data []byte

//go:embed files/errors/501.txt
var ErrorBanner501Data []byte

//go:embed files/errors/502.txt
var ErrorBanner502Data []byte

//go:embed files/errors/503.txt
var ErrorBanner503Data []byte

//go:embed files/errors/504.txt
var ErrorBanner504Data []byte

//go:embed files/errors/505.txt
var ErrorBanner505Data []byte

//go:embed files/errors/506.txt
var ErrorBanner506Data []byte

//go:embed files/errors/507.txt
var ErrorBanner507Data []byte

//go:embed files/errors/508.txt
var ErrorBanner508Data []byte

//go:embed files/errors/510.txt
var ErrorBanner510Data []byte

//go:embed files/errors/511.txt
var ErrorBanner511Data []byte

