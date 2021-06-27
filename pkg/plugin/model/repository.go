package model

import (
	"encoding/json"

	"github.com/evalsocket/policyreport-octant-plugin/pkg/plugin/view"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (r *Repository) GetPolicyReportV1(reports *view.PolicyReports) error {
	unstructuredList, err := r.client.List(reports.Ctx, store.Key{
		APIVersion: "wgpolicyk8s.io/v1alpha1",
		Kind:       "PolicyReport",
	})
	if err != nil {
		return err
	}
	if len(unstructuredList.Items) == 0 {
		return err
	}
	err = r.structure(unstructuredList, &reports.Source.PolicyReportsv1)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetPolicyReport(reports *view.PolicyReports) error {
	unstructuredList, err := r.client.List(reports.Ctx, store.Key{
		APIVersion: "wgpolicyk8s.io/v1alpha2",
		Kind:       "PolicyReport",
	})
	if err != nil {
		return err
	}
	if len(unstructuredList.Items) == 0 {
		return err
	}
	err = r.structure(unstructuredList, &reports.Source.ClusterPolicyReports)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetPolicyReports(reports *view.PolicyReports) error {
	err := r.GetPolicyReportV1(reports)
	if err != nil {
		return err
	}

	err = r.GetPolicyReport(reports)
	if err != nil {
		return err
	}
	MergePolicyReport(reports)
	return nil
}

func (r *Repository) GetClusterPolicyReports(reports *view.PolicyReports) error {

	err := r.GetClusterPolicyReportV1(reports)
	if err != nil {
		return err
	}

	err = r.GetClusterPolicyReport(reports)
	if err != nil {
		return err
	}
	MergeClusterPolicyReport(reports)
	return nil
}

func (r *Repository) GetClusterPolicyReportV1(reports *view.PolicyReports) error {
	unstructuredList, err := r.client.List(reports.Ctx, store.Key{
		APIVersion: "wgpolicyk8s.io/v1alpha1",
		Kind:       "ClusterPolicyReport",
	})
	if err != nil {
		return err
	}
	if len(unstructuredList.Items) == 0 {
		return err
	}
	err = r.structure(unstructuredList, &reports.Source.ClusterPolicyReportsv1)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetClusterPolicyReport(reports *view.PolicyReports) error {
	unstructuredList, err := r.client.List(reports.Ctx, store.Key{
		APIVersion: "wgpolicyk8s.io/v1alpha2",
		Kind:       "ClusterPolicyReport",
	})
	if err != nil {
		return err
	}
	if len(unstructuredList.Items) == 0 {
		return err
	}
	err = r.structure(unstructuredList, &reports.Source.ClusterPolicyReports)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) structure(m json.Marshaler, v interface{}) (err error) {
	b, err := m.MarshalJSON()
	if err != nil {
		return
	}
	err = json.Unmarshal(b, v)
	return
}

func MergePolicyReport(reports *view.PolicyReports) {
	for _, v := range reports.Source.PolicyReportsv1.Items {
		for _, r := range v.Results {
			for _, s := range r.Resources {
				reports.Results = append(reports.Results, view.SingleReport{
					Version: "v1alpha1",
					Scope:   "PolicyReport",
					Subject: s,
					Engine:  v.Labels["engine"],
					Result: v1alpha2.PolicyReportResult{
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
					}})
			}
		}
	}

	for _, v := range reports.Source.PolicyReports.Items {
		for _, r := range v.Results {
			for _, s := range r.Subjects {
				reports.Results = append(reports.Results, view.SingleReport{
					Version: "v1alpha2",
					Scope:   "PolicyReport",
					Subject: s,
					Engine:  v.Labels["engine"],
					Result:  *r})
			}
		}
	}
}

func MergeClusterPolicyReport(reports *view.PolicyReports) {
	for _, v := range reports.Source.ClusterPolicyReportsv1.Items {
		for _, r := range v.Results {
			for _, s := range r.Resources {
				reports.Results = append(reports.Results, view.SingleReport{
					Version: "v1alpha1",
					Scope:   "ClusterPolicyReport",
					Subject: s,
					Engine:  v.Labels["engine"],
					Result: v1alpha2.PolicyReportResult{
						Source:          "",
						Policy:          r.Policy,
						Rule:            r.Rule,
						Category:        r.Category,
						Severity:        v1alpha2.PolicyResultSeverity(r.Severity),
						Timestamp:       v1.Timestamp{},
						Result:          v1alpha2.PolicyResult(r.Status),
						Scored:          r.Scored,
						SubjectSelector: r.ResourceSelector,
						Description:     r.Message,
						Properties:      r.Data,
					}})
			}

		}
	}

	for _, v := range reports.Source.PolicyReports.Items {
		for _, r := range v.Results {
			for _, s := range r.Subjects {
				reports.Results = append(reports.Results, view.SingleReport{
					Version: "v1alpha2",
					Scope:   "ClusterPolicyReport",
					Subject: s,
					Engine:  v.Labels["engine"],
					Result:  *r})
			}
		}
	}
}
