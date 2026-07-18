package maps

import(
	"qb/build"
	"qb/policies/vc"
	"qb/policies/llvm/clang/vc"
)

var DIFF_HDR_PROVIDERS = map[string]func(*qb.BuildState, *vc.FileState)(vc.FileDiff, qb.BuildError){
	"Compile": clang_vc.DiffHeaders,
	"Link": nil,
}
var DIFF_SRC_PROVIDERS = map[string]func(*qb.BuildState, *vc.FileState)(vc.FileDiff, qb.BuildError){
	"Compile": clang_vc.DiffSources,
	"Link": nil,
}
var DIFF_OUT_PROVIDERS = map[string]func(*qb.BuildState, *vc.FileState)(vc.ObjectDiff, qb.BuildError){
	"Compile": clang_vc.DiffOutForAll,
	"Link": clang_vc.DiffOutForAny,
}

