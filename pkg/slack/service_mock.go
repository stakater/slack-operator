package slack

import (
	"github.com/go-logr/logr"
	"github.com/slack-go/slack"
	"github.com/stakater/slack-operator/pkg/slack/mock"
)

var mockSlackService *SlackService

//NewMockService creates a mock service with SlackTestServer
func NewMockService(log logr.Logger) *SlackService {

	if mockSlackService == nil {
		testServer := mock.InitSlackTestServer()
		go testServer.Start()

		log.Info("Starting Test Server", "url", testServer.GetAPIURL())

		opts := slack.OptionAPIURL(testServer.GetAPIURL())

		mockSlackService = &SlackService{
			api: slack.New("apitoken", opts),
			log: log.WithName("SlackService"),
		}
	}

	return mockSlackService
}
