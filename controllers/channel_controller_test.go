package controllers

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	slackv1alpha1 "github.com/stakater/slack-operator/api/v1alpha1"
	"github.com/stakater/slack-operator/pkg/slack/mock"
	slackMock "github.com/stakater/slack-operator/pkg/slack/mock"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("ChannelController", func() {

	var channelName string

	BeforeEach(func() {
		channelName = util.RandSeq(10)
	})

	AfterEach(func() {
		util.TryDeleteChannel(channelName, ns)
	})

	Describe("Creating SlackChannel resource", func() {
		Context("With required fields", func() {
			It("should set status.ID to public channel ID", func() {
				_ = util.CreateChannel(channelName, false, "", "", []string{mock.ExistingUserEmail}, ns)
				channel := util.GetChannel(channelName, ns)

				Expect(channel.Status.ID).To(Equal(slackMock.PublicConversationID))
				Expect(len(channel.Status.Conditions)).To(Equal(1))
				Expect(channel.Status.Conditions[0].Reason).To(Equal("Successful"))
			})
		})

		Context("With private field true", func() {
			It("should set status.ID status to private channel ID and not set status.Error", func() {
				_ = util.CreateChannel(channelName, true, "", "", []string{mock.ExistingUserEmail}, ns)
				channel := util.GetChannel(channelName, ns)

				Expect(channel.Status.ID).To(Equal(slackMock.PrivateConversationID))
				Expect(len(channel.Status.Conditions)).To(Equal(1))
				Expect(channel.Status.Conditions[0].Reason).To(Equal("Successful"))
			})
		})

		Context("With description", func() {
			It("should set channel description", func() {
				description := "my description"

				_ = util.CreateChannel(channelName, true, "", description, []string{mock.ExistingUserEmail}, ns)
				channel := util.GetChannel(channelName, ns)

				Expect(channel.Spec.Description).To(Equal(description))
				Expect(len(channel.Status.Conditions)).To(Equal(1))
				Expect(channel.Status.Conditions[0].Reason).To(Equal("Successful"))
			})
		})

		Context("With topic", func() {
			It("should set channel topic", func() {
				topic := "topic of the channel"

				_ = util.CreateChannel(channelName, true, topic, "", []string{mock.ExistingUserEmail}, ns)
				channel := util.GetChannel(channelName, ns)

				Expect(channel.Spec.Topic).To(Equal(topic))
				Expect(len(channel.Status.Conditions)).To(Equal(1))
				Expect(channel.Status.Conditions[0].Reason).To(Equal("Successful"))
			})
		})

		Context("With user emails", func() {
			It("should set success condition when user exists", func() {

				_ = util.CreateChannel(channelName, true, "", "", []string{mock.ExistingUserEmail}, ns)
				channel := util.GetChannel(channelName, ns)

				Expect(len(channel.Status.Conditions)).To(Equal(1))
				Expect(channel.Status.Conditions[0].Reason).To(Equal("Successful"))
			})

			It("should set error condition when user does not exists", func() {

				_ = util.CreateChannel(channelName, true, "", "", []string{"nonexistent@slack.com"}, ns)
				channel := util.GetChannel(channelName, ns)

				Expect(len(channel.Status.Conditions)).To(Equal(1))
				Expect(channel.Status.Conditions[0].Reason).To(Equal("Failed"))
				Expect(channel.Status.Conditions[0].Message).To(Equal("users_not_found"))
			})
		})
	})

	Describe("Updating SlackChannel resource", func() {
		Context("With new name", func() {
			It("should assign new name to channel", func() {
				_ = util.CreateChannel(channelName, false, "", "", []string{mock.ExistingUserEmail}, ns)
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

				Expect(updatedChannel.Spec.Name).To(Equal(newName))

				Expect(len(channel.Status.Conditions)).To(Equal(1))
				Expect(channel.Status.Conditions[0].Reason).To(Equal("Successful"))
			})
		})
	})

	Describe("Deleting SlackChannel resource", func() {
		Context("When Channel on slack was created", func() {
			It("should remove resource and delete channel ", func() {
				_ = util.CreateChannel(channelName, false, "", "", []string{mock.ExistingUserEmail}, ns)
				channel := util.GetChannel(channelName, ns)

				Expect(channel.Status.ID).ToNot(BeEmpty())
				Expect(len(channel.Status.Conditions)).To(Equal(1))
				Expect(channel.Status.Conditions[0].Reason).To(Equal("Successful"))

				util.DeleteChannel(channelName, ns)

				channelObject := &slackv1alpha1.Channel{}
				err := k8sClient.Get(ctx, types.NamespacedName{Name: channelName, Namespace: ns}, channelObject)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("When Channel on slack was not created", func() {
			It("should remove resource ", func() {
				_ = util.CreateChannel(mock.InvalidChannelName, false, "", "", []string{mock.ExistingUserEmail}, ns)
				channel := util.GetChannel(mock.InvalidChannelName, ns)

				Expect(len(channel.Status.Conditions)).To(Equal(1))
				Expect(channel.Status.Conditions[0].Reason).To(Equal("Failed"))

				util.DeleteChannel(mock.InvalidChannelName, ns)

				channelObject := &slackv1alpha1.Channel{}
				err := k8sClient.Get(ctx, types.NamespacedName{Name: channelName, Namespace: ns}, channelObject)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
