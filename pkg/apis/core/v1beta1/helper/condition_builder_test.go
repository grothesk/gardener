// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package helper_test

import (
	"time"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	. "github.com/gardener/gardener/pkg/apis/core/v1beta1/helper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Builder", func() {
	const (
		conditionType = gardencorev1beta1.ConditionType("Test")
		// re-decalared so the underlying constant is not changed
		unknowStatus       = gardencorev1beta1.ConditionStatus("Unknown")
		fooStatus          = gardencorev1beta1.ConditionStatus("Foo")
		bazReason          = "Baz"
		fubarMessage       = "FuBar"
		unitializedMessage = `The condition has been initialized but its semantic check has not been performed yet.`
		initializedReason  = "ConditionInitialized"
	)

	var (
		defaultTime     metav1.Time
		defaultTimeFunc func() metav1.Time
		codes           = []gardencorev1beta1.ErrorCode{
			gardencorev1beta1.ErrorInfraDependencies,
		}
	)

	BeforeEach(func() {
		defaultTime = metav1.NewTime(time.Unix(2, 2))
		defaultTimeFunc = func() metav1.Time {
			return defaultTime
		}
	})

	Describe("#NewConditionBuilder", func() {
		It("should return error if condition type is empty", func() {
			bldr, err := NewConditionBuilder("")

			Expect(bldr).To(BeNil())
			Expect(err).To(HaveOccurred())
		})

		It("should return not empty builder on success", func() {
			bldr, err := NewConditionBuilder("Foo")

			Expect(bldr).ToNot(BeNil())
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("#Build", func() {
		var (
			result  gardencorev1beta1.Condition
			updated bool
			bldr    ConditionBuilder
		)

		JustBeforeEach(func() {
			bldr, _ = NewConditionBuilder(conditionType)
		})

		Context("empty condition", func() {
			JustBeforeEach(func() {
				result, updated = bldr.WithNowFunc(defaultTimeFunc).Build()
			})

			It("should mark the result as updated", func() {
				Expect(updated).To(BeTrue())
			})

			It("should return correct result", func() {
				Expect(result).To(Equal(gardencorev1beta1.Condition{
					Type:               conditionType,
					Status:             unknowStatus,
					LastTransitionTime: defaultTime,
					LastUpdateTime:     defaultTime,
					Reason:             initializedReason,
					Message:            unitializedMessage,
				}))
			})
		})

		Context("#WithStatus", func() {
			JustBeforeEach(func() {
				result, updated = bldr.
					WithNowFunc(defaultTimeFunc).
					WithStatus(fooStatus).
					Build()
			})

			It("should mark the result as updated", func() {
				Expect(updated).To(BeTrue())
			})

			It("should return correct result", func() {
				Expect(result).To(Equal(gardencorev1beta1.Condition{
					Type:               conditionType,
					Status:             fooStatus,
					LastTransitionTime: defaultTime,
					LastUpdateTime:     defaultTime,
					Reason:             initializedReason,
					Message:            unitializedMessage,
				}))
			})
		})

		Context("#WithReason", func() {
			JustBeforeEach(func() {
				result, updated = bldr.
					WithNowFunc(defaultTimeFunc).
					WithReason(bazReason).
					Build()
			})

			It("should mark the result as updated", func() {
				Expect(updated).To(BeTrue())
			})

			It("should return correct result", func() {
				Expect(result).To(Equal(gardencorev1beta1.Condition{
					Type:               conditionType,
					Status:             unknowStatus,
					LastTransitionTime: defaultTime,
					LastUpdateTime:     defaultTime,
					Reason:             bazReason,
					Message:            unitializedMessage,
				}))
			})
		})

		Context("#WithMessage", func() {
			JustBeforeEach(func() {
				result, updated = bldr.
					WithNowFunc(defaultTimeFunc).
					WithMessage(fubarMessage).
					Build()
			})

			It("should mark the result as updated", func() {
				Expect(updated).To(BeTrue())
			})

			It("should return correct result", func() {
				Expect(result).To(Equal(gardencorev1beta1.Condition{
					Type:               conditionType,
					Status:             unknowStatus,
					LastTransitionTime: defaultTime,
					LastUpdateTime:     defaultTime,
					Reason:             initializedReason,
					Message:            fubarMessage,
				}))
			})
		})

		Context("#WithCodes", func() {
			JustBeforeEach(func() {
				result, updated = bldr.
					WithNowFunc(defaultTimeFunc).
					WithCodes(codes...).
					Build()
			})

			It("should mark the result as updated", func() {
				Expect(updated).To(BeTrue())
			})

			It("should return correct result", func() {
				Expect(result).To(Equal(gardencorev1beta1.Condition{
					Type:               conditionType,
					Status:             unknowStatus,
					LastTransitionTime: defaultTime,
					LastUpdateTime:     defaultTime,
					Reason:             initializedReason,
					Message:            unitializedMessage,
					Codes:              codes,
				}))
			})
		})

		Context("#WithOldCondition", func() {
			JustBeforeEach(func() {
				result, updated = bldr.
					WithNowFunc(defaultTimeFunc).
					WithOldCondition(gardencorev1beta1.Condition{
						Type:               conditionType,
						Status:             fooStatus,
						LastTransitionTime: metav1.NewTime(time.Unix(10, 0)),
						LastUpdateTime:     metav1.NewTime(time.Unix(11, 0)),
						Reason:             bazReason,
						Message:            fubarMessage,
						Codes:              codes,
					}).
					Build()
			})

			It("should mark the result as not updated", func() {
				Expect(updated).To(BeFalse())
			})

			It("should return correct result", func() {
				Expect(result).To(Equal(gardencorev1beta1.Condition{
					Type:               conditionType,
					Status:             fooStatus,
					LastTransitionTime: metav1.NewTime(time.Unix(10, 0)),
					LastUpdateTime:     metav1.NewTime(time.Unix(11, 0)),
					Reason:             bazReason,
					Message:            fubarMessage,
					Codes:              codes,
				}))
			})
		})

		Context("Full override", func() {
			JustBeforeEach(func() {
				result, updated = bldr.
					WithNowFunc(defaultTimeFunc).
					WithStatus("SomeNewStatus").
					WithMessage("Some message").
					WithReason("SomeNewReason").
					WithCodes(codes...).
					WithOldCondition(gardencorev1beta1.Condition{
						Type:               conditionType,
						Status:             fooStatus,
						LastTransitionTime: metav1.NewTime(time.Unix(10, 0)),
						LastUpdateTime:     metav1.NewTime(time.Unix(11, 0)),
						Reason:             bazReason,
						Message:            fubarMessage,
						Codes:              []gardencorev1beta1.ErrorCode{gardencorev1beta1.ErrorInfraQuotaExceeded},
					}).
					Build()
			})

			It("should mark the result as updated", func() {
				Expect(updated).To(BeTrue())
			})

			It("should return correct result", func() {
				Expect(result).To(Equal(gardencorev1beta1.Condition{
					Type:               conditionType,
					Status:             gardencorev1beta1.ConditionStatus("SomeNewStatus"),
					LastTransitionTime: defaultTime,
					LastUpdateTime:     defaultTime,
					Reason:             "SomeNewReason",
					Message:            "Some message",
					Codes:              codes,
				}))
			})
		})
	})
})
