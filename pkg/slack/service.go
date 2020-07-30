package slack

import (
	"github.com/go-logr/logr"
	"github.com/slack-go/slack"
)

// Service structure
type Service struct {
	Log      logr.Logger
	APIToken string
}

var api *slack.Client

func (s *Service) init() {
	if api == nil {
		api = slack.New(s.APIToken)
	}
}

// CreateChannel creates a public or private channel on slack with the given name
func (s *Service) CreateChannel(name string, isPrivate bool) (*ChannelID, error) {
	s.Log.Info("Creating Slack Channel", "name", name, "isPrivate", isPrivate)
	s.init()

	channel, err := api.CreateConversation(name, isPrivate)

	if err != nil {
		s.Log.Error(err, "Error Creating channel", "name", name)
		return nil, err
	}

	s.Log.V(1).Info("Created Slack Channel", "channel", channel)

	channelID := &ChannelID{
		value: channel.ID,
	}
	return channelID, nil
}

// SetDescription sets description/"purpose" of the slack channel
func (s *Service) SetDescription(channelID *ChannelID, description string) error {
	log := s.Log.WithValues("channelID", channelID.Get())
	s.init()

	log.V(1).Info("Setting Description of the Slack Channel")

	_, err := api.SetPurposeOfConversation(channelID.Get(), description)

	if err != nil {
		log.Error(err, "Error setting description of the channel")
		return err
	}
	return nil
}

// SetTopic sets "topic" of the slack channel
func (s *Service) SetTopic(channelID *ChannelID, topic string) error {
	log := s.Log.WithValues("channelID", channelID.Get())
	s.init()

	log.V(1).Info("Setting Topic of the Slack Channel")

	_, err := api.SetTopicOfConversation(channelID.Get(), topic)

	if err != nil {
		log.Error(err, "Error setting topic of the channel")
		return err
	}
	return nil
}

// RenameChannel renames the slack channel
func (s *Service) RenameChannel(channelID *ChannelID, name string) error {
	log := s.Log.WithValues("channelID", channelID.Get())
	s.init()

	log.V(1).Info("Renaming Slack Channel", "name", name)

	_, err := api.RenameConversation(channelID.Get(), name)

	if err != nil {
		log.Error(err, "Error renaming channel")
		return err
	}
	return nil
}

// InviteUsers invites users to the slack channel
func (s *Service) InviteUsers(channelID *ChannelID, userEmails []string) error {
	log := s.Log.WithValues("channelID", channelID.Get())
	s.init()

	for _, email := range userEmails {
		user, err := api.GetUserByEmail(email)

		if err != nil {
			s.Log.Error(err, "Error getting user by email")
			return err
		}

		log.V(1).Info("Inviting user to Slack Channel", "userID", user.ID)
		_, err = api.InviteUsersToConversation(channelID.value, user.ID)

		if err != nil {
			log.Error(err, "Error Inviting user to channel", "userID", user.ID)
			return err
		}
	}
	return nil
}
