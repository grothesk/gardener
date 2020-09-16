// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package keys_test

import (
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/gardener/gardener/pkg/client/kubernetes/clientmap/internal"
	"github.com/gardener/gardener/pkg/client/kubernetes/clientmap/keys"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Keys", func() {

	It("#ForGarden", func() {
		key := keys.ForGarden().(internal.GardenClientSetKey)
		Expect(key.Key()).To(Equal("garden"))
	})

	It("#ForSeed", func() {
		seed := &gardencorev1beta1.Seed{
			ObjectMeta: metav1.ObjectMeta{
				Name: "seed-eu1",
			},
		}
		key := keys.ForSeed(seed).(internal.SeedClientSetKey)
		Expect(key.Key()).To(Equal(seed.Name))
		Expect(key).To(BeEquivalentTo(seed.Name))
	})

	It("#ForSeedWithName", func() {
		name := "seed-eu1"
		key := keys.ForSeedWithName(name).(internal.SeedClientSetKey)
		Expect(key.Key()).To(Equal(name))
		Expect(key).To(BeEquivalentTo(name))
	})

	It("#ForShoot", func() {
		shoot := &gardencorev1beta1.Shoot{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "inception",
				Namespace: "core",
			},
		}
		key := keys.ForShoot(shoot).(internal.ShootClientSetKey)
		Expect(key.Key()).To(Equal(shoot.Namespace + "/" + shoot.Name))
		Expect(key.Namespace).To(Equal(shoot.Namespace))
		Expect(key.Name).To(Equal(shoot.Name))
	})

	It("#ForShootWithNamespacedName", func() {
		name := "inception"
		namespace := "core"
		key := keys.ForShootWithNamespacedName(namespace, name).(internal.ShootClientSetKey)
		Expect(key.Key()).To(Equal(namespace + "/" + name))
		Expect(key.Namespace).To(Equal(namespace))
		Expect(key.Name).To(Equal(name))
	})

	It("#ForPlant", func() {
		plant := &gardencorev1beta1.Plant{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "water-me",
				Namespace: "core",
			},
		}
		key := keys.ForPlant(plant).(internal.PlantClientSetKey)
		Expect(key.Key()).To(Equal(plant.Namespace + "/" + plant.Name))
		Expect(key.Namespace).To(Equal(plant.Namespace))
		Expect(key.Name).To(Equal(plant.Name))
	})

	It("#ForPlantWithNamespacedName", func() {
		name := "water-me"
		namespace := "core"
		key := keys.ForPlantWithNamespacedName(namespace, name).(internal.PlantClientSetKey)
		Expect(key.Key()).To(Equal(namespace + "/" + name))
		Expect(key.Namespace).To(Equal(namespace))
		Expect(key.Name).To(Equal(name))
	})

})
