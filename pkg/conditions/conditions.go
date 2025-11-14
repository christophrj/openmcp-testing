package conditions

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2"
	"sigs.k8s.io/e2e-framework/klient/k8s/resources"
	apimachineryconditions "sigs.k8s.io/e2e-framework/klient/wait/conditions"
)

type Conditions struct {
	apimachineryconditions.Condition
	resources *resources.Resources
}

func New(r *resources.Resources) *Conditions {
	return &Conditions{Condition: *apimachineryconditions.New(r), resources: r}
}

func (c *Conditions) ClusterProviderConditionMatch(name string, conditionType string, conditionStatus v1.ConditionStatus) wait.ConditionWithContextFunc {
	return func(ctx context.Context) (done bool, err error) {
		klog.Infof("Waiting for cluster provider %s", name)
		cl, err := dynamic.NewForConfig(c.resources.GetConfig())
		if err != nil {
			return false, err
		}
		res := cl.Resource(schema.GroupVersionResource{
			Group:    "openmcp.cloud",
			Version:  "v1alpha1",
			Resource: "clusterproviders",
		})
		providerObject, err := res.Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		result := checkCondition(providerObject, conditionType, conditionStatus)
		return result, nil
	}
}

func checkCondition(unstruc *unstructured.Unstructured, desiredType string, desiredStatus v1.ConditionStatus) bool {
	conditions, ok, err := unstructured.NestedSlice(unstruc.UnstructuredContent(), "status", "conditions")
	if err != nil {
		klog.Infof("Could not extract conditions of (%s) %s, %s", unstruc.GroupVersionKind().String(), unstruc.GetName(), err.Error())
		return false
	} else if !ok {
		klog.Infof("Object (%s) %s doesnt have conditions", unstruc.GroupVersionKind().String(), unstruc.GetName())
		return false
	}
	status := ""
	message := ""
	for _, condition := range conditions {
		c := condition.(map[string]interface{})
		curType := c["type"]
		if curType == desiredType {
			status = c["status"].(string)
			msg, convertible := c["message"].(string)
			if convertible {
				message = msg
			}
		}
	}
	matchedConditionStatus := status == string(desiredStatus)
	matchedConditionReason := true
	klog.Infof("Object (%s) %s, condition: %s: %s, matched: %t, message: %s",
		unstruc.GroupVersionKind().String(),
		unstruc.GetName(),
		desiredType,
		status,
		matchedConditionStatus,
		message,
	)
	return matchedConditionStatus && matchedConditionReason
}
