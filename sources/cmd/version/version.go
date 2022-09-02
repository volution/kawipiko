

package version


import "bytes"
import "fmt"
import "io"
import "os"


import . "github.com/volution/kawipiko/lib/common"
import . "github.com/volution/kawipiko/embedded"




func Version (_executableName string, _executable string, _stream io.Writer) (error) {
	
	if _executable == "<os.Executable>" {
		if _executable_0, _error := os.Executable (); _error == nil {
			_executable = _executable_0
		} else {
			return _error
		}
	}
	
	_buffer := bytes.NewBuffer (nil)
	
	fmt.Fprintf (_buffer, "* tool          : %s\n", _executableName)
	fmt.Fprintf (_buffer, "* version       : %s\n", BUILD_VERSION)
	if _executable != "" {
		fmt.Fprintf (_buffer, "* executable    : %s\n", _executable)
	}
	fmt.Fprintf (_buffer, "* build target  : %s, %s-%s, %s, %s\n", BUILD_TARGET, BUILD_TARGET_OS, BUILD_TARGET_ARCH, BUILD_COMPILER_VERSION, BUILD_COMPILER_TYPE)
	fmt.Fprintf (_buffer, "* build number  : %s, %s\n", BUILD_NUMBER, BUILD_TIMESTAMP)
	fmt.Fprintf (_buffer, "* code & issues : %s\n", PROJECT_URL)
	fmt.Fprintf (_buffer, "* sources git   : %s\n", BUILD_GIT_HASH)
	fmt.Fprintf (_buffer, "* sources hash  : %s\n", BUILD_SOURCES_HASH)
	fmt.Fprintf (_buffer, "* uname node    : %s\n", UNAME_NODE)
	fmt.Fprintf (_buffer, "* uname system  : %s, %s, %s\n", UNAME_SYSTEM, UNAME_RELEASE, UNAME_MACHINE)
	fmt.Fprintf (_buffer, "* uname hash    : %s\n", UNAME_FINGERPRINT)
	
	if _, _error := _buffer.WriteTo (_stream); _error != nil {
		return _error
	}
	
	return nil
}




func Main (_executableName string, _special string) () {
	
	switch _special {
		
		case "version" :
			if _error := Version (_executableName, "<os.Executable>", os.Stdout); _error != nil {
				AbortError (_error, "[74bfa815]  unexpected error!")
			}
		
		case "sources.md5" :
			if _, _error := os.Stdout.Write (StringToBytes (BuildSourcesMd5)); _error != nil {
				AbortError (_error, "[c8f6b6c9]  unexpected error!")
			}
		
		case "sources.cpio" :
			if _, _error := os.Stdout.Write (BuildSourcesCpioGz); _error != nil {
				AbortError (_error, "[292d8865]  unexpected error!")
			}
		
		case "sbom.txt" :
			if _, _error := os.Stdout.Write (StringToBytes (SbomTxt)); _error != nil {
				AbortError (_error, "[e8f9540a]  unexpected error!")
			}
		
		case "sbom.json" :
			if _, _error := os.Stdout.Write (StringToBytes (SbomJson)); _error != nil {
				AbortError (_error, "[9f9eb7ff]  unexpected error!")
			}
		
		default :
			panic ("[65a6f35f]")
	}
}

