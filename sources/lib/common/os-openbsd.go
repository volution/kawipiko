
//go:build openbsd


package common


import "syscall"




func SysSetrlimitMemory (_limitMemory uint) (error) {
	{
		_limitMb := _limitMemory
		_limit := syscall.Rlimit {
				Cur : uint64 (_limitMb) * 1024 * 1024,
				Max : uint64 (_limitMb) * 1024 * 1024,
			}
		if _error := syscall.Setrlimit (syscall.RLIMIT_DATA, &_limit); _error != nil {
			return _error
		}
	}
	return nil
}




func SysSetrlimitDescriptors (_limitDescriptors uint) (error) {
	{
		_limit := syscall.Rlimit {
				Cur : uint64 (_limitDescriptors),
				Max : uint64 (_limitDescriptors),
			}
		if _error := syscall.Setrlimit (syscall.RLIMIT_NOFILE, &_limit); _error != nil {
			return _error
		}
	}
	return nil
}

