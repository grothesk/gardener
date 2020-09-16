// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package controllerinstallation_test

import (
	"testing"

	"github.com/gardener/gardener/pkg/apis/core"
	"github.com/gardener/gardener/pkg/registry/core/controllerinstallation"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
)

func TestControllerInstallation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry ControllerInstallation Suite")
}

var _ = Describe("ToSelectableFields", func() {
	It("should return correct fields", func() {
		result := controllerinstallation.ToSelectableFields(newControllerInstallation())

		Expect(result).To(HaveLen(3))
		Expect(result.Has("metadata.name")).To(BeTrue())
		Expect(result.Get("metadata.name")).To(Equal("test"))
		Expect(result.Has(core.RegistrationRefName)).To(BeTrue())
		Expect(result.Get(core.RegistrationRefName)).To(Equal("baz"))
		Expect(result.Has(core.SeedRefName)).To(BeTrue())
		Expect(result.Get(core.SeedRefName)).To(Equal("qux"))
	})
})

var _ = Describe("GetAttrs", func() {
	It("should return error when object is not ControllerInstallation", func() {
		_, _, err := controllerinstallation.GetAttrs(&core.ControllerRegistration{})
		Expect(err).To(HaveOccurred())
	})

	It("should return correct result", func() {
		ls, fs, err := controllerinstallation.GetAttrs(newControllerInstallation())

		Expect(err).NotTo(HaveOccurred())
		Expect(ls).To(HaveLen(1))
		Expect(ls.Get("foo")).To(Equal("bar"))
		Expect(fs.Get(core.SeedRefName)).To(Equal("qux"))
	})
})

var _ = Describe("MatchControllerInstallation", func() {
	It("should return correct predicate", func() {
		ls, _ := labels.Parse("app=test")
		fs := fields.OneTermEqualSelector(core.SeedRefName, "foo")

		result := controllerinstallation.MatchControllerInstallation(ls, fs)

		Expect(result.Label).To(Equal(ls))
		Expect(result.Field).To(Equal(fs))
	})
})

func newControllerInstallation() *core.ControllerInstallation {
	return &core.ControllerInstallation{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "test",
			Labels: map[string]string{"foo": "bar"},
		},
		Spec: core.ControllerInstallationSpec{
			RegistrationRef: corev1.ObjectReference{
				Name: "baz",
			},
			SeedRef: corev1.ObjectReference{
				Name: "qux",
			},
		},
	}
}
