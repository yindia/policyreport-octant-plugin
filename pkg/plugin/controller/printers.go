package controller

import (
	"errors"

	"github.com/evalsocket/policyreport-octant-plugin/pkg/plugin/model"
	"github.com/evalsocket/policyreport-octant-plugin/pkg/plugin/view"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/wg-policy-prototypes/policy-report/api/v1alpha1"
	"sigs.k8s.io/wg-policy-prototypes/policy-report/api/v1alpha2"
)

// ResourceTabPrinter is called when Octant wants to add new tab for the underlying resource.
func ResourceTabPrinter(request *service.PrintRequest) (tab plugin.TabResponse, err error) {
	if request.Object == nil {
		err = errors.New("request object is nil")
		return
	}

	return ReportTabPrinter(request)

}

func ReportTabPrinter(request *service.PrintRequest) (plugin.TabResponse, error) {
	flexLayout := component.NewFlexLayout("Policy Report")
	repository := model.NewRepository(request.DashboardClient)
	obj := request.Object.DeepCopyObject()
	subject, err := getSubject(obj)
	if err != nil {
		return plugin.TabResponse{}, nil
	}

	reports := &view.PolicyReports{
		Subject: subject,
		Request: request,
		Ctx: request.Context(),
		Source: view.PolicyReportSource{
			PolicyReports:          v1alpha2.PolicyReportList{},
			PolicyReportsv1:        v1alpha1.PolicyReportList{},
			ClusterPolicyReports:   v1alpha2.ClusterPolicyReportList{},
			ClusterPolicyReportsv1: v1alpha1.ClusterPolicyReportList{},
		},
	}
	if err := repository.GetPolicyReports(reports); err != nil {
		return plugin.TabResponse{}, nil
	}

	filter := filterReports(reports)

	flexLayout.AddSections([]component.FlexLayoutSection{
		{
			{
				Width: component.WidthHalf,
				View:  view.CreateQuadrant("Policy Report", filter),
			},
			{
				Width: component.WidthHalf,
				View:  view.PrintFixes(filter, reports.Subject, "High Severity"),
			},
			{
				Width: component.WidthFull,
				View:  view.PrintPolicyReportTable(filter, reports.Subject, "Policy Report"),
			},
		},
	}...)

	tab := component.NewTabWithContents(*flexLayout)
	return plugin.TabResponse{Tab: tab}, nil
}

func getSubject(obj runtime.Object) (*corev1.ObjectReference, error) {
	accessor := meta.NewAccessor()
	kind, err := accessor.Kind(obj)
	if err != nil {
		return &corev1.ObjectReference{}, err
	}

	name, err := accessor.Name(obj)
	if err != nil {
		return &corev1.ObjectReference{}, err
	}
	namespace, _ := accessor.Namespace(obj)
	if obj.GetObjectKind().GroupVersionKind().Kind == "Namespace" {
		return &corev1.ObjectReference{
			Kind:      kind,
			Name:      name,
			Namespace: "",
		}, err
	}
	return &corev1.ObjectReference{
		Kind:      kind,
		Name:      name,
		Namespace: namespace,
	}, err
}

func getAnalytics(status v1alpha2.PolicyResult, report *view.ReportView) *view.ReportView {
	switch status {
	case v1alpha2.PolicyResult("pass"):
		report.Analytics.Pass++
		return report
	case v1alpha2.PolicyResult("warn"):
		report.Analytics.Warn++
		return report
	case v1alpha2.PolicyResult("error"):
		report.Analytics.Error++
		return report
	case v1alpha2.PolicyResult("fail"):
		report.Analytics.Fail++
		return report
	case v1alpha2.PolicyResult("skip"):
		report.Analytics.Skip++
		return report
	}
	return report
}

func filterReports(reports *view.PolicyReports) *view.ReportView {
	filterResult := &view.ReportView{
		Analytics: &v1alpha2.PolicyReportSummary{
			Pass:  0,
			Fail:  0,
			Error: 0,
			Warn:  0,
			Skip:  0,
		},
	}
	if reports.Subject.Kind == "Namespace" {
		for _, r := range reports.Results {
			if r.Subject.Namespace == reports.Subject.Name {
				filterResult.Reports = append(filterResult.Reports, r)
				filterResult = getAnalytics(r.Result.Result, filterResult)
				if r.Result.Severity == v1alpha2.PolicyResultSeverity("high") {
					filterResult.HighSeverity = append(filterResult.HighSeverity, r)
				}
			}
		}
		return filterResult
	}
	for _, r := range reports.Results {
		if r.Subject.Kind == reports.Subject.Kind && r.Subject.Name == reports.Subject.Name && r.Subject.Namespace == reports.Subject.Namespace {
			filterResult.Reports = append(filterResult.Reports, r)
			filterResult = getAnalytics(r.Result.Result, filterResult)
			if r.Result.Severity == v1alpha2.PolicyResultSeverity("high") {
				filterResult.HighSeverity = append(filterResult.HighSeverity, r)
			}
		}
	}
	return filterResult
}
