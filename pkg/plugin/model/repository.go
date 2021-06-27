package model

import (
	"context"
	"encoding/json"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/wg-policy-prototypes/policy-report/api/v1alpha1"
	"sigs.k8s.io/wg-policy-prototypes/policy-report/api/v1alpha2"

	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
)

const (
	ClusterKind = "Cluster"
)

// Repository defines methods for accessing Kubernetes object.
type Repository struct {
	client service.Dashboard
}

// NewRepository constructs new Kubernetes objects repository
// with the specified Kubernetes client provided by Octant extensions API.
func NewRepository(client service.Dashboard) *Repository {
	return &Repository{
		client: client,
	}
}

// PolicyReport
type PolicyReport struct {
	Name   string
	Report v1alpha2.PolicyReport
}

func (r *Repository) GetPolicyReport(ctx context.Context) []v1alpha2.PolicyReportResult {
	v1, err := r.GetPolicyReportV1(ctx)
	if err != nil {
		return []v1alpha2.PolicyReportResult{}
	}
	v2, err := r.GetPolicyReportV2(ctx)
	if err != nil {
		return []v1alpha2.PolicyReportResult{}
	}
	results := MergePolicyReport(v1.Items, v2.Items)
	return results
}

func (r *Repository) GetClusterPolicyReport(ctx context.Context) []v1alpha2.PolicyReportResult {
	v1, err := r.GetClusterPolicyReportV1(ctx)
	if err != nil {
		return []v1alpha2.PolicyReportResult{}
	}
	v2, err := r.GetClusterPolicyReportV2(ctx)
	if err != nil {
		return []v1alpha2.PolicyReportResult{}
	}
	results := MergeClusterPolicyReport(v1.Items, v2.Items)
	return results
}

func (r *Repository) GetClusterPolicyReportV1(ctx context.Context) (reports v1alpha1.ClusterPolicyReportList, err error) {
	unstructuredList, err := r.client.List(ctx, store.Key{
		APIVersion: "wgpolicyk8s.io/v1alpha1",
		Kind:       "ClusterPolicyReport",
	})
	if err != nil {
		return
	}
	if len(unstructuredList.Items) == 0 {
		return
	}
	var reportList v1alpha1.ClusterPolicyReportList
	err = r.structure(unstructuredList, &reportList)
	if err != nil {
		return
	}

	reports = reportList
	return
}

func (r *Repository) GetClusterPolicyReportV2(ctx context.Context) (reports v1alpha2.ClusterPolicyReportList, err error) {
	unstructuredList, err := r.client.List(ctx, store.Key{
		APIVersion: "wgpolicyk8s.io/v1alpha2",
		Kind:       "ClusterPolicyReport",
	})
	if err != nil {
		return
	}
	if len(unstructuredList.Items) == 0 {
		return
	}
	var reportList v1alpha2.ClusterPolicyReportList
	err = r.structure(unstructuredList, &reportList)
	if err != nil {
		return
	}

	reports = reportList
	return
}

func (r *Repository) GetPolicyReportV1(ctx context.Context) (report v1alpha1.PolicyReportList, err error) {
	unstructuredList, err := r.client.List(ctx, store.Key{
		APIVersion: "wgpolicyk8s.io/v1alpha1",
		Kind:       "PolicyReport",
	})
	if err != nil {
		return
	}
	if len(unstructuredList.Items) == 0 {
		return
	}
	var reportList v1alpha1.PolicyReportList
	err = r.structure(unstructuredList, &reportList)
	if err != nil {
		return
	}

	return reportList, nil
}

func (r *Repository) GetPolicyReportV2(ctx context.Context) (report v1alpha2.PolicyReportList, err error) {
	unstructuredList, err := r.client.List(ctx, store.Key{
		APIVersion: "wgpolicyk8s.io/v1alpha2",
		Kind:       "PolicyReport",
	})
	if err != nil {
		return
	}
	if len(unstructuredList.Items) == 0 {
		return
	}
	var reportList v1alpha2.PolicyReportList
	err = r.structure(unstructuredList, &reportList)
	if err != nil {
		return
	}

	return reportList, nil
}

func (r *Repository) structure(m json.Marshaler, v interface{}) (err error) {
	b, err := m.MarshalJSON()
	if err != nil {
		return
	}
	err = json.Unmarshal(b, v)
	return
}

func MergePolicyReport(policyReportv1 []v1alpha1.PolicyReport, policyReportv2 []v1alpha2.PolicyReport) []v1alpha2.PolicyReportResult {
	reports := []v1alpha2.PolicyReportResult{}
	for _, v := range policyReportv1 {
		for _, r := range v.Results {
			reports = append(reports, v1alpha2.PolicyReportResult{
				Source:          "",
				Policy:          r.Policy,
				Rule:            r.Rule,
				Category:        r.Category,
				Severity:        v1alpha2.PolicyResultSeverity(r.Severity),
				Timestamp:       v1.Timestamp{},
				Result:          v1alpha2.PolicyResult(r.Status),
				Scored:          r.Scored,
				Subjects:        r.Resources,
				SubjectSelector: r.ResourceSelector,
				Description:     r.Message,
				Properties:      r.Data,
			})
		}
	}
	//TODO: // Remove common policy if any
	//for _,v := range policyReportv2 {
	//	for _,r := range v.Results {
	//		reports = append(reports,*r)
	//	}
	//}
	return reports
}

func MergeClusterPolicyReport(policyReportv1 []v1alpha1.ClusterPolicyReport, policyReportv2 []v1alpha2.ClusterPolicyReport) []v1alpha2.PolicyReportResult {
	reports := []v1alpha2.PolicyReportResult{}
	for _, v := range policyReportv1 {
		for _, r := range v.Results {
			reports = append(reports, v1alpha2.PolicyReportResult{
				Source:          "",
				Policy:          r.Policy,
				Rule:            r.Rule,
				Category:        r.Category,
				Severity:        v1alpha2.PolicyResultSeverity(r.Severity),
				Timestamp:       v1.Timestamp{},
				Result:          v1alpha2.PolicyResult(r.Status),
				Scored:          r.Scored,
				Subjects:        r.Resources,
				SubjectSelector: r.ResourceSelector,
				Description:     r.Message,
				Properties:      r.Data,
			})
		}

	}
	//for _,v := range policyReportv2 {
	//	for _,r := range v.Results {
	//		reports = append(reports,*r)
	//	}
	//}
	return reports
}
