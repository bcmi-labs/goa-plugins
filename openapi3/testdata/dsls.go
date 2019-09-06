package testdata

import (
	. "goa.design/goa/v3/dsl"
	_ "goa.design/plugins/v3/docs"
)

var SimpleDSL = func() {
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
	var _ = API("test", func() {
		Server("test", func() {
			Host("localhost", func() {
				URI("https://goa.design")
			})
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
