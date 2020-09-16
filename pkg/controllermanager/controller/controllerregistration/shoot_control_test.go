// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package controllerregistration

import (
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Controller", func() {
	var (
		queue *fakeQueue
		c     *Controller

		seedName = "seed"
	)

	BeforeEach(func() {
		queue = &fakeQueue{}
		c = &Controller{
			controllerRegistrationSeedQueue: queue,
		}
	})

	Describe("#shootAdd", func() {
		It("should do nothing because object is not a Shoot", func() {
			obj := &gardencorev1beta1.CloudProfile{}

			c.shootAdd(obj)

			Expect(queue.Len()).To(BeZero())
		})

		It("should do nothing because the seedName is nil", func() {
			obj := &gardencorev1beta1.Shoot{}

			c.shootAdd(obj)

			Expect(queue.Len()).To(BeZero())
		})

		It("should add the object to the queue", func() {
			obj := &gardencorev1beta1.Shoot{
				Spec: gardencorev1beta1.ShootSpec{
					SeedName: &seedName,
				},
			}

			c.shootAdd(obj)

			Expect(queue.Len()).To(Equal(1))
			Expect(queue.items[0]).To(Equal(seedName))
		})
	})

	Describe("#shootUpdate", func() {
		It("should do nothing because old object is not a Shoot", func() {
			oldObj := &gardencorev1beta1.CloudProfile{}
			newObj := &gardencorev1beta1.Shoot{}

			c.shootUpdate(oldObj, newObj)

			Expect(queue.Len()).To(BeZero())
		})

		It("should do nothing because new object is not a Shoot", func() {
			oldObj := &gardencorev1beta1.Shoot{}
			newObj := &gardencorev1beta1.CloudProfile{}

			c.shootUpdate(oldObj, newObj)

			Expect(queue.Len()).To(BeZero())
		})

		It("should do nothing because nothing changed", func() {
			oldObj := &gardencorev1beta1.Shoot{}
			newObj := &gardencorev1beta1.Shoot{}

			c.shootUpdate(oldObj, newObj)

			Expect(queue.Len()).To(BeZero())
		})

		It("should add the new obj to the queue because seed name changed", func() {
			oldObj := &gardencorev1beta1.Shoot{}
			newObj := &gardencorev1beta1.Shoot{
				Spec: gardencorev1beta1.ShootSpec{
					SeedName: &seedName,
				},
			}

			c.shootUpdate(oldObj, newObj)

			Expect(queue.Len()).To(Equal(1))
			Expect(queue.items[0]).To(Equal(seedName))
		})

		It("should add the new obj to the queue because workers changed", func() {
			oldObj := &gardencorev1beta1.Shoot{
				Spec: gardencorev1beta1.ShootSpec{
					SeedName: &seedName,
				},
			}
			newObj := oldObj.DeepCopy()
			newObj.Spec.Provider.Workers = []gardencorev1beta1.Worker{{}}

			c.shootUpdate(oldObj, newObj)

			Expect(queue.Len()).To(Equal(1))
			Expect(queue.items[0]).To(Equal(seedName))
		})

		It("should add the new obj to the queue because extensions changed", func() {
			oldObj := &gardencorev1beta1.Shoot{
				Spec: gardencorev1beta1.ShootSpec{
					SeedName: &seedName,
				},
			}
			newObj := oldObj.DeepCopy()
			newObj.Spec.Extensions = []gardencorev1beta1.Extension{{}}

			c.shootUpdate(oldObj, newObj)

			Expect(queue.Len()).To(Equal(1))
			Expect(queue.items[0]).To(Equal(seedName))
		})

		It("should add the new obj to the queue because dns changed", func() {
			oldObj := &gardencorev1beta1.Shoot{
				Spec: gardencorev1beta1.ShootSpec{
					SeedName: &seedName,
				},
			}
			newObj := oldObj.DeepCopy()
			newObj.Spec.DNS = &gardencorev1beta1.DNS{}

			c.shootUpdate(oldObj, newObj)

			Expect(queue.Len()).To(Equal(1))
			Expect(queue.items[0]).To(Equal(seedName))
		})

		It("should add the new obj to the queue because networking type changed", func() {
			oldObj := &gardencorev1beta1.Shoot{
				Spec: gardencorev1beta1.ShootSpec{
					SeedName: &seedName,
				},
			}
			newObj := oldObj.DeepCopy()
			newObj.Spec.Networking.Type = "type"

			c.shootUpdate(oldObj, newObj)

			Expect(queue.Len()).To(Equal(1))
			Expect(queue.items[0]).To(Equal(seedName))
		})

		It("should add the new obj to the queue because provider type changed", func() {
			oldObj := &gardencorev1beta1.Shoot{
				Spec: gardencorev1beta1.ShootSpec{
					SeedName: &seedName,
				},
			}
			newObj := oldObj.DeepCopy()
			newObj.Spec.Provider.Type = "type"

			c.shootUpdate(oldObj, newObj)

			Expect(queue.Len()).To(Equal(1))
			Expect(queue.items[0]).To(Equal(seedName))
		})
	})

	Describe("#shootDelete", func() {
		It("should do nothing because object is not a Shoot", func() {
			obj := &gardencorev1beta1.CloudProfile{}

			c.shootDelete(obj)

			Expect(queue.Len()).To(BeZero())
		})

		It("should do nothing because the seedName is nil", func() {
			obj := &gardencorev1beta1.Shoot{}

			c.shootDelete(obj)

			Expect(queue.Len()).To(BeZero())
		})

		It("should add the object to the queue", func() {
			obj := &gardencorev1beta1.Shoot{
				Spec: gardencorev1beta1.ShootSpec{
					SeedName: &seedName,
				},
			}

			c.shootDelete(obj)

			Expect(queue.Len()).To(Equal(1))
			Expect(queue.items[0]).To(Equal(seedName))
		})
	})
})
