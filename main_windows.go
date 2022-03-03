package main

import (
	"errors"
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

func main() {
	var pFactory *IDXGIFactory
	if err := CreateDXGIFactory(&IID_IDXGIFactory, (*unsafe.Pointer)(unsafe.Pointer(&pFactory))); err != nil {
		fmt.Printf("Create DXGIFactory failed: %v\n", err)
		os.Exit(1)
	}

	for numAdapter := uint32(0); ; numAdapter++ {
		var pAdapter *IDXGIAdapter
		if err := pFactory.EnumAdapters(numAdapter, &pAdapter); err != nil {
			if errors.Is(err, DXGI_ERROR_NOT_FOUND) {
				break
			}
			fmt.Printf("EnumAdapters failed: %v\n", err)
			os.Exit(1)
		}
		// fmt.Printf("EnumAdapters(%d): %v\n", numAdapter, uintptr(unsafe.Pointer(pAdapter)))

		var desc DXGI_ADAPTER_DESC
		if err := pAdapter.GetDesc(&desc); err != nil {
			fmt.Printf("Adapter(%d).GetDesc failed: %v\n", numAdapter, err)
			os.Exit(1)
		}

		fmt.Printf("%d: %v\n", numAdapter, windows.UTF16ToString(desc.Description[:]))
		fmt.Printf("%d: Vendor 0x%08x, Device 0x%08x\n", numAdapter, desc.VendorId, desc.DeviceId)
		fmt.Printf("%d: Subsystem 0x%08x Revision 0x%08x\n", numAdapter, desc.SubSysId, desc.Revision)
		fmt.Printf("%d: Dedicated Video Memory: %d\n", numAdapter, desc.DedicatedVideoMemory)
		fmt.Printf("%d: Dedicated System Memory: %d\n", numAdapter, desc.DedicatedSystemMemory)
		fmt.Printf("%d: Shared System Memory: %d\n", numAdapter, desc.SharedSystemMemory)
		fmt.Printf("%d: LUID 0x%08x:0x%08x\n", numAdapter, desc.AdapterLuid.HighPart, desc.AdapterLuid.LowPart)
	}
}
