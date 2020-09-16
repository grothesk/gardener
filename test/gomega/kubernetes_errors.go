// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package gomega

import (
	"fmt"

	"github.com/onsi/gomega/format"
)

type kubernetesErrors struct {
	checkFunc func(error) bool
	message   string
}

func (k *kubernetesErrors) Match(actual interface{}) (success bool, err error) {
	// is purely nil?
	if actual == nil {
		return false, nil
	}

	actualErr, actualOk := actual.(error)
	if !actualOk {
		return false, fmt.Errorf("expected an error-type.  got:\n%s", format.Object(actual, 1))
	}

	return k.checkFunc(actualErr), nil
}

func (k *kubernetesErrors) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, fmt.Sprintf("to be %s error", k.message))
}
func (k *kubernetesErrors) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, fmt.Sprintf("to not be %s error", k.message))
}
