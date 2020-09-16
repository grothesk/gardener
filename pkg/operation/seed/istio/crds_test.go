// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package istio_test

import (
	"context"

	cr "github.com/gardener/gardener/pkg/chartrenderer"
	"github.com/gardener/gardener/pkg/client/kubernetes"
	"github.com/gardener/gardener/pkg/operation/botanist/component"
	. "github.com/gardener/gardener/pkg/operation/seed/istio"
	. "github.com/gardener/gardener/test/gomega"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/version"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("#CRDs", func() {
	var (
		ctx context.Context
		c   client.Client
		crd component.DeployWaiter
	)

	BeforeEach(func() {
		ctx = context.TODO()

		s := runtime.NewScheme()
		Expect(apiextensionsv1beta1.AddToScheme(s)).NotTo(HaveOccurred())
		Expect(apiextensionsv1beta1.AddToScheme(s)).NotTo(HaveOccurred())

		c = fake.NewFakeClientWithScheme(s)

		mapper := meta.NewDefaultRESTMapper([]schema.GroupVersion{apiextensionsv1beta1.SchemeGroupVersion})
		mapper.Add(apiextensionsv1beta1.SchemeGroupVersion.WithKind("CustomResourceDefinition"), meta.RESTScopeRoot)

		renderer := cr.NewWithServerVersion(&version.Info{})

		ca := kubernetes.NewChartApplier(renderer, kubernetes.NewApplier(c, mapper))
		Expect(ca).NotTo(BeNil(), "should return chart applier")

		crd = NewIstioCRD(ca, chartsRootPath, c)

	})

	JustBeforeEach(func() {
		deprecatedCRDs := []apiextensionsv1beta1.CustomResourceDefinition{
			{ObjectMeta: metav1.ObjectMeta{Name: "meshpolicies.authentication.istio.io"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "policies.authentication.istio.io"}},
		}

		for _, deprecated := range deprecatedCRDs {
			Expect(c.Create(ctx, &deprecated)).ToNot(HaveOccurred())
		}

		Expect(crd.Deploy(ctx)).ToNot(HaveOccurred(), "istio crd deploy succeeds")
	})

	DescribeTable("CRD is deployed",
		func(crdName string) {
			Expect(c.Get(
				ctx,
				client.ObjectKey{Name: crdName},
				&apiextensionsv1beta1.CustomResourceDefinition{},
			)).ToNot(HaveOccurred())
		},
		Entry("DestinationRule", "destinationrules.networking.istio.io"),
		Entry("EnvoyFilter", "envoyfilters.networking.istio.io"),
		Entry("Gateways", "gateways.networking.istio.io"),
		Entry("ServiceEntry", "serviceentries.networking.istio.io"),
		Entry("Sidecar", "sidecars.networking.istio.io"),
		Entry("VirtualServices", "virtualservices.networking.istio.io"),
		Entry("AuthorizationPolicy", "authorizationpolicies.security.istio.io"),
		Entry("PeerAuthentication", "peerauthentications.security.istio.io"),
		Entry("RequestAuthentications", "requestauthentications.security.istio.io"),
		Entry("WorkloadEntries", "workloadentries.networking.istio.io"),
		// TODO (mvladev): Entries bellow should be moved to unused CRDs table when
		// they are no longer used by future versions of istio.
		Entry("HTTPAPISpec (DEPRECATED, but needed)", "httpapispecs.config.istio.io"),
		Entry("QuotaSpecBinding (DEPRECATED, but needed)", "quotaspecbindings.config.istio.io"),
		Entry("HTTPAPISpecBinding (DEPRECATED, but needed)", "httpapispecbindings.config.istio.io"),
		Entry("QuotaSpec (DEPRECATED, but needed)", "quotaspecs.config.istio.io"),
		Entry("ClusterRBACConfig (DEPRECATED, but needed)", "clusterrbacconfigs.rbac.istio.io"),
		Entry("RBACConfig (DEPRECATED, but needed)", "rbacconfigs.rbac.istio.io"),
		Entry("ServiceRole (DEPRECATED, but needed)", "serviceroles.rbac.istio.io"),
		Entry("ServiceRoleBindings (DEPRECATED, but needed)", "servicerolebindings.rbac.istio.io"),
	)

	DescribeTable("unused CRDs are not deployed",
		func(crdName string) {
			Expect(c.Get(
				ctx,
				client.ObjectKey{Name: crdName},
				&apiextensionsv1beta1.CustomResourceDefinition{},
			)).To(BeNotFoundError())
		},

		Entry("AttributeManifsts", "attributemanifests.config.istio.io"),
		Entry("Handlers", "handlers.config.istio.io"),
		Entry("Instances", "instances.config.istio.io"),
		Entry("Rules", "rules.config.istio.io"),
		Entry("MeshPolicy", "meshpolicies.authentication.istio.io"),
		Entry("Policy", "policies.authentication.istio.io"),
	)
})
