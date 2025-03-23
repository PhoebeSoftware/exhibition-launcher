//go:build windows

package utils

import (
	"syscall"
	"unsafe"
)

type DiskStatus struct {
	All  uint64
	Free uint64
	Used uint64
}

func DiskUsage(path string) (disk DiskStatus) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getDiskFreeSpaceEx := kernel32.NewProc("GetDiskFreeSpaceExW")

	var lpFreeBytesAvailable, lpTotalNumberOfBytes, lpTotalNumberOfFreeBytes uint64

	pathPtr, _ := syscall.UTF16PtrFromString(path)
	ret, _, _ := getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)),
	)
	if ret == 0 {
		return
	}

	disk.All = lpTotalNumberOfBytes
	disk.Free = lpTotalNumberOfFreeBytes
	disk.Used = disk.All - disk.Free
	return
}

