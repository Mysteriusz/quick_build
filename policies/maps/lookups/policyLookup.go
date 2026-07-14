package lookup

import(
	"fmt"

	. "qb/policies"
	. "qb/policies/maps"
)

/*
	Get the policy config by file name and policy name
*/
func QBPolicyLookup(_policy_alias string) (policy QB_PolicyInfo, res bool){
	policy,found := POLICY_INFO_LOOKUP[_policy_alias]
	if !found{
		fmt.Printf("Policy file '%s' not found\n", _policy_alias)
		return
	}

	return policy, true
}

