package tfnode

import "github.com/hashicorp/hcl2/hcl"

// Root is terraform root node that parsed.
type Root struct {
	Resources []Resource `hcl:"resource,block" json:"resource"`
	Providers []Provider `hcl:"provider,block" json:"provider"`
	Variables []Variable `hcl:"variable,block" json:"variable"`
	Locals    []Local    `hcl:"local,block" json:"local"`
	DataList  []Data     `hcl:"data,block" json:"data"`
	Remain    hcl.Body   `hcl:",remain" json:"-"` // extra fields
}

// https://www.terraform.io/docs/configuration/resources.html
type Resource struct {
	Label  string   `hcl:"name,label" json:"name"`
	Type   string   `hcl:"type,label" json:"type"`
	Remain hcl.Body `hcl:",remain" json:"-"`
}

type Provider struct {
	Label  string   `hcl:"name,label" json:"name"`
	Remain hcl.Body `hcl:",remain" json:"-"`
}

type Variable struct {
	Label  string   `hcl:"name,label" json:"name"`
	Remain hcl.Body `hcl:",remain" json:"-"`
}

// Local represents terraform local values block.
// https://www.terraform.io/docs/configuration/locals.html
type Local struct {
	Label  string   `hcl:"name,label" json:"name"`
	Remain hcl.Body `hcl:",remain" json:"-"`
}

// Data represents terraform data source block.
// https://www.terraform.io/docs/configuration/data-sources.html
type Data struct {
	Label  string   `hcl:"name,label"`
	Remain hcl.Body `hcl:",remain" json:"-"`
}
