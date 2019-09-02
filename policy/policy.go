package policy

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/go-multierror"
	"github.com/laqiiz/tfpolicy/tfnode"
	"log"
	"reflect"
	"regexp"
)

type Policy interface {
	Accept(tfnode.Root) bool
	Violate(tfnode.Root) error
}

type RegexRule struct {
	Selector       NodeSelector
	RegexCondition *regexp.Regexp
	Message        string
}

type NodeSelector struct {
	HCLNodePath string
}

func (t NodeSelector) Select(tn tfnode.Root) (interface{}, error) {
	b, err := json.Marshal(tn)
	if err != nil {
		return nil, err
	}

	var input interface{}
	if err := json.Unmarshal(b, &input); err != nil {
		return nil, err
	}
	return jsonpath.Get(t.HCLNodePath, input)
}

func (r RegexRule) Accept(tn tfnode.Root) bool {
	s, err := r.Selector.Select(tn)
	if err != nil {
		log.Println("[ERROR] selector error: " + err.Error())
		return false
	}

	return s != nil
}

func (r RegexRule) Violate(tn tfnode.Root) error {
	var merr *multierror.Error

	s, err := r.Selector.Select(tn)
	if err != nil {
		return err
	}

	fmt.Printf("[DEBUG] selector result: %v\n",  s)

	v := reflect.ValueOf(s)
	switch v.Kind() {
	case reflect.Bool:
		return fmt.Errorf("unsupported type: %v", s)
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		return fmt.Errorf("unsupported type: %v", s)
	case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		return fmt.Errorf("unsupported type: %v", s)
	case reflect.Float32, reflect.Float64:
	case reflect.String:
		return fmt.Errorf("unsupported type: %v", s)
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			elm := v.Index(i)

			if !r.RegexCondition.MatchString(fmt.Sprintf("%v", elm)) {
				msg := fmt.Sprintf("%s: path='%s',v='%v'", r.Message, r.Selector.HCLNodePath, elm)
				merr = multierror.Append(merr, errors.New(msg))
			}
		}
	case reflect.Map:
		// Key Only Check
		for _, elm := range v.MapKeys() {
			if !r.RegexCondition.MatchString(fmt.Sprintf("%v", elm)) {
				msg := fmt.Sprintf("%s: path='%s',v='%v'", r.Message, r.Selector.HCLNodePath, elm)
				merr = multierror.Append(merr, errors.New(msg))
			}
		}
	default:
		return fmt.Errorf("unsupported type: %v", s)
	}

	return merr
}
