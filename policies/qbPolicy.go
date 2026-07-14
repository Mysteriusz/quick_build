package policies

import(
	. "qb/io"
	. "qb/build"
	. "qb/policies/vc"
)

type QB_Capabilities struct{
	VersionControl 		bool
} 
type QB_PolicyInfo interface{
	/*
		Static struct that defines capabilities
		supported by the policy
	*/
	GetCapabilities() QB_Capabilities
	/*
		Formatted file that defines all the policies
	*/
	GetFile() *QB_File
	/*
		Execute policy on the build state object
	*/
	Run(_state *QB_BuildState) bool

	VC_PolicyInt
}

