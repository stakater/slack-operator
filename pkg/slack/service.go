package slack

//TODO: in pkg or in root?
import (
	"github.com/go-logr/logr"
	"github.com/slack-go/slack"
)

// Service interface
type Service interface {
	CreateChannel(string, bool) (*string, error)
	SetDescription(string, string) (*slack.Channel, error)
	SetTopic(string, string) (*slack.Channel, error)
	RenameChannel(string, string) (*slack.Channel, error)
	ArchiveChannel(string) error
	InviteUsers(string, []string) error
}

// SlackService structure
type SlackService struct {
	log logr.Logger
	api *slack.Client
}

// New creates a new SlackService
func New(APIToken string, logger logr.Logger) *SlackService {
	return &SlackService{
		api: slack.New(APIToken),
		log: logger,
	}
}

// CreateChannel creates a public or private channel on slack with the given name
func (s *SlackService) CreateChannel(name string, isPrivate bool) (*string, error) {
	s.log.Info("Creating Slack Channel", "name", name, "isPrivate", isPrivate)

	channel, err := s.api.CreateConversation(name, isPrivate)
	if err != nil {
		s.log.Error(err, "Error Creating channel", "name", name)
		return nil, err
	}

	s.log.V(1).Info("Created Slack Channel", "channel", channel)

	return &channel.ID, nil
}

// SetDescription sets description/"purpose" of the slack channel
func (s *SlackService) SetDescription(channelID string, description string) (*slack.Channel, error) {
	log := s.log.WithValues("channelID", channelID)

	channel, err := s.api.GetConversationInfo(channelID, false)

	if err != nil {
		log.Error(err, "Error fetching channel")
		return nil, err
	}

	if channel.Purpose.Value == description {
		return channel, nil
	}

	log.V(1).Info("Setting Description of the Slack Channel")

	channel, err = s.api.SetPurposeOfConversation(channelID, description)

	if err != nil {
		log.Error(err, "Error setting description of the channel")
		return nil, err
	}
	return channel, nil
}

// SetTopic sets "topic" of the slack channel
func (s *SlackService) SetTopic(channelID string, topic string) (*slack.Channel, error) {
	log := s.log.WithValues("channelID", channelID)

	channel, err := s.api.GetConversationInfo(channelID, false)

	if err != nil {
		log.Error(err, "Error fetching channel")
		return nil, err
	}

	if channel.Topic.Value == topic {
		return channel, nil
	}

	log.V(1).Info("Setting Topic of the Slack Channel")

	channel, err = s.api.SetTopicOfConversation(channelID, topic)

	if err != nil {
		log.Error(err, "Error setting topic of the channel")
		return nil, err
	}
	return channel, nil
}

// RenameChannel renames the slack channel
func (s *SlackService) RenameChannel(channelID string, newName string) (*slack.Channel, error) {
	log := s.log.WithValues("channelID", channelID)

	channel, err := s.api.GetConversationInfo(channelID, false)

	if err != nil {
		log.Error(err, "Error fetching channel")
		return nil, err
	}
	if channel.Name == newName {
		return channel, nil
	}

	log.V(1).Info("Renaming Slack Channel", "newName", newName)

	channel, err = s.api.RenameConversation(channelID, newName)

	if err != nil {
		log.Error(err, "Error renaming channel")
		return nil, err
	}
	return channel, nil
}

// ArchiveChannel archives the slack channel
func (s *SlackService) ArchiveChannel(channelID string) error {
	log := s.log.WithValues("channelID", channelID)

	log.V(1).Info("Archiving channel")
	err := s.api.ArchiveConversation(channelID)

	if err != nil {
		log.Error(err, "Error archiving channel")
		return err
	}

	return nil
}

// InviteUsers invites users to the slack channel
func (s *SlackService) InviteUsers(channelID string, userEmails []string) error {
	log := s.log.WithValues("channelID", channelID)

	for _, email := range userEmails {
		user, err := s.api.GetUserByEmail(email)

		if err != nil {
			log.Error(err, "Error getting user by email")
			return err
		}

		log.V(1).Info("Inviting user to Slack Channel", "userID", user.ID)
		_, err = s.api.InviteUsersToConversation(channelID, user.ID)

		if err != nil {
			log.Error(err, "Error Inviting user to channel", "userID", user.ID)
			return err
		}
	}
	return nil
}
