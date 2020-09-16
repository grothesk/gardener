// SPDX-FileCopyrightText: 2019 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	mockmanager "github.com/gardener/gardener/pkg/mock/controller-runtime/manager"
	mockwebhook "github.com/gardener/gardener/pkg/mock/gardener/extensions/webhook"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gomodules.xyz/jsonpatch/v2"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var logger = log.Log.WithName("controlplane-webhook-test")

var _ = Describe("Handler", func() {
	const (
		name      = "foo"
		namespace = "default"
	)

	var (
		ctrl *gomock.Controller
		mgr  *mockmanager.MockManager

		objTypes = []runtime.Object{&corev1.Service{}}
		svc      = &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: name},
		}

		req admission.Request
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		// Build scheme
		scheme := runtime.NewScheme()
		_ = corev1.AddToScheme(scheme)

		// Create mock manager
		mgr = mockmanager.NewMockManager(ctrl)
		mgr.EXPECT().GetScheme().Return(scheme)

		req = admission.Request{
			AdmissionRequest: admissionv1beta1.AdmissionRequest{
				Kind:      metav1.GroupVersionKind{Group: "", Version: "v1", Kind: "Service"},
				Name:      name,
				Namespace: namespace,
				Operation: admissionv1beta1.Create,
				Object:    runtime.RawExtension{Raw: encode(svc)},
			},
		}

	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("#Handle", func() {
		It("should return an allowing response if the resource wasn't changed by mutator", func() {
			// Create mock mutator
			mutator := mockwebhook.NewMockMutator(ctrl)
			mutator.EXPECT().Mutate(context.TODO(), svc, nil).Return(nil)

			// Create handler
			h, err := NewBuilder(mgr, logger).WithMutator(mutator, objTypes...).Build()
			Expect(err).NotTo(HaveOccurred())
			Expect(err).NotTo(HaveOccurred())

			// Call Handle and check response
			resp := h.Handle(context.TODO(), req)
			Expect(resp).To(Equal(admission.Response{
				AdmissionResponse: admissionv1beta1.AdmissionResponse{
					Allowed: true,
					Result: &metav1.Status{
						Code: 200,
					},
				},
			}))
		})

		It("should return an allowing response if the resource wasn't changed by mutator and it's update", func() {
			// Create mock mutator
			mutator := mockwebhook.NewMockMutator(ctrl)

			oldSvc := svc.DeepCopy()
			oldSvc.ObjectMeta.Generation = 2

			mutator.EXPECT().Mutate(context.TODO(), svc, oldSvc).Return(nil)

			// Create handler
			h, err := NewBuilder(mgr, logger).WithMutator(mutator, objTypes...).Build()
			Expect(err).NotTo(HaveOccurred())
			Expect(err).NotTo(HaveOccurred())

			req.AdmissionRequest.Operation = admissionv1beta1.Update
			req.AdmissionRequest.OldObject = runtime.RawExtension{Raw: encode(oldSvc)}

			// Call Handle and check response
			resp := h.Handle(context.TODO(), req)
			Expect(resp).To(Equal(admission.Response{
				AdmissionResponse: admissionv1beta1.AdmissionResponse{
					Allowed: true,
					Result: &metav1.Status{
						Code: 200,
					},
				},
			}))
		})

		It("should return a patch response if the resource was changed by mutator", func() {
			// Create mock mutator
			mutator := mockwebhook.NewMockMutator(ctrl)
			mutator.EXPECT().Mutate(context.TODO(), svc, nil).DoAndReturn(func(ctx context.Context, obj, oldOjb runtime.Object) error {
				accessor, _ := meta.Accessor(obj)
				accessor.SetAnnotations(map[string]string{"foo": "bar"})
				return nil
			})

			// Create handler
			h, err := NewBuilder(mgr, logger).WithMutator(mutator, objTypes...).Build()
			Expect(err).NotTo(HaveOccurred())
			Expect(err).NotTo(HaveOccurred())

			// Call Handle and check response
			resp := h.Handle(context.TODO(), req)
			pt := admissionv1beta1.PatchTypeJSONPatch
			Expect(resp).To(Equal(admission.Response{
				Patches: []jsonpatch.JsonPatchOperation{
					{
						Operation: "add",
						Path:      "/metadata/annotations",
						Value:     map[string]interface{}{"foo": "bar"},
					},
				},
				AdmissionResponse: admissionv1beta1.AdmissionResponse{
					Allowed:   true,
					PatchType: &pt,
				},
			}))
		})

		It("should return an error response if the mutator returned an error", func() {
			// Create mock mutator
			mutator := mockwebhook.NewMockMutator(ctrl)
			mutator.EXPECT().Mutate(context.TODO(), svc, nil).Return(errors.New("test error"))

			// Create handler
			h, err := NewBuilder(mgr, logger).WithMutator(mutator, objTypes...).Build()
			Expect(err).NotTo(HaveOccurred())
			Expect(err).NotTo(HaveOccurred())

			// Call Handle and check response
			resp := h.Handle(context.TODO(), req)
			Expect(resp).To(Equal(admission.Response{
				AdmissionResponse: admissionv1beta1.AdmissionResponse{
					Allowed: false,
					Result: &metav1.Status{
						Code:    http.StatusBadRequest,
						Message: "test error",
					},
				},
			}))
		})
	})
})

func encode(obj runtime.Object) []byte {
	data, _ := json.Marshal(obj)
	return data
}
