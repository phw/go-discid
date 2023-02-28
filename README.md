# Go bindings for MusicBrainz libdiscid
[![PkgGoDev](https://pkg.go.dev/badge/go.uploadedlobster.com/discid)](https://pkg.go.dev/go.uploadedlobster.com/discid)

**This project has been moved to https://git.sr.ht/~phw/go-discid**

## About
discid provides Go bindings for the MusicBrainz DiscID library [libdiscid](http://musicbrainz.org/doc/libdiscid).
It allows calculating DiscIDs (MusicBrainz and freedb) for Audio CDs. Additionally
the library can extract the MCN/UPC/EAN and the ISRCs from disc.

## Requirements
* libdiscid >= 0.6.0

## Usage

```go
package main

import (
	"fmt"
	"log"

	"go.uploadedlobster.com/discid"
)

func main() {
  // Specifying the device is optional. If set to an empty string a platform
  // specific default will be used.
  disc, err := discid.ReadFeatures("", discid.FeatureRead|discid.FeatureMcn)
  if err != nil {
    log.Fatal(err)
  }
  defer disc.Close()
  fmt.Printf("Disc ID: %v\n", disc.Id())
  fmt.Printf("MCN    : %v\n", disc.Mcn())
}
```

See the [API documentation](https://pkg.go.dev/go.uploadedlobster.com/discid) for details.

## Contribute
The source code for discid-sys is available on
[SourceHut](https://git.sr.ht/~phw/go-discid).

Please report any issues on the
[issue tracker](https://todo.sr.ht/~phw/discid-bindings).

Patches can be submitted to the [mailing list](https://lists.sr.ht/~phw/musicbrainz).
You can clone the repository directly on SourceHut and submit your changes
with the "Prepare patchset" button. Please see SourceHut's
[documentation for sending patches upstream](https://man.sr.ht/git.sr.ht/#sending-patches-upstream)
for details.

## License
discid Copyright (c) 2020-2023 by Philipp Wolfer <ph.wolfer@gmail.com>

discid is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

See [LICENSE](./LICENSE) for details.
