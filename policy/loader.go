package policy

import "regexp"

func Load(path string) []Policy {

	p1 := RegexRule{
		Selector:       NodeSelector{"$.resource[*].type"},
		RegexCondition: regexp.MustCompile("resource label must match pattern [a-zA-Z]+-[a-zA-Z]+"),
		Message:        "resource label must match pattern [a-zA-Z]+-[a-zA-Z]+",
	}

	return []Policy{
		p1,
	}
}
