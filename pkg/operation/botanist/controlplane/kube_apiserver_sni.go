// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package controlplane

import (
	"context"
	"path/filepath"

	"github.com/gardener/gardener/pkg/client/kubernetes"
	"github.com/gardener/gardener/pkg/operation/botanist/component"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
)

// kubeAPIServiceValues configure the kube-apiserver service SNI.
type KubeAPIServerSNIValues struct {
	Hosts                    []string  `json:"hosts,omitempty"`
	Name                     string    `json:"name,omitempty"`
	NamespaceUID             types.UID `json:"namespaceUID,omitempty"`
	ApiserverClusterIP       string    `json:"apiserverClusterIP,omitempty"`
	IstioIngressNamespace    string    `json:"istioIngressNamespace,omitempty"`
	EnableKonnectivityTunnel bool      `json:"enableKonnectivityTunnel,omitempty"`
}

// NewKubeAPIServerSNI creates a new instance of DeployWaiter which deploys Istio resources for
// kube-apiserver SNI access.
func NewKubeAPIServerSNI(
	values *KubeAPIServerSNIValues,
	namespace string,
	applier kubernetes.ChartApplier,
	chartsRootPath string,

) component.DeployWaiter {
	if values == nil {
		values = &KubeAPIServerSNIValues{}
	}

	return &kubeAPIServerSNI{
		ChartApplier: applier,
		chartPath:    filepath.Join(chartsRootPath, "seed-controlplane", "charts", "kube-apiserver-sni"),
		values:       values,
		namespace:    namespace,
	}
}

type kubeAPIServerSNI struct {
	values    *KubeAPIServerSNIValues
	namespace string
	kubernetes.ChartApplier
	chartPath string
}

func (k *kubeAPIServerSNI) Deploy(ctx context.Context) error {
	return k.Apply(
		ctx,
		k.chartPath,
		k.namespace,
		k.values.Name,
		kubernetes.Values(k.values),
	)
}

func (k *kubeAPIServerSNI) Destroy(ctx context.Context) error {
	return k.Delete(
		ctx,
		k.chartPath,
		k.namespace,
		k.values.Name,
		kubernetes.Values(k.values),
		kubernetes.TolerateErrorFunc(meta.IsNoMatchError),
	)
}

func (k *kubeAPIServerSNI) Wait(ctx context.Context) error        { return nil }
func (k *kubeAPIServerSNI) WaitCleanup(ctx context.Context) error { return nil }
