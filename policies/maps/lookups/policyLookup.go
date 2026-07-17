package lookup

import(
	"qb/policies"
	"qb/policies/maps"
)

/*
	Get the policy config by file name and policy name
*/
func PolicyLookup(_policy_alias string) (policies.PolicyInfo, bool){
	policy,found := maps.POLICY_INFO_LOOKUP[_policy_alias]
	return policy, found
}

