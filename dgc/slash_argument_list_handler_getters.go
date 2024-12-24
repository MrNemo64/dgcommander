package dgc

// Boolean

func (args *slashCommandArgumentList) GetRequiredBool(name string) bool {
	arg, found := args.values[name]
	if !found {
		panic(ErrMissingRequiredArgument.withArgs("boolean", name))
	}
	if value, ok := arg.(bool); ok {
		return value
	}
	panic(ErrArgumentIsWrongType.withArgs("boolean", name, arg, arg))
}

func (args *slashCommandArgumentList) GetBool(name string) (value bool, found bool) {
	arg, found := args.values[name]
	if !found {
		return false, false
	}
	if value, ok := arg.(bool); ok {
		return value, true
	}
	panic(ErrArgumentIsWrongType.withArgs("boolean", name, arg, arg))
}

func (args *slashCommandArgumentList) GetBoolOrDefault(name string, def bool) bool {
	if value, found := args.GetBool(name); found {
		return value
	}
	return def
}

// Number

func (args *slashCommandArgumentList) GetRequiredNumber(name string) float64 {
	arg, found := args.values[name]
	if !found {
		panic(ErrMissingRequiredArgument.withArgs("number", name))
	}
	if value, ok := arg.(float64); ok {
		return value
	}
	panic(ErrArgumentIsWrongType.withArgs("number", name, arg, arg))
}

func (args *slashCommandArgumentList) GetNumber(name string) (value float64, found bool) {
	arg, found := args.values[name]
	if !found {
		return 0, false
	}
	if value, ok := arg.(float64); ok {
		return value, true
	}
	panic(ErrArgumentIsWrongType.withArgs("number", name, arg, arg))
}

func (args *slashCommandArgumentList) GetNumberOrDefault(name string, def float64) float64 {
	if value, found := args.GetNumber(name); found {
		return value
	}
	return def
}
