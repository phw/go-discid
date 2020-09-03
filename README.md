# Go bindings for MusicBrainz libdiscid
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/phw/go-discid?label=package%20version)](https://github.com/phw/go-discid/releases)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/phw/go-discid)](https://pkg.go.dev/github.com/phw/go-discid)
[![GitHub license](https://img.shields.io/github/license/phw/go-discid)](https://github.com/phw/go-discid/blob/master/LICENSE)

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

	"github.com/phw/go-discid"
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

See the [API documentation](https://pkg.go.dev/github.com/phw/go-discid) for details.

## Contribute
The source code for discid is available on
[GitHub](https://github.com/phw/go-discid).

Please report any issues on the
[issue tracker](https://github.com/phw/go-discid/issues).

## License
discid Copyright (c) 2020 by Philipp Wolfer <ph.wolfer@gmail.com>

discid is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

See LICENSE for details.
