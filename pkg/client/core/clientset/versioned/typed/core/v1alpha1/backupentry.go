/*
SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
SPDX-License-Identifier: Apache-2.0
*/

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"
	scheme "github.com/gardener/gardener/pkg/client/core/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// BackupEntriesGetter has a method to return a BackupEntryInterface.
// A group's client should implement this interface.
type BackupEntriesGetter interface {
	BackupEntries(namespace string) BackupEntryInterface
}

// BackupEntryInterface has methods to work with BackupEntry resources.
type BackupEntryInterface interface {
	Create(ctx context.Context, backupEntry *v1alpha1.BackupEntry, opts v1.CreateOptions) (*v1alpha1.BackupEntry, error)
	Update(ctx context.Context, backupEntry *v1alpha1.BackupEntry, opts v1.UpdateOptions) (*v1alpha1.BackupEntry, error)
	UpdateStatus(ctx context.Context, backupEntry *v1alpha1.BackupEntry, opts v1.UpdateOptions) (*v1alpha1.BackupEntry, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.BackupEntry, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.BackupEntryList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.BackupEntry, err error)
	BackupEntryExpansion
}

// backupEntries implements BackupEntryInterface
type backupEntries struct {
	client rest.Interface
	ns     string
}

// newBackupEntries returns a BackupEntries
func newBackupEntries(c *CoreV1alpha1Client, namespace string) *backupEntries {
	return &backupEntries{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the backupEntry, and returns the corresponding backupEntry object, and an error if there is any.
func (c *backupEntries) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.BackupEntry, err error) {
	result = &v1alpha1.BackupEntry{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("backupentries").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of BackupEntries that match those selectors.
func (c *backupEntries) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.BackupEntryList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.BackupEntryList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("backupentries").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested backupEntries.
func (c *backupEntries) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("backupentries").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a backupEntry and creates it.  Returns the server's representation of the backupEntry, and an error, if there is any.
func (c *backupEntries) Create(ctx context.Context, backupEntry *v1alpha1.BackupEntry, opts v1.CreateOptions) (result *v1alpha1.BackupEntry, err error) {
	result = &v1alpha1.BackupEntry{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("backupentries").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(backupEntry).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a backupEntry and updates it. Returns the server's representation of the backupEntry, and an error, if there is any.
func (c *backupEntries) Update(ctx context.Context, backupEntry *v1alpha1.BackupEntry, opts v1.UpdateOptions) (result *v1alpha1.BackupEntry, err error) {
	result = &v1alpha1.BackupEntry{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("backupentries").
		Name(backupEntry.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(backupEntry).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *backupEntries) UpdateStatus(ctx context.Context, backupEntry *v1alpha1.BackupEntry, opts v1.UpdateOptions) (result *v1alpha1.BackupEntry, err error) {
	result = &v1alpha1.BackupEntry{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("backupentries").
		Name(backupEntry.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(backupEntry).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the backupEntry and deletes it. Returns an error if one occurs.
func (c *backupEntries) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("backupentries").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *backupEntries) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("backupentries").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched backupEntry.
func (c *backupEntries) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.BackupEntry, err error) {
	result = &v1alpha1.BackupEntry{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("backupentries").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
