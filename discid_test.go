// Copyright (C) 2020-2023 Philipp Wolfer <ph.wolfer@gmail.com>
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
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uploadedlobster.com/discid"
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

func TestReadInvalidDevice(t *testing.T) {
	_, err := discid.Read("notadevice")
	if err == nil {
		t.Errorf("Expected error for accessing invalid device")
	}
}

func ExampleRead() {
	disc, err := discid.Read("") // Read from default device
	if err != nil {
		log.Fatal(err)
	}
	defer disc.Close()
	fmt.Printf("Disc ID: %v\n", disc.Id())
}

func ExampleReadFeatures() {
	// Read TOC and MCN from the disc in /dev/cdrom
	disc, err := discid.ReadFeatures("/dev/cdrom", discid.FeatureRead|discid.FeatureMcn)
	if err != nil {
		log.Fatal(err)
	}
	defer disc.Close()
	fmt.Printf("Disc ID: %v\n", disc.Id())
	fmt.Printf("MCN    : %v\n", disc.Mcn())
}

func TestPut(t *testing.T) {
	assert := assert.New(t)
	first := 1
	offsets := []int{
		206535, 150, 18901, 39738, 59557, 79152, 100126, 124833, 147278, 166336, 182560,
	}
	disc, err := discid.Put(first, offsets)
	if err != nil {
		t.Fatal(err)
	}
	defer disc.Close()
	assert.Equal("Wn8eRBtfLDfM0qjYPdxrz.Zjs_U-", disc.Id())
	assert.Equal(disc.Id(), fmt.Sprint(disc))
	assert.Equal("830abf0a", disc.FreedbId())
	assert.Equal(1, disc.FirstTrackNum())
	assert.Equal(10, disc.LastTrackNum())
	assert.Equal(206535, disc.Sectors())
	assert.Equal(
		"1 10 206535 150 18901 39738 59557 79152 100126 124833 147278 166336 182560",
		disc.TocString())
	assert.Equal(
		"http://musicbrainz.org/cdtoc/attach?id=Wn8eRBtfLDfM0qjYPdxrz.Zjs_U-&tracks=10&toc=1+10+206535+150+18901+39738+59557+79152+100126+124833+147278+166336+182560",
		disc.SubmissionUrl())
	for i := disc.FirstTrackNum(); i <= disc.LastTrackNum(); i++ {
		track := disc.Track(i)
		offset := offsets[track.Number]
		next := 0
		if track.Number < disc.LastTrackNum() {
			next = track.Number + 1
		}
		length := offsets[next] - offset
		assert.Equal(i, track.Number)
		assert.Equal(offset, track.Offset)
		assert.Equal(length, track.Sectors)
		assert.Equal("", track.Isrc)
	}
}

func TestPutFirstTrackLargerOne(t *testing.T) {
	assert := assert.New(t)
	first := 3
	offsets := []int{
		206535, 150, 18901, 39738, 59557, 79152, 100126, 124833, 147278, 166336, 182560,
	}
	disc, err := discid.Put(first, offsets[0:])
	if err != nil {
		t.Fatal(err)
	}
	defer disc.Close()
	assert.Equal("ByBKvJM1hBL7XtvsPyYtIjlX0Bw-", disc.Id())
	assert.Equal(3, disc.FirstTrackNum())
	assert.Equal(12, disc.LastTrackNum())
	assert.Equal(206535, disc.Sectors())
}

func TestPutTooManyOffsets(t *testing.T) {
	first := 1
	offsets := [101]int{}
	disc, err := discid.Put(first, offsets[0:])
	assert.Empty(t, disc)
	assert.NotEmpty(t, err)
	if err.Error() != "Illegal track limits" {
		t.Errorf("Expected error \"Illegal track limits\"")
	}
}

func TestPutTooManyTracks(t *testing.T) {
	// First track number is 82, 19 tracks
	// => last track number would be 100, but 99 is max.
	first := 82
	offsets := [20]int{}
	disc, err := discid.Put(first, offsets[0:])
	assert.Empty(t, disc)
	assert.NotEmpty(t, err)
	if err.Error() != "Illegal track limits" {
		t.Errorf("Expected error \"Illegal track limits\"")
	}
}

