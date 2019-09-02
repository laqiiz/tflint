package linter

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/laqiiz/tfpolicy/policy"
	"github.com/laqiiz/tfpolicy/tfnode"
)

//var regex = regexp.MustCompile("[a-zA-Z]+-[a-zA-Z]+")
//var advice = "resource label must match pattern [a-zA-Z]+-[a-zA-Z]+"

type HCLRequest struct {
	FilePath    string
	DisplayPath string
	Body        []byte
}

type Linter struct {
	Policies []policy.Policy
}

func (l Linter) Validate(req HCLRequest) error {
	var merr *multierror.Error

	f, diags := hclparse.NewParser().ParseHCL(req.Body, req.FilePath, )
	if diags != nil {
		return diags
	}

	var root tfnode.Root
	if diags := gohcl.DecodeBody(f.Body, nil, &root); diags != nil {
		return diags
	}

	for _, p := range l.Policies {
		if !p.Accept(root) {
			continue
		}

		violate := p.Violate(root)
		if violate == nil {
			continue
		}

		// TODO wrap and add error message like below
		// msg := fmt.Sprintf("[ERROR] %s: %s: label='%s',type='%s'", req.DisplayPath, advice, v.Label, v.Type)
		merr = multierror.Append(merr, violate)
	}

	return merr
}
