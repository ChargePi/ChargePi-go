## Testing the client

## Running unit tests on software components

Testing the functionality of the client can be done manually, however, we can test some software components of the
client with unit tests. Unit test will compare the actual output of a function to the desired output. We parse
arguments/input to the function, and we expect a certain output/result. If the function does not produce a desired
result, the test will fail, and we know that we have a bug in our code (or test).

1. Copy the files from config directory and create test certificates:

   ```bash
   cd test/
   cp -r ../config/ .
   chmod +x create-test-certs.sh
   ./create-test-certs.sh
   ```

2. Run the tests:

    ```bash
   go test -v
   ```

### Writing unit tests

ChargePi-go uses [testify](https://github.com/stretchr/testify) to test the code. Check out their docs on how to write
tests.

An example test for a [connector](../../internal/components/connector/connector_test.go):

```golang
package test

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xBlaz3kx/ChargePi-go/chargepoint"
	"github.com/xBlaz3kx/ChargePi-go/data"
	"github.com/xBlaz3kx/ChargePi-go/hardware"
	"testing"
)

func TestConnector_ReserveConnector(t *testing.T) {
	require := require.New(t)

	connector1, err := chargepoint.NewConnector(1, 1, "Schuko",
		hardware.NewRelay(25, false), nil, false, 15)
	require.NoError(err)

	//ok case
	err = connector1.ReserveConnector(1)
	require.NoError(err)

	// connector already reserved
	err = connector1.ReserveConnector(2)
	require.Error(err)

	err = connector1.RemoveReservation()
	assert.NoError(t, err)

	// invalid connector status
	connector1.SetStatus(core.ChargePointStatusCharging, core.NoError)
	err = connector1.ReserveConnector(2)
	require.Error(err)
}
```

### Limits

As of now, we can only test certain segments of the client: AuthCache, Sessions, Connector and Settings. Hardware
testing must be done manually. If we wanted to test the client as a whole, it would be necessary to mock the central
system's responses and connectivity ([example](https://github.com/lorenzodonini/ocpp-go/tree/master/ocpp1.6_test)). 
