package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/laqiiz/tflint/tf"
	"github.com/pkg/errors"
)

/**
c.f.
- syntax
  - https://www.terraform.io/docs/configuration/syntax.html
  - https://github.com/hashicorp/hcl2/blob/master/hcl/hclsyntax/spec.md
  - https://buildmedia.readthedocs.org/media/pdf/hcl/guide/hcl.pdf
- GoDoc
  - https://godoc.org/github.com/hashicorp/hcl
- http://jen20.com/2015/09/07/using-hcl-part-1.html
- https://qiita.com/TsuyoshiUshio@github/items/ba4a2101be4784253cc5

- good
  - https://github.com/hashicorp/hcl2
  - https://github.com/hashicorp/hcl2/blob/master/guide/go_decoding_gohcl.rst
*/

var regex = regexp.MustCompile("[a-zA-Z]+-[a-zA-Z]+")
var advice = "resource label must match pattern [a-zA-Z]+-[a-zA-Z]+"

func main() {
	log.Println("Starting audit...")

	currentDir := pwd()
	log.Printf("dir=%s", currentDir)

	err := DirWalk(currentDir, CheckProcess)
	if merr, ok := err.(*multierror.Error); ok {
		if merr.Len() > 0 {
			log.Fatalf("%+v", merr)
		}
	}
	log.Println("Audit done.")
}

func pwd() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

// CheckProcess is style check main logic.
// This func is used by DirWalk like callback.
func CheckProcess(path string) error {
	// multiple error
	var result *multierror.Error

	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		log.Fatal(diags.Error())
	}

	var root tf.Root
	if diags := gohcl.DecodeBody(f.Body, nil, &root); diags.HasErrors() {
		log.Fatal(diags.Error())
	}
	log.Printf("root=%+v\n", root)

	for _, v := range root.Variables {
		// check label name style
		if !regex.MatchString(v.Label) {
			msg := fmt.Sprintf("unmatch format %s", v.Label)
			result = multierror.Append(result, errors.New(msg))
		}
	}

	for _, v := range root.Resources {
		// check label name style
		if !regex.MatchString(v.Label) {
			alternativePath := filepath.Clean(strings.Replace(path, pwd(), "", -1))
			msg := fmt.Sprintf("[ERROR] %s: %s: label='%s',type='%s'", alternativePath, advice, v.Label, v.Type)
			result = multierror.Append(result, errors.New(msg))
		}
	}
	return result.ErrorOrNil()
}

// DirWalk is Recursive directory search and returns terraform file paths.
func DirWalk(dir string, callback func(filePath string) error) error {
	// multiple error
	var merr *multierror.Error

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, file := range files {
		if file.IsDir() {
			if err := DirWalk(filepath.Join(dir, file.Name()), callback); err != nil {
				// TODO To separate system error and check style error
				merr = multierror.Append(merr, err)
			}
			continue
		}

		if strings.HasSuffix(file.Name(), ".tf") || strings.HasSuffix(file.Name(), ".tf.json") {
			// TODO tfvars
			// > The Terraform language uses configuration files that are named with the .tf file extension.
			// > There is also a JSON-based variant of the language that is named with the .tf.json file extension.
			// > https://www.terraform.io/docs/configuration/#code-organization
			if err := callback(filepath.Join(dir, file.Name())); err != nil {
				// TODO To separate system error and check style error
				merr = multierror.Append(merr, err)

			}
		}
	}
	return merr.ErrorOrNil()
}
