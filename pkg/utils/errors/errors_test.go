// SPDX-FileCopyrightText: 2018 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package errors_test

import (
	"fmt"
	"testing"

	errorsmock "github.com/gardener/gardener/pkg/mock/gardener/utils/errors"
	utilerrors "github.com/gardener/gardener/pkg/utils/errors"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestErrors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Errors Suite")
}

var _ = Describe("Errors", func() {
	var (
		err1, err2 error
	)
	BeforeEach(func() {
		err1 = fmt.Errorf("error 1")
		err2 = fmt.Errorf("error 2")
	})

	Describe("#WithSuppressed", func() {
		It("should return nil if the error is nil", func() {
			Expect(utilerrors.WithSuppressed(nil, err2)).To(BeNil())
		})

		It("should return the error if the suppressed error is nil", func() {
			Expect(utilerrors.WithSuppressed(err1, nil)).To(BeIdenticalTo(err1))
		})

		It("should return an error with cause and suppressed equal to the given errors", func() {
			err := utilerrors.WithSuppressed(err1, err2)

			Expect(errors.Cause(err)).To(BeIdenticalTo(err1))
			Expect(utilerrors.Suppressed(err)).To(BeIdenticalTo(err2))
		})
	})

	Describe("#Suppressed", func() {
		It("should retrieve the suppressed error", func() {
			Expect(utilerrors.Suppressed(utilerrors.WithSuppressed(err1, err2))).To(BeIdenticalTo(err2))
		})

		It("should retrieve nil if the error doesn't have a suppressed error", func() {
			Expect(utilerrors.Suppressed(err1)).To(BeNil())
		})
	})

	Context("withSuppressed", func() {
		Describe("#Error", func() {
			It("should return an error message consisting of the two errors", func() {
				Expect(utilerrors.WithSuppressed(err1, err2).Error()).To(Equal("error 1, suppressed: error 2"))
			})
		})

		Describe("#Format", func() {
			It("should correctly format the error in verbose mode", func() {
				Expect(fmt.Sprintf("%+v", utilerrors.WithSuppressed(err1, err2))).
					To(Equal("error 1\nsuppressed: error 2"))
			})
		})
	})

	Describe("Error context", func() {
		It("Should panic with duplicate error IDs", func() {
			defer func() {
				_ = recover()
			}()

			errorContext := utilerrors.NewErrorContext("Test context", nil)
			errorContext.AddErrorID("ID1")
			errorContext.AddErrorID("ID1")
			Fail("Panic should have occurred")
		})
	})

	Describe("Error handling", func() {
		var (
			errorContext *utilerrors.ErrorContext
			ctrl         *gomock.Controller
		)

		BeforeEach(func() {
			errorContext = utilerrors.NewErrorContext("Test context", nil)
			ctrl = gomock.NewController(GinkgoT())
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		It("Should update the error context", func() {
			errID := "x1"
			Expect(utilerrors.HandleErrors(errorContext,
				nil,
				nil,
				utilerrors.ToExecute(errID, func() error {
					return nil
				}),
			)).To(Succeed())
			Expect(errorContext.HasErrorWithID(errID)).To(BeTrue())
		})

		It("Should call default failure handler", func() {
			errorID := "x1"
			errorText := fmt.Sprintf("Error in %s", errorID)
			expectedErr := utilerrors.WithID(errorID, fmt.Errorf("%s failed (%s)", errorID, errorText))
			err := utilerrors.HandleErrors(errorContext,
				nil,
				nil,
				utilerrors.ToExecute(errorID, func() error {
					return fmt.Errorf(errorText)
				}),
			)

			Expect(err).To(Equal(expectedErr))
		})

		It("Should call failure handler on fail", func() {
			errID := "x1"
			errorText := "Error from task"
			expectedErr := fmt.Errorf("Got %s %s", errID, errorText)
			err := utilerrors.HandleErrors(errorContext,
				nil,
				func(errorID string, err error) error {
					return fmt.Errorf(fmt.Sprintf("Got %s %s", errorID, err))
				},
				utilerrors.ToExecute(errID, func() error {
					return fmt.Errorf(errorText)
				}),
			)

			Expect(err).To(Equal(expectedErr))
		})

		It("Should return a cancelError when manually canceled", func() {
			errID := "x1"
			err := utilerrors.HandleErrors(errorContext,
				nil,
				nil,
				utilerrors.ToExecute(errID, func() error {
					return utilerrors.Cancel()
				}),
			)

			Expect(utilerrors.WasCanceled(errors.Cause(err))).To(BeTrue())
		})

		It("Should stop execution on error", func() {
			expectedErr := fmt.Errorf("Err1")
			f1 := errorsmock.NewMockTaskFunc(ctrl)
			f2 := errorsmock.NewMockTaskFunc(ctrl)
			f3 := errorsmock.NewMockTaskFunc(ctrl)

			f1.EXPECT().Do(errorContext).Return("x1", nil)
			f2.EXPECT().Do(errorContext).Return("x2", expectedErr)
			f3.EXPECT().Do(errorContext).Times(0)

			err := utilerrors.HandleErrors(errorContext,
				nil,
				func(errorID string, err error) error {
					return err
				},
				f1,
				f2,
				f3,
			)

			Expect(err).To(Equal(expectedErr))
		})

		It("Should call success handler on error resolution", func() {
			errID := "x2"
			errorContext := utilerrors.NewErrorContext("Check success handler", []string{errID})
			Expect(utilerrors.HandleErrors(errorContext,
				func(errorID string) error {
					return nil
				},
				nil,
				utilerrors.ToExecute("x1", func() error {
					return nil
				}),
				utilerrors.ToExecute(errID, func() error {
					return nil
				}),
			)).To(Succeed())
		})

		It("Should execute methods sequentially in the specified order", func() {
			f1 := errorsmock.NewMockTaskFunc(ctrl)
			f2 := errorsmock.NewMockTaskFunc(ctrl)
			f3 := errorsmock.NewMockTaskFunc(ctrl)

			gomock.InOrder(
				f1.EXPECT().Do(errorContext).Return("x1", nil),
				f2.EXPECT().Do(errorContext).Return("x2", nil),
				f3.EXPECT().Do(errorContext).Return("x3", nil),
			)

			err := utilerrors.HandleErrors(errorContext,
				nil,
				func(errorID string, err error) error {
					return err
				},
				f1,
				f2,
				f3,
			)

			Expect(err).To(Succeed())
		})
	})
})

var _ = Describe("Multierrors", func() {
	var (
		allErrs    *multierror.Error
		err1, err2 error
	)

	BeforeEach(func() {
		err1 = fmt.Errorf("error 1")
		err2 = fmt.Errorf("error 2")
	})

	Describe("#NewErrorFormatFuncWithPrefix", func() {
		BeforeEach(func() {
			allErrs = &multierror.Error{
				ErrorFormat: utilerrors.NewErrorFormatFuncWithPrefix("prefix"),
			}
		})

		It("should format a multierror correctly if it contains 1 error", func() {
			allErrs.Errors = []error{err1}
			Expect(allErrs.Error()).To(Equal("prefix: 1 error occurred: error 1"))
		})
		It("should format a multierror correctly if it contains multiple errors", func() {
			allErrs.Errors = []error{err1, err2}
			Expect(allErrs.Error()).To(Equal("prefix: 2 errors occurred: [error 1, error 2]"))
		})
	})
})
