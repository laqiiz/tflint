package linter

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/laqiiz/tfpolicy/tfnode"
	"log"
	"regexp"
)

var regex = regexp.MustCompile("[a-zA-Z]+-[a-zA-Z]+")
var advice = "resource label must match pattern [a-zA-Z]+-[a-zA-Z]+"

type HCLRequest struct {
	FilePath    string
	DisplayPath string
	Body        []byte
}

func Validate(req HCLRequest) error {
	var merr *multierror.Error

	f, diags := hclparse.NewParser().ParseHCL(req.Body, req.FilePath, )
	if diags != nil {
		return diags
	}

	var root tfnode.Root
	if diags := gohcl.DecodeBody(f.Body, nil, &root); diags != nil {
		return diags
	}

	log.Printf("root=%+v\n", root)

	for _, v := range root.Variables {
		// check label name style
		if !regex.MatchString(v.Label) {
			msg := fmt.Sprintf("unmatch format %s", v.Label)
			merr = multierror.Append(merr, errors.New(msg))
		}
	}

	for _, v := range root.Resources {
		// check label name style
		if !regex.MatchString(v.Label) {
			msg := fmt.Sprintf("[ERROR] %s: %s: label='%s',type='%s'", req.DisplayPath, advice, v.Label, v.Type)
			merr = multierror.Append(merr, errors.New(msg))
		}
	}

	return merr
}
