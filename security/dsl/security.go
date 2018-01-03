package dsl

import (
	goadesign "goa.design/goa/design"
	"goa.design/goa/eval"
	"goa.design/goa/http/dsl"
	"goa.design/plugins/security/design"
)

// Description sets the expression description.
//
// Description must appear in API, Docs, Type, Attribute, BasicAuthSecurity,
// APIKeySecurity, OAuth2Security or JWTSecurity.
//
// Description accepts one arguments: the description string.
//
// Example:
//
//    API("adder", func() {
//        Description("Adder API")
//    })
//
func Description(d string) {
	switch expr := eval.Current().(type) {
	case *design.SchemeExpr:
		expr.Description = d
	default:
		dsl.Description(d)
	}
}

// BasicAuthSecurity defines a basic authentication security scheme.
//
// BasicAuthSecurity is a top level DSL.
//
// BasicAuthSecurity takes a name as first argument and an optional DSL as
// second argument.
//
// Example:
//
//     var Basic = BasicAuthSecurity("password", func() {
//         Description("Use your own password!")
//     })
//
func BasicAuthSecurity(name string, dsl ...func()) *design.SchemeExpr {
	if _, ok := eval.Current().(eval.TopExpr); !ok {
		eval.IncompatibleDSL()
		return nil
	}

	if securitySchemeRedefined(name) {
		return nil
	}

	expr := &design.SchemeExpr{
		Kind:       design.BasicAuthKind,
		SchemeName: name,
	}

	if len(dsl) != 0 {
		if !eval.Execute(dsl[0], expr) {
			return nil
		}
	}

	design.Root.Schemes = append(design.Root.Schemes, expr)

	return expr
}

// APIKeySecurity defines an API key security scheme where a key must be
// provided by the client to perform authorization.
//
// APIKeySecurity is a top level DSL.
//
// APIKeySecurity takes a name as first argument and an optional DSL as
// second argument.
//
// Example:
//
//    var APIKey = APIKeySecurity("key", func() {
//          Description("Shared secret")
//    })
//
func APIKeySecurity(name string, dsl ...func()) *design.SchemeExpr {
	if _, ok := eval.Current().(eval.TopExpr); !ok {
		eval.IncompatibleDSL()
		return nil
	}

	if securitySchemeRedefined(name) {
		return nil
	}

	expr := &design.SchemeExpr{
		Kind:       design.APIKeyKind,
		SchemeName: name,
	}

	if len(dsl) != 0 {
		if !eval.Execute(dsl[0], expr) {
			return nil
		}
	}

	design.Root.Schemes = append(design.Root.Schemes, expr)

	return expr
}

// OAuth2Security defines an OAuth2 security scheme. The DSL provided as second
// argument defines the specific flow supported by the scheme, one of
// ImplicitFlow, PasswordFlow, ClientCredentialsFlow or AuthorizationCodeFlow.
// The DSL also defines the scopes that may be associated with the incoming
// request tokens.
//
// OAuth2Security is a top level DSL.
//
// OAuth2Security takes a name as first argument and a DSL as second argument.
//
// Example:
//
//    var OAuth2 = OAuth2Security("googAuth", func() {
//        ImplicitFlow("/authorization")
//
//        Scope("api:write", "Write acess")
//        Scope("api:read", "Read access")
//    })
//
func OAuth2Security(name string, dsl ...func()) *design.SchemeExpr {
	if _, ok := eval.Current().(eval.TopExpr); !ok {
		eval.IncompatibleDSL()
		return nil
	}

	if securitySchemeRedefined(name) {
		return nil
	}

	expr := &design.SchemeExpr{
		SchemeName: name,
		Kind:       design.OAuth2Kind,
	}

	if len(dsl) != 0 {
		if !eval.Execute(dsl[0], expr) {
			return nil
		}
	}

	design.Root.Schemes = append(design.Root.Schemes, expr)

	return expr
}

// JWTSecurity defines an HTTP security scheme where a JWT is passed in the
// request Authorization header as a bearer token to perform auth. This scheme
// supports defining scopes that endpoint may require to authorize the request.
// The scheme also supports specifying a token URL used to retrieve token
// values.
//
// Since scopes are not compatible with the Swagger specification, the swagger
// generator inserts comments in the description of the different elements on
// which they are defined.
//
// JWTSecurity is a top level DSL.
//
// JWTSecurity takes a name as first argument and an optional DSL as second
// argument.
//
// Example:
//
//    var JWT = JWTSecurity("jwt", func() {
//        Scope("my_system:write", "Write to the system")
//        Scope("my_system:read", "Read anything in there")
//    })
//
func JWTSecurity(name string, dsl ...func()) *design.SchemeExpr {
	if _, ok := eval.Current().(eval.TopExpr); !ok {
		eval.IncompatibleDSL()
		return nil
	}

	if securitySchemeRedefined(name) {
		return nil
	}

	expr := &design.SchemeExpr{
		SchemeName: name,
		Kind:       design.JWTKind,
	}

	if len(dsl) != 0 {
		if !eval.Execute(dsl[0], expr) {
			return nil
		}
	}

	design.Root.Schemes = append(design.Root.Schemes, expr)

	return expr
}

