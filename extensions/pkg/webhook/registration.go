// SPDX-FileCopyrightText: 2019 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// RegisterWebhooks registers the given webhooks in the Kubernetes cluster targeted by the provided manager.
func RegisterWebhooks(ctx context.Context, mgr manager.Manager, namespace, providerName string, servicePort int, mode, url string, caBundle []byte, webhooks []*Webhook) (webhooksToRegisterSeed []admissionregistrationv1beta1.MutatingWebhook, webhooksToRegisterShoot []admissionregistrationv1beta1.MutatingWebhook, err error) {
	var (
		fail                             = admissionregistrationv1beta1.Fail
		ignore                           = admissionregistrationv1beta1.Ignore
		mutatingWebhookConfigurationSeed = &admissionregistrationv1beta1.MutatingWebhookConfiguration{ObjectMeta: metav1.ObjectMeta{Name: "gardener-extension-" + providerName}}
	)

	for _, webhook := range webhooks {
		var rules []admissionregistrationv1beta1.RuleWithOperations
		for _, t := range webhook.Types {
			rule, err := buildRule(mgr, t)
			if err != nil {
				return nil, nil, err
			}
			rules = append(rules, *rule)
		}

		webhookToRegister := admissionregistrationv1beta1.MutatingWebhook{
			Name:              fmt.Sprintf("%s.%s.extensions.gardener.cloud", webhook.Name, strings.TrimPrefix(providerName, "provider-")),
			NamespaceSelector: webhook.Selector,
			Rules:             rules,
		}

		switch webhook.Target {
		case TargetSeed:
			webhookToRegister.FailurePolicy = &fail
			webhookToRegister.ClientConfig = buildClientConfigFor(webhook, namespace, providerName, servicePort, mode, url, caBundle)
			webhooksToRegisterSeed = append(webhooksToRegisterSeed, webhookToRegister)
		case TargetShoot:
			webhookToRegister.FailurePolicy = &ignore
			webhookToRegister.ClientConfig = buildClientConfigFor(webhook, namespace, providerName, servicePort, ModeURLWithServiceName, url, caBundle)
			webhooksToRegisterShoot = append(webhooksToRegisterShoot, webhookToRegister)
		default:
			return nil, nil, fmt.Errorf("invalid webhook target: %s", webhook.Target)
		}
	}

	if len(webhooksToRegisterSeed) > 0 {
		c, err := getClient(mgr)
		if err != nil {
			return nil, nil, err
		}

		if _, err := controllerutil.CreateOrUpdate(ctx, c, mutatingWebhookConfigurationSeed, func() error {
			mutatingWebhookConfigurationSeed.Webhooks = webhooksToRegisterSeed
			return nil
		}); err != nil {
			return nil, nil, err
		}
	}

	return webhooksToRegisterSeed, webhooksToRegisterShoot, nil
}

// buildRule creates and returns a RuleWithOperations for the given object type.
func buildRule(mgr manager.Manager, t runtime.Object) (*admissionregistrationv1beta1.RuleWithOperations, error) {
	// Get GVK from the type
	gvk, err := apiutil.GVKForObject(t, mgr.GetScheme())
	if err != nil {
		return nil, errors.Wrapf(err, "could not get GroupVersionKind from object %v", t)
	}

	// Get REST mapping from GVK
	mapping, err := mgr.GetRESTMapper().RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get REST mapping from GroupVersionKind '%s'", gvk.String())
	}

	// Create and return RuleWithOperations
	return &admissionregistrationv1beta1.RuleWithOperations{
		Operations: []admissionregistrationv1beta1.OperationType{
			admissionregistrationv1beta1.Create,
			admissionregistrationv1beta1.Update,
		},
		Rule: admissionregistrationv1beta1.Rule{
			APIGroups:   []string{gvk.Group},
			APIVersions: []string{gvk.Version},
			Resources:   []string{mapping.Resource.Resource},
		},
	}, nil
}

func buildClientConfigFor(webhook *Webhook, namespace, providerName string, servicePort int, mode, url string, caBundle []byte) admissionregistrationv1beta1.WebhookClientConfig {
	path := "/" + webhook.Path

	clientConfig := admissionregistrationv1beta1.WebhookClientConfig{
		CABundle: caBundle,
	}

	switch mode {
	case ModeURL:
		url := fmt.Sprintf("https://%s%s", url, path)
		clientConfig.URL = &url
	case ModeURLWithServiceName:
		url := fmt.Sprintf("https://gardener-extension-%s.%s:%d%s", providerName, namespace, servicePort, path)
		clientConfig.URL = &url
	case ModeService:
		clientConfig.Service = &admissionregistrationv1beta1.ServiceReference{
			Namespace: namespace,
			Name:      "gardener-extension-" + providerName,
			Path:      &path,
		}
	}

	return clientConfig
}
