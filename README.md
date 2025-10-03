# domain-normalizer

Normalize and parse domain names from messy input. Cleans whitespace and stray symbols, preserves label boundaries, converts Unicode to ASCII (IDNA/punycode), and extracts eTLD+1 (e.g. `co.uk`), SLD, and subdomain.

## Install

```shell
go get github.com/pixel365/domain-normalizer
```

## Examples

**Normalize a free-form string**

```go
package main

import (
	"fmt"

	"github.com/pixel365/domain-normalizer"
)

func main() {
	fmt.Println(normalizer.Normalize("  SomE...DOMa in..com!!!  "))
	// Output: some.doma-in.com
}
```

**Parse (IDNA + public suffix)**

```go
package main

import (
	"fmt"

	"github.com/pixel365/domain-normalizer"
)

func main() {
	d, err := normalizer.Parse("a.bücher.de")
	if err != nil { panic(err) }

	fmt.Println(d.Normalized)   // a.bücher.de
	fmt.Println(d.ASCII)        // a.xn--bcher-kva.de
	fmt.Println(d.Tld)          // de
	fmt.Println(d.Registerable) // bücher.de
	fmt.Println(d.Sld)          // bücher
	fmt.Println(d.SubDomain)    // a
}
```


**Multipart public suffix (`co.uk`)**

```go
package main

import (
	"fmt"

	"github.com/pixel365/domain-normalizer"
)

func main() {
	d, _ := normalizer.Parse("x.example.co.uk")
	fmt.Println(d.Tld)          // co.uk
	fmt.Println(d.Registerable) // example.co.uk
	fmt.Println(d.SubDomain)    // x
}
```
