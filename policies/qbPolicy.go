package policies

import(
	"qb/qbio"
	"qb/build"
	"qb/policies/vc"
)

type Capabilities struct{
	VersionControl 		bool
} 
type PolicyInfo interface{
	/*
		Static struct that defines capabilities
		supported by the policy
	*/
	GetCapabilities() Capabilities
	/*
		Formatted file that defines all the policies
	*/
	GetFile() *qbio.File
	/*
		Execute policy on the build state object
	*/
	Run(_state *qb.BuildState) bool

	vc.PolicyInt
}

