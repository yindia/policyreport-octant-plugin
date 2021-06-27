package controller

import (
	"github.com/evalsocket/policyreport-octant-plugin/pkg/plugin/model"
	"github.com/evalsocket/policyreport-octant-plugin/pkg/plugin/view"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"sigs.k8s.io/wg-policy-prototypes/policy-report/api/v1alpha1"
	"sigs.k8s.io/wg-policy-prototypes/policy-report/api/v1alpha2"
)

func rootHandler(request service.Request) (component.ContentResponse, error) {
	rootView, err := buildRootViewForRequest(request)
	if err != nil {
		return component.EmptyContentResponse, err
	}
	response := component.NewContentResponse(nil)
	response.Add(rootView)
	return *response, nil
}

func buildRootViewForRequest(request service.Request) (*component.FlexLayout, error) {
	flexLayout := component.NewFlexLayout("Report")
	repository := model.NewRepository(request.DashboardClient())

	reports := &view.PolicyReports{
		Subject:       nil,
		Ctx:           request.Context(),
		ModuleRequest: &request,
		Source: view.PolicyReportSource{
			PolicyReports:          v1alpha2.PolicyReportList{},
			PolicyReportsv1:        v1alpha1.PolicyReportList{},
			ClusterPolicyReports:   v1alpha2.ClusterPolicyReportList{},
			ClusterPolicyReportsv1: v1alpha1.ClusterPolicyReportList{},
		},
	}

	if err := repository.GetPolicyReports(reports); err != nil {
		return flexLayout, err
	}

	if err := repository.GetClusterPolicyReports(reports); err != nil {
		return flexLayout, err
	}

	filterPRResult := filterModuleReport(reports,"PolicyReport")
	filterPRResult.Reports = reports.Results
	filterCPRResult := filterModuleReport(reports,"ClusterPolicyReport")
	filterCPRResult.Reports = reports.Results

	flexLayout.AddSections([]component.FlexLayoutSection{
		{
			{
				Width: component.WidthHalf,
				View:  view.CreateQuadrant("Policy Report", filterPRResult),
			},
			{
				Width: component.WidthHalf,
				View:  view.CreateQuadrant("Cluster Policy Report", filterCPRResult),
			},
			{
				Width: component.WidthFull,
				View:  view.PrintPolicyReportTable(filterPRResult, nil, "Policy Report"),
			},
			{
				Width: component.WidthFull,
				View:  view.PrintPolicyReportTable(filterCPRResult, nil, "Cluster Policy Report"),
			},
		},
	}...)

	return flexLayout, nil
}

func filterModuleReport(reports *view.PolicyReports,scope string) *view.ReportView {
	filterResult := &view.ReportView{
		Analytics: &v1alpha2.PolicyReportSummary{
			Pass:  0,
			Fail:  0,
			Error: 0,
			Warn:  0,
			Skip:  0,
		},
	}
	for _, r := range reports.Results {
		if scope == r.Scope {
			filterResult = getAnalytics(r.Result.Result, filterResult)
		}
	}
	return filterResult
}
