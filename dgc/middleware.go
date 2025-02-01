package dgc

import "errors"

type MiddlewareCallChainError struct{ DGCError }

func (e MiddlewareCallChainError) New(errs ...error) MiddlewareCallChainError {
	base := e.DGCError.withArgs(errors.Join(errs...))
	return MiddlewareCallChainError{DGCError: base}
}

func (e MiddlewareCallChainError) Values() (errs []error) {
	if len(e.args) != 4 {
		panic("MiddlewareCallChainError has not had its values specified yet")
	}
	return e.args[0].(interface{ Unwrap() []error }).Unwrap()
}

func (e MiddlewareCallChainError) Is(target error) bool {
	_, ok := target.(MiddlewareCallChainError)
	return ok
}

var (
	ErrMiddlewareCallChain = MiddlewareCallChainError{DGCError: makeError("errors returned by the middlewares: %w")}
)

func newMiddlewareChain[T any](value *T, middlewares []func(*T, func()) error, handler func() error) middlewareChain[T] {
	return middlewareChain[T]{
		value:       value,
		middlewares: middlewares,
		handler:     handler,
	}
}

type middlewareChain[T any] struct {
	value                *T
	middlewares          []func(*T, func()) error
	handler              func() error
	index                int
	allMiddlewaresCalled bool
	errs                 []error
}

func (mc *middlewareChain[T]) startChain() error {
	mc.index = -1
	mc.errs = nil
	mc.allMiddlewaresCalled = false
	mc.next()
	if len(mc.errs) > 0 {
		return ErrMiddlewareCallChain.New(mc.errs...)
	}
	return nil
}

func (mc *middlewareChain[T]) next() {
	mc.index++
	if mc.index >= len(mc.middlewares) {
		mc.allMiddlewaresCalled = true
		if err := mc.handler(); err != nil {
			mc.errs = append(mc.errs, err)
		}
		return
	}
	middle := mc.middlewares[mc.index]
	if err := middle(mc.value, mc.next); err != nil {
		mc.errs = append(mc.errs, err)
	}
}
