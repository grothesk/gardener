// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package operations

import (
	"context"
	"fmt"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	gardenerutils "github.com/gardener/gardener/pkg/utils"
	"github.com/gardener/gardener/test/framework"

	"github.com/onsi/ginkgo"
	g "github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Shoot container runtime testing", func() {

	f := framework.NewShootFramework(nil)

	f.Beta().Serial().CIt("should add worker pool with containerd", func(ctx context.Context) {
		shoot := f.Shoot

		if len(shoot.Spec.Provider.Workers) == 0 {
			ginkgo.Skip("at least one worker pool is required in the test shoot")
		}

		worker := shoot.Spec.Provider.Workers[0]
		// containerD is supported only with Ubuntu OS for now.
		// TODO: adapt/remove this when containerD is available on other OS.
		if worker.Machine.Image.Name != "ubuntu" {
			ginkgo.Skip("worker with machine image 'ubuntu' is required")
		}

		containerdWorker := worker.DeepCopy()

		allowedCharacters := "0123456789abcdefghijklmnopqrstuvwxyz"
		id, err := gardenerutils.GenerateRandomStringFromCharset(3, allowedCharacters)
		framework.ExpectNoError(err)

		containerdWorker.Name = fmt.Sprintf("test-%s", id)
		containerdWorker.Maximum = 1
		containerdWorker.Minimum = 1
		containerdWorker.CRI = &gardencorev1beta1.CRI{
			Name:              extensionsv1alpha1.CRINameContainerD,
			ContainerRuntimes: nil,
		}

		shoot.Spec.Provider.Workers = append(shoot.Spec.Provider.Workers, *containerdWorker)

		ginkgo.By("adding containerd worker pool")

		defer func(ctx context.Context, workerPoolName string) {
			ginkgo.By("removing containerd worker pool after test execution")
			err := f.UpdateShoot(ctx, func(s *gardencorev1beta1.Shoot) error {
				var workers []gardencorev1beta1.Worker
				for _, current := range s.Spec.Provider.Workers {
					if current.Name == workerPoolName {
						continue
					}
					workers = append(workers, current)
				}
				s.Spec.Provider.Workers = workers
				return nil
			})
			framework.ExpectNoError(err)
		}(ctx, containerdWorker.Name)

		err = f.UpdateShoot(ctx, func(s *gardencorev1beta1.Shoot) error {
			s.Spec.Provider.Workers = shoot.Spec.Provider.Workers
			return nil
		})
		framework.ExpectNoError(err)

		// check the node labels of the worker pool to contain containerd label
		nodeList, err := framework.GetAllNodesInWorkerPool(ctx, f.ShootClient, &containerdWorker.Name)
		framework.ExpectNoError(err)
		g.Expect(len(nodeList.Items)).To(g.Equal(int(containerdWorker.Minimum)))

		for _, node := range nodeList.Items {
			value, found := node.Labels[extensionsv1alpha1.CRINameWorkerLabel]
			g.Expect(found).To(g.BeTrue())
			g.Expect(value).To(g.Equal(extensionsv1alpha1.CRINameContainerD))
		}

		// deploy root pod
		rootPodExecutor := framework.NewRootPodExecutor(f.Logger, f.ShootClient, &nodeList.Items[0].Name, "kube-system")

		// check the configuration on the host
		initializerServiceCommand := fmt.Sprintf("systemctl is-active %s", "containerd-initializer")
		executeCommand(ctx, rootPodExecutor, initializerServiceCommand, "active")

		containerdServiceCommand := fmt.Sprintf("systemctl is-active %s", "containerd")
		executeCommand(ctx, rootPodExecutor, containerdServiceCommand, "active")

		// check that config.toml is configured
		checkConfigurationCommand := "cat /etc/systemd/system/containerd.service.d/11-exec_config.conf | grep 'usr/bin/containerd --config=/etc/containerd/config.toml' |  echo $?"
		executeCommand(ctx, rootPodExecutor, checkConfigurationCommand, "0")

		// check that config.toml exists
		checkConfigCommand := "[ -f /etc/containerd/config.toml ] && echo 'found' || echo 'Not found'"
		executeCommand(ctx, rootPodExecutor, checkConfigCommand, "found")
	}, scaleWorkerTimeout)
})

// executeCommand executes a command on the host and checks the returned result
func executeCommand(ctx context.Context, rootPodExecutor framework.RootPodExecutor, command, expected string) {
	response, err := rootPodExecutor.Execute(ctx, command)
	framework.ExpectNoError(err)
	g.Expect(response).ToNot(g.BeNil())
	g.Expect(string(response)).To(g.Equal(fmt.Sprintf("%s\n", expected)))
}
