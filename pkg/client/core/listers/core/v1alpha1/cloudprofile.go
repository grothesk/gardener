/*
SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
SPDX-License-Identifier: Apache-2.0
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// CloudProfileLister helps list CloudProfiles.
type CloudProfileLister interface {
	// List lists all CloudProfiles in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.CloudProfile, err error)
	// Get retrieves the CloudProfile from the index for a given name.
	Get(name string) (*v1alpha1.CloudProfile, error)
	CloudProfileListerExpansion
}

// cloudProfileLister implements the CloudProfileLister interface.
type cloudProfileLister struct {
	indexer cache.Indexer
}

// NewCloudProfileLister returns a new CloudProfileLister.
func NewCloudProfileLister(indexer cache.Indexer) CloudProfileLister {
	return &cloudProfileLister{indexer: indexer}
}

// List lists all CloudProfiles in the indexer.
func (s *cloudProfileLister) List(selector labels.Selector) (ret []*v1alpha1.CloudProfile, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.CloudProfile))
	})
	return ret, err
}

// Get retrieves the CloudProfile from the index for a given name.
func (s *cloudProfileLister) Get(name string) (*v1alpha1.CloudProfile, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("cloudprofile"), name)
	}
	return obj.(*v1alpha1.CloudProfile), nil
}
