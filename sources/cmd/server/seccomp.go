
package server


import "log"

import "github.com/volution/kawipiko/lib/seccomp"

import . "github.com/volution/kawipiko/lib/common"




func seccompApplyPhase1 () () {
	log.Printf ("[ii] [d53cf86e]  [seccomp.]  applying Linux seccomp filter (phase 1)...\n")
	if _error := seccomp.ApplyServer (); _error != nil {
		AbortError (_error, "[58d1492b]  failed to apply Linux seccomp filter (phase 1)!")
	}
}


func seccompApplyPhase2 () () {
	log.Printf ("[ii] [a338ddaf]  [seccomp.]  applying Linux seccomp filter (phase 2)...\n")
	if _error := seccomp.ApplyServer (); _error != nil {
		AbortError (_error, "[68283e68]  failed to apply Linux seccomp filter (phase 2)!")
	}
}


func seccompApplyPhase3 () () {
	log.Printf ("[ii] [a319ff21]  [seccomp.]  applying Linux seccomp filter (phase 3)...\n")
	if _error := seccomp.ApplyServer (); _error != nil {
		AbortError (_error, "[7c5a0f44]  failed to apply Linux seccomp filter (phase 3)!")
	}
}




var seccompSupported = seccomp.Supported

