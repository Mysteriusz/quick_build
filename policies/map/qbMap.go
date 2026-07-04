package policies

import(
	"fmt"

	. "qb/policies"
	. "qb/policies/llvm/ar"
	. "qb/policies/llvm/clang"
)

var LLVM_POLICY_INFO_LOOKUP map[string]QB_PolicyInfo = map[string]QB_PolicyInfo{
	"llvm-clang": &Clang_Policy{
		PATH: "./policies/llvm/clang.toml",
		CAPS: QB_Capabilities{
			VersionControl: true,
		},
	},
	"llvm-ar": &Ar_Policy{
		PATH: "./policies/llvm/ar.toml",
		CAPS: QB_Capabilities{
			VersionControl: false,
		},
	},
}


/*
	Get the policy config by file name and policy name
*/
func QBPolicyLookup(_policy_alias string) (policy QB_PolicyInfo, res bool){
	policy,found := LLVM_POLICY_INFO_LOOKUP[_policy_alias]
	if !found{
		fmt.Printf("Policy file '%s' not found\n", _policy_alias)
		return
	}

	return policy, true
}

