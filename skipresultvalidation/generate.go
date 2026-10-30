// Package skipresultvalidation is a goa plugin that stops generated server
// endpoints from validating results before returning them.
//
// Since goa v3.31, the endpoints of methods that return a ResultType validate
// the projected result (views.ValidateXxx) before returning it. Results that
// do not satisfy the design validations (e.g. Format, Pattern or Enum) then
// fail at run time, while goa v3.30 and earlier returned them as is. goa v3.31
// to v3.33 return the validation error as is, and goa v3.34 and later report
// it as a server fault. This plugin removes the call from the generated
// endpoints and leaves the OpenAPI specification and the views package
// unchanged.
//
// Import the package for its side effect in the design package:
//
//	import _ "github.com/ikawaha/goaplugins/skipresultvalidation"
package skipresultvalidation

import (
	"fmt"
	"regexp"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
)

// sectionName is the name of the section that renders a service endpoint
// method in the generated endpoints.go files.
const sectionName = "endpoint-method"

var (
	// validationCall matches the result validation in the endpoint method
	// template of goa v3.31 and later. goa v3.31 to v3.33 return the
	// validation error as is, while goa v3.34 and later wrap it in a server
	// fault preceded by a comment.
	validationCall = regexp.MustCompile(`(?m)^\t*if err := \{\{ \.ViewedResult\.ViewsPkg \}\}\.\{\{ \.ViewedResult\.Validate\.Declaration\.Name \}\}\(vres\); err != nil \{\n(?:\t*//.*\n)*\t*return nil, (?:err|goa\.NewServiceError\(err, "fault", false, false, true\))\n\t*\}\n`)
	// validationRef matches any reference to the result validation in the
	// template. It detects template changes that validationCall no longer
	// matches.
	validationRef = regexp.MustCompile(`\.ViewedResult\.Validate\b`)
)

func init() {
	codegen.RegisterPlugin("skipresultvalidation", "gen", nil, Generate)
}

// Generate removes the result validation from the endpoint method sections.
// It does nothing for goa versions whose templates do not validate results
// (v3.30 and earlier) and returns an error if the template references the
// result validation in a form it does not recognize.
func Generate(_ string, _ []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, f := range files {
		for _, s := range f.SectionTemplates {
			if s.Name != sectionName {
				continue
			}
			refs := len(validationRef.FindAllStringIndex(s.Source, -1))
			if refs == 0 {
				continue
			}
			if calls := len(validationCall.FindAllStringIndex(s.Source, -1)); calls != refs {
				return nil, fmt.Errorf("skipresultvalidation: %s: found %d result validation calls but %d references in the %q template; the goa template may have changed", f.Path, calls, refs, sectionName)
			}
			s.Source = validationCall.ReplaceAllString(s.Source, "")
		}
	}
	return files, nil
}
