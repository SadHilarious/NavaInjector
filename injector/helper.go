package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/gentlemanautomaton/winproc/psapi"
	"golang.org/x/sys/windows"
)

//go:embed Milim.dll
var milimDll []byte

//go:embed nava.dll
var navaDll []byte

var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")

	VirtualAllocEx     = kernel32.NewProc("VirtualAllocEx")
	VirtualFreeEx      = kernel32.NewProc("VirtualFreeEx")
	CreateRemoteThread = kernel32.NewProc("CreateRemoteThread")
	TerminateThread    = kernel32.NewProc("TerminateThread")

	GetModuleHandle = kernel32.NewProc("GetModuleHandleW")
	GetProcAddress  = kernel32.NewProc("GetProcAddress")

	injectedPIDs = make(map[uint32]bool)
	mInjected    sync.Mutex
)

func findProcessSync(name string) uint32 {
	ch := make(chan uint32, 1)
	go findProcess(name, ch)
	return <-ch
}

func backgroundInject(targetName, dllPath string, interval time.Duration) {
	fmt.Printf("[Nava] Watching for %s . . . \n", targetName)

	for {
		pids, err := findProcessesByName(targetName)
		if err != nil {
			fmt.Printf("[Nava::Warn] Scan failed: %v\n", err)
			time.Sleep(interval)
			continue
		}

		for _, pid := range pids {
			if !isInjected(pid) {
				fmt.Printf("[Nava] Detected %s (PID: %d), injecting...\n", targetName, pid)

				if err := inject(pid, dllPath); err != nil {
					fmt.Printf("[Nava::Error] Inject failed for PID %d: %v\n", pid, err)
					continue
				}

				markInjected(pid)
				fmt.Printf("[Nava] Injected into PID %d\n", pid)
			}
		}

		cleanDp(pids)
		time.Sleep(interval)
	}
}

func findProcessesByName(targetName string) ([]uint32, error) {
	var pids []uint32
	targetName = strings.ToLower(targetName)

	snapshot, err := psapi.CreateSnapshot(psapi.SnapAll, 0)
	if err != nil {
		return nil, err
	}
	defer syscall.CloseHandle(snapshot)

	proc, err := psapi.FirstProcess(snapshot)
	for err == nil {
		if strings.Contains(strings.ToLower(proc.Name()), targetName) {
			pids = append(pids, proc.ProcessID)
		}
		proc, err = psapi.NextProcess(snapshot)
	}
	return pids, nil
}

func cleanDp(currentPIDs []uint32) {
	mInjected.Lock()
	defer mInjected.Unlock()

	currentSet := make(map[uint32]bool)
	for _, pid := range currentPIDs {
		currentSet[pid] = true
	}

	for pid := range injectedPIDs {
		if !currentSet[pid] {
			fmt.Printf("[Nava] Process %d exited removed from tracking\n", pid)
			delete(injectedPIDs, pid)
		}
	}
}

func isInjected(pid uint32) bool {
	mInjected.Lock()
	defer mInjected.Unlock()
	return injectedPIDs[pid]
}

func markInjected(pid uint32) {
	mInjected.Lock()
	defer mInjected.Unlock()
	injectedPIDs[pid] = true
}

func extractDll(data []byte, destPath string) error {
	if dir := filepath.Dir(destPath); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return os.WriteFile(destPath, data, 0755)
}

func findProcess(targetProcess string, targetPid chan<- uint32) {
	defer close(targetPid)

	var process psapi.ProcessEntry
	process.Size = uint32(reflect.TypeFor[psapi.ProcessEntry]().Size())

	for {
		snapshot, errSnap := psapi.CreateSnapshot(psapi.SnapAll, 0)
		if errSnap != nil {
			return
		}

		process, _ = psapi.FirstProcess(snapshot)

		var e error
		for e == nil {
			if strings.Contains(strings.ToLower(process.Name()), targetProcess) {
				syscall.CloseHandle(snapshot)
				targetPid <- process.ProcessID
			}

			process, e = psapi.NextProcess(snapshot)
		}
	}

}

