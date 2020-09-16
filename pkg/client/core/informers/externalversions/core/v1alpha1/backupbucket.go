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

// BackupBucketInformer provides access to a shared informer and lister for
// BackupBuckets.
type BackupBucketInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.BackupBucketLister
}

type backupBucketInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewBackupBucketInformer constructs a new informer for BackupBucket type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewBackupBucketInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredBackupBucketInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredBackupBucketInformer constructs a new informer for BackupBucket type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredBackupBucketInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CoreV1alpha1().BackupBuckets().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CoreV1alpha1().BackupBuckets().Watch(context.TODO(), options)
			},
		},
		&corev1alpha1.BackupBucket{},
		resyncPeriod,
		indexers,
	)
}

func (f *backupBucketInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredBackupBucketInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *backupBucketInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&corev1alpha1.BackupBucket{}, f.defaultInformer)
}

func (f *backupBucketInformer) Lister() v1alpha1.BackupBucketLister {
	return v1alpha1.NewBackupBucketLister(f.Informer().GetIndexer())
}
