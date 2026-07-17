package clang

/*	

FIELDS:
	'Function' -> determines what execution flow the policy uses
	Acceptable 'Function' for Clang are:
		- Compile 			-> ClangCompileFromState
		- Link 				-> ClangLinkFromState

*/

type PolicyConfig struct{
	Function	string 	`toml:"function"`
	Vars		map[string]any `toml:"vars"`
}

const OUT_DEP string = "dependency_file"
const OUT_SRC string = "source_file"


