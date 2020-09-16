// SPDX-FileCopyrightText: 2019 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package validation_test

import (
	"github.com/gardener/gardener/pkg/apis/core"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"

	. "github.com/gardener/gardener/pkg/apis/core/validation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("validation", func() {
	var controllerInstallation *core.ControllerInstallation

	BeforeEach(func() {
		controllerInstallation = &core.ControllerInstallation{
			ObjectMeta: metav1.ObjectMeta{
				Name: "extension-abc",
			},
			Spec: core.ControllerInstallationSpec{
				RegistrationRef: corev1.ObjectReference{
					Name:            "extension",
					ResourceVersion: "1",
				},
				SeedRef: corev1.ObjectReference{
					Name:            "aws",
					ResourceVersion: "1",
				},
			},
		}
	})

	Describe("#ValidateControllerInstallation", func() {
		It("should forbid empty ControllerInstallation resources", func() {
			errorList := ValidateControllerInstallation(&core.ControllerInstallation{})

			Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeRequired),
				"Field": Equal("metadata.name"),
			})), PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeRequired),
				"Field": Equal("spec.registrationRef.name"),
			})), PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeRequired),
				"Field": Equal("spec.seedRef.name"),
			}))))
		})

		It("should allow valid ControllerInstallation resources", func() {
			errorList := ValidateControllerInstallation(controllerInstallation)

			Expect(errorList).To(BeEmpty())
		})
	})

	Describe("#ValidateControllerInstallationUpdate", func() {
		It("should prevent updating anything if deletion time stamp is set", func() {
			now := metav1.Now()

			newControllerInstallation := prepareControllerInstallationForUpdate(controllerInstallation)
			controllerInstallation.DeletionTimestamp = &now
			newControllerInstallation.DeletionTimestamp = &now
			newControllerInstallation.Spec.RegistrationRef.APIVersion = "another-api-version"

			errorList := ValidateControllerInstallationUpdate(newControllerInstallation, controllerInstallation)

			Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeInvalid),
				"Field": Equal("spec"),
			}))))
		})

		It("should prevent updating immutable fields", func() {
			newControllerInstallation := prepareControllerInstallationForUpdate(controllerInstallation)
			newControllerInstallation.Spec.RegistrationRef.Name = "another-name"
			newControllerInstallation.Spec.RegistrationRef.ResourceVersion = "2"
			newControllerInstallation.Spec.SeedRef.Name = "another-name"
			newControllerInstallation.Spec.SeedRef.ResourceVersion = "2"

			errorList := ValidateControllerInstallationUpdate(newControllerInstallation, controllerInstallation)

			Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeInvalid),
				"Field": Equal("spec.registrationRef.name"),
			})), PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":  Equal(field.ErrorTypeInvalid),
				"Field": Equal("spec.seedRef.name"),
			}))))
		})
	})
})

func prepareControllerInstallationForUpdate(obj *core.ControllerInstallation) *core.ControllerInstallation {
	newObj := obj.DeepCopy()
	newObj.ResourceVersion = "1"
	return newObj
}