func inject(pid uint32, dllpath string) error {
	hProc, e := windows.OpenProcess(
		windows.PROCESS_VM_OPERATION|
			windows.PROCESS_VM_WRITE|
			windows.PROCESS_CREATE_THREAD|
			windows.PROCESS_QUERY_INFORMATION, false, pid)
	if e != nil {
		return e
	}
	defer syscall.CloseHandle(syscall.Handle(hProc))

	kernel32Ptr, _ := syscall.UTF16PtrFromString("kernel32.dll")
	loadLibPtr, _ := syscall.BytePtrFromString("LoadLibraryA")

	hMod, _, _ := GetModuleHandle.Call(uintptr(unsafe.Pointer(kernel32Ptr)))
	if hMod == 0 {
		return fmt.Errorf("GetModuleHandle call -> kernel32.dll failed")
	}

	loadLibAddr, _, _ := GetProcAddress.Call(hMod, uintptr(unsafe.Pointer(loadLibPtr)))
	if loadLibAddr == 0 {
		return fmt.Errorf("GetProcAddress call -> LoadLibraryA failed")
	}

	dllPathBytes := []byte(dllpath + "\x00")
	addr, _, _ := VirtualAllocEx.Call(uintptr(hProc), 0, uintptr(len(dllPathBytes)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)
	if addr == 0 {
		return windows.GetLastError()
	}
	defer VirtualFreeEx.Call(uintptr(hProc), addr, 0, windows.MEM_RELEASE)

	var bytes uintptr
	e = windows.WriteProcessMemory(hProc, addr, &dllPathBytes[0], uintptr(len(dllPathBytes)), &bytes)

	hThread, _, e := CreateRemoteThread.Call(uintptr(hProc), 0, 0, loadLibAddr, addr, 0, 0)
	if hThread == 0 {
		return windows.GetLastError()
	}
	defer syscall.CloseHandle(syscall.Handle(hThread))

	syscall.WaitForSingleObject(syscall.Handle(hThread), windows.INFINITE)
	syscall.CloseHandle(syscall.Handle(hThread))
	return nil
}

func enablePriv() {
	var hToken syscall.Token
	var luid windows.LUID
	currProc, _ := syscall.GetCurrentProcess()
	e := syscall.OpenProcessToken(currProc, syscall.TOKEN_ADJUST_PRIVILEGES|syscall.TOKEN_QUERY, &hToken)
	if e != nil {
		fmt.Println("[Nava::Error] OpenProcessToken failed: ", e)
		return
	}
	defer syscall.CloseHandle(syscall.Handle(hToken))

	if e = windows.LookupPrivilegeValue(nil, windows.StringToUTF16Ptr("SeDebugPrivilege"), &luid); e != nil {
		syscall.CloseHandle(syscall.Handle(hToken))
		fmt.Println("[Nava::Error] LookupPrivilege failed", e.Error())
		return
	}

	tokenPriv := windows.Tokenprivileges{
		PrivilegeCount: 1,
		Privileges: [1]windows.LUIDAndAttributes{
			{
				Luid:       luid,
				Attributes: windows.SE_PRIVILEGE_ENABLED,
			},
		},
	}

	privSize := uint32(unsafe.Sizeof(tokenPriv))
	if e = windows.AdjustTokenPrivileges(windows.Token(hToken), false, &tokenPriv, privSize, nil, nil); e != nil {
		fmt.Println("[Nava::Error] AdjustTokenPrivileges failed: ", e)
		return
	}

	if err, ok := windows.GetLastError().(syscall.Errno); ok && err == windows.ERROR_NOT_ALL_ASSIGNED {
		fmt.Println("[Nava::Error] Not all privilege assigned successfully :(")
		return
	}

	fmt.Println("[Nava] Success enable privilege")
}
