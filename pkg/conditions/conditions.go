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
	"k8s.io/klog/v2"
	fwconditions "sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

type Conditions struct {
	fwconditions.Condition
	cfg *envconf.Config
	gvr schema.GroupVersionResource
}

func New(cfg *envconf.Config, gvr schema.GroupVersionResource) *Conditions {
	return &Conditions{
		Condition: *fwconditions.New(cfg.Client().Resources()),
		cfg:       cfg,
		gvr:       gvr,
	}
}

// Match returns true if the conditionType of an object matches the conditionStatus.
// If an object is not found, the condition is not satisfied and no error is returned.
func (c *Conditions) Match(ref types.NamespacedName, conditionType string, conditionStatus v1.ConditionStatus) wait.ConditionWithContextFunc {
	return func(ctx context.Context) (done bool, err error) {
		klog.Infof("Waiting for condition: %s %s", c.gvr, ref)
		obj, err := resources.GetObject(ctx, c.cfg, ref, c.gvr)
		if err != nil {
			return false, ignoreNotFound(err)
		}
		return checkCondition(obj, conditionType, conditionStatus), nil
	}
}

// Status returns true if the status key of an object matches the status value.
// If an object is not found, the condition is not satisfied and no error is returned.
func (c *Conditions) Status(ref types.NamespacedName, key string, value string) wait.ConditionWithContextFunc {
	return func(ctx context.Context) (done bool, err error) {
		klog.Infof("Waiting for status: %s %s", c.gvr, ref)
		obj, err := resources.GetObject(ctx, c.cfg, ref, c.gvr)
		if err != nil {
			return false, ignoreNotFound(err)
		}
		status, found, err := unstructured.NestedMap(obj.Object, "status")
		if err != nil {
			return false, err
		}
		if !found {
			return false, nil
		}
		return status[key] == value, nil
	}
}

// Deleted returns true if an object is not found.
func (c *Conditions) Deleted(ref types.NamespacedName) wait.ConditionWithContextFunc {
	return func(ctx context.Context) (done bool, err error) {
		klog.Infof("Waiting for deletion: %s %s", c.gvr, ref)
		_, err = resources.GetObject(ctx, c.cfg, ref, c.gvr)
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
