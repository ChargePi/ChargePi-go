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

### Writing unit tests in Go

More about testing in Golang [here](https://golang.org/doc/tutorial/add-a-test).

An example of a function test in Golang:

```golang
package test

import "testing"

func TestFunction(t *testing.T) {
	type args struct {
		testArgument string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "PositiveTestCase",
			args: args{
				testArgument: "123",
			},
			want: true,
		}, {
			name: "NegativeTestCase",
			args: args{
				testArgument: "1234",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := functionWeAreTestingThatReturnsBool(tt.args.testArgument); got != tt.want {
				t.Errorf("functionWeAreTestingThatReturnsBool() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

### Limits

As of now, we can only test certain segments of the client: AuthCache, Sessions, Connector and Settings. Hardware
testing must be done manually. If we wanted to test the client as a whole, it would be necessary to mock the central
system's responses and connectivity (todo, [example](https://github.com/lorenzodonini/ocpp-go/tree/master/ocpp1.6_test))
. 