// Security defines authentication requirements to access an API, a service or a
// service endpoint.
//
// The requirement refers to one or more OAuth2Security, BasicAuthSecurity,
// APIKeySecurity or JWTSecurity security scheme. If the schemes include a
// OAuth2Security or JWTSecurity scheme then required scopes may be listed by
// name in the Security DSL. All the listed schemes must be validated by the
// client for the request to be authorized. Security may appear multiple times
// in the same scope in which case the client may validate any one of the
// requirements for the request to be authorized.
//
// Security must appear in a API, Service or Method expression.
//
// Security accepts an arbitrary number of security schemes as argument
// specified by name or by reference and an optional DSL function as last
// argument.
//
// Examples:
//
//    var _ = API("calc", func() {
//        // All API endpoints are secured via basic auth by default.
//        Security(BasicAuth)
//    })
//
//    var _ = Service("calculator", func() {
//        // Override default API security requirements. Accept either basic
//        // auth or OAuth2 access token with "api:read" scope.
//        Security(BasicAuth)
//        Security("oauth2", func() {
//            Scope("api:read")
//        })
//
//        Method("add", func() {
//            Description("Add two operands")
//
//            // Override default service security requirements. Require
//            // both basic auth and OAuth2 access token with "api:write"
//            // scope.
//            Security(BasicAuth, "oauth2", func() {
//                Scope("api:write")
//            })
//
//            Payload(Operands)
//            Error(ErrBadRequest, ErrorResult)
//        })
//
//        Method("health-check", func() {
//            Description("Check health")
//
//            // Remove need for authorization for this endpoint.
//            NoSecurity()
//
//            Payload(Operands)
//            Error(ErrBadRequest, ErrorResult)
//        })
//    })
//
func Security(args ...interface{}) {
	var dsl func()
	{
		if d, ok := args[len(args)-1].(func()); ok {
			args = args[:len(args)-1]
			dsl = d
		}
	}

	var schemes []*design.SchemeExpr
	{
		schemes = make([]*design.SchemeExpr, len(args))
		for i, arg := range args {
			switch val := arg.(type) {
			case string:
				for _, s := range design.Root.Schemes {
					if s.SchemeName == val {
						schemes[i] = s
						break
					}
				}
				if schemes[i] == nil {
					eval.ReportError("security scheme %q not found", val)
					return
				}
			case *design.SchemeExpr:
				schemes[i] = val
			default:
				eval.InvalidArgError("security scheme or security scheme name", val)
				return
			}
		}
	}

	security := &design.SecurityExpr{Schemes: schemes}
	if dsl != nil {
		if !eval.Execute(dsl, security) {
			return
		}
	}

	current := eval.Current()
	switch actual := current.(type) {
	case *goadesign.MethodExpr:
		sec := &design.EndpointSecurityExpr{SecurityExpr: security, Method: actual}
		design.Root.EndpointSecurity = append(design.Root.EndpointSecurity, sec)
	case *goadesign.ServiceExpr:
		sec := &design.ServiceSecurityExpr{SecurityExpr: security, Service: actual}
		design.Root.ServiceSecurity = append(design.Root.ServiceSecurity, sec)
	case *goadesign.APIExpr:
		design.Root.APISecurity = append(design.Root.APISecurity, security)
	default:
		eval.IncompatibleDSL()
		return
	}
}

// NoSecurity removes the need for an endpoint to perform authorization.
//
// NoSecurity must appear in Method.
func NoSecurity() {
	security := &design.SecurityExpr{
		Schemes: []*design.SchemeExpr{
			&design.SchemeExpr{Kind: design.NoKind},
		},
	}

	current := eval.Current()
	switch actual := current.(type) {
	case *goadesign.MethodExpr:
		sec := &design.EndpointSecurityExpr{SecurityExpr: security, Method: actual}
		design.Root.EndpointSecurity = append(design.Root.EndpointSecurity, sec)
	default:
		eval.IncompatibleDSL()
		return
	}
}

