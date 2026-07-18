package qb

import(
	"errors"
)

type BuildError struct{
	Objects	[]string
	Pipe	string
	Err 	error
}

func(BuildError) None() BuildError{
	return BuildError{}
}
func(_err BuildError) New(_state *BuildState, _msg string, _objs ...string) BuildError{
	_err.Pipe = _state.CurrentPipe().CommandPolicyName
	_err.Err = errors.New(_msg)
	_err.Objects = _objs
	return _err
}

func(_err BuildError) NilArgument(_state *BuildState) BuildError{
		return BuildError{}.New(_state, "Nil arugment provided")
}

func(err BuildError) Check() bool{
	return err.Err != nil
}

