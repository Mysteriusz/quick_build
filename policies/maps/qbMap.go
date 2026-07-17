package maps

import(
	"qb/policies"

	"qb/policies/llvm/ar"
	"qb/policies/llvm/clang"
	"qb/policies/shell"
)

var POLICY_INFO_LOOKUP map[string]policies.PolicyInfo = map[string]policies.PolicyInfo{
	"llvm-clang": &clang.Policy{
		PATH: "./policies/llvm/clang.toml",
		CAPS: policies.Capabilities{
			VersionControl: true,
		},
	},
	"llvm-ar": &ar.Policy{
		PATH: "./policies/llvm/ar.toml",
		CAPS: policies.Capabilities{
			VersionControl: true,
		},
	},
	"shell-exec": &shell.Policy{
		PATH: "./policies/shell/exec.toml",
		CAPS: policies.Capabilities{
			VersionControl: false,
		},
	},
}

