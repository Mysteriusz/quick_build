package build

func qbRefMapVar(_state *QB_BuildState, fld string) []string{
	switch(fld){
	case "STATE.WORKING_SET":
		return _state.WorkingSet.StringArray()
	case "STATE.GET_SOURCES":
		return _state.GatherAllSources().AllPaths()
	case "STATE.GET_HEADERS":
		return _state.GatherAllHeaders().AllPaths()
	case "STATE.SOURCE_DIRECTORY":
		return []string{_state.Config.SourceDirectory}
	case "STATE.OUTPUT_DIRECTORY":
		return []string{_state.Config.OutputDirectory}
	case "STATE.NAME":
		return []string{_state.Config.Name}
	default:
		return nil
	}
}

func QBRefConvertFld(_state *QB_BuildState, fld string) []string{
	if fld[0] != '$'{
		return []string{}
	}
	return qbRefMapVar(_state, fld[1:])
}

