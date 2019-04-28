package main

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

// Directories is for parse
type Directories struct {
	Name        string `hcl:",key"`
	Description string
}

// Config is for
type Config struct {
	Variable []Directories
}

/// This sample test read/write from HCL file.
func main() {

	type Directories struct {
		Name        string `hcl:",key"`
		Description string
	}

	type Config struct {
		Variable []Directories
	}

	//  var conf map[string]map[string]string
	var conf Config
	b, err := ioutil.ReadFile("testdata/azure.tf")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
	if err := hcl.Decode(&conf, string(b)); err != nil {
		panic(err)
	}
	fmt.Println(conf.Variable[0].Description)
	fmt.Println(conf.Variable[0].Name)
}