// Username defines the attribute used to provide the username to an endpoint
// secured with basic authentication. The parameters and usage of Username are
// the same as the goa DSL Attribute function.
//
// The generated code produced by goa uses the value of the corresponding
// payload field to compute the basic authentication Authorization header value.
//
// Username must appear in Payload or Type.
//
// Example:
//
//    Method("login", func() {
//        Security(Basic)
//        Payload(func() {
//            Username("user", String)
//            Password("pass", String)
//        })
//        HTTP(func() {
//            // The "Authorization" header is defined implicitly.
//            POST("/login")
//        })
//    })
//
func Username(name string, args ...interface{}) {
	args = useDSL(args, func() { dsl.Metadata("security:username") })
	dsl.Attribute(name, args...)
}

// Password defines the attribute used to provide the password to an endpoint
// secured with basic authentication. The parameters and usage of Password are
// the same as the goa DSL Attribute function.
//
// The generated code produced by goa uses the value of the corresponding
// payload field to compute the basic authentication Authorization header value.
//
// Password must appear in Payload or Type.
//
// Example:
//
//    Method("login", func() {
//        Security(Basic)
//        Payload(func() {
//            Username("user", String)
//            Password("pass", String)
//        })
//        HTTP(func() {
//            // The "Authorization" header is defined implicitly.
//            POST("/login")
//        })
//    })
//
func Password(name string, args ...interface{}) {
	args = useDSL(args, func() { dsl.Metadata("security:password") })
	dsl.Attribute(name, args...)
}

// APIKey defines the attribute used to provide the API key to an endpoint
// secured with API keys. The parameters and usage of APIKey are the same as the
// goa DSL Attribute function except that it accepts an extra first argument
// corresponding to the name of the API key security scheme.
//
// The generated code produced by goa uses the value of the corresponding
// payload field to set the API key value.
//
// APIKey must appear in Payload or Type.
//
// Example:
//
//    Method("secured_read", func() {
//        Security(APIKeyAuth)
//        Payload(func() {
//            APIKey("api_key", "key", String, "API key used to perform authorization")
//            Required("key")
//        })
//        Result(String)
//        HTTP(func() {
//            GET("/")
//            Param("key:k") // Provide the key as a query string param "k"
//        })
//    })
//
//    Method("secured_write", func() {
//        Security(APIKeyAuth)
//        Payload(func() {
//            APIKey("api_key", "key", String, "API key used to perform authorization")
//            Attribute("data", String, "Data to be written")
//            Required("key", "data")
//        })
//        HTTP(func() {
//            POST("/")
//            Header("key:Authorization") // Provide the key as Authorization header (default)
//        })
//    })
//
func APIKey(scheme, name string, args ...interface{}) {
	args = useDSL(args, func() { dsl.Metadata("security:apikey:"+scheme, scheme) })
	dsl.Attribute(name, args...)
}

// AccessToken defines the attribute used to provide the access token to an
// endpoint secured with OAuth2. The parameters and usage of AccessToken are the
// same as the goa DSL Attribute function.
//
// The generated code produced by goa uses the value of the corresponding
// payload field to initialize the Authorization header.
//
// AccessToken must appear in Payload or Type.
//
// Example:
//
//    Method("secured", func() {
//        Security(OAuth2)
//        Payload(func() {
//            AccessToken("token", String, "OAuth2 access token used to perform authorization")
//            Required("token")
//        })
//        Result(String)
//        HTTP(func() {
//            // The "Authorization" header is defined implicitly.
//            GET("/")
//        })
//    })
//
func AccessToken(name string, args ...interface{}) {
	args = useDSL(args, func() { dsl.Metadata("security:accesstoken") })
	dsl.Attribute(name, args...)
}

// Token defines the attribute used to provide the JWT to an endpoint secured
// via JWT. The parameters and usage of Token are the same as the goa DSL
// Attribute function.
//
// The generated code produced by goa uses the value of the corresponding
// payload field to initialize the Authorization header.
//
// Example:
//
//    Method("secured", func() {
//        Security(JWT)
//        Payload(func() {
//            Token("token", String, "JWT token used to perform authorization")
//            Required("token")
//        })
//        Result(String)
//        HTTP(func() {
//            // The "Authorization" header is defined implicitly.
//            GET("/")
//        })
//    })
//
func Token(name string, args ...interface{}) {
	args = useDSL(args, func() { dsl.Metadata("security:token") })
	dsl.Attribute(name, args...)
}

