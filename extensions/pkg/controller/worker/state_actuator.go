// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"
	"encoding/json"

	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	workerhelper "github.com/gardener/gardener/extensions/pkg/controller/worker/helper"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"

	machinev1alpha1 "github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type genericStateActuator struct {
	logger logr.Logger

	client client.Client
}

// NewStateActuator creates a new Actuator that reconciles Worker's State subresource
// It provides a default implementation that allows easier integration of providers.
func NewStateActuator(logger logr.Logger) StateActuator {
	return &genericStateActuator{logger: logger.WithName("worker-state-actuator")}
}

func (a *genericStateActuator) InjectClient(client client.Client) error {
	a.client = client
	return nil
}

// Reconcile update the Worker state with the latest.
func (a *genericStateActuator) Reconcile(ctx context.Context, worker *extensionsv1alpha1.Worker) error {
	copyOfWorker := worker.DeepCopy()
	if err := a.updateWorkerState(ctx, copyOfWorker); err != nil {
		return errors.Wrapf(err, "failed to update the state in worker status")
	}

	return nil
}

func (a *genericStateActuator) updateWorkerState(ctx context.Context, worker *extensionsv1alpha1.Worker) error {
	state, err := a.getWorkerState(ctx, worker.Namespace)
	if err != nil {
		return err
	}
	rawState, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return extensionscontroller.TryUpdateStatus(ctx, retry.DefaultBackoff, a.client, worker, func() error {
		worker.Status.State = &runtime.RawExtension{Raw: rawState}
		return nil
	})
}

func (a *genericStateActuator) getWorkerState(ctx context.Context, namespace string) (*State, error) {
	existingMachineDeployments := &machinev1alpha1.MachineDeploymentList{}
	if err := a.client.List(ctx, existingMachineDeployments, client.InNamespace(namespace)); err != nil {
		return nil, err
	}

	machineSets, err := a.getExistingMachineSetsMap(ctx, namespace)
	if err != nil {
		return nil, err
	}

	machines, err := a.getExistingMachinesMap(ctx, namespace)
	if err != nil {
		return nil, err
	}

	workerState := &State{
		MachineDeployments: make(map[string]*MachineDeploymentState),
	}
	for _, deployment := range existingMachineDeployments.Items {
		machineDeploymentState := &MachineDeploymentState{}

		machineDeploymentState.Replicas = deployment.Spec.Replicas

		machineDeploymentMachineSets, ok := machineSets[deployment.Name]
		if !ok {
			continue
		}
		addMachineSetToMachineDeploymentState(machineDeploymentMachineSets, machineDeploymentState)

		for _, machineSet := range machineDeploymentMachineSets {
			currentMachines := append(machines[machineSet.Name], machines[deployment.Name]...)
			if len(currentMachines) <= 0 {
				continue
			}

			for index := range currentMachines {
				addMachineToMachineDeploymentState(&currentMachines[index], machineDeploymentState)
			}
		}

		workerState.MachineDeployments[deployment.Name] = machineDeploymentState
	}

	return workerState, nil
}

// getExistingMachineSetsMap returns a map of existing MachineSets as values and their owners as keys
func (a *genericStateActuator) getExistingMachineSetsMap(ctx context.Context, namespace string) (map[string][]machinev1alpha1.MachineSet, error) {
	existingMachineSets := &machinev1alpha1.MachineSetList{}
	if err := a.client.List(ctx, existingMachineSets, client.InNamespace(namespace)); err != nil {
		return nil, err
	}

	return workerhelper.BuildOwnerToMachineSetsMap(existingMachineSets.Items), nil
}

// getExistingMachinesMap returns a map of the existing Machines as values and the name of their owner
// no matter of being machineSet or MachineDeployment. If a Machine has a ownerRefernce the key(owner)
// will be the MachineSet if not the key will be the name of the MachineDeployment which is stored as
// a lable. We assume that there is no MachineDeployment and MachineSet with the same names.
func (a *genericStateActuator) getExistingMachinesMap(ctx context.Context, namespace string) (map[string][]machinev1alpha1.Machine, error) {
	existingMachines := &machinev1alpha1.MachineList{}
	if err := a.client.List(ctx, existingMachines, client.InNamespace(namespace)); err != nil {
		return nil, err
	}

	return workerhelper.BuildOwnerToMachinesMap(existingMachines.Items), nil
}

func addMachineSetToMachineDeploymentState(machineSets []machinev1alpha1.MachineSet, machineDeploymentState *MachineDeploymentState) {
	if len(machineSets) < 1 || machineDeploymentState == nil {
		return
	}

	//remove redundant data from the machine set
	for index := range machineSets {
		machineSet := &machineSets[index]
		machineSet.ObjectMeta = metav1.ObjectMeta{
			Name:        machineSet.Name,
			Namespace:   machineSet.Namespace,
			Annotations: machineSet.Annotations,
			Labels:      machineSet.Labels,
		}
		machineSet.OwnerReferences = nil
		machineSet.Status = machinev1alpha1.MachineSetStatus{}
	}

	machineDeploymentState.MachineSets = machineSets
}

func addMachineToMachineDeploymentState(machine *machinev1alpha1.Machine, machineDeploymentState *MachineDeploymentState) {
	if machine == nil || machineDeploymentState == nil {
		return
	}

	//remove redundant data from the machine
	machine.ObjectMeta = metav1.ObjectMeta{
		Name:        machine.Name,
		Namespace:   machine.Namespace,
		Annotations: machine.Annotations,
		Labels:      machine.Labels,
	}
	machine.OwnerReferences = nil
	machine.Status = machinev1alpha1.MachineStatus{
		Node: machine.Status.Node,
	}

	machineDeploymentState.Machines = append(machineDeploymentState.Machines, *machine)
}
