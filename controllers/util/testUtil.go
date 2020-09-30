package util

import (
	"context"
	"math/rand"
	"time"

	"github.com/onsi/ginkgo"
	ginko "github.com/onsi/ginkgo"
	slackv1alpha1 "github.com/stakater/slack-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// TestUtil contains necessary objects required to perform operations during tests
type TestUtil struct {
	ctx       context.Context
	k8sClient client.Client
	r         reconcile.Reconciler
}

// New creates new TestUtil
func New(ctx context.Context, k8sClient client.Client, r reconcile.Reconciler) *TestUtil {
	return &TestUtil{
		ctx:       ctx,
		k8sClient: k8sClient,
		r:         r,
	}
}

// CreateChannel creates and submits a Slack Channel object to the kubernetes server
func (t *TestUtil) CreateChannel(name string, isPrivate bool, topic string, description string, users []string, namespace string) *slackv1alpha1.Channel {
	channelObject := t.CreateSlackChannelObject(name, isPrivate, topic, description, users, namespace)
	err := t.k8sClient.Create(t.ctx, channelObject)

	if err != nil {
		ginko.Fail(err.Error())
	}

	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: namespace}}

	_, err = t.r.Reconcile(req)
	if err != nil {
		ginko.Fail(err.Error())
	}

	return channelObject
}

// GetChannel fetches a channel object from kubernetes
func (t *TestUtil) GetChannel(name string, namespace string) *slackv1alpha1.Channel {
	channelObject := &slackv1alpha1.Channel{}
	err := t.k8sClient.Get(t.ctx, types.NamespacedName{Name: name, Namespace: namespace}, channelObject)

	if err != nil {
		ginko.Fail(err.Error())
	}

	return channelObject
}

// DeleteChannel deletes the channel resource
func (t *TestUtil) DeleteChannel(name string, namespace string) {
	channelObject := &slackv1alpha1.Channel{}
	err := t.k8sClient.Get(t.ctx, types.NamespacedName{Name: name, Namespace: namespace}, channelObject)

	if err != nil {
		ginko.Fail(err.Error())
	}

	err = t.k8sClient.Delete(t.ctx, channelObject)

	if err != nil {
		ginko.Fail(err.Error())
	}

	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: namespace}}

	_, err = t.r.Reconcile(req)
	if err != nil {
		ginko.Fail(err.Error())
	}
}

// TryDeleteChannel - Tries to delete channel if it exists, does not fail on any error
func (t *TestUtil) TryDeleteChannel(name string, namespace string) {
	channelObject := &slackv1alpha1.Channel{}
	_ = t.k8sClient.Get(t.ctx, types.NamespacedName{Name: name, Namespace: namespace}, channelObject)
	_ = t.k8sClient.Delete(t.ctx, channelObject)
	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: namespace}}
	_, _ = t.r.Reconcile(req)
}

// CreateNamespace creates a namespace in the kubernetes server
func (t *TestUtil) CreateNamespace(name string) {
	namespaceObject := t.CreateNamespaceObject(name)
	err := t.k8sClient.Create(t.ctx, namespaceObject)

	if err != nil {
		ginkgo.Fail(err.Error())
	}
}

// CreateSlackChannelObject creates a slack channel custom resource object
func (t *TestUtil) CreateSlackChannelObject(name string, isPrivate bool, topic string, description string, users []string, namespace string) *slackv1alpha1.Channel {
	return &slackv1alpha1.Channel{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: slackv1alpha1.ChannelSpec{
			Name:        name,
			Private:     isPrivate,
			Topic:       topic,
			Description: description,
			Users:       users,
		},
	}
}

// CreateNamespaceObject creates a namespace object
func (t *TestUtil) CreateNamespaceObject(name string) *v1.Namespace {
	return &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

// RandSeq Generates a letter sequence with `n` characters
func (t *TestUtil) RandSeq(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")

	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
