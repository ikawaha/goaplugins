package skipresultvalidation_test

import (
	"strings"
	"testing"

	"goa.design/goa/v3/codegen"

	"github.com/ikawaha/goaplugins/skipresultvalidation"
)

// endpointMethodV331 is an excerpt of the endpoint method template of goa
// v3.31.1 (codegen/service/templates/service_endpoint_method.go.tpl) with
// the three result validations.
const endpointMethodV331 = `{{- if .ViewedResult }}
		vres := {{ $.ViewedResult.Init.Declaration.Name }}(res, view)
		if err := {{ .ViewedResult.ViewsPkg }}.{{ .ViewedResult.Validate.Declaration.Name }}(vres); err != nil {
			return nil, err
		}
		return vres, nil
	{{- else if .ViewedResult }}
	vres := {{ $.ViewedResult.Init.Declaration.Name }}(res, view)
	if err := {{ .ViewedResult.ViewsPkg }}.{{ .ViewedResult.Validate.Declaration.Name }}(vres); err != nil {
		return nil, err
	}
	return vres, nil
{{- else if .ViewedResult }}
	vres := {{ $.ViewedResult.Init.Declaration.Name }}(res, view)
	if err := {{ .ViewedResult.ViewsPkg }}.{{ .ViewedResult.Validate.Declaration.Name }}(vres); err != nil {
		return nil, err
	}
	return vres, nil
{{- end }}
`

// endpointMethodV330 is an excerpt of the same template in goa v3.30.0,
// which does not validate results.
const endpointMethodV330 = `{{- if .ViewedResult }}
	vres := {{ $.ViewedResult.Init.Declaration.Name }}(res, view)
	return vres, nil
{{- end }}
`

func generate(t *testing.T, name, source string) (string, error) {
	t.Helper()
	s := &codegen.SectionTemplate{Name: name, Source: source}
	files := []*codegen.File{{Path: "gen/svc/endpoints.go", SectionTemplates: []*codegen.SectionTemplate{s}}}
	_, err := skipresultvalidation.Generate("", nil, files)
	return s.Source, err
}

func TestGenerate(t *testing.T) {
	t.Run("removes the result validations", func(t *testing.T) {
		got, err := generate(t, "endpoint-method", endpointMethodV331)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(got, "Validate") {
			t.Errorf("result validation remains:\n%s", got)
		}
		if n := strings.Count(got, "return vres, nil"); n != 3 {
			t.Errorf("got %d returns, want 3:\n%s", n, got)
		}
	})
	t.Run("does nothing without result validation", func(t *testing.T) {
		got, err := generate(t, "endpoint-method", endpointMethodV330)
		if err != nil {
			t.Fatal(err)
		}
		if got != endpointMethodV330 {
			t.Errorf("source changed:\n%s", got)
		}
	})
	t.Run("ignores other sections", func(t *testing.T) {
		got, err := generate(t, "service", endpointMethodV331)
		if err != nil {
			t.Fatal(err)
		}
		if got != endpointMethodV331 {
			t.Errorf("source changed:\n%s", got)
		}
	})
	t.Run("fails on an unknown form of the result validation", func(t *testing.T) {
		changed := strings.Replace(endpointMethodV331, "return nil, err", "return nil, fmt.Errorf(\"invalid result: %w\", err)", 1)
		if _, err := generate(t, "endpoint-method", changed); err == nil {
			t.Error("want an error")
		}
	})
}
