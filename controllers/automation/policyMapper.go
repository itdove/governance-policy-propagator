// Copyright Contributors to the Open Cluster Management project

package automation

import (
	"context"

	policiesv1 "github.com/open-cluster-management/governance-policy-propagator/api/v1"
	policyv1beta1 "github.com/open-cluster-management/governance-policy-propagator/api/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func policyMapper(c client.Client) handler.MapFunc {
	return func(obj client.Object) []reconcile.Request {
		policy := obj.(*policiesv1.Policy)
		var result []reconcile.Request
		policyAutomationList := &policyv1beta1.PolicyAutomationList{}
		err := c.List(context.TODO(), policyAutomationList, &client.ListOptions{Namespace: policy.GetNamespace()})
		if err != nil {
			return nil
		}
		found := false
		policyAutomation := policyv1beta1.PolicyAutomation{}
		for _, policyAutomationTemp := range policyAutomationList.Items {
			if policyAutomationTemp.Spec.PolicyRef == policy.GetName() {
				found = true
				policyAutomation = policyAutomationTemp
				break
			}
		}
		if found {
			if policyAutomation.Spec.Mode == "scan" {
				// scan mode, do not queue
			} else if policyAutomation.Spec.Mode == "once" {
				request := reconcile.Request{NamespacedName: types.NamespacedName{
					Name:      policyAutomation.GetName(),
					Namespace: policyAutomation.GetNamespace(),
				}}
				result = append(result, request)
			}
		}
		return result
	}
}
