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
	var err error
	flexLayout := component.NewFlexLayout("Policy Report")
	repository := model.NewRepository(request.DashboardClient)

	results := &v1alpha2.PolicyReportSummary{
		Pass:  0,
		Fail:  0,
		Error: 0,
		Warn:  0,
		Skip:  0,
	}

	obj := request.Object.DeepCopyObject()
	kind, name, namespace, err := getResourceInfo(obj)
	if err != nil {
		return plugin.TabResponse{}, nil
	}

	policyReport := repository.GetPolicyReport(request.Context())
	subject := &corev1.ObjectReference{
		Kind:      kind,
		Name:      name,
		Namespace: namespace,
	}

	filteredReports := []v1alpha2.PolicyReportResult{}
	filteredRedReports := []v1alpha2.PolicyReportResult{}
	for _, r := range policyReport {
		for _, s := range r.Subjects {
			if s.Kind == kind && s.Name == name && s.Namespace == namespace {
				filteredReports = append(filteredReports, r)
				getAnalytics(r.Result, results)
				if r.Severity == v1alpha2.PolicyResultSeverity("high") {
					filteredRedReports = append(filteredRedReports, r)
				}
			}
		}
	}

	d := component.NewDonutChart()
	d.Config = component.DonutChartConfig{
			Size: component.DonutChartSizeSmall,
			Segments: []component.DonutSegment{
				{
					Count: results.Pass,
					Status: component.NodeStatusOK,
				},
			},
			Labels: component.DonutChartLabels{
				Singular: "Policy Report",
			},
	}


	flexLayout.AddSections([]component.FlexLayoutSection{
		{
			{
				Width: component.WidthHalf,
				View: view.PrintFixes(filteredRedReports,"High Severity"),
			},
			{
				Width: component.WidthHalf,
				View:  view.CreateQuadrant("Policy Report", results),
			},
			{
				Width: component.WidthFull,
				View:  view.PrintPolicyReportTable(filteredReports, subject, "Policy Report"),
			},
		},
	}...)

	tab := component.NewTabWithContents(*flexLayout)
	return plugin.TabResponse{Tab: tab}, nil
}

func getResourceInfo(obj runtime.Object) (string, string, string, error) {
	accessor := meta.NewAccessor()
	kind, err := accessor.Kind(obj)
	if err != nil {
		return "", "", "", err
	}

	name, err := accessor.Name(obj)
	if err != nil {
		return "", "", "", err
	}

	namespace, err := accessor.Namespace(obj)
	if err != nil {
		return "", "", "", err
	}
	return kind, name, namespace, nil
}

func getAnalytics(status v1alpha2.PolicyResult, result *v1alpha2.PolicyReportSummary) {
	switch status {
	case v1alpha2.PolicyResult("pass"):
		result.Pass++
		return
	case v1alpha2.PolicyResult("warn"):
		result.Warn++
		return
	case v1alpha2.PolicyResult("error"):
		result.Error++
		return
	case v1alpha2.PolicyResult("fail"):
		result.Fail++
		return
	case v1alpha2.PolicyResult("skip"):
		result.Skip++
		return
	}
}
