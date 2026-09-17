package policies

import(
	"qb/build"
)

type Capabilities struct{
	VersionControl 		bool
} 
type PolicyInfoInt interface{
	GetCapabilities() Capabilities
	GetFile(_search_directory string) *PolicyFile
	/*
		Execute policy on the build state object
	*/
	Run(_state *qb.BuildState) qb.BuildError
}

