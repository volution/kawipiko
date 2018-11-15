

package common


import "log"
import "os"
import "regexp"




func AbortError (_error error, _message string) () {
	LogError (_error, _message)
	log.Printf ("[!!] [89251d36]  aborting!\n")
	os.Exit (1)
}




func LogError (_error error, _message string) () {
	
	if _message != "#" {
		if (_message == "") && (_error != nil) {
			_message = "[70d7e7c6]  unexpected error encountered!";
		}
		if _message != "" {
			log.Printf ("[ee] %s\n", _message)
		}
	}
	
	if _error != nil {
		_errorString := _error.Error ()
		if logErrorMessageProper.MatchString (_errorString) {
			log.Printf ("[ee] %s\n", _errorString)
		} else {
			log.Printf ("[ee] [c776ae31]  %q\n", _errorString)
			log.Printf ("[ee] [ddd6baae]  %#v\n", _error)
		}
	}
}


var logErrorMessageProper = regexp.MustCompile (`\A\[[0-9a-f]{8}\] [^\n]+\z`)

