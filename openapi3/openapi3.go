package openapi3

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/expr"
)

func NewV3(r *expr.RootExpr) (openapi3.Swagger, error) {
	return openapi3.Swagger{
		OpenAPI: "3.0.0",
		Info:    info(r.API),
		Servers: servers(r.API.Servers),
		Paths:   paths(r),
	}, nil
}

func info(api *expr.APIExpr) openapi3.Info {
	version := "unversioned"
	if api.Version != "" {
		version = api.Version
	}
	return openapi3.Info{
		Title:          api.Title,
		Description:    api.Description,
		TermsOfService: api.TermsOfService,
		Contact:        contact(api.Contact),
		License:        license(api.License),
		Version:        version,
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

func servers(s []*expr.ServerExpr) []*openapi3.Server {
	if s == nil {
		return nil
	}

	servers := []*openapi3.Server{}

	for _, server := range s {
		for _, host := range server.Hosts {
			for _, uri := range host.URIs {
				servers = append(servers, &openapi3.Server{
					URL:         string(uri),
					Description: host.Description,
				})
			}

		}
	}

	return servers
}

func paths(r *expr.RootExpr) map[string]*openapi3.PathItem {
	paths := map[string]*openapi3.PathItem{}

	for _, service := range r.API.HTTP.Services {
		for _, endpoint := range service.HTTPEndpoints {
			for _, route := range endpoint.Routes {
				path, ok := paths[route.Path]
				if !ok {
					path = &openapi3.PathItem{}
				}

				operation := operation(service, endpoint, route)
				path.Get = operation

				paths[route.Path] = path
			}

			// paths[endpoint] = &openapi3.PathItem{
			// 	Get: operation("get", r),
			// }
		}
	}

	return paths
}

func operation(s *expr.HTTPServiceExpr, e *expr.HTTPEndpointExpr, r *expr.RouteExpr) *openapi3.Operation {
	params := paramsFromExpr(e.Params, r.Path)
	// params = append(params, paramsFromHeaders(e)...)

	responses := map[string]*openapi3.ResponseRef{}

	return &openapi3.Operation{
		OperationID: fmt.Sprintf("%s#%s", s.Name(), e.Name()),
		Description: r.Endpoint.Description(),
		Parameters:  params,
		Responses:   responses,
	}
}

func paramsFromExpr(params *expr.MappedAttributeExpr, path string) []*openapi3.ParameterRef {
	if params == nil {
		return nil
	}
	var (
		res       []*openapi3.ParameterRef
		wildcards = expr.ExtractHTTPWildcards(path)
		i         = 0
	)
	_ = codegen.WalkMappedAttr(params, func(n, pn string, required bool, at *expr.AttributeExpr) error {
		in := "query"
		for _, w := range wildcards {
			if n == w {
				in = "path"
				required = true
				break
			}
		}

		param := paramFor(at, pn, in, required)
		res = append(res, param)
		i++
		return nil
	})
	return res
}

func paramFor(at *expr.AttributeExpr, name, in string, required bool) *openapi3.ParameterRef {
	p := &openapi3.ParameterRef{
		Value: &openapi3.Parameter{
			In:          in,
			Name:        name,
			Description: at.Description,
			Required:    required,
		},
	}

	if expr.IsArray(at.Type) {
		true := true
		p.Value.Explode = &true
	}

	return p
}
