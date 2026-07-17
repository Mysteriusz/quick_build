package maps

import(
	"qb/build"
	"qb/policies/vc"
	"qb/policies/llvm/clang/vc"
)

var OUTPUT_DIFF_FUNCS = map[string]func(*qb.BuildState, *vc.FileState)(vc.ObjectDiff){
	"Compile": clang_vc.DiffOutForAll,
	"Link": clang_vc.DiffOutForAny,
}

