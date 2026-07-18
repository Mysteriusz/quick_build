package maps

import(
	"qb/policies"

	"qb/policies/llvm/ar"
	"qb/policies/llvm/clang"
	"qb/policies/shell"
)

var POLICY_INFO_LOOKUP map[string]policies.PolicyInfoInt = map[string]policies.PolicyInfoInt{
	"llvm-clang": &clang.PolicyInfo{},
	"llvm-ar": &ar.PolicyInfo{},
	"shell-exec": &shell.PolicyInfo{},
}

