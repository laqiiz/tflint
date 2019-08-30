package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/laqiiz/tfpolicy/linter"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

/**
c.f.
- syntax
  - https://www.terraform.io/docs/configuration/syntax.html
  - https://github.com/hashicorp/hcl2/blob/master/hcl/hclsyntax/spec.md
- good
  - https://github.com/hashicorp/hcl2
  - https://github.com/hashicorp/hcl2/blob/master/guide/go_decoding_gohcl.rst
*/

type Filter struct {
	Suffix []string
}

func (f Filter) Accept(path string) bool {
	for _, v := range f.Suffix {
		if strings.HasSuffix(path, v) {
			return true
		}
	}
	return false
}

func main() {

	dir := flag.String("d", ".", "d is target directory path")
	flag.StringVar(dir, "dir", ".", "dir is target directory path")
	flag.Parse()

	filter := Filter{
		// > The Terraform language uses configuration files that are named with the .tf file extension.
		// > There is also a JSON-based variant of the language that is named with the .tf.json file extension.
		// > https://www.terraform.io/docs/configuration/#code-organization
		// TODO tfvars
		Suffix: []string{".tf", ".tf.json"},
	}

	err := dirWalk(*dir, filter, linter.Validate)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%+v", err)
	}

}

func dirWalk(dir string, filter Filter, handler func(req linter.HCLRequest) error) error {
	var merr *multierror.Error

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, file := range files {
		if file.IsDir() {
			if err := dirWalk(filepath.Join(dir, file.Name()), filter, handler); err != nil {
				// TODO To separate system error and check style error
				merr = multierror.Append(merr, err)
			}
			continue
		}

		if filter.Accept(file.Name()) {
			src, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
			if err != nil {
				// unknown error, logic miss?
				return err
			}

			request := linter.HCLRequest{
				FilePath:    file.Name(),
				DisplayPath: file.Name(),
				Body:        src,
			}

			if err := handler(request); err != nil {
				// TODO To separate system error and check style error
				merr = multierror.Append(merr, err)
			}
		}
	}
	return merr
}
