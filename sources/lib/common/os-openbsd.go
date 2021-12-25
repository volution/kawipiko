
//go:build openbsd


package common


import "syscall"




func SysSetrlimit (_limitMemory uint) (error) {
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

