// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package genericactuator

import (
	"context"
	"fmt"

	"github.com/gardener/gardener/extensions/pkg/controller"
	gardencorev1beta1helper "github.com/gardener/gardener/pkg/apis/core/v1beta1/helper"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"

	machinev1alpha1 "github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"
	"github.com/pkg/errors"
)

// Migrate removes all machine related resources (e.g. MachineDeployments, MachineClasses, MachineClassSecrets, MachineSets and Machines)
// without waiting for machine-controller-manager to do that. Before removal it ensures that the MCM is deleted.
func (a *genericActuator) Migrate(ctx context.Context, worker *extensionsv1alpha1.Worker, cluster *controller.Cluster) error {
	logger := a.logger.WithValues("worker", kutil.KeyFromObject(worker), "operation", "migrate")

	workerDelegate, err := a.delegateFactory.WorkerDelegate(ctx, worker, cluster)
	if err != nil {
		return errors.Wrap(err, "could not instantiate actuator context")
	}

	// Make sure machine-controller-manager is deleted before deleting the machines.
	if err := a.deleteMachineControllerManager(ctx, logger, worker); err != nil {
		return errors.Wrap(err, "failed deleting machine-controller-manager")
	}

	if err := a.waitUntilMachineControllerManagerIsDeleted(ctx, logger, worker.Namespace); err != nil {
		return errors.Wrap(err, "failed deleting machine-controller-manager")
	}

	if err := a.shallowDeleteAllObjects(ctx, logger, worker.Namespace, &machinev1alpha1.MachineList{}); err != nil {
		return errors.Wrap(err, "shallow deletion of all machine failed")
	}

	if err := a.shallowDeleteAllObjects(ctx, logger, worker.Namespace, &machinev1alpha1.MachineSetList{}); err != nil {
		return errors.Wrap(err, "shallow deletion of all machineSets failed")
	}

	if err := a.shallowDeleteAllObjects(ctx, logger, worker.Namespace, &machinev1alpha1.MachineDeploymentList{}); err != nil {
		return errors.Wrap(err, "shallow deletion of all machineDeployments failed")
	}

	if err := a.shallowDeleteAllObjects(ctx, logger, worker.Namespace, workerDelegate.MachineClassList()); err != nil {
		return errors.Wrap(err, "cleaning up machine classes failed")
	}

	if err := a.shallowDeleteMachineClassSecrets(ctx, logger, worker.Namespace, nil); err != nil {
		return errors.Wrap(err, "cleaning up machine class secrets failed")
	}

	// Wait until all machine resources have been properly deleted.
	if err := a.waitUntilMachineResourcesDeleted(ctx, logger, worker, workerDelegate); err != nil {
		return gardencorev1beta1helper.DetermineError(err, fmt.Sprintf("Failed while waiting for all machine resources to be deleted: '%s'", err.Error()))
	}

	return nil
}
