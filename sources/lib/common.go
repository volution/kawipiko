

package lib


import "flag"
import "fmt"
import "log"
import "os"
import "regexp"




func Main (_main func () (error)) () {
	
	log.SetFlags (0)
	
	if _error := _main (); _error == nil {
		os.Exit (0)
	} else {
		AbortError (_error, "#")
	}
}




func FlagsParse (_flags *flag.FlagSet, _argumentsMin uint, _argumentsMax uint) ([]string) {
	
	_arguments := os.Args[1:]
	
	if _error := _flags.Parse (_arguments); _error != nil {
		AbortError (_error, fmt.Sprintf ("[8fae7a93]  failed parsing arguments:  `%v`!", _arguments))
	}
	
	_flagsNArg := uint (_flags.NArg ())
	if _argumentsMin == _argumentsMax {
		if _flagsNArg != _argumentsMin {
			AbortError (nil, fmt.Sprintf ("[c451f1f8]  expected exactly `%d` positional arguments!", _argumentsMin))
		}
	} else {
		if _flagsNArg < _argumentsMin {
			AbortError (nil, fmt.Sprintf ("[c451f1f8]  expected at least `%d` positional arguments!", _argumentsMin))
		}
		if _flagsNArg > _argumentsMax {
			AbortError (nil, fmt.Sprintf ("[fa0a8c22]  expected at most `%d` positional arguments!", _argumentsMax))
		}
	}
	
	return _flags.Args ()
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
		if _matches, _matchesError := regexp.MatchString (`^\[[0-9a-f]{8}\] [^\n]+$`, _errorString); _matchesError == nil {
			if _matches {
				log.Printf ("[ee] %s\n", _errorString)
			} else {
				log.Printf ("[ee] [c776ae31]  %q\n", _errorString)
				log.Printf ("[ee] [ddd6baae]  %#v\n", _error)
			}
		} else {
			log.Printf ("[ee] [609a0410]  %q\n", _errorString)
			log.Printf ("[ee] [2ddce4bf]  %#v\n", _error)
		}
	}
}


func AbortError (_error error, _message string) () {
	LogError (_error, _message)
	log.Printf ("[!!] [89251d36]  aborting!\n")
	os.Exit (1)
}

