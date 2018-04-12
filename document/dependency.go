package document

type Dependency struct {
	Dep            string
	Governor       int    `json:"governor"`
	GovernorGloss  string `json:"governorGloss"`
	Dependent      int    `json:"dependent"`
	DependentGloss string `json:"dependentGloss"`
}
