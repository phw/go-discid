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

package discid_test

import (
	"fmt"
	"strings"
	"testing"

	"musicbrainz.org/discid"
)

func TestDefaultDevice(t *testing.T) {
	device := discid.DefaultDevice()
	if device == "" {
		t.Errorf("TestDefaultDevice() is empty; expected device name")
	}
}

func ExampleDefaultDevice() {
	fmt.Printf("Default device: %v\n", discid.DefaultDevice())
}

func TestVersion(t *testing.T) {
	version := discid.Version()
	if !strings.HasPrefix(version, "libdiscid") {
		t.Errorf("Version() = %v; expected starting with \"libdiscid\"", version)
	}
}

func ExampleVersion() {
	fmt.Printf("Version: %v\n", discid.Version())
}

func TestHasFeature(t *testing.T) {
	result := discid.HasFeature(discid.FeatureRead)
	if !result {
		t.Errorf("HasFeature() = %v; expected true", result)
	}
}

func ExampleHasFeature() {
	if discid.HasFeature(discid.FeatureIsrc) {
		fmt.Println("ISRC support available")
	}
}

func ExampleRead() {
	disc := discid.Read("") // Read from default device
	defer disc.Close()
	fmt.Printf("Disc ID: %v\n", disc.Id())
}

func ExampleReadFeatures() {
	// Read TOC and MCN from the disc in /dev/cdrom
	disc := discid.ReadFeatures("/dev/cdrom", discid.FeatureRead|discid.FeatureMcn)
	defer disc.Close()
	fmt.Printf("Disc ID: %v\n", disc.Id())
	fmt.Printf("MCN    : %v\n", disc.Mcn())
}
