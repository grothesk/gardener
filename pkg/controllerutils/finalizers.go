// SPDX-FileCopyrightText: 2018 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package controllerutils

import (
	"context"
	"fmt"
	"time"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EnsureFinalizer ensure the <finalizer> is present for the object.
func EnsureFinalizer(ctx context.Context, c client.Client, obj kutil.Object, finalizer string) error {
	if err := kutil.TryUpdate(ctx, retry.DefaultBackoff, c, obj, func() error {
		finalizers := sets.NewString(obj.GetFinalizers()...)
		finalizers.Insert(finalizer)
		obj.SetFinalizers(finalizers.UnsortedList())
		return nil
	}); err != nil {
		return fmt.Errorf("could not ensure %q finalizer: %+v", finalizer, err)
	}
	return nil
}

// RemoveGardenerFinalizer removes the gardener finalizer from the object.
func RemoveGardenerFinalizer(ctx context.Context, c client.Client, obj kutil.Object) error {
	return RemoveFinalizer(ctx, c, obj, gardencorev1beta1.GardenerName)
}

// RemoveFinalizer removes the <finalizer> from the object.
func RemoveFinalizer(ctx context.Context, c client.Client, obj kutil.Object, finalizer string) error {
	if err := kutil.TryUpdate(ctx, retry.DefaultBackoff, c, obj, func() error {
		finalizers := sets.NewString(obj.GetFinalizers()...)
		finalizers.Delete(finalizer)
		obj.SetFinalizers(finalizers.UnsortedList())
		return nil
	}); client.IgnoreNotFound(err) != nil {
		return fmt.Errorf("could not remove %q finalizer: %+v", finalizer, err)
	}

	// Wait until the above modifications are reflected in the cache to prevent unwanted reconcile
	// operations (sometimes the cache is not synced fast enough).
	pollerCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return wait.PollImmediateUntil(time.Second, func() (bool, error) {
		err := c.Get(ctx, kutil.KeyFromObject(obj), obj)
		if apierrors.IsNotFound(err) {
			return true, nil
		}
		if err != nil {
			return false, err
		}
		if !HasFinalizer(obj, finalizer) {
			return true, nil
		}
		return false, nil
	}, pollerCtx.Done())
}

// HasFinalizer checks whether the given obj has the given finalizer.
func HasFinalizer(obj metav1.Object, finalizer string) bool {
	return sets.NewString(obj.GetFinalizers()...).Has(finalizer)
}
