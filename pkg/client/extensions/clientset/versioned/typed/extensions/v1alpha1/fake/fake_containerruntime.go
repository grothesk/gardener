/*
SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
SPDX-License-Identifier: Apache-2.0
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeContainerRuntimes implements ContainerRuntimeInterface
type FakeContainerRuntimes struct {
	Fake *FakeExtensionsV1alpha1
	ns   string
}

var containerruntimesResource = schema.GroupVersionResource{Group: "extensions.gardener.cloud", Version: "v1alpha1", Resource: "containerruntimes"}

var containerruntimesKind = schema.GroupVersionKind{Group: "extensions.gardener.cloud", Version: "v1alpha1", Kind: "ContainerRuntime"}

// Get takes name of the containerRuntime, and returns the corresponding containerRuntime object, and an error if there is any.
func (c *FakeContainerRuntimes) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ContainerRuntime, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(containerruntimesResource, c.ns, name), &v1alpha1.ContainerRuntime{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ContainerRuntime), err
}

// List takes label and field selectors, and returns the list of ContainerRuntimes that match those selectors.
func (c *FakeContainerRuntimes) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ContainerRuntimeList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(containerruntimesResource, containerruntimesKind, c.ns, opts), &v1alpha1.ContainerRuntimeList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ContainerRuntimeList{ListMeta: obj.(*v1alpha1.ContainerRuntimeList).ListMeta}
	for _, item := range obj.(*v1alpha1.ContainerRuntimeList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested containerRuntimes.
func (c *FakeContainerRuntimes) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(containerruntimesResource, c.ns, opts))

}

// Create takes the representation of a containerRuntime and creates it.  Returns the server's representation of the containerRuntime, and an error, if there is any.
func (c *FakeContainerRuntimes) Create(ctx context.Context, containerRuntime *v1alpha1.ContainerRuntime, opts v1.CreateOptions) (result *v1alpha1.ContainerRuntime, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(containerruntimesResource, c.ns, containerRuntime), &v1alpha1.ContainerRuntime{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ContainerRuntime), err
}

// Update takes the representation of a containerRuntime and updates it. Returns the server's representation of the containerRuntime, and an error, if there is any.
func (c *FakeContainerRuntimes) Update(ctx context.Context, containerRuntime *v1alpha1.ContainerRuntime, opts v1.UpdateOptions) (result *v1alpha1.ContainerRuntime, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(containerruntimesResource, c.ns, containerRuntime), &v1alpha1.ContainerRuntime{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ContainerRuntime), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeContainerRuntimes) UpdateStatus(ctx context.Context, containerRuntime *v1alpha1.ContainerRuntime, opts v1.UpdateOptions) (*v1alpha1.ContainerRuntime, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(containerruntimesResource, "status", c.ns, containerRuntime), &v1alpha1.ContainerRuntime{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ContainerRuntime), err
}

// Delete takes name of the containerRuntime and deletes it. Returns an error if one occurs.
func (c *FakeContainerRuntimes) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(containerruntimesResource, c.ns, name), &v1alpha1.ContainerRuntime{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeContainerRuntimes) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(containerruntimesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ContainerRuntimeList{})
	return err
}

// Patch applies the patch and returns the patched containerRuntime.
func (c *FakeContainerRuntimes) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ContainerRuntime, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(containerruntimesResource, c.ns, name, pt, data, subresources...), &v1alpha1.ContainerRuntime{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ContainerRuntime), err
}
