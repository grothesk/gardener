// SPDX-FileCopyrightText: 2018 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package logger_test

import (
	"fmt"
	"os"

	. "github.com/gardener/gardener/pkg/logger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("logger", func() {
	Describe("logger", func() {
		AfterEach(func() {
			Logger = nil
		})

		Describe("#NewLogger", func() {
			It("should return a pointer to a Logger object ('info' level)", func() {
				logger := NewLogger("info")

				Expect(logger.Out).To(Equal(os.Stderr))
				Expect(logger.Level).To(Equal(logrus.InfoLevel))
				Expect(Logger).To(Equal(logger))
			})

			It("should return a pointer to a Logger object ('debug' level)", func() {
				logger := NewLogger("debug")

				Expect(logger.Out).To(Equal(os.Stderr))
				Expect(logger.Level).To(Equal(logrus.DebugLevel))
				Expect(Logger).To(Equal(logger))
			})

			It("should return a pointer to a Logger object ('error' level)", func() {
				logger := NewLogger("error")

				Expect(logger.Out).To(Equal(os.Stderr))
				Expect(logger.Level).To(Equal(logrus.ErrorLevel))
				Expect(Logger).To(Equal(logger))
			})
		})

		Describe("#NewShootLogger", func() {
			It("should return an Entry object with additional fields (w/o operationID)", func() {
				logger := NewLogger("info")
				namespace := "core"
				name := "shoot01"

				shootLogger := NewShootLogger(logger, name, namespace)

				Expect(shootLogger.Data).To(HaveKeyWithValue("shoot", fmt.Sprintf("%s/%s", namespace, name)))
			})
		})

		Describe("#NewFieldLogger", func() {
			It("should return an Entry object with additional fields", func() {
				logger := NewLogger("info")
				key := "foo"
				value := "bar"

				fieldLogger := NewFieldLogger(logger, key, value)

				Expect(fieldLogger.Data).To(HaveKeyWithValue(key, value))
			})
		})

		Describe("#AddWriter", func() {
			It("should return a pointer to a Test Logger object ('info' level)", func() {
				logger := AddWriter(NewLogger(""), GinkgoWriter)
				Expect(logger.Out).To(Equal(GinkgoWriter))
				Expect(logger.Level).To(Equal(logrus.InfoLevel))
				Expect(Logger).To(Equal(logger))
			})
		})

	})
})
