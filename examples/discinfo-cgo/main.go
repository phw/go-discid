// This example uses cgo to interface with libdiscid directly instead of using go-discid.
// Compare this to the much simpler use of go-discid in the discinfo example.

package main

// #cgo LDFLAGS: -ldiscid
// #include <stdlib.h>
// #include "discid/discid.h"
import "C"
import (
	"fmt"
	"log"
	"unsafe"
)

func main() {
	fmt.Printf("Version       : %v\n", C.GoString(C.discid_get_version_string()))
	fmt.Printf("Default device: %v\n", C.GoString(C.discid_get_default_device()))
	disc := C.discid_new()
	defer C.discid_free(disc)
	// Read from device
	device := C.CString("/dev/cdrom")
	defer C.free(unsafe.Pointer(device))
	status := C.discid_read_sparse(disc, device, C.DISCID_FEATURE_READ|C.DISCID_FEATURE_MCN)
	if status == 0 {
		err := C.discid_get_error_msg(disc)
		log.Fatal(err)
	}
	discid := C.GoString(C.discid_get_id(disc))
	fmt.Printf("Disc ID       : %v\n", discid)
	mcn := C.GoString(C.discid_get_mcn(disc))
	fmt.Printf("MCN           : %v\n", mcn)
}
