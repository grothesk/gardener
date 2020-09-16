// SPDX-FileCopyrightText: 2019 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package genericactuator_test

import (
	"context"

	"github.com/gardener/gardener/extensions/pkg/controller/backupentry"
	"github.com/gardener/gardener/extensions/pkg/controller/backupentry/genericactuator"
	mockgenericactuator "github.com/gardener/gardener/pkg/mock/gardener/extensions/controller/backupentry/genericactuator"

	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
)

const (
	providerSecretName      = "backupprovider"
	providerSecretNamespace = "garden"
	shootTechnicalID        = "shoot--foo--bar"
	shootUID                = "asd234-asd-34"
	bucketName              = "test-bucket"
)

var _ = Describe("Actuator", func() {
	var (
		ctrl *gomock.Controller

		be = &extensionsv1alpha1.BackupEntry{
			ObjectMeta: metav1.ObjectMeta{
				Name: shootTechnicalID + "--" + shootUID,
			},
			Spec: extensionsv1alpha1.BackupEntrySpec{
				BucketName: bucketName,
				SecretRef: corev1.SecretReference{
					Name:      providerSecretName,
					Namespace: providerSecretNamespace,
				},
			},
		}

		backupProviderSecretData = map[string][]byte{
			"foo":        []byte("bar"),
			"bucketName": []byte(bucketName),
		}

		beSecret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      providerSecretName,
				Namespace: providerSecretNamespace,
			},
			Data: map[string][]byte{
				"foo": []byte("bar"),
			},
		}

		etcdBackupSecretData = map[string][]byte{
			"bucketName": []byte(bucketName),
			"foo":        []byte("bar"),
		}

		etcdBackupSecretKey = client.ObjectKey{Namespace: shootTechnicalID, Name: genericactuator.EtcdBackupSecretName}
		etcdBackupSecret    = &corev1.Secret{
			TypeMeta: metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
			ObjectMeta: metav1.ObjectMeta{
				Name:            genericactuator.EtcdBackupSecretName,
				Namespace:       shootTechnicalID,
				ResourceVersion: "1",
			},
			Data: etcdBackupSecretData,
		}

		seedNamespace = &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: shootTechnicalID,
			},
		}

		logger = log.Log.WithName("test")
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("#Reconcile", func() {
		var (
			client client.Client
			a      backupentry.Actuator
		)

		Context("seed namespace exist", func() {

			BeforeEach(func() {
				// Create fake client
				client = fakeclient.NewFakeClientWithScheme(scheme.Scheme, seedNamespace, beSecret)
			})

			It("should create secrets", func() {
				// Create mock values provider
				backupEntryDelegate := mockgenericactuator.NewMockBackupEntryDelegate(ctrl)
				backupEntryDelegate.EXPECT().GetETCDSecretData(context.TODO(), be, backupProviderSecretData).Return(etcdBackupSecretData, nil)

				// Create actuator
				a = genericactuator.NewActuator(backupEntryDelegate, logger)
				err := a.(inject.Client).InjectClient(client)
				Expect(err).NotTo(HaveOccurred())

				// Call Reconcile method and check the result
				err = a.Reconcile(context.TODO(), be)
				Expect(err).NotTo(HaveOccurred())

				deployedSecret := &corev1.Secret{}
				err = client.Get(context.TODO(), etcdBackupSecretKey, deployedSecret)
				Expect(err).NotTo(HaveOccurred())
				Expect(deployedSecret).To(Equal(etcdBackupSecret))
			})
		})

		Context("seed namespace does not exist", func() {
			BeforeEach(func() {
				// Create fake client
				client = fakeclient.NewFakeClientWithScheme(scheme.Scheme, beSecret)
			})

			It("should not create secrets", func() {
				// Create mock values provider
				backupEntryDelegate := mockgenericactuator.NewMockBackupEntryDelegate(ctrl)

				// Create actuator
				a = genericactuator.NewActuator(backupEntryDelegate, logger)
				err := a.(inject.Client).InjectClient(client)
				Expect(err).NotTo(HaveOccurred())

				// Call Reconcile method and check the result
				err = a.Reconcile(context.TODO(), be)
				Expect(err).NotTo(HaveOccurred())

				deployedSecret := &corev1.Secret{}
				err = client.Get(context.TODO(), etcdBackupSecretKey, deployedSecret)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
