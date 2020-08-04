package slack

import (
	"testing"

	"github.com/bombsimon/logrusr"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slacktest"
)

var slackService *SlackService

func newMockService() *SlackService {

	if slackService == nil {
		logger := logrus.New()
		log := logrusr.NewLogger(logger).WithName("SlackServiceTest")

		testServer := slacktest.NewTestServer()
		go testServer.Start()

		log.Info("Starting Test Server", "url", testServer.GetAPIURL())

		opts := slack.OptionAPIURL(testServer.GetAPIURL())

		slackService = &SlackService{
			api: slack.New("apitoken", opts),
			log: log,
		}
	}

	return slackService
}

func TestSlackService_CreateChannel(t *testing.T) {
	s := newMockService()

	//TODO: mock not implemented in slacktest
	id, err := s.CreateChannel("my-channel", false)

	if err != nil {
		panic(err)
	} else {
		s.log.Info("id", *id)
	}
}
