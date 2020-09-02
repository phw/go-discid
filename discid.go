// Copyright (C) 2020 Philipp Wolfer <ph.wolfer@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// discid allows you to calculate MusicBrainz or FreeDB disc IDs for audio CDs.
//
// discid provides bindings to the MusicBrainz libdiscid (https://musicbrainz.org/doc/libdiscid)
// library. In addition to calculating the disc IDs you can also get advanced data from the
// audio CD such as MCN (media catalogue number) and per-track ISRCs.
//
// To get started see the documentation and examples of discid.Read, discid.ReadFeatures,
// discid.Put and discid.Parse.
//
// Details about the use and calculation of disc IDs can be found at the MusicBrainz
// disc ID documentation (https://musicbrainz.org/doc/Disc_ID).
//
// The source code of this library is available on GitHub (https://github.com/phw/go-discid)
// under the terms of the GNU Lesser General Public License version 3 or later.
package discid

// #cgo LDFLAGS: -ldiscid
// #include <stdlib.h>
// #include "discid/discid.h"
import "C"
import "unsafe"

type Feature uint

const (
	FeatureRead = C.DISCID_FEATURE_READ
	FeatureMcn  = C.DISCID_FEATURE_MCN
	FeatureIsrc = C.DISCID_FEATURE_ISRC
)

// Holds information about a read disc (TOC, MCN, ISRCs).
//
// Use discid.Read, discid.ReadFeatures, discid.Put or discid.Parse
// to initialize an instance of Disc.
//
// Use the Close method to free the allocated resources after use, e.g.:
//   disc := discid.Read("") // Read from default device
//   defer disc.Close()
type Disc struct {
	handle *C.DiscId
}

// Return the name of the default disc drive for this operating system.
//
// The default device is system dependent, e.g. "/dev/cdrom" on Linux and "D:" on Windows.
func DefaultDevice() string {
	device := C.discid_get_default_device()
	return C.GoString(device)
}

// Return version information about libdiscid.
//
// The returned string will be e.g. "libdiscid 0.6.2".
func Version() string {
	version := C.discid_get_version_string()
	return C.GoString(version)
}

// Read the disc in the given CD-ROM/DVD-ROM drive extracting only the TOC.
//
// This function reads the disc in the drive specified by the given device
// identifier. If the device is an empty string, the default device, as
// returned by discid.DefaultDevice, is used.
//
// This function will only read the TOC, hence only the disc ID itself will be
// available. Use discid::ReadFeatures if you want to read also MCN and ISRCs.
func Read(device string) Disc {
	return ReadFeatures(device, FeatureRead)
}

func ReadFeatures(device string, features Feature) Disc {
	handle := C.discid_new()
	var c_device *C.char = nil
	if device != "" {
		c_device = C.CString(device)
	}
	C.discid_read_sparse(handle, c_device, C.uint(features))
	C.free(unsafe.Pointer(c_device))
	return Disc{handle}
}

// Release the memory allocated for the Disc object.
func (disc Disc) Close() {
	C.discid_free(disc.handle)
}

// Returns the MusicBrainz disc ID.
func (disc Disc) Id() string {
	id := C.discid_get_id(disc.handle)
	return C.GoString(id)
}

// Return the Media Catalogue Number (MCN) for the disc.
//
// This is essentially an EAN (= UPC with 0 prefix).
func (disc Disc) Mcn() string {
	mcn := C.discid_get_mcn(disc.handle)
	return C.GoString(mcn)
}
