package slack

import (
	"testing"

	"github.com/stakater/slack-operator/pkg/slack/mock"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var log = zap.New()

func TestSlackService_CreateChannel_shouldCreatePublicChannel_whenPrivateIsFalse(t *testing.T) {
	s := NewMockService(log)

	id, err := s.CreateChannel("my-channel", false)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, mock.PublicConversationID, *id)
}

func TestSlackService_CreateChannel_shouldCreatePrivateChannel_whenPrivateIsTrue(t *testing.T) {
	s := NewMockService(log)

	id, err := s.CreateChannel("my-channel", true)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, mock.PrivateConversationID, *id)
}

func TestSlackService_CreateChannel_shouldThrowError_whenChannelWithSameNameExists(t *testing.T) {
	s := NewMockService(log)

	_, err := s.CreateChannel(mock.NameTakenConversationName, true)

	assert.EqualError(t, err, "name_taken")
}

func TestSlackService_SetDescription_shouldSetPurpose(t *testing.T) {
	s := NewMockService(log)

	channel, err := s.SetDescription(mock.PublicConversationID, "myDescription")
	assert.NoError(t, err)
	assert.Equal(t, "myDescription", channel.Purpose.Value)
}

func TestSlackService_SetTopic_shouldSetTopic(t *testing.T) {
	s := NewMockService(log)
	channel, err := s.SetTopic(mock.PublicConversationID, "myTopic")
	assert.NoError(t, err)
	assert.Equal(t, "myTopic", channel.Topic.Value)
}

func TestSlackService_RenameChannel_shouldSetNewName(t *testing.T) {
	s := NewMockService(log)
	channel, err := s.RenameChannel(mock.PublicConversationID, "new-channel")
	assert.NoError(t, err)
	assert.Equal(t, "new-channel", channel.Name)
}

func TestSlackService_ArchiveChannel_shouldArchiveChannel(t *testing.T) {
	s := NewMockService(log)
	err := s.ArchiveChannel(mock.PublicConversationID)
	assert.NoError(t, err)
}

func TestSlackService_ArchiveChannel_shouldThrowError_whenChannelNotFound(t *testing.T) {
	s := NewMockService(log)
	err := s.ArchiveChannel(mock.NotFoundConversationID)
	assert.EqualError(t, err, "channel_not_found")
}

func TestSlackService_InviteUsers_shouldSendUserInvites_whenUserExists(t *testing.T) {
	s := NewMockService(log)
	err := s.InviteUsers(mock.PublicConversationID, []string{mock.ExistingUserEmail})
	assert.NoError(t, err)
}

func TestSlackService_InviteUsers_shouldThowError_whenUserDoesNotExists(t *testing.T) {
	s := NewMockService(log)
	err := s.InviteUsers(mock.PublicConversationID, []string{"spengler@ghostbusters.example.com"})
	assert.EqualError(t, err, "users_not_found")
}
