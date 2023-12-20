package deploy

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	listersappsv1 "k8s.io/client-go/listers/apps/v1"

	"github.com/openshift/cluster-resource-override-admission-operator/pkg/asset"
	"github.com/openshift/cluster-resource-override-admission-operator/pkg/ensurer"
	operatorruntime "github.com/openshift/cluster-resource-override-admission-operator/pkg/runtime"
)

func NewDaemonSetInstall(lister listersappsv1.DaemonSetLister,
	oc operatorruntime.OperandContext,
	asset *asset.Asset,
	deployment *ensurer.DaemonSetEnsurer,
	kubeclient kubernetes.Interface) Interface {
	return &daemonset{
		lister:     lister,
		context:    oc,
		asset:      asset,
		deployment: deployment,
		kubeclient: kubeclient,
	}
}

type daemonset struct {
	lister     listersappsv1.DaemonSetLister
	context    operatorruntime.OperandContext
	asset      *asset.Asset
	deployment *ensurer.DaemonSetEnsurer
	kubeclient kubernetes.Interface
}

func (d *daemonset) Name() string {
	return d.asset.DaemonSet().Name()
}

func (d *daemonset) IsAvailable() (available bool, err error) {
	name := d.asset.DaemonSet().Name()
	current, err := d.lister.DaemonSets(d.context.WebhookNamespace()).Get(name)
	if err != nil {
		return
	}

	available, err = GetDaemonSetStatus(current)
	return
}

func (d *daemonset) Get() (object runtime.Object, accessor metav1.Object, err error) {
	name := d.asset.DaemonSet().Name()
	object, err = d.lister.DaemonSets(d.context.WebhookNamespace()).Get(name)
	if err != nil {
		return
	}

	accessor, err = meta.Accessor(object)
	return
}

func (d *daemonset) Ensure(parent, child Applier) (current runtime.Object, accessor metav1.Object, err error) {
	desired := d.asset.DaemonSet().New()

	controPlaneNodeLabelKey := "node-role.kubernetes.io/master"
	nodes, err := d.kubeclient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{LabelSelector: "node-role.kubernetes.io/control-plane"})
	if err != nil {
		return nil, nil, err
	}
	if len(nodes.Items) != 0 {
		controPlaneNodeLabelKey = "node-role.kubernetes.io/control-plane"
	}
	desired.Spec.Template.Spec.NodeSelector = map[string]string{controPlaneNodeLabelKey: ""}

	if parent != nil {
		parent.Apply(desired)
	}
	if child != nil {
		child.Apply(&desired.Spec.Template)
	}

	current, err = d.deployment.Ensure(desired)
	if err != nil {
		return
	}

	accessor, err = meta.Accessor(current)
	return
}
