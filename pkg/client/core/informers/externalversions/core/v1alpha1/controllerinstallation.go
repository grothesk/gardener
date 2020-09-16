/*
SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
SPDX-License-Identifier: Apache-2.0
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	corev1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"
	versioned "github.com/gardener/gardener/pkg/client/core/clientset/versioned"
	internalinterfaces "github.com/gardener/gardener/pkg/client/core/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/gardener/gardener/pkg/client/core/listers/core/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ControllerInstallationInformer provides access to a shared informer and lister for
// ControllerInstallations.
type ControllerInstallationInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.ControllerInstallationLister
}

type controllerInstallationInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewControllerInstallationInformer constructs a new informer for ControllerInstallation type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewControllerInstallationInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredControllerInstallationInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredControllerInstallationInformer constructs a new informer for ControllerInstallation type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredControllerInstallationInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CoreV1alpha1().ControllerInstallations().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CoreV1alpha1().ControllerInstallations().Watch(context.TODO(), options)
			},
		},
		&corev1alpha1.ControllerInstallation{},
		resyncPeriod,
		indexers,
	)
}

func (f *controllerInstallationInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredControllerInstallationInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *controllerInstallationInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&corev1alpha1.ControllerInstallation{}, f.defaultInformer)
}

func (f *controllerInstallationInformer) Lister() v1alpha1.ControllerInstallationLister {
	return v1alpha1.NewControllerInstallationLister(f.Informer().GetIndexer())
}
