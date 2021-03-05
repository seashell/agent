package structs

import "github.com/mitchellh/hashstructure/v2"

// DragoConfiguration :
type DragoConfiguration struct {
	Name    string
	DataDir string
	Servers []string
	Secret  string
	Meta    map[string]string
}

// Hash returns a unique hash of the struct
func (c *DragoConfiguration) Hash() uint64 {
	hash, err := hashstructure.Hash(c, hashstructure.FormatV2, nil)
	if err != nil {
		panic(err)
	}
	return hash
}

// NomadConfiguration :
type NomadConfiguration struct {
	Name                string
	DataDir             string
	InterfaceName       string
	InterfaceAddress    string
	PublicInterfaceName string
	Meta                map[string]string
}

// Hash returns a unique hash of the struct
func (c *NomadConfiguration) Hash() uint64 {
	hash, err := hashstructure.Hash(c, hashstructure.FormatV2, nil)
	if err != nil {
		panic(err)
	}
	return hash
}

// ConsulConfiguration :
type ConsulConfiguration struct {
	Name        string
	DataDir     string
	BindAddress string
	RetryJoin   string
	Meta        map[string]string
}

// Hash returns a unique hash of the struct
func (c *ConsulConfiguration) Hash() uint64 {
	hash, err := hashstructure.Hash(c, hashstructure.FormatV2, nil)
	if err != nil {
		panic(err)
	}
	return hash
}
