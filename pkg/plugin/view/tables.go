package view

import (
	"fmt"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/wg-policy-prototypes/policy-report/api/v1alpha2"
)

func GetProperties(r v1alpha2.PolicyReportResult) *component.Selectors {
	var data []component.Selector
	for k, v := range r.Properties {
		data = append(data, component.NewLabelSelector(k, v))
	}
	s := component.NewSelectors(data)
	return s
}

func PrintPolicyReportTable(results []v1alpha2.PolicyReportResult, s *corev1.ObjectReference, title string) *component.Table {
	table := component.NewTableWithRows(title, title, policyReportCol, []component.TableRow{})
	for _, r := range results {
		if s == nil && len(r.Subjects) > 0 {
			s.Name = r.Subjects[0].Name
			s.Kind = r.Subjects[0].Kind
			s.Namespace = r.Subjects[0].Namespace
		}
		table.Add(component.TableRow{
			"Rule":        component.NewText(r.Rule),
			"Policy":      component.NewText(r.Policy),
			"Description": component.NewText(r.Description),
			"Result":      component.NewText(string(r.Result)),
			"Severity":    component.NewText(string(r.Severity)),
			"Properties":  GetProperties(r),
			"Kind":        component.NewText(s.Kind),
			"APIVersion":  component.NewText(s.APIVersion),
		})
	}
	return table
}

func PrintFixes(results []v1alpha2.PolicyReportResult, title string) *component.Table {
	table := component.NewTableWithRows(title, title, fixHighSeverity, []component.TableRow{})
	for _, r := range results {
		table.Add(component.TableRow{
			"Rule":        component.NewText(r.Rule),
			"Policy":      component.NewText(r.Policy),
			"Result":      component.NewText(string(r.Result)),
			"Severity":    component.NewText(string(r.Severity)),
		})
	}
	return table
}

func CreateQuadrant(title string, results *v1alpha2.PolicyReportSummary) *component.Quadrant {
	quadrant := component.NewQuadrant(title)
	quadrant.Set(component.QuadNW, "Pass", fmt.Sprintf("%v", results.Pass))
	quadrant.Set(component.QuadNE, "Error", fmt.Sprintf("%v", results.Error))
	quadrant.Set(component.QuadSE, "Warn", fmt.Sprintf("%v", results.Warn))
	quadrant.Set(component.QuadSW, "Fail", fmt.Sprintf("%v", results.Fail))
	return quadrant
}
