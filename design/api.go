package design

import (
	. "goa.design/goa/v3/dsl"
	cors "goa.design/plugins/v3/cors/dsl"
)

var _ = API("foo", func() {
	Title("Foo Service")

	cors.Origin("*.example.com", func() {
		cors.Headers("X-Api-Version", "X-Shared-Secret")
		cors.MaxAge(100)
		cors.Credentials()
	})
})

var _ = Service("foo", func() {
	Method("foo1", func() {
		Payload(Int)
		Result(Int)
		HTTP(func() {
			POST("/foo1")
		})
	})
	Method("foo2", func() {
		Payload(Int)
		Result(Int)
		HTTP(func() {
			POST("/foo2")
		})
	})
	Method("foo3", func() {
		Payload(Int)
		Result(Int)
		HTTP(func() {
			POST("/foo3")
		})
	})

	Method("fooOptions", func() {
		HTTP(func() {
			OPTIONS("/foo1")
			Response(StatusOK)
		})
	})
})