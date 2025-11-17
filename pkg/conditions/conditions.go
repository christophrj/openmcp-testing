package conditions

import (
	"context"

	"github.com/christophrj/openmcp-testing/pkg/resources"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	fwresources "sigs.k8s.io/e2e-framework/klient/k8s/resources"
	fwconditions "sigs.k8s.io/e2e-framework/klient/wait/conditions"
)

type Conditions struct {
	fwconditions.Condition
	cfg *rest.Config
}

func New(r *fwresources.Resources) *Conditions {
	return &Conditions{Condition: *fwconditions.New(r), cfg: r.GetConfig()}
}

func (c *Conditions) ClusterProviderConditionMatch(name string, conditionType string, conditionStatus v1.ConditionStatus) wait.ConditionWithContextFunc {
	return func(ctx context.Context) (done bool, err error) {
		klog.Infof("Waiting for cluster provider %s", name)
		obj, err := c.retrieveObject(ctx, types.NamespacedName{
			Name: name,
		}, resources.ClusterproviderGVR)
		if err != nil {
			return false, err
		}
		return checkCondition(obj, conditionType, conditionStatus), nil
	}
}

func (c *Conditions) ClusterConditionMatch(ref types.NamespacedName, conditionType string, conditionStatus v1.ConditionStatus) wait.ConditionWithContextFunc {
	return func(ctx context.Context) (done bool, err error) {
		klog.Infof("Waiting for cluster: %s", ref)
		obj, err := c.retrieveObject(ctx, ref, resources.ClusterGVR)
		if err != nil {
			return false, ignoreNotFound(err)
		}
		return checkCondition(obj, conditionType, conditionStatus), nil
	}
}

func (c *Conditions) ClusterDelete(ref types.NamespacedName) wait.ConditionWithContextFunc {
	return func(ctx context.Context) (done bool, err error) {
		klog.Infof("Waiting for cluster deletion: %s", ref)
		_, err = c.retrieveObject(ctx, ref, resources.ClusterGVR)
		if err != nil && !errors.IsNotFound(err) {
			return false, err
		}
		return err != nil && errors.IsNotFound(err), nil
	}
}

func ignoreNotFound(err error) error {
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

func (c *Conditions) retrieveObject(ctx context.Context, ref types.NamespacedName, gvr schema.GroupVersionResource) (*unstructured.Unstructured, error) {
	cl, err := resources.New(c.cfg)
	if err != nil {
		return nil, err
	}
	return cl.GetObject(ctx, ref, gvr)
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
