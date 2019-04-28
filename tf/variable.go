package tf

import "github.com/hashicorp/hcl2/hcl"

// Root is terraform root node that parsed.
type Root struct {
	Resources []Resource `hcl:"resource,block"`
	Providers []Provider `hcl:"provider,block"`
	Variables []Variable `hcl:"variable,block"`
	Locals    []Local    `hcl:"local,block"`
	DataList  []Data     `hcl:"data,block"`
	Remain    hcl.Body   `hcl:",remain"` // extra fields
}

// https://www.terraform.io/docs/configuration/resources.html
type Resource struct {
	Label  string   `hcl:"name,label"`
	Type   string   `hcl:"type,label"`
	Remain hcl.Body `hcl:",remain"`
}

type Provider struct {
	Label  string   `hcl:"name,label"`
	Remain hcl.Body `hcl:",remain"`
}

type Variable struct {
	Label  string   `hcl:"name,label"`
	Remain hcl.Body `hcl:",remain"`
}

// Local represents terraform local values block.
// https://www.terraform.io/docs/configuration/locals.html
type Local struct {
	Label  string   `hcl:"name,label"`
	Remain hcl.Body `hcl:",remain"`
}

// Data represents terraform data source block.
// https://www.terraform.io/docs/configuration/data-sources.html
type Data struct {
	Label  string   `hcl:"name,label"`
	Remain hcl.Body `hcl:",remain"`
}
