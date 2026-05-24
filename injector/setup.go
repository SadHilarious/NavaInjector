package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func setup() string {

	temp := os.TempDir()
	nDllPath := filepath.Join(temp, "nava.dll")
	mDllPath := filepath.Join(temp, "milim.dll")

	fmt.Printf("[Nava] Extracting %s\n", nDllPath)
	if err := extractDll(navaDll, nDllPath); err != nil {
		fmt.Printf("[Nava::Error] Failed to extract nava: %v\n", err)
		return ""
	}

	fmt.Printf("[Nava] Extracting %s\n", mDllPath)
	if err := extractDll(milimDll, mDllPath); err != nil {
		fmt.Printf("[Nava::Error] Failed to extract milim: %v\n", err)
		return ""
	}

	return nDllPath
}
