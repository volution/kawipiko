

package common




const MimeTypeText = "text/plain; charset=utf-8"
const MimeTypeCsv = "text/csv; charset=utf-8"

const MimeTypeHtml = "text/html; charset=utf-8"
const MimeTypeCss = "text/css; charset=utf-8"
const MimeTypeJs = "application/javascript; charset=utf-8"
const MimeTypeJson = "application/json; charset=utf-8"

const MimeTypeXml = "application/xml; charset=utf-8"
const MimeTypeXslt = "application/xslt+xml; charset=utf-8"
const MimeTypeXhtml = "application/xhtml+xml; charset=utf-8"

// NOTE:  Based on: https://developer.mozilla.org/en-US/docs/Web/Media/Formats/Image_types
const MimeTypeSvg = "image/svg+xml; charset=utf-8"
const MimeTypePng = "image/png"
const MimeTypeJpeg = "image/jpeg"
const MimeTypeWebp = "image/webp"
const MimeTypeGif = "image/gif"
const MimeTypeIco = "image/x-icon"
const MimeTypeBmp = "image/bmp"
const MimeTypeTiff = "image/tiff"
const MimeTypeApng = "image/apng"

const MimeTypeAvif = "image/avif"
const MimeTypeAvifSequence = "image/avif-sequence"
const MimeTypeHeic = "image/heic"
const MimeTypeHeicSequence = "image/heic-sequence"
const MimeTypeHeif = "image/heif"
const MimeTypeHeifSequence = "image/heif-sequence"

const MimeTypeOtf = "font/otf"
const MimeTypeTtf = "font/ttf"
const MimeTypeWoff = "font/woff"
const MimeTypeWoff2 = "font/woff2"

const MimeTypeWav = "audio/wav"
const MimeTypeMidi = "audio/midi"
const MimeTypeAac = "audio/aac"
const MimeTypeOpus = "audio/opus"
const MimeTypeAudioMpeg = "audio/mpeg"
const MimeTypeAudioWebm = "audio/webm"
const MimeTypeAudioOgg = "audio/ogg"

const MimeTypeAvi = "video/x-msvideo"
const MimeTypeMp4 = "video/mp4"
const MimeTypeVideoMpeg = "video/mpeg"
const MimeTypeVideoWebm = "video/webm"
const MimeTypeVideoOgg = "video/ogg"

const MimeTypePdf = "application/pdf"
const MimeTypePs = "application/postscript"
const MimeTypeIcs = "text/calendar"

const MimeTypeZip = "application/zip"

const MimeTypeRaw = "application/octet-stream"


var MimeTypesByExtension = map[string]string {
		
		// NOTE:  Based on: https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
		
		"txt" : MimeTypeText,
		"csv" : MimeTypeCsv,
		"tsv" : MimeTypeCsv,
		
		"html" : MimeTypeHtml,
		"css" : MimeTypeCss,
		"js" : MimeTypeJs,
		"mjs" : MimeTypeJs,
		"json" : MimeTypeJson,
		
		"xml" : MimeTypeXml,
		"xslt" : MimeTypeXml,
		"xhtml" : MimeTypeXhtml,
		
		"svg" : MimeTypeSvg,
		"png" : MimeTypePng,
		"jpeg" : MimeTypeJpeg,
		"jpg" : MimeTypeJpeg,
		"webp" : MimeTypeWebp,
		"gif" : MimeTypeGif,
		"ico" : MimeTypeIco,
		"cur" : MimeTypeIco,
		"bmp" : MimeTypeBmp,
		"tiff" : MimeTypeTiff,
		"tif" : MimeTypeTiff,
		"apng" : MimeTypeApng,
		
		"avif" : MimeTypeAvif,
		"avifs" : MimeTypeAvifSequence,
		"heic" : MimeTypeHeic,
		"heics" : MimeTypeHeicSequence,
		"heif" : MimeTypeHeif,
		"heifs" : MimeTypeHeifSequence,
		
		"otf" : MimeTypeOtf,
		"ttf" : MimeTypeTtf,
		"woff" : MimeTypeWoff,
		"woff2" : MimeTypeWoff2,
		
		"wav" : MimeTypeWav,
		"mid" : MimeTypeMidi,
		"midi" : MimeTypeMidi,
		"aac" : MimeTypeAac,
		"opus" : MimeTypeOpus,
		"mp3" : MimeTypeAudioMpeg,
		"weba" : MimeTypeAudioWebm,
		"oga" : MimeTypeAudioOgg,
		
		"avi" : MimeTypeAvi,
		"mp4" : MimeTypeMp4,
		"mpeg" : MimeTypeVideoMpeg,
		"webm" : MimeTypeVideoWebm,
		"ogv" : MimeTypeVideoOgg,
		
		"pdf" : MimeTypePdf,
		"ps" : MimeTypePs,
		"ics" : MimeTypeIcs,
		
		"zip" : MimeTypeZip,
		
	}

