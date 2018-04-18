package document

// A Dependency representes the result of the DependencyParseAnnotator
// https://stanfordnlp.github.io/CoreNLP/depparse.html
type Dependency struct {
	Dep            string
	Governor       int    `json:"governor"`
	GovernorGloss  string `json:"governorGloss"`
	Dependent      int    `json:"dependent"`
	DependentGloss string `json:"dependentGloss"`
}
