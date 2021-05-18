package goa_cors_override

import (
	"context"
	"log"

	foo2 "experiments/goa-cors-override/gen/foo"
)

// foo service example implementation.
// The example methods log the requests and return zero values.
type foosrvc struct {
	logger *log.Logger
}

// NewFoo returns the foo service implementation.
func NewFoo(logger *log.Logger) foo2.Service {
	return &foosrvc{logger}
}

// Foo1 implements foo1.
func (s *foosrvc) Foo1(ctx context.Context, p int) (res int, err error) {
	s.logger.Print("foo.foo1")
	return
}

// Foo2 implements foo2.
func (s *foosrvc) Foo2(ctx context.Context, p int) (res int, err error) {
	s.logger.Print("foo.foo2")
	return
}

// Foo3 implements foo3.
func (s *foosrvc) Foo3(ctx context.Context, p int) (res int, err error) {
	s.logger.Print("foo.foo3")
	return
}

// FooOptions implements fooOptions.
func (s *foosrvc) FooOptions(ctx context.Context) (err error) {
	s.logger.Print("foo.fooOptions")
	return
}
