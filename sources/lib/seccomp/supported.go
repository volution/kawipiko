

//go:build linux && amd64 && seccomp




package seccomp


import "fmt"
import "syscall"


import "github.com/seccomp/libseccomp-golang"




func init () {
	
	Supported = true
	
	if false {
		
		_filter, _error := seccomp.NewFilter (seccomp.ActLog)
		if _error != nil {
			panic (_error)
		}
		
		if _error = _filter.Load (); _error != nil {
			panic (_error)
		}
	}
}




func AllowOnlySyscalls (_syscalls []string) (error) {
	
	_architectures := []seccomp.ScmpArch {
			seccomp.ArchNative,
			seccomp.ArchX86,
			seccomp.ArchX32,
			seccomp.ArchAMD64,
		}
	
	_fallbackAction := seccomp.ActKill
	switch 0 {
		case 1 :
			_fallbackAction = seccomp.ActErrno.SetReturnCode (int16 (syscall.EPERM))
		case 2 :
			_fallbackAction = seccomp.ActLog
	}
	
	_filter, _error := seccomp.NewFilter (_fallbackAction)
	if _error != nil {
		return _error
	}
	
	for _, _architecture := range _architectures {
		if _error := _filter.AddArch (_architecture); _error != nil {
			return _error
		}
	}
	
	for _, _syscall := range _syscalls {
		
		var _sc_syscall seccomp.ScmpSyscall
		switch {
			
			case _syscall[0] == '!' :
				continue
			
			default :
				if _sc_syscall_0, _error := seccomp.GetSyscallFromNameByArch (_syscall, seccomp.ArchNative); _error == nil {
					_sc_syscall = _sc_syscall_0
				} else {
					return fmt.Errorf ("[5cf9cd60]  failed resolving syscall `%s`:  %w", _syscall, _error)
				}
		}
		
		if _error := _filter.AddRule (_sc_syscall, seccomp.ActAllow); _error != nil {
			return _error
		}
	}
	
	if _error = _filter.Load (); _error != nil {
		return _error
	}
	
	return nil
}

