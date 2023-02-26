## 🖥️ Display hardware

All displays must implement the `Display` interface. It is recommended that you implement the interface in a new file named
after the model of the display/LCD in the `hardware/display` package. 

Then you should add a **constant** named after the **model** of the display in the `display` file in the package and add a switch case with the implementation and the
necessary logic that returns a pointer to the struct.

```golang
package display

import (
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/notifications"
)

const (
	DriverHD44780 = "hd44780"
)

// Display is an abstraction layer for concrete implementation of a display.
type Display interface {
	DisplayMessage(message notifications.Message)
	Cleanup()
	Clear()
	GetType() string
}

```