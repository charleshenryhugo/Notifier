package slackNotify

import (
	"fmt"
	"net/http"
	"notifier/consts"
	"strings"
)

func buildPayload(channelID, title, text, userName, iconEmoji string) string {
	return `{` + `"username": "` + userName +
		`", "icon_emoji": "` + iconEmoji +
		`", "text": "` + title +
		`", "channel": "` + channelID +
		`", "attachments":[{"text": "` + text +
		`"}]` + `}`
}

func postMsgWebhooks(hookURLs []string, title, text string, userName, iconEmoji string) consts.ERR {
	for _, hurl := range hookURLs {
		if err := postMsgWebhook(hurl, title, text, userName, iconEmoji); err != consts.NIL {
			return err
		}
	}
	fmt.Println("(If the post is [HTTP 200 OK] but you did not receive any notification, please check the webhook urls)")
	return consts.NIL
}

//PostMsgWebhook post a message to the default hookURL channel
func postMsgWebhook(hookURL string, title, text string, userName, iconEmoji string) consts.ERR {
	return postMsgWebhookWithChannel(hookURL, "", title, text, userName, iconEmoji)
}

func postMsgWebhookWithChannels(hookURL string, channelIDs []string, title, text string, userName, iconEmoji string) consts.ERR {
	for _, chID := range channelIDs {
		if err := postMsgWebhookWithChannel(hookURL, chID, title, text, userName, iconEmoji); err != consts.NIL {
			return err
		}
	}
	fmt.Println("(If the post is sucessfully[HTTP 200 OK] but you did not receive any notification, please check the webhook urls)")
	return consts.NIL
}

//PostMsgWebhookWithChannel post a message to the default hookURL channel or to the channel specified by  para:"channel"
func postMsgWebhookWithChannel(hookURL string, channelID, title, text string, userName, iconEmoji string) consts.ERR {
	//build a complete message with attatchments
	payload := buildPayload(channelID, title, text, userName, iconEmoji)
	body := strings.NewReader(payload)
	req, err := http.NewRequest("POST", hookURL, body)
	//fmt.Println("req:", req)
	if err != nil {
		fmt.Println("Please check you network connection and try again.")
		return consts.REQ_FAIL
	}
	resp, err := http.DefaultClient.Do(req)
	//check the response status code. (default: 200 OK)
	switch resp.StatusCode {
	case 400:
		fmt.Println("[HTTP 400 BAD REQUEST]. The payload you sent can not be understood: " + payload)
		return consts.INVALID_PAYLOAD
	case 403:
		fmt.Println("[HTTP 403 FORBIDDEN]. The team associated with your posting has some kind of restriction on the webhook posting in this context")
		return consts.ACTION_FORBID
	case 404:
		fmt.Println("[HTTP 404 NOT FOUND]. Invalid Webhook or channel ID.\nPlease check the target channel \"" + channelID +
			"\" or Webhook url: " + hookURL)
		return consts.CHL_NOT_FOUND
	case 410:
		fmt.Println("[HTTP 410 GONE]. The channel \"" + channelID + "\" has been archived and doesn't accept further messages, even from your incoming webhook")
		return consts.CHL_ARCHIVED
	case 500:
		fmt.Println("[HTTP 500 SERVER ERR]. Something strange and unusual happened that was likely not your fault at all.")
		return consts.ROLLUP_ERROR
	}
	resp.Body.Close()

	fmt.Println("[HTTP 200 OK]. Message posted successfully with webhook url: " + hookURL)
	return consts.NIL
}
