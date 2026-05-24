package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("[Nava] Welcome . . .")
	sebParent := "safeexambrowser.exe"
	sebClient := "safeexambrowser.client"

	fmt.Println("[Nava] Running setup")
	nDll := setup()
	enablePriv()

	fmt.Println("[Nava] Starting . . .")
	sebPid := findProcessSync(sebParent)

	if sebPid != 0 {
		fmt.Println("[Nava] Found parent pid")

		time.Sleep(100 * time.Millisecond)
		if e := inject(sebPid, nDll); e != nil {
			fmt.Println("[Nava::Error] Failed to inject nava in parent process", e.Error())
		} else {
			markInjected(sebPid)
			if e := callExportedFunction(sebPid, nDll, "Nava"); e != nil {
				fmt.Println("[Nava::Error] Failed to call Nava", e.Error())
			}
		}
	}

	fmt.Println("[Nava::AutoReInject] Starting watcher")

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		backgroundInject(
			sebClient,
			nDll,
			1*time.Second,
			func(pid uint32) {
				// onsuccess
				fmt.Printf("[Nava::AutoReInject] Calling Nava on PID %d\n", pid)
				time.Sleep(300 * time.Millisecond)
				if err := callExportedFunction(pid, nDll, "Nava"); err != nil {
					fmt.Printf("[Nava::AutoReInject::Error] Call Nava failed: %v\n", err)
				}
			},
		)
	}()

	select {}

}
