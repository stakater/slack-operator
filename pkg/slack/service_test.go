package slack

import (
	"testing"

	"github.com/stakater/slack-operator/pkg/slack/mock"
	"github.com/stretchr/testify/assert"
)

func TestSlackService_CreateChannel_withPublic_shouldCreatePublicChannel(t *testing.T) {
	s := NewMockService()

	id, err := s.CreateChannel("my-channel", false)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, mock.PublicConversationID, *id)
}

func TestSlackService_CreateChannel_withPrivate_shouldCreatePrivateChannel(t *testing.T) {
	s := NewMockService()

	id, err := s.CreateChannel("my-channel", true)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, mock.PrivateConversationID, *id)
}

func TestSlackService_SetDescription_shouldSetPurpose(t *testing.T) {
	s := NewMockService()

	channel, err := s.SetDescription(mock.PublicConversationID, "myDescription")
	assert.NoError(t, err)
	assert.Equal(t, "myDescription", channel.Purpose.Value)
}

func TestSlackService_SetTopic_shouldSetTopic(t *testing.T) {
	s := NewMockService()
	channel, err := s.SetTopic(mock.PublicConversationID, "myTopic")
	assert.NoError(t, err)
	assert.Equal(t, "myTopic", channel.Topic.Value)
}

func TestSlackService_RenameChannel_shouldSetNewName(t *testing.T) {
	s := NewMockService()
	channel, err := s.RenameChannel(mock.PublicConversationID, "new-channel")
	assert.NoError(t, err)
	assert.Equal(t, "new-channel", channel.Name)
}

func TestSlackService_InviteUsers_shouldSendUserInvites(t *testing.T) {
	s := NewMockService()
	err := s.InviteUsers(mock.PublicConversationID, []string{"spengler@ghostbusters.example.com"})
	assert.NoError(t, err)
}
