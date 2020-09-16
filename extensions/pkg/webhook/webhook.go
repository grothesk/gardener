// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	// TargetSeed defines that the webhook is to be installed in the seed.
	TargetSeed = "seed"
	// TargetShoot defines that the webhook is to be installed in the shoot.
	TargetShoot = "shoot"

	// ValidatorName is a common name for a validation webhook.
	ValidatorName = "validator"
	// ValidatorPath is a common path for a validation webhook.
	ValidatorPath = "/webhooks/validate"
)

// Webhook is the specification of a webhook.
type Webhook struct {
	Name     string
	Kind     string
	Provider string
	Path     string
	Target   string
	Types    []runtime.Object
	Webhook  *admission.Webhook
	Handler  http.Handler
	Selector *metav1.LabelSelector
}

type Args struct {
	Provider   string
	Name       string
	Path       string
	Predicates []predicate.Predicate
	Validators map[Validator][]runtime.Object
	Mutators   map[Mutator][]runtime.Object
}

// New creates a new Webhook with the given args.
func New(mgr manager.Manager, args Args) (*Webhook, error) {
	logger := log.Log.WithName(args.Name).WithValues("provider", args.Provider)

	// Create handler
	builder := NewBuilder(mgr, logger)

	for val, objs := range args.Validators {
		builder.WithValidator(val, objs...)
	}

	for mut, objs := range args.Mutators {
		builder.WithMutator(mut, objs...)
	}

	builder.WithPredicates(args.Predicates...)

	handler, err := builder.Build()
	if err != nil {
		return nil, err
	}

	// Create webhook
	logger.Info("Creating webhook")

	return &Webhook{
		Path:    args.Path,
		Webhook: &admission.Webhook{Handler: handler},
	}, nil
}
