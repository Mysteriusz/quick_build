package policy_map

import(
	. "qb/policies"

	. "qb/policies/llvm/ar"
	. "qb/policies/llvm/clang"
	. "qb/policies/shell"
)

var POLICY_INFO_LOOKUP map[string]QB_PolicyInfo = map[string]QB_PolicyInfo{
	"llvm-clang": &Clang_Policy{
		PATH: "./policies/llvm/clang.toml",
		CAPS: QB_Capabilities{
			VersionControl: true,
		},
	},
	"llvm-ar": &Ar_Policy{
		PATH: "./policies/llvm/ar.toml",
		CAPS: QB_Capabilities{
			VersionControl: true,
		},
	},
	"shell-exec": &Shell_Policy{
		PATH: "./policies/shell/exec.toml",
		CAPS: QB_Capabilities{
			VersionControl: false,
		},
	},
}

