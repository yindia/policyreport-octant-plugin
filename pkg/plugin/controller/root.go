package controller

import (
	"github.com/evalsocket/policyreport-octant-plugin/pkg/plugin/model"
	"github.com/evalsocket/policyreport-octant-plugin/pkg/plugin/view"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
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
	repository := model.NewRepository(request.DashboardClient())
	flexLayout := component.NewFlexLayout("Policy Report")

	policyReports := repository.GetPolicyReport(request.Context())
	clusterPolicyReport := repository.GetClusterPolicyReport(request.Context())

	flexLayout.AddSections([]component.FlexLayoutSection{
		{
			{
				Width: component.WidthFull,
				View:  view.PrintPolicyReportTable(policyReports, nil, "Policy Report"),
			},
			{
				Width: component.WidthFull,
				View:  view.PrintPolicyReportTable(clusterPolicyReport, nil, "Cluster Policy Report"),
			},
		},
	}...)

	return flexLayout, nil
}
