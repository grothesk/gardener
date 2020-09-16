// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package internal_test

import (
	"context"
	"fmt"

	"github.com/gardener/gardener/pkg/client/kubernetes"
	"github.com/gardener/gardener/pkg/client/kubernetes/clientmap"
	"github.com/gardener/gardener/pkg/client/kubernetes/clientmap/internal"
	"github.com/gardener/gardener/pkg/client/kubernetes/clientmap/keys"
	fakeclientset "github.com/gardener/gardener/pkg/client/kubernetes/fake"
	"github.com/gardener/gardener/pkg/logger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("GardenClientMap", func() {
	var (
		ctx        context.Context
		cm         clientmap.ClientMap
		key        clientmap.ClientSetKey
		factory    *internal.GardenClientSetFactory
		restConfig *rest.Config
	)

	BeforeEach(func() {
		ctx = context.TODO()
		key = keys.ForGarden()

		restConfig = &rest.Config{}
		factory = &internal.GardenClientSetFactory{
			RESTConfig: restConfig,
		}
		cm = internal.NewGardenClientMap(factory, logger.NewNopLogger())
	})

	Context("#GetClient", func() {
		It("should fail if ClientSetKey type is unsupported", func() {
			key = fakeKey{}
			cs, err := cm.GetClient(ctx, key)
			Expect(cs).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("unsupported ClientSetKey")))
		})

		It("should fail if NewClientSetWithConfig fails", func() {
			fakeErr := fmt.Errorf("fake")
			internal.NewClientSetWithConfig = func(fns ...kubernetes.ConfigFunc) (i kubernetes.Interface, err error) {
				return nil, fakeErr
			}

			cs, err := cm.GetClient(ctx, key)
			Expect(err).To(MatchError(fmt.Sprintf("error creating new ClientSet for key %q: fake", key.Key())))
			Expect(cs).To(BeNil())
		})

		It("should correctly construct a new ClientSet", func() {
			fakeCS := fakeclientset.NewClientSetBuilder().WithRESTConfig(restConfig).Build()
			internal.NewClientSetWithConfig = func(fns ...kubernetes.ConfigFunc) (i kubernetes.Interface, err error) {
				Expect(fns).To(kubernetes.ConsistOfConfigFuncs(
					kubernetes.WithRESTConfig(restConfig),
					kubernetes.WithClientOptions(client.Options{Scheme: kubernetes.GardenScheme}),
				))
				return fakeCS, nil
			}

			cs, err := cm.GetClient(ctx, key)
			Expect(err).NotTo(HaveOccurred())
			Expect(cs).To(BeIdenticalTo(fakeCS))
			Expect(cs.RESTConfig()).To(Equal(restConfig))
		})
	})
})
