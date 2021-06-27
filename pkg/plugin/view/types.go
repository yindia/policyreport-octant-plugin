package view

import (
	"context"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/wg-policy-prototypes/policy-report/api/v1alpha1"
	"sigs.k8s.io/wg-policy-prototypes/policy-report/api/v1alpha2"
)

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
	namespaceReportCol = []component.TableCol{
		{Name: "Rule", Accessor: "Rule"},
		{Name: "Policy", Accessor: "Policy"},
		{Name: "Resource", Accessor: "Resource"},
		{Name: "Result", Accessor: "Result"},
		{Name: "Severity", Accessor: "Severity"},
		{Name: "Kind", Accessor: "Kind"},
	}
)

type ReportView struct {
	TargetSubject *corev1.ObjectReference
	HighSeverity  []SingleReport
	Reports       []SingleReport
	Analytics     *v1alpha2.PolicyReportSummary
}

type PolicyReports struct {
	// Subject points out to the resource where we triggered the event
	Subject       *corev1.ObjectReference
	Request       *service.PrintRequest
	ModuleRequest *service.Request
	Scope         string
	Ctx           context.Context
	Results       []SingleReport
	Source        PolicyReportSource
}

type PolicyReportSource struct {
	PolicyReports          v1alpha2.PolicyReportList
	PolicyReportsv1        v1alpha1.PolicyReportList
	ClusterPolicyReports   v1alpha2.ClusterPolicyReportList
	ClusterPolicyReportsv1 v1alpha1.ClusterPolicyReportList
}

type SingleReport struct {
	Version string
	Scope   string
	Engine  string
	Subject *corev1.ObjectReference
	Result  v1alpha2.PolicyReportResult
}
