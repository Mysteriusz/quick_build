package policies

import(
	. "qb/io"
	. "qb/build"
	. "qb/policies/version_control"
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
		Begin policy defined version check 
		transaction on the build state object

		The function should but isn`t required to check
		for Version control capability of it`s policy info

		IMPORTANT!
		This should require 'GetCapabilities' to return an object
		with field 'QB_Capabilities.VersionControl' == true

		Eles the not_updated value should always be 0 (false)
	*/
	BeginVersionControl(_state *QB_BuildState) (not_updated bool, _vc_state VC_FileState)
	/*
		Finish and save the version check transaction

		The function should but isn`t required to check
		for Version control capability of it`s policy info

		IMPORTANT!
		This should require 'GetCapabilities' to return an object
		with field 'QB_Capabilities.VersionControl' == true
	*/
	EndVersionControl(_state *QB_BuildState, _vc_state *VC_FileState)
}

