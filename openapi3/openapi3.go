package openapi3

import (
	"github.com/getkin/kin-openapi/openapi3"
	"goa.design/goa/v3/expr"
)

func NewV3(r *expr.RootExpr) (openapi3.Swagger, error) {
	return openapi3.Swagger{
		OpenAPI: "3.0.0",
		Info:    info(r.API),
	}, nil
}

func info(api *expr.APIExpr) openapi3.Info {
	return openapi3.Info{
		Title:          api.Title,
		Description:    api.Description,
		TermsOfService: api.TermsOfService,
		Contact:        contact(api.Contact),
		License:        license(api.License),
		Version:        api.Version,
	}
}

func contact(c *expr.ContactExpr) *openapi3.Contact {
	if c != nil {
		return &openapi3.Contact{
			Name:  c.Name,
			URL:   c.URL,
			Email: c.Email,
		}
	}

	return nil
}

func license(l *expr.LicenseExpr) *openapi3.License {
	if l != nil {
		return &openapi3.License{
			Name: l.Name,
			URL:  l.URL,
		}
	}

	return nil
}
