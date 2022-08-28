

package archiver


import "strings"




var IndexNames = []string {
		
		"_index.html", "_index.htm",
		"_index.xhtml", "_index.xht",
		"_index.txt",
		"_index.json",
		"_index.xml",
		
		"index.html", "index.htm",
		"index.xhtml", "index.xht",
		"index.txt",
		"index.json",
		"index.xml",
	}


var StripSuffixes = []string {
		".html", ".htm",
		".xhtml", ".xht",
		".txt",
	}


var SkipPrefixes = []string {
		".",
		"#",
	}

var SkipSuffixes = []string {
		"~",
		"#",
		".log",
		".tmp",
		".temp",
		".lock",
	}

var SkipInfixes = []string {
		"#",
	}

var SkipNames = []string {
		"Thumbs.db",
		".DS_Store",
	}




func ShouldSkipName (_name string) (bool) {
	for _, _skipName := range SkipNames {
		if _skipName == _name {
			return true
		}
	}
	for _, _skipPrefix := range SkipPrefixes {
		if strings.HasPrefix (_name, _skipPrefix) {
			return true
		}
	}
	for _, _skipSuffix := range SkipSuffixes {
		if strings.HasSuffix (_name, _skipSuffix) {
			return true
		}
	}
	for _, _skipInfix := range SkipInfixes {
		if strings.Contains (_name, _skipInfix) {
			return true
		}
	}
	return false
}

