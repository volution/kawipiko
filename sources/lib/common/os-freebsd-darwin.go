
//go:build freebsd || darwin


package common


import "fmt"
import "os"
import "syscall"




func SysStatExtract (_stat os.FileInfo) (_dev uint64, _inode uint64, _size uint64, _timestamp uint64, _error error) {
	if _stat, _ok := _stat.Sys () .(*syscall.Stat_t); _ok {
		_dev = uint64 (_stat.Dev)
		_inode = uint64 (_stat.Ino)
		_size = uint64 (_stat.Size)
		_timestamp = (uint64 (_stat.Mtimespec.Sec) * 1000000) + (uint64 (_stat.Mtimespec.Nsec) / 1000)
	} else {
		_error = fmt.Errorf ("[7739aed3]  failed `stat`-ing`!")
	}
	return
}