// Scope has two uses: in JWTSecurity or OAuth2Security it defines a scope
// supported by the scheme. In Security it lists required scopes.
//
// Scope must appear in Security, JWTSecurity or OAuth2Security.
//
// Scope accepts one or two arguments: the first argument is the scope name and
// when used in JWTSecurity or OAuth2Security the second argument is a
// description.
//
// Example:
//
//    var JWT = JWTSecurity("JWT", func() {
//        Scope("api:read", "Read access") // Defines a scope
//        Scope("api:write", "Write access")
//    })
//
//    Method("secured", func() {
//        Security(JWT, func() {
//            Scope("api:read") // Required scope for auth
//        })
//    })
//
func Scope(name string, desc ...string) {
	switch current := eval.Current().(type) {
	case *design.SecurityExpr:
		if len(desc) >= 1 {
			eval.ReportError("too many arguments")
			return
		}
		current.Scopes = append(current.Scopes, name)
	case *design.SchemeExpr:
		if len(desc) > 1 {
			eval.ReportError("too many arguments")
			return
		}
		d := "no description"
		if len(desc) == 1 {
			d = desc[0]
		}
		current.Scopes = append(current.Scopes,
			&design.ScopeExpr{Name: name, Description: d})
	default:
		eval.IncompatibleDSL()
	}
}

// AuthorizationCodeFlow defines an authorizationCode OAuth2 flow as described
// in section 1.3.1 of RFC 6749.
//
// AuthorizationCodeFlow must be used in OAuth2Security.
//
// AuthorizationCodeFlow accepts three arguments: the authorization, token and
// refresh URLs.
func AuthorizationCodeFlow(authorizationURL, tokenURL, refreshURL string) {
	current, ok := eval.Current().(*design.SchemeExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	if current.Kind != design.OAuth2Kind {
		eval.ReportError("cannot specify flow for non-oauth2 security scheme.")
		return
	}
	current.Flows = append(current.Flows, &design.FlowExpr{
		Kind:             design.AuthorizationCodeFlowKind,
		AuthorizationURL: authorizationURL,
		TokenURL:         tokenURL,
		RefreshURL:       refreshURL,
	})
}

// ImplicitFlow defines an implicit OAuth2 flow as described in section 1.3.2
// of RFC 6749.
//
// ImplicitFlow must be used in OAuth2Security.
//
// ImplicitFlow accepts two arguments: the authorization and refresh URLs.
func ImplicitFlow(authorizationURL, refreshURL string) {
	current, ok := eval.Current().(*design.SchemeExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	if current.Kind != design.OAuth2Kind {
		eval.ReportError("cannot specify flow for non-oauth2 security scheme.")
		return
	}
	current.Flows = append(current.Flows, &design.FlowExpr{
		Kind:             design.ImplicitFlowKind,
		AuthorizationURL: authorizationURL,
		RefreshURL:       refreshURL,
	})
}

// PasswordFlow defines an Resource Owner Password Credentials OAuth2 flow as
// described in section 1.3.3 of RFC 6749.
//
// PasswordFlow must be used in OAuth2Security.
//
// PasswordFlow accepts two arguments: the token and refresh URLs.
func PasswordFlow(tokenURL, refreshURL string) {
	current, ok := eval.Current().(*design.SchemeExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	if current.Kind != design.OAuth2Kind {
		eval.ReportError("cannot specify flow for non-oauth2 security scheme.")
		return
	}
	current.Flows = append(current.Flows, &design.FlowExpr{
		Kind:       design.PasswordFlowKind,
		TokenURL:   tokenURL,
		RefreshURL: refreshURL,
	})
}

// ClientCredentialsFlow defines an clientCredentials OAuth2 flow as described
// in section 1.3.4 of RFC 6749.
//
// ClientCredentialsFlow must be used in OAuth2Security.
//
// ClientCredentialsFlow accepts two arguments: the token and refresh URLs.
func ClientCredentialsFlow(tokenURL, refreshURL string) {
	current, ok := eval.Current().(*design.SchemeExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	if current.Kind != design.OAuth2Kind {
		eval.ReportError("cannot specify flow for non-oauth2 security scheme.")
		return
	}
	current.Flows = append(current.Flows, &design.FlowExpr{
		Kind:       design.ClientCredentialsFlowKind,
		TokenURL:   tokenURL,
		RefreshURL: refreshURL,
	})
}

func securitySchemeRedefined(name string) bool {
	for _, s := range design.Root.Schemes {
		if s.SchemeName == name {
			eval.ReportError("cannot redefine security scheme with name %q", name)
			return true
		}
	}
	return false
}

// useDSL modifies the Attribute function to use the given function as DSL,
// merging it with any pre-exsiting DSL.
func useDSL(args []interface{}, d func()) []interface{} {
	ds, ok := args[len(args)-1].(func())
	if ok {
		newdsl := func() { ds(); d() }
		args = append(args[:len(args)-1], newdsl)
	} else {
		args = append(args, d)
	}
	return args
}
