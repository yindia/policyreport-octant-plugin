package view

import "github.com/vmware-tanzu/octant/pkg/view/component"

var (
	fixHighSeverity = []component.TableCol{
		{Name: "Rule", Accessor: "Rule"},
		{Name: "Policy", Accessor: "Policy"},
		{Name: "Result", Accessor: "Result"},
		{Name: "Severity", Accessor: "Severity"},
	}
	policyReportCol = []component.TableCol{
		{Name: "Rule", Accessor: "Rule"},
		{Name: "Policy", Accessor: "Policy"},
		{Name: "Description", Accessor: "Description"},
		{Name: "Result", Accessor: "Result"},
		{Name: "Severity", Accessor: "Severity"},
		{Name: "Kind", Accessor: "Kind"},
		{Name: "Properties", Accessor: "Properties"},
	}
)
