
package server


import "log"

import "github.com/volution/kawipiko/lib/seccomp"

import . "github.com/volution/kawipiko/lib/common"




var _seccompBaseSyscalls = []string {
		
		"brk",
		"mmap",
		"munmap",
		"madvise",
		"mprotect",
		
		"clone3",
		"getpid",
		"gettid",
		"tgkill",
		"exit_group",
		"sched_yield",
		"nanosleep",
		
		"sigaltstack",
		"rt_sigaction",
		"rt_sigprocmask",
		"rt_sigreturn",
		"restart_syscall",
		
		"futex",
		"set_robust_list",
		"rseq",
		
	}


// NOTE:  While serving.
var _seccompPhase3Syscalls = append ([]string {
		
		"accept4",
		"close",
		"getsockname",
		"getpeername",
		"getsockopt",
		"setsockopt",
		
		"read",
		"write",
		
		"pread64",
		"pwrite64",
		
		"recvmmsg",
		"sendmsg",
		
		"epoll_ctl",
		"epoll_pwait",
		
		"getrandom",
		
		"getrusage",
		
	}, _seccompBaseSyscalls ...)


// NOTE:  While listening.
var _seccompPhase2Syscalls = append ([]string {
		
		"socket",
		"bind",
		"listen",
		
		"pipe2",
		"fcntl",
		
		"epoll_create1",
		
		"seccomp",
		"prctl",
		
	}, _seccompPhase3Syscalls ...)


// NOTE:  While loading.
var _seccompPhase1Syscalls = append ([]string {
		
		"openat",
		"fstat",
		"newfstatat",
		
		"mmap",
		
		"setrlimit",
		
		"seccomp",
		"prctl",
		
	}, _seccompPhase2Syscalls ...)




func seccompApplyPhase1 () () {
	seccompApplied = true
	log.Printf ("[ii] [d53cf86e]  [seccomp.]  applying Linux seccomp filter (phase 1)...\n")
	if _error := seccomp.AllowOnlySyscalls (_seccompPhase1Syscalls); _error != nil {
		AbortError (_error, "[58d1492b]  failed to apply Linux seccomp filter (phase 1)!")
	}
}


func seccompApplyPhase2 () () {
	seccompApplied = true
	log.Printf ("[ii] [a338ddaf]  [seccomp.]  applying Linux seccomp filter (phase 2)...\n")
	if _error := seccomp.AllowOnlySyscalls (_seccompPhase2Syscalls); _error != nil {
		AbortError (_error, "[68283e68]  failed to apply Linux seccomp filter (phase 2)!")
	}
}


func seccompApplyPhase3 () () {
	seccompApplied = true
	log.Printf ("[ii] [a319ff21]  [seccomp.]  applying Linux seccomp filter (phase 3)...\n")
	if _error := seccomp.AllowOnlySyscalls (_seccompPhase3Syscalls); _error != nil {
		AbortError (_error, "[7c5a0f44]  failed to apply Linux seccomp filter (phase 3)!")
	}
}




var seccompApplied = false
var seccompSupported = seccomp.Supported

