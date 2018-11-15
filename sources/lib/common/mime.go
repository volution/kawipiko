

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


const MimeTypeRaw = "application/octet-stream"


var MimeTypesByExtension = map[string]string {
		
		"txt" : MimeTypeText,
		"csv" : MimeTypeCsv,
		"tsv" : MimeTypeCsv,
		
		"html" : MimeTypeHtml,
		"css" : MimeTypeCss,
		"js" : MimeTypeJs,
		"json" : MimeTypeJson,
		
		"xml" : MimeTypeXml,
		"xslt" : MimeTypeXml,
		"xhtml" : MimeTypeXhtml,
		
	}

