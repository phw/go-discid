// This example uses go-discid to query details about the disc in the default device.

package main

import (
	"fmt"
	"log"

	"github.com/phw/go-discid"
)

func main() {
	fmt.Printf("Version       : %v\n", discid.Version())
	fmt.Printf("Default device: %v\n", discid.DefaultDevice())
	// Read from default device
	disc, err := discid.ReadFeatures("", discid.FeatureAll)
	if err != nil {
		log.Fatal(err)
	}
	defer disc.Close()
	fmt.Printf("Disc ID       : %v\n", disc.Id())
	fmt.Printf("FreeDB ID     : %v\n", disc.FreedbId())
	fmt.Printf("TOC           : %v\n", disc.TocString())
	fmt.Printf("MCN           : %v\n", disc.Mcn())
	fmt.Printf("First track   : %v\n", disc.FirstTrackNum())
	fmt.Printf("Last track    : %v\n", disc.LastTrackNum())
	fmt.Printf("Sectors       : %v\n\n", disc.Sectors())

	for n := disc.FirstTrackNum(); n <= disc.LastTrackNum(); n++ {
		track := disc.Track(n)
		fmt.Printf("Track #%v:\n", track.Number)
		fmt.Printf("    ISRC   : %v\n", track.Isrc)
		fmt.Printf("    Offset : %v\n", track.Offset)
		fmt.Printf("    Sectors: %v\n", track.Sectors)
	}
}
