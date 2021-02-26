## Seashell API Client

This directory contains the `api` package which aims at providing programmatic access to Seashell's HTTP API.

### Documentation

...

### Usage

```go
package main

import "github.com/seashell/agent/api"

func main() {
	// Get a new client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}   
	
	// Get a handle to the networks API
	networks := client.Networks()

	// Create a new network
	n := &api.Network{
		Name: "my-new-network",
		IPAddressRange: "10.1.1.0/24"
	}

	id, err := networks.Create(context.Background(), n)
	if err != nil {
		panic(err)
	}

    	...
}
```

To run this example, start a Seashell server:

```
seashell agent --server
```

Copy the code above into a file such as `main.go`, and run it.

After running the code, you can also view the values in the Seashell UI on your local machine at http://localhost:8080/ui/
