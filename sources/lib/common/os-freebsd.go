
//go:build freebsd


package common


import "syscall"




func SysSetrlimitMemory (_limitMemory uint) (error) {
	{
		_limitMb := (2 * _limitMemory) + (1 * 1024)
		_limit := syscall.Rlimit {
				Cur : int64 (_limitMb) * 1024 * 1024,
				Max : int64 (_limitMb) * 1024 * 1024,
			}
		if _error := syscall.Setrlimit (syscall.RLIMIT_AS, &_limit); _error != nil {
			return _error
		}
	}
	{
		_limitMb := _limitMemory
		_limit := syscall.Rlimit {
				Cur : int64 (_limitMb) * 1024 * 1024,
				Max : int64 (_limitMb) * 1024 * 1024,
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
				Cur : int64 (_limitDescriptors),
				Max : int64 (_limitDescriptors),
			}
		if _error := syscall.Setrlimit (syscall.RLIMIT_NOFILE, &_limit); _error != nil {
			return _error
		}
	}
	return nil
}

