

//go:build !(linux && amd64 && seccomp)




package seccomp




func AllowOnlySyscalls (_syscalls []string) (error) {
	return nil
}