func ExamplePut() {
	first := 1
	offsets := []int{
		242457, 150, 44942, 61305, 72755, 96360, 130485, 147315, 164275, 190702, 205412, 220437,
	}
	disc, err := discid.Put(first, offsets)
	if err != nil {
		log.Fatal(err)
	}
	defer disc.Close()
	fmt.Println(disc.Id())
	// Output: lSOVc5h6IXSuzcamJS1Gp4_tRuA-
}

func ExampleParse() {
	toc := "1 11 242457 150 44942 61305 72755 96360 130485 147315 164275 190702 205412 220437"
	disc, err := discid.Parse(toc)
	if err != nil {
		log.Fatal(err)
	}
	defer disc.Close()
	fmt.Println(disc.Id())
	// Output: lSOVc5h6IXSuzcamJS1Gp4_tRuA-
}

func TestParseMinimal(t *testing.T) {
	assert := assert.New(t)
	toc := "1 1 44942 150"
	disc, err := discid.Parse(toc)
	if err != nil {
		t.Fatal(err)
	}
	defer disc.Close()
	assert.Equal("ANJa4DGYN_ktpzOwvVPtcjwP7mE-", disc.Id())
	assert.Equal(toc, disc.TocString())
}

func TestParseFirstTrackNotOne(t *testing.T) {
	assert := assert.New(t)
	toc := "3 12 242457 150 18901 39738 59557 79152 100126 124833 147278 166336 182560"
	disc, err := discid.Parse(toc)
	if err != nil {
		t.Fatal(err)
	}
	defer disc.Close()
	assert.Equal("fC1yNbC5bVjbvphqlAY9JyYoWEY-", disc.Id())
	assert.Equal(toc, disc.TocString())
}

func TestParseNaN(t *testing.T) {
	toc := "1 2 242457 150 a"
	_, err := discid.Parse(toc)
	if assert.Error(t, err) {
		if err.(*strconv.NumError).Err != strconv.ErrSyntax {
			t.Errorf("Expected strconv.ErrSyntax, got \"%v\"", err)
		}
	}
}

func TestParseInvalidEmpty(t *testing.T) {
	toc := ""
	_, err := discid.Parse(toc)
	if assert.Error(t, err) {
		if err.(*strconv.NumError).Err != strconv.ErrSyntax {
			t.Errorf("Expected strconv.ErrSyntax, got \"%v\"", err)
		}
	}
}

func TestParseTooManyOffsets(t *testing.T) {
	assert := assert.New(t)
	toc := "1 2 242457 150 200 300"
	_, err := discid.Parse(toc)
	assert.Error(err)
	assert.Equal("TOC string contains too many offsets (max. 100)", err.Error())
}

func TestParseTooManyOffsetsTotal(t *testing.T) {
	assert := assert.New(t)
	indexes := [103]string{"1", "99", "20000"}
	for i := 3; i < len(indexes); i++ {
		indexes[i] = strconv.Itoa(i * 100)
	}
	toc := strings.Join(indexes[:], " ")
	_, err := discid.Parse(toc)
	assert.Error(err)
	assert.Equal("TOC string contains too many offsets (max. 100)", err.Error())
}

func TestParseInvalidMissingOffsets(t *testing.T) {
	assert := assert.New(t)
	toc := "1 2 242457 150"
	_, err := discid.Parse(toc)
	assert.Error(err)
	assert.Equal("Number of offsets 1 does not match track count 2", err.Error())
}

func TestParseInvalidNotEnoughElements(t *testing.T) {
	assert := assert.New(t)
	toc := "1"
	_, err := discid.Parse(toc)
	assert.Error(err)
	assert.Equal("Invalid TOC string \"1\"", err.Error())
}

func TestTrackOutOfRange(t *testing.T) {
	assert := assert.New(t)
	first := 1
	offsets := []int{
		206535, 150, 18901, 39738, 59557, 79152, 100126, 124833, 147278, 166336, 182560,
	}
	disc, err := discid.Put(first, offsets)
	if err != nil {
		t.Fatal(err)
	}
	defer disc.Close()
	assert.Panics(func() { disc.Track(disc.FirstTrackNum() - 1) })
	assert.NotPanics(func() { disc.Track(disc.FirstTrackNum()) })
	assert.NotPanics(func() { disc.Track(disc.LastTrackNum()) })
	assert.Panics(func() { disc.Track(disc.LastTrackNum() + 1) })
}
