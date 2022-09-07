
package main


import "fmt"
import "os"


import "github.com/volution/kawipiko/cmd/server"
import "github.com/volution/kawipiko/cmd/archiver"
import "github.com/volution/kawipiko/cmd/version"




func main () {
	
	if len (os.Args) >= 2 {
		
		_command := os.Args[1]
		os.Args = os.Args[1:]
		
		switch _command {
			
			case "version", "--version", "-v" :
				version.Main ("kawipiko-wrapper", "version")
				return
			
			case "--sources-md5" :
				version.Main ("kawipiko-wrapper", "sources.md5")
				return
			
			case "--sources-cpio" :
				version.Main ("kawipiko-wrapper", "sources.cpio")
				return
			
			case "--sbom-text", "--sbom-txt", "--sbom" :
				version.Main ("kawipiko-wrapper", "sbom.txt")
				return
			
			case "--sbom-json" :
				version.Main ("kawipiko-wrapper", "sbom.json")
				return
			
			case "server" :
				server.Main ()
				return
			
			case "archiver" :
				archiver.Main ()
				return
			
			default :
				fmt.Fprintf (os.Stderr, "[!!] [44887671] unknown command: `%s`;  aborting!\n", _command)
		}
		
	} else {
		fmt.Fprintf (os.Stderr, "[!!] [3628f38a]  expected command: `server`, `archiver`, or `version`;  aborting!\n")
	}
}

