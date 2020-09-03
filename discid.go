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

// Platform dependent feature
//
// The platform dependent features are currently discid.FeatureRead,
// discid.FeatureMcn and discid.FeatureIsrc.
//
// See the libdiscid feature matrix (https://musicbrainz.org/doc/libdiscid#Feature_Matrix)
// for a list of supported features per platform.
type Feature uint

const (
	// Read TOC from disc
	FeatureRead = C.DISCID_FEATURE_READ
	// Read MCN from disc
	FeatureMcn = C.DISCID_FEATURE_MCN
	// Read ISRCs from disc
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

// Check if a certain feature is implemented on the current platform.
//
// This only works for single features, not bit masks with multiple features.
//
// See the libdiscid feature matrix (https://musicbrainz.org/doc/libdiscid#Feature_Matrix)
// for a list of supported features per platform.
func HasFeature(feature Feature) bool {
	result := C.discid_has_feature(uint32(feature))
	return result == 1
}

// Read the disc in the given CD-ROM/DVD-ROM drive extracting only the TOC.
//
// This function reads the disc in the drive specified by the given device
// identifier. If the device is an empty string, the default device, as
// returned by discid.DefaultDevice, is used.
//
// This function will only read the TOC, hence only the disc ID itself will be
// available. Use discid::ReadFeatures if you want to read also MCN and ISRCs.
func Read(device string) (disc Disc, err error) {
	return ReadFeatures(device, FeatureRead)
}

// Read the disc in the given CD-ROM/DVD-ROM drive with additional features.
//
// This function is similar to disc.Read but allows to read information about
// MCN and per-track ISRCs in addition to the normal TOC data.
//
// The parameter features accepts a bitwise combination of values.
// discid.FeatureRead is always implied, so it is not necessary to specify it.
//
// Reading MCN and ISRCs is not available on all platforms. You can use the
// has_feature function to check if a specific feature is available. Passing
// unsupported features here will just be ignored.
//
// Note that reading MCN and ISRC data is significantly slower than just
// reading the TOC, so only request the features you actually need.
func ReadFeatures(device string, features Feature) (disc Disc, err error) {
	handle := C.discid_new()
	disc = Disc{handle}
	var c_device *C.char = nil
	if device != "" {
		c_device = C.CString(device)
		defer C.free(unsafe.Pointer(c_device))
	}
	var status = C.discid_read_sparse(handle, c_device, C.uint(features))
	if status == 0 {
		err = disc
	}
	return
}

// Provides the TOC of a known CD.
//
// This function may be used if the TOC has been read earlier and you want to calculate
// the disc ID afterwards, without accessing the disc drive.
//
// first is the track number of the first track (1-99).
// The offsets parameter points to an array which contains the track offsets for each track.
// The first element, offsets[0], is the leadout track. It must contain the total number of
// sectors on the disc. offsets must not be longer than 100 elements (leadout + 99 tracks).
func Put(first int, offsets []int) (disc Disc, err error) {
	last := first + len(offsets) - 2
	handle := C.discid_new()
	disc = Disc{handle}
	// libdiscid always expects an array of 100 integers, no matter the track count.
	var c_offsets [100]C.int
	for i := 0; i < len(offsets); i++ {
		c_offsets[i] = C.int(offsets[i])
	}
	var status = C.discid_put(handle, C.int(first), C.int(last), &c_offsets[0])
	if status == 0 {
		err = disc
	}
	return
}

// Release the memory allocated for the Disc object.
func (disc Disc) Close() {
	C.discid_free(disc.handle)
}

// Return a human-readable error message.
//
// This function may only be used if disc.Read failed. The returned error
// message is only valid as long as the Disc object exists.
func (disc Disc) Error() string {
	err := C.discid_get_error_msg(disc.handle)
	return C.GoString(err)
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
