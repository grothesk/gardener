/*
SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
SPDX-License-Identifier: Apache-2.0
*/

// Code generated by lister-gen. DO NOT EDIT.

package internalversion

import (
	core "github.com/gardener/gardener/pkg/apis/core"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ShootStateLister helps list ShootStates.
type ShootStateLister interface {
	// List lists all ShootStates in the indexer.
	List(selector labels.Selector) (ret []*core.ShootState, err error)
	// ShootStates returns an object that can list and get ShootStates.
	ShootStates(namespace string) ShootStateNamespaceLister
	ShootStateListerExpansion
}

// shootStateLister implements the ShootStateLister interface.
type shootStateLister struct {
	indexer cache.Indexer
}

// NewShootStateLister returns a new ShootStateLister.
func NewShootStateLister(indexer cache.Indexer) ShootStateLister {
	return &shootStateLister{indexer: indexer}
}

// List lists all ShootStates in the indexer.
func (s *shootStateLister) List(selector labels.Selector) (ret []*core.ShootState, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*core.ShootState))
	})
	return ret, err
}

// ShootStates returns an object that can list and get ShootStates.
func (s *shootStateLister) ShootStates(namespace string) ShootStateNamespaceLister {
	return shootStateNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ShootStateNamespaceLister helps list and get ShootStates.
type ShootStateNamespaceLister interface {
	// List lists all ShootStates in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*core.ShootState, err error)
	// Get retrieves the ShootState from the indexer for a given namespace and name.
	Get(name string) (*core.ShootState, error)
	ShootStateNamespaceListerExpansion
}

// shootStateNamespaceLister implements the ShootStateNamespaceLister
// interface.
type shootStateNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ShootStates in the indexer for a given namespace.
func (s shootStateNamespaceLister) List(selector labels.Selector) (ret []*core.ShootState, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*core.ShootState))
	})
	return ret, err
}

// Get retrieves the ShootState from the indexer for a given namespace and name.
func (s shootStateNamespaceLister) Get(name string) (*core.ShootState, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(core.Resource("shootstate"), name)
	}
	return obj.(*core.ShootState), nil
}
