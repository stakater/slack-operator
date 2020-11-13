package mock

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

	"github.com/slack-go/slack/slacktest"
)

// InitSlackTestServer initializes mock server for slack api
func InitSlackTestServer() *slacktest.Server {

	testServer := slacktest.NewTestServer(
		func(c slacktest.Customize) {
			c.Handle("/conversations.info", conversationInfoHandler)
		},
		func(c slacktest.Customize) {
			c.Handle("/conversations.create", createConversationHandler)
		},
		func(c slacktest.Customize) {
			c.Handle("/conversations.setTopic", setConversationTopicHandler)
		},
		func(c slacktest.Customize) {
			c.Handle("/conversations.setPurpose", setConversationPurposeHandler)
		},
		func(c slacktest.Customize) {
			c.Handle("/conversations.rename", renameConversationHandler)
		},
		func(c slacktest.Customize) {
			c.Handle("/conversations.archive", archiveConversationHandler)
		},
		func(c slacktest.Customize) {
			c.Handle("/conversations.invite", inviteConversationHandler)
		},
		func(c slacktest.Customize) {
			c.Handle("/users.lookupByEmail", usersLookupByEmailHandler)
		},
		func(c slacktest.Customize) {
			c.Handle("/conversations.members", getMembersInConversationHandler)
		},
		func(c slacktest.Customize) {
			c.Handle("/conversations.kick", kickMemberFromConversationHandler)
		},
	)

	return testServer
}

// handle conversations.info
func conversationInfoHandler(w http.ResponseWriter, r *http.Request) {

	channelID := extractParamValue(r, "channel")

	var responseJSON string
	if channelID == PublicConversationID {
		responseJSON = publicConversationJSON
	} else if channelID == PrivateConversationID {
		responseJSON = privateConversationJSON
	}

	_, _ = w.Write([]byte(responseJSON))
}

// handle conversations.create
func createConversationHandler(w http.ResponseWriter, r *http.Request) {

	isPrivate := extractParamValue(r, "is_private")
	channelName := extractParamValue(r, "name")

	var responseJSON string

	if channelName == NameTakenConversationName {
		responseJSON = conversationNameTakenJSON
	} else if isPrivate == "true" {
		responseJSON = privateConversationJSON
	} else {
		responseJSON = publicConversationJSON
	}

	_, _ = w.Write([]byte(responseJSON))
}

// handle conversations.setTopic
func setConversationTopicHandler(w http.ResponseWriter, r *http.Request) {
	topic := extractParamValue(r, "topic")
	_, _ = w.Write([]byte(getConversationTopicResponse(topic)))
}

// handle conversations.setPurpose
func setConversationPurposeHandler(w http.ResponseWriter, r *http.Request) {
	purpose := extractParamValue(r, "purpose")
	_, _ = w.Write([]byte(getConversationPurposeResponse(purpose)))
}

// handle conversations.rename
func renameConversationHandler(w http.ResponseWriter, r *http.Request) {
	newName := extractParamValue(r, "name")
	_, _ = w.Write([]byte(getConversationNameResponse(newName)))
}

// handle conversations.archive
func archiveConversationHandler(w http.ResponseWriter, r *http.Request) {
	channelID := extractParamValue(r, "channel")

	response := ""
	if channelID == NotFoundConversationID {
		response = getConversationArchiveChannelNotFoundRespose()
	} else {
		response = getConversationArchiveResponse()
	}
	_, _ = w.Write([]byte(response))
}

// handle conversations.invite
func inviteConversationHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(inviteConversationJSON))
}

// handle conversations.members
func getMembersInConversationHandler(w http.ResponseWriter, r *http.Request) {
	channelID := extractParamValue(r, "channel")

	response := ""
	if channelID == NotFoundConversationID {
		response = getConversationNotFoundResponse()
	} else {
		response = getMembersInConversationResponse()
	}

	_, _ = w.Write([]byte(response))
}

// handle conversations.kick
func kickMemberFromConversationHandler(w http.ResponseWriter, r *http.Request) {
	channelID := extractParamValue(r, "channel")
	userID := extractParamValue(r, "user")

	response := ""
	if channelID == NotFoundConversationID {
		response = getConversationNotFoundResponse()
	} else if userID == NotFoundUserID {
		response = getUserNotFoundResponse()
	} else {
		response = kickUserFromConversationSuccessResponse()
	}

	_, _ = w.Write([]byte(response))
}

// handle users.lookupByEmail
func usersLookupByEmailHandler(w http.ResponseWriter, r *http.Request) {
	email := extractParamValue(r, "email")

	userJSON := ""
	if email == url.QueryEscape(ExistingUserEmail) {
		userJSON = fmt.Sprintf(templateUserJSON, ExistingUserEmail)
	} else {
		userJSON = userNotFoundJSON
	}

	_, _ = w.Write([]byte(userJSON))
}

func extractParamValue(r *http.Request, key string) string {
	buf, bodyErr := ioutil.ReadAll(r.Body)
	if bodyErr != nil {
		return ""
	}

	rdr1, _ := ioutil.ReadAll(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
	r.Body = rdr2

	bodyParams := string(rdr1)
	re := regexp.MustCompile(key + "=(.*?)(&|$)")
	match := re.FindStringSubmatch(bodyParams)

	if len(match) > 1 {
		return match[1]
	}

	return ""
}
