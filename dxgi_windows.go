package main

import (
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

// mksyscall_windows.go borrowed from hcsshim, as it has HRESULT support, and lxnwin added.
//go:generate go run mksyscall_windows.go -lxnwin -output zsyscalls_windows.go dxgi_windows.go

//sys CreateDXGIFactory(riid *win.IID, ppFactory *unsafe.Pointer) (hr error) = dxgi.CreateDXGIFactory?

var (
	// https://github.com/apitrace/dxsdk/blob/master/Include/dxgi.h#L1984
	IID_IDXGIFactory = win.IID{Data1: 0x7b7166ec, Data2: 0x21c7, Data3: 0x44ae, Data4: [8]byte{0xb2, 0x1a, 0xc9, 0xae, 0x32, 0x1a, 0xe3, 0x69}}

	DXGI_ERROR_NOT_FOUND = syscall.Errno(0x887A0002)
)

// https://github.com/apitrace/dxsdk/blob/master/Include/dxgi.h#L2019-L2098

type IDXGIFactoryVtbl struct {
	QueryInterface          uintptr
	AddRef                  uintptr
	Release                 uintptr
	SetPrivateData          uintptr
	SetPrivateDataInterface uintptr
	GetPrivateData          uintptr
	GetParent               uintptr
	EnumAdapters            uintptr
	MakeWindowAssociation   uintptr
	GetWindowAssociation    uintptr
	CreateSwapChain         uintptr
	CreateSoftwareAdapter   uintptr
}

type IDXGIFactory struct {
	LpVtbl *IDXGIFactoryVtbl
}

// TODO: Fill out other APIs? I can't believe someone hasn't automated at least the IUnknown ones yet.

func (obj *IDXGIFactory) EnumAdapters(adapter uint32, ppAdapter **IDXGIAdapter) (hr error) {
	r0, _, _ := syscall.Syscall(obj.LpVtbl.EnumAdapters, 3, uintptr(unsafe.Pointer(obj)), uintptr(adapter), uintptr(unsafe.Pointer(ppAdapter)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

// https://github.com/apitrace/dxsdk/blob/master/Include/dxgi.h#L1279-L1350
type IDXGIAdapterVtbl struct {
	QueryInterface          uintptr
	AddRef                  uintptr
	Release                 uintptr
	SetPrivateData          uintptr
	SetPrivateDataInterface uintptr
	GetPrivateData          uintptr
	GetParent               uintptr
	EnumOutputs             uintptr
	GetDesc                 uintptr
	CheckInterfaceSupport   uintptr
}

type IDXGIAdapter struct {
	LpVtbl *IDXGIAdapterVtbl
}

func (obj *IDXGIAdapter) GetDesc(pDesc *DXGI_ADAPTER_DESC) (hr error) {
	r0, _, _ := syscall.Syscall(obj.LpVtbl.GetDesc, 2, uintptr(unsafe.Pointer(obj)), uintptr(unsafe.Pointer(pDesc)), 0)
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}
		hr = syscall.Errno(r0)
	}
	return
}

// https://github.com/apitrace/dxsdk/blob/master/Include/dxgi.h#L186-L190
type LUID struct {
	LowPart  uint32 // DWORD
	HighPart int32  // LONG
}

// https://github.com/apitrace/dxsdk/blob/master/Include/dxgi.h#L196-L207
type DXGI_ADAPTER_DESC struct {
	Description           [128]uint16
	VendorId              uint32
	DeviceId              uint32
	SubSysId              uint32
	Revision              uint32
	DedicatedVideoMemory  uintptr // size_t
	DedicatedSystemMemory uintptr // size_t
	SharedSystemMemory    uintptr // size_t
	AdapterLuid           LUID
}
