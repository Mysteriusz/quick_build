package clang

import(
	. "qb/build"
	. "qb/policies/vc"

	"qb/policies/llvm/clang/vc"
)

var OUTPUT_DIFF_FUNCS = map[string]func(*QB_BuildState, *VC_FileState)(VC_ObjectDiff){
	"Compile": clang.ClangVCDiffOutForAll,
	"Link": clang.ClangVCDiffOutForAny,
}

