package policies

import(
	"qb/build"
)

type Capabilities struct{
	VersionControl 		bool
} 
type PolicyInfoInt interface{
	GetCapabilities() Capabilities
	GetFile() *PolicyFile
	/*
		Execute policy on the build state object
	*/
	Run(_state *qb.BuildState) qb.BuildError
}

