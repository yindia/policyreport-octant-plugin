package view

import (
	"fmt"

	"github.com/vmware-tanzu/octant/pkg/view/component"
	corev1 "k8s.io/api/core/v1"
)

func GetProperties(r SingleReport) *component.Selectors {
	var data []component.Selector
	for k, v := range r.Result.Properties {
		data = append(data, component.NewLabelSelector(k, v))
	}
	s := component.NewSelectors(data)
	return s
}

func setRow(table *component.Table, r SingleReport,isNamespace bool) {
	if isNamespace {
		table.Add(component.TableRow{
			"Rule":        component.NewText(r.Result.Rule),
			"Policxwy":    component.NewText(r.Result.Policy),
			"Description": component.NewText(r.Result.Description),
			"Result":      component.NewText(string(r.Result.Result)),
			"Severity":    component.NewText(string(r.Result.Severity)),
			"Kind":        component.NewText(r.Subject.Kind),
			"Properties":  GetProperties(r),
			"Resource":    component.NewText(r.Subject.Name),
			"Namespace": component.NewText(r.Subject.Namespace),
		})
		return
	}
	table.Add(component.TableRow{
		"Rule":        component.NewText(r.Result.Rule),
		"Policxwy":    component.NewText(r.Result.Policy),
		"Description": component.NewText(r.Result.Description),
		"Result":      component.NewText(string(r.Result.Result)),
		"Severity":    component.NewText(string(r.Result.Severity)),
		"Kind":        component.NewText(r.Subject.Kind),
		"Properties":  GetProperties(r),
		"Resource":    component.NewText(r.Subject.Name),
	})
	return
}

func PrintPolicyReportTable(results *ReportView, s *corev1.ObjectReference, title string) *component.Table {
	cols := policyReportCol
	var isModule = true
	if s == nil {
		cols = append(cols, component.TableCol{
			Name: "Resource", Accessor: "Resource",
		})
	} else {
		isModule = false
		if s.Kind == "Namespace" {
			cols = append(cols, component.TableCol{
				Name: "Resource", Accessor: "Resource",
			})
			cols = append(cols, component.TableCol{
				Name: "Namespace", Accessor: "Namespace",
			})
		}
	}
	table := component.NewTableWithRows(title, title, cols, []component.TableRow{})
	for _, r := range results.Reports {
		if isModule {
			setRow(table, r,false)
			continue
		}

		if s.Kind == "Namespace" {
			setRow(table, r,true)
			continue
		}

		table.Add(component.TableRow{
			"Rule":        component.NewText(r.Result.Rule),
			"Policy":      component.NewText(r.Result.Policy),
			"Description": component.NewText(r.Result.Description),
			"Result":      component.NewText(string(r.Result.Result)),
			"Severity":    component.NewText(string(r.Result.Severity)),
			"Kind":        component.NewText(r.Subject.Kind),
			"Properties":  GetProperties(r),
		})
	}
	return table
}

func PrintFixes(results *ReportView, s *corev1.ObjectReference, title string) *component.Table {
	cols := fixHighSeverity
	if s.Kind == "Namespace" {
		cols = append(cols, component.TableCol{
			Name: "Resource", Accessor: "Resource",
		})
	}
	table := component.NewTableWithRows(title, title, cols, []component.TableRow{})
	for _, r := range results.HighSeverity {
		if s.Kind == "Namespace" {
			table.Add(component.TableRow{
				"Rule":     component.NewText(r.Result.Rule),
				"Policy":   component.NewText(r.Result.Policy),
				"Result":   component.NewText(string(r.Result.Result)),
				"Severity": component.NewText(string(r.Result.Severity)),
				"Resource": component.NewText(string(r.Subject.Name)),
			})
		}
		table.Add(component.TableRow{
			"Rule":     component.NewText(r.Result.Rule),
			"Policy":   component.NewText(r.Result.Policy),
			"Result":   component.NewText(string(r.Result.Result)),
			"Severity": component.NewText(string(r.Result.Severity)),
		})
	}
	return table
}

func CreateQuadrant(title string, results *ReportView) *component.Quadrant {
	quadrant := component.NewQuadrant(title)
	quadrant.Set(component.QuadNW, "Pass", fmt.Sprintf("%v", results.Analytics.Pass))
	quadrant.Set(component.QuadNE, "Error", fmt.Sprintf("%v", results.Analytics.Error))
	quadrant.Set(component.QuadSE, "Warn", fmt.Sprintf("%v", results.Analytics.Warn))
	quadrant.Set(component.QuadSW, "Fail", fmt.Sprintf("%v", results.Analytics.Fail))
	return quadrant
}
