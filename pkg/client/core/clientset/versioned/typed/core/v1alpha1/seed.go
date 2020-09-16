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

// SeedsGetter has a method to return a SeedInterface.
// A group's client should implement this interface.
type SeedsGetter interface {
	Seeds() SeedInterface
}

// SeedInterface has methods to work with Seed resources.
type SeedInterface interface {
	Create(ctx context.Context, seed *v1alpha1.Seed, opts v1.CreateOptions) (*v1alpha1.Seed, error)
	Update(ctx context.Context, seed *v1alpha1.Seed, opts v1.UpdateOptions) (*v1alpha1.Seed, error)
	UpdateStatus(ctx context.Context, seed *v1alpha1.Seed, opts v1.UpdateOptions) (*v1alpha1.Seed, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.Seed, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.SeedList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Seed, err error)
	SeedExpansion
}

// seeds implements SeedInterface
type seeds struct {
	client rest.Interface
}

// newSeeds returns a Seeds
func newSeeds(c *CoreV1alpha1Client) *seeds {
	return &seeds{
		client: c.RESTClient(),
	}
}

// Get takes name of the seed, and returns the corresponding seed object, and an error if there is any.
func (c *seeds) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Seed, err error) {
	result = &v1alpha1.Seed{}
	err = c.client.Get().
		Resource("seeds").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Seeds that match those selectors.
func (c *seeds) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.SeedList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.SeedList{}
	err = c.client.Get().
		Resource("seeds").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested seeds.
func (c *seeds) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("seeds").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a seed and creates it.  Returns the server's representation of the seed, and an error, if there is any.
func (c *seeds) Create(ctx context.Context, seed *v1alpha1.Seed, opts v1.CreateOptions) (result *v1alpha1.Seed, err error) {
	result = &v1alpha1.Seed{}
	err = c.client.Post().
		Resource("seeds").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(seed).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a seed and updates it. Returns the server's representation of the seed, and an error, if there is any.
func (c *seeds) Update(ctx context.Context, seed *v1alpha1.Seed, opts v1.UpdateOptions) (result *v1alpha1.Seed, err error) {
	result = &v1alpha1.Seed{}
	err = c.client.Put().
		Resource("seeds").
		Name(seed.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(seed).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *seeds) UpdateStatus(ctx context.Context, seed *v1alpha1.Seed, opts v1.UpdateOptions) (result *v1alpha1.Seed, err error) {
	result = &v1alpha1.Seed{}
	err = c.client.Put().
		Resource("seeds").
		Name(seed.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(seed).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the seed and deletes it. Returns an error if one occurs.
func (c *seeds) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("seeds").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *seeds) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("seeds").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched seed.
func (c *seeds) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Seed, err error) {
	result = &v1alpha1.Seed{}
	err = c.client.Patch(pt).
		Resource("seeds").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
