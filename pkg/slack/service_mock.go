package slack

import (
	"github.com/bombsimon/logrusr"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/stakater/slack-operator/pkg/slack/mock"
)

var mockSlackService *SlackService

//NewMockService creates a mock service with SlackTestServer
func NewMockService() *SlackService {

	if mockSlackService == nil {
		logger := logrus.New()
		log := logrusr.NewLogger(logger).WithName("SlackTestServer")

		testServer := mock.InitSlackTestServer()
		go testServer.Start()

		log.Info("Starting Test Server", "url", testServer.GetAPIURL())

		opts := slack.OptionAPIURL(testServer.GetAPIURL())

		mockSlackService = &SlackService{
			api: slack.New("apitoken", opts),
			log: log,
		}
	}

	return mockSlackService
}