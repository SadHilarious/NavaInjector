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

	fmt.Println("[Nava] Ready . . .")
	sebPid := findProcessSync(sebParent)

	if sebPid != 0 {
		fmt.Println("[Nava] Found parent pid")

		time.Sleep(100 * time.Millisecond)
		if e := inject(sebPid, nDll); e != nil {
			fmt.Println("[Nava::Error] Failed to inject nava in parent process", e.Error())
		} else {
			markInjected(sebPid)
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
			500*time.Millisecond,
		)
	}()

	select {}

}
