package jiraapiexporter

type IssueCounts map[string]map[string]map[string]float64

func NewIssueCounts() IssueCounts {
	return make(map[string]map[string]map[string]float64)
}

func (c IssueCounts) Count(project string, release string, issuetype string) {
	proj, ok := c[project]
	if !ok {
		proj = make(map[string]map[string]float64)
		c[project] = proj
	}

	rel, ok := proj[release]
	if !ok {
		rel = make(map[string]float64)
		c[project][release] = rel
	}

	rel[issuetype] += 1
}
