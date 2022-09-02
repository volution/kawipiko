
package main


import "fmt"
import "os"


import "github.com/volution/kawipiko/cmd/server"
import "github.com/volution/kawipiko/cmd/archiver"
import "github.com/volution/kawipiko/cmd/version"




func main () () {
	
	if len (os.Args) >= 2 {
		
		_command := os.Args[1]
		os.Args = os.Args[1:]
		
		switch _command {
			
			case "version", "--version", "-v" :
				version.Main ("kawipiko-wrapper")
			
			case "server" :
				server.Main ()
			
			case "archiver" :
				archiver.Main ()
			
			default :
				fmt.Fprintf (os.Stderr, "[!!] [44887671] unknown command: `%s`;  aborting!\n", _command)
		}
		
	} else {
		fmt.Fprintf (os.Stderr, "[!!] [3628f38a]  expected command: `server`, `archiver`, or `version`;  aborting!\n")
	}
}

