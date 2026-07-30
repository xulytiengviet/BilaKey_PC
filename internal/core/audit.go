package core

// CVNSSAudit is a stable, machine-readable summary of the embedded rule graph.
type CVNSSAudit struct {
	RuleVersion          string `json:"rule_version"`
	BaseRows             int    `json:"base_rows"`
	PatchEntries         int    `json:"patch_entries"`
	Codes                int    `json:"codes"`
	AmbiguityGroups      int    `json:"ambiguity_groups"`
	CanonicalPolicies    int    `json:"canonical_policies"`
	CriticalCollisions   int    `json:"critical_collisions"`
	MaxCandidatesPerCode int    `json:"max_candidates_per_code"`
	SilentOverwrite      int    `json:"silent_reverse_overwrite"`
}

func AuditCVNSS() CVNSSAudit {
	maxCandidates := 1
	for _, candidates := range cvnssCandidates {
		if len(candidates) > maxCandidates {
			maxCandidates = len(candidates)
		}
	}
	return CVNSSAudit{
		RuleVersion:          CVNSSRuleVersion,
		BaseRows:             CVNSSBaseRows,
		PatchEntries:         CVNSSPatchEntries,
		Codes:                len(cvnssToCQN),
		AmbiguityGroups:      CVNSSAmbiguityGroups,
		CanonicalPolicies:    CVNSSCanonicalPolicies,
		CriticalCollisions:   CVNSSCriticalCollisionGroups,
		MaxCandidatesPerCode: maxCandidates,
		SilentOverwrite:      0,
	}
}
