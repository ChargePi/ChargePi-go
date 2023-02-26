## 💳 Reader hardware

All readers must implement the `Reader` interface. It is recommended that you implement the interface in a new file
named after the model of the reader in the `hardware/reader` package.

Then you should add a **constant** named after the **model** of the reader in the `reader` file in the package.

```golang
package reader

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
)

// Supported readers - by libnfc
const (
	PN532  = "PN532"
	ACR122 = "ACR122"
	PN533  = "PN533"
	BR500  = "BR500"
	R502   = "R502"
)

// Reader is an abstraction for an RFID/NFC tag reader.
type Reader interface {
	ListenForTags(ctx context.Context)
	Cleanup()
	Reset()
	GetTagChannel() <-chan string
	GetType() string
}
```