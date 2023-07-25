package test_test

import (
	"fmt"
	"github.com/ghodss/yaml"
	"os"
	"testing"
)

func TestDemo(t *testing.T) {
	file, err := os.ReadFile("conf.yml")
	if err != nil {
		t.Error(err)
	}
	data := T{}
	if err := yaml.Unmarshal(file, &data); err != nil {
		t.Error(err)
	}
	for _, v := range data.Rules {
		fmt.Println(v)
	}
}

type T struct {
	Global struct {
		AllowRPC       bool          `json:"AllowRPC"`
		ApiVersion     string        `json:"ApiVersion"`
		Date           string        `json:"Date"`
		DisableRules   []interface{} `json:"DisableRules"`
		EnableRules    []interface{} `json:"EnableRules"`
		MaxLogInput    int           `json:"MaxLogInput"`
		MaxRegexRuleID int           `json:"MaxRegexRuleID"`
		Mode           string        `json:"Mode"`
	} `json:"Global"`
	MaskRules []struct {
		IgnoreCharSet string   `json:"IgnoreCharSet,omitempty"`
		IgnoreKind    []string `json:"IgnoreKind,omitempty"`
		Length        int      `json:"Length,omitempty"`
		MaskType      string   `json:"MaskType"`
		Offset        int      `json:"Offset,omitempty"`
		Padding       int      `json:"Padding,omitempty"`
		Reverse       bool     `json:"Reverse,omitempty"`
		RuleName      string   `json:"RuleName"`
		Value         string   `json:"Value,omitempty"`
	} `json:"MaskRules"`
	Rules []struct {
		CnName      string `json:"CnName"`
		Description string `json:"Description"`
		Detect      struct {
			KDict []string      `json:"KDict,omitempty"`
			KReg  []interface{} `json:"KReg,omitempty"`
			VDict []string      `json:"VDict,omitempty"`
			VReg  []string      `json:"VReg,omitempty"`
		} `json:"Detect"`
		EnName  string `json:"EnName"`
		ExtInfo struct {
			CnGroup string `json:"CnGroup"`
			EnGroup string `json:"EnGroup"`
		} `json:"ExtInfo"`
		Filter struct {
			BAlgo []string      `json:"BAlgo"`
			BDict []interface{} `json:"BDict,omitempty"`
			BReg  []interface{} `json:"BReg,omitempty"`
		} `json:"Filter,omitempty"`
		InfoType string `json:"InfoType"`
		Level    string `json:"Level"`
		Mask     string `json:"Mask,omitempty"`
		RuleID   int    `json:"RuleID"`
		Verify   struct {
			CDict []string      `json:"CDict,omitempty"`
			CReg  []interface{} `json:"CReg,omitempty"`
			VAlgo []string      `json:"VAlgo,omitempty"`
		} `json:"Verify,omitempty"`
		GroupName string `json:"GroupName,omitempty"`
	} `json:"Rules"`
}
