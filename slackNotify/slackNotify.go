package slackNotify

import (
	"log"
	"notifier/consts"
	"notifier/parsers"
	"strings"

	"github.com/nlopes/slack"
)

//parse tokens from "notifyrc.xml"
func getToken(ntf parsers.SlackNotifier) string {
	if ntf.Type == "slack" && ntf.State == true {
		return ntf.Token
	}
	return ""
}

//get all channels using your token
func getSlackChannels(token string) (channels []slack.Channel, err error) {
	api := slack.New(token)
	channels, err = api.GetChannels(true)
	return channels, err
}

//get all group users using your token
func getSlackUsers(token string) (users []slack.User, err error) {
	api := slack.New(token)
	users, err = api.GetUsers()
	return users, err
}

//build a slack attachment for slack message parameter and return it
//(to be extended......)
func buildAttachment(title, pretext, text string) slack.Attachment {
	attachment := slack.Attachment{
		Title:   title,
		Pretext: pretext,
		Text:    text,
	}

	return attachment
}

//build a slack message parameter and return it
func buildMessageParameters(attachment slack.Attachment, ntf parsers.SlackNotifier) slack.PostMessageParameters {
	params := slack.PostMessageParameters{}
	params.Attachments = []slack.Attachment{attachment}
	params.AsUser = ntf.AsUser
	params.Username = ntf.UserName
	params.IconEmoji = ":" + ntf.IconEmoji + ":"
	return params
}

//send message to channels using your token parsed from SlackNotifier
func postMsgChannels(ntf parsers.SlackNotifier, channelIDs []string, msgTitle, attachTitle, attachPretext, attachText string) ([]string, string, consts.ERR) {
	token := ntf.Token
	if token == "" {
		log.Println("Your slack token is invalid, please check that.")
		return []string{}, "", consts.SLK_TOKEN_INVAL
	}
	api := slack.New(token)
	msgAttachment := buildAttachment(attachTitle, attachPretext, attachText)
	params := buildMessageParameters(msgAttachment, ntf)

	var (
		timestamp string
		err       error
	)
	for _, channelID := range channelIDs {
		_, timestamp, err = api.PostMessage(channelID, msgTitle, params)
		if err != nil {
			log.Println(err)
			//return exact ERR code using the err string info
			if strings.Contains(err.Error(), "auth") {
				log.Println("Your slack token is invalid, please check that.")
				return channelIDs, timestamp, consts.SLK_TOKEN_INVAL
			} else if strings.Contains(err.Error(), "dial tcp: lookup slack.com: no such host") {
				log.Println("You may lose Internet connection or be refused by remote host.",
					"Try fixing your network and send again")
				return channelIDs, timestamp, consts.SLK_SVR_CONN_ERR
			} else if strings.Contains(err.Error(), "channel_not_found") {
				log.Println("Try checking this slack user(or channel):", channelID, " and send again")
				return channelIDs, timestamp, consts.SLK_CHL_ERR
			}
			return channelIDs, timestamp, consts.SLK_CHL_ERR
		}
		log.Println("slack userID(channelID): ", channelID, " posted successfully")
	}

	return channelIDs, timestamp, consts.SUCCESS
}

//send message to users using your token parsed from SlackNotifier
func postMsgUsers(ntf parsers.SlackNotifier, userIDs []string, msgTitle string, attachment slack.Attachment) ([]string, string, consts.ERR) {
	return postMsgChannels(ntf, userIDs, msgTitle,
		attachment.Title, attachment.Pretext, attachment.Text)
}

//SlackNotify (to []string, subject, msg string, ntfs Notifiers)
//post a notification with subject and message provided with parameters
//to the slack userIDs(ChannelIDs) stored in(to []string)
//ChannelID and UserID are both available
func SlackNotify(to []string, subject, msg string, ntfs parsers.Notifiers) ([]string, string, consts.ERR) {
	ntf := ntfs.SlackNotifier
	if ntf.State == true {
		switch strings.ToLower(ntf.Type) {
		case "slack":
			if len(to) == 0 {
				return []string{}, "", consts.SLK_NOTGT
			}
			attachment := slack.Attachment{Text: msg}
			return postMsgUsers(ntf, to, subject, attachment)
		case "slackwebhook":
			IconEmoji := ":" + ntf.IconEmoji + ":"
			//post to all channelIDs stored in slacklistfile only when there is just one webhook url
			if len(ntf.WebhookURLs) == 1 && len(to) > 0 {
				return to, "",
					postMsgWebhookWithChannels(ntf.WebhookURLs[0], to, subject, msg, ntf.UserName, IconEmoji)
			}
			return to, "", postMsgWebhooks(ntf.WebhookURLs, subject, msg, ntf.UserName, IconEmoji)
		}
	}
	return []string{}, "", consts.SLK_INVAL
}
