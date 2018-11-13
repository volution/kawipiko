

package common


import "flag"
import "fmt"
import "log"
import "os"




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
			AbortError (nil, fmt.Sprintf ("[b2ac515b]  expected exactly `%d` positional arguments!", _argumentsMin))
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

