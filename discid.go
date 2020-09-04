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
import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

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
	// Read with all features
	FeatureAll = FeatureRead | FeatureMcn | FeatureIsrc
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

// Holds information about a single track
type Track struct {
	// Track number (1-99) of the track
	Number int
	// Start offset in sectors
	Offset int
	// Track length in sectors
	Sectors int
	// ISRC for this track (might be empty).
	//
	// This will only bet set if discid.ReadFeatures` is called with discid.FeatureIsrc.
	Isrc string
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
	d := Disc{C.discid_new()}
	var c_device *C.char = nil
	if device != "" {
		c_device = C.CString(device)
		defer C.free(unsafe.Pointer(c_device))
	}
	var status = C.discid_read_sparse(d.handle, c_device, C.uint(features))
	if status == 0 {
		defer d.Close()
		err = errors.New(d.ErrorMessage())
	} else {
		disc = d
	}
	return
}

// Provides the TOC of a known CD.
//
// This function may be used if the TOC has been read earlier and you want to calculate
// the disc ID afterwards, without accessing the disc drive.
//
// first is the track number of the first track (1-99).
// The offsets parameter is an array which contains the track offsets for each track.
// The first element, offsets[0], is the leadout track. It must contain the total number of
// sectors on the disc. offsets must not be longer than 100 elements (leadout + 99 tracks).
func Put(first int, offsets []int) (disc Disc, err error) {
	last := first + len(offsets) - 2
	d := Disc{C.discid_new()}
	// libdiscid always expects an array of 100 integers, no matter the track count.
	var c_offsets [100]C.int
	c_offsets[0] = C.int(offsets[0])
	for i, n := range offsets[1:] {
		track := i + first
		if track > 99 {
			break
		}
		c_offsets[track] = C.int(n)
	}
	var status = C.discid_put(d.handle, C.int(first), C.int(last), &c_offsets[0])
	if status == 0 {
		defer d.Close()
		err = errors.New(d.ErrorMessage())
	} else {
		disc = d
	}
	return
}

// Parses a TOC string and returns a Disc instance for it.
//
// The TOC string provided here must have the same format as returned by Disc.TocString.
//
// This function can be used if you already have a TOC string like e.g.
// "1 11 242457 150 44942 61305 72755 96360 130485 147315 164275 190702 205412 220437".
func Parse(toc string) (disc Disc, err error) {
	first := 0
	last := 0
	var offsets [100]int
	var i int
	var part string
	for i, part = range strings.Split(toc, " ") {
		parsedInt, e := strconv.Atoi(part)
		if e != nil {
			err = e
			return
		}
		if i == 0 {
			first = parsedInt
		} else if i == 1 {
			last = parsedInt
		} else {
			if i > (last+2) || i > 99+2 {
				err = errors.New("TOC string contains too many offsets (max. 100)")
				return
			}
			offsets[i-2] = parsedInt
		}
	}

	if i < 2 || first < 1 || last < 1 || last > 99 {
		msg := fmt.Sprintf("Invalid TOC string \"%v\"", toc)
		err = errors.New(msg)
		return
	}

	offsetCount := i - 2
	trackCount := last - first + 1
	if offsetCount < trackCount {
		msg := fmt.Sprintf("Number of offsets %v does not match track count %v",
			offsetCount, trackCount)
		err = errors.New(msg)
		return
	}

	return Put(first, offsets[0:trackCount+1])
}

// Release the memory allocated for the Disc object.
func (d Disc) Close() {
	C.discid_free(d.handle)
}

// Return a human-readable error message.
//
// This function may only be used if discid.Read failed.
func (d Disc) ErrorMessage() string {
	err := C.discid_get_error_msg(d.handle)
	return C.GoString(err)
}

// String representation of the disc, same as Id()
func (d Disc) String() string {
	return d.Id()
}

// Returns the MusicBrainz disc ID.
func (d Disc) Id() string {
	id := C.discid_get_id(d.handle)
	return C.GoString(id)
}

// Returns the FreeDB disc ID.
func (d Disc) FreedbId() string {
	id := C.discid_get_freedb_id(d.handle)
	return C.GoString(id)
}

// Return a string representing CD Table Of Contents (TOC).
//
// The TOC string is a list of integers separated by a single space character.
//
// Example: 1 11 242457 150 44942 61305 72755 96360 130485 147315 164275 190702 205412 220437
//
// The integers represent (in order):
//
// - First track number (normally one)
//
// - Last track number
//
// - Lead-out track offset
//
// - Up to 99 frame offsets
func (d Disc) TocString() string {
	toc := C.discid_get_toc_string(d.handle)
	return C.GoString(toc)
}

// An URL for submitting the DiscID to MusicBrainz.
func (d Disc) SubmissionUrl() string {
	url := C.discid_get_submission_url(d.handle)
	return C.GoString(url)
}

// The number of the first track on this disc.
func (d Disc) FirstTrackNum() int {
	return int(C.discid_get_first_track_num(d.handle))
}

// The number of the last track on this disc.
func (d Disc) LastTrackNum() int {
	return int(C.discid_get_last_track_num(d.handle))
}

// The length of the disc in sectors.
func (d Disc) Sectors() int {
	return int(C.discid_get_sectors(d.handle))
}

// Return the Media Catalogue Number (MCN) for the disc, if present.
//
// This is essentially an EAN (= UPC with 0 prefix).
func (d Disc) Mcn() string {
	mcn := C.discid_get_mcn(d.handle)
	return C.GoString(mcn)
}

// Return the Media Catalogue Number (MCN) for the disc, if present.
//
// This is essentially an EAN (= UPC with 0 prefix).
func (d Disc) Track(number int) Track {
	first := d.FirstTrackNum()
	last := d.LastTrackNum()
	if number < first || number > last {
		err := fmt.Sprintf(
			"track number out of bounds: given %v, expected between %v and %v",
			number, first, last)
		panic(err)
	}
	n := C.int(number)
	return Track{
		number,
		int(C.discid_get_track_offset(d.handle, n)),
		int(C.discid_get_track_length(d.handle, n)),
		C.GoString(C.discid_get_track_isrc(d.handle, n)),
	}
}
