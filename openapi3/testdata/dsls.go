package testdata

import (
	. "goa.design/goa/v3/dsl"
	_ "goa.design/plugins/v3/docs"
)

var FullDSL = func() {
	var _ = API("test", func() {
		Title("test api")
		Description("an api to test openapi3")
		TermsOfService("https://example.com/tos")

		Contact(func() {
			Name("name")
			URL("https://example.com")
			Email("test@test.com")
		})

		License(func() {
			Name("license")
			URL("https://example.com/license")
		})

		Server("test", func() {
			Host("localhost", func() {
				URI("https://goa.design")
			})
		})
	})

	var PayloadT = Type("Payload", func() {
		Attribute("string", String, func() {
			Example("")
		})
	})
	var ResultT = Type("Result", func() {
		Attribute("string", String, func() {
			Example("")
		})
	})

	Service("testService", func() {
		Method("testEndpoint", func() {
			Payload(PayloadT)
			Result(ResultT)
			HTTP(func() {
				GET("/")
			})
		})
	})
}
