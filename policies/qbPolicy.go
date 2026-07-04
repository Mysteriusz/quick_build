package policies

import(
	. "qb/io"
	. "qb/build"
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
	/*
		Execute policy defined version control check
		on the build state object

		The function should but isn`t required to check
		for Version control capability of it`s policy info

		IMPORTANT!
		This requires 'GetCapabilities' to return an object
		with field 'QB_Capabilities.VersionControl' == true
	*/
	RunVersionControl(_state *QB_BuildState) (not_updated bool)
}

