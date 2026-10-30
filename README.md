# goaplugins

A collection of plugins for [goa](https://goa.design) v3.

## skipresultvalidation

Stops generated server endpoints from validating results before returning them.

Since goa v3.31, the endpoints of methods that return a `ResultType` validate the projected result before returning it.
Results that do not satisfy the design validations (e.g. `Format`, `Pattern` or `Enum`) then fail at run time, while goa v3.30 and earlier returned them as is.
How the failure is reported depends on the goa version:

- goa v3.31 to v3.33 return the validation error as is. The default HTTP error encoder reports it as `400 Bad Request`.
- goa v3.34 and later report it as a server fault. The default HTTP error encoder reports it as `500 Internal Server Error`.

This plugin removes the validation call from the generated endpoints. The OpenAPI specification and the views package are unchanged.

```go
import (
	_ "github.com/ikawaha/goaplugins/skipresultvalidation"
	. "goa.design/goa/v3/dsl"
)
```

- With goa v3.30 and earlier, the plugin does nothing because the endpoints do not validate results.
- The plugin edits the endpoint template of goa. If the template references the result validation in a form the plugin does not recognize, code generation fails instead of silently keeping the validation.
- Generated clients still validate responses, as in goa v3.30 and earlier.
