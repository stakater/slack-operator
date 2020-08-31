package controllers

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	slackMock "github.com/stakater/slack-operator/pkg/slack/mock"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("ChannelController", func() {

	var channelName string

	BeforeEach(func() {
		channelName = util.RandSeq(10)
	})

	Describe("Creating SlackChannel resource", func() {
		Context("With required fields", func() {
			It("should set status.ID to public channel ID", func() {
				_ = util.CreateChannel(channelName, false, "", "", []string{}, ns)
				channel := util.GetChannel(channelName, ns)

				//Todo: after util
				//Expect(channel.Status.Conditions.GetCondition("ReconcileError")).To(BeEmpty())
				Expect(channel.Status.ID).To(Equal(slackMock.PublicConversationID))
			})
		})

		Context("With private field true", func() {
			It("should set status.ID status to private channel ID and not set status.Error", func() {
				_ = util.CreateChannel(channelName, true, "", "", []string{}, ns)
				channel := util.GetChannel(channelName, ns)

				//Expect(channel.Status.Error).To(BeEmpty())
				Expect(channel.Status.ID).To(Equal(slackMock.PrivateConversationID))
			})
		})

		Context("With description", func() {
			It("should set channel description", func() {
				description := "my description"

				_ = util.CreateChannel(channelName, true, "", description, []string{}, ns)
				channel := util.GetChannel(channelName, ns)

				//Expect(channel.Status.Error).To(BeEmpty())
				Expect(channel.Spec.Description).To(Equal(description))
			})
		})

		Context("With topic", func() {
			It("should set channel topic", func() {
				topic := "topic of the channel"

				_ = util.CreateChannel(channelName, true, topic, "", []string{}, ns)
				channel := util.GetChannel(channelName, ns)

				//Expect(channel.Status.Error).To(BeEmpty())
				Expect(channel.Spec.Topic).To(Equal(topic))
			})
		})
	})

	Describe("Updating SlackChannel resource", func() {
		Context("With new name", func() {
			It("should assign new name to channel", func() {
				_ = util.CreateChannel(channelName, false, "", "", []string{}, ns)
				channel := util.GetChannel(channelName, ns)

				newName := "old-channel-new-name"
				channel.Spec.Name = newName
				err := k8sClient.Update(ctx, channel)

				if err != nil {
					Fail(err.Error())
				}

				req := reconcile.Request{NamespacedName: types.NamespacedName{Name: channelName, Namespace: ns}}
				_, err = r.Reconcile(req)
				if err != nil {
					Fail(err.Error())
				}

				updatedChannel := util.GetChannel(channelName, ns)

				//Expect(updatedChannel.Status.Error).To(BeEmpty())
				Expect(updatedChannel.Spec.Name).To(Equal(newName))
			})
		})
	})
})
