package dgc

type ArgumentHasInvalidValueError struct{ DGCError }

func (e ArgumentHasInvalidValueError) New(name string, foundValue any, expectedValueType string) ArgumentHasInvalidValueError {
	base := e.DGCError.withArgs(name, foundValue, foundValue, expectedValueType)
	return ArgumentHasInvalidValueError{DGCError: base}
}

func (e ArgumentHasInvalidValueError) Values() (name string, foundValue any, expectedValueType string) {
	if len(e.args) != 4 {
		panic("ArgumentHasInvalidValueError has not had its values specified yet")
	}
	return e.args[0].(string), e.args[1], e.args[3].(string)
}

func (e ArgumentHasInvalidValueError) Is(target error) bool {
	_, ok := target.(ArgumentHasInvalidValueError)
	return ok
}

type ArgumentHasNoValueError struct{ DGCError }

func (e ArgumentHasNoValueError) New(name string) ArgumentHasNoValueError {
	base := e.DGCError.withArgs(name)
	return ArgumentHasNoValueError{DGCError: base}
}

func (e ArgumentHasNoValueError) Values() (name string) {
	if len(e.args) != 1 {
		panic("ArgumentHasNoValueError has not had its values specified yet")
	}
	return e.args[0].(string)
}

func (e ArgumentHasNoValueError) Is(target error) bool {
	_, ok := target.(ArgumentHasNoValueError)
	return ok
}
