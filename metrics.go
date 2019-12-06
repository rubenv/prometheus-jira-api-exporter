package jiraapiexporter

import "github.com/prometheus/client_golang/prometheus"

var issuesMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "jira",
		Subsystem: "issues",
		Name:      "total",
		Help:      "Issue counts",
	},
	[]string{
		"projectkey",
		"project",
		"release",
		"type",
	},
)

func init() {
	prometheus.MustRegister(issuesMetric)
}
