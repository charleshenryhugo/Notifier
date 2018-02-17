package main

import (
	"log"
	"strings"

	"github.com/nlopes/slack"
)

//parse tokens from "notifyrc.xml"
func getToken(ntf SlackNotifier) string {
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
func buildMessageParameters(attachment slack.Attachment, ntf SlackNotifier) slack.PostMessageParameters {
	params := slack.PostMessageParameters{}
	params.Attachments = []slack.Attachment{attachment}
	params.AsUser = ntf.AsUser
	params.Username = ntf.UserName
	params.IconEmoji = ":" + ntf.IconEmoji + ":"
	return params
}

//send message to channels using your token parsed from SlackNotifier
func postMsgChannels(ntf SlackNotifier, channelIDs []string, msgTitle, attachTitle, attachPretext, attachText string) ([]string, string, ERR) {
	token := ntf.Token
	if token == "" {
		log.Println("Your slack token is invalid, please check that.")
		return []string{}, "", SLK_TOKEN_INVAL
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
				return channelIDs, timestamp, SLK_TOKEN_INVAL
			} else if strings.Contains(err.Error(), "dial tcp: lookup slack.com: no such host") {
				log.Println("You may lose Internet connection or be refused by remote host.",
					"Try fixing your network and send again")
				return channelIDs, timestamp, SLK_SVR_CONN_ERR
			} else if strings.Contains(err.Error(), "channel_not_found") {
				log.Println("Try checking this slack user(or channel):", channelID, " and send again")
				return channelIDs, timestamp, SLK_CHL_ERR
			}
			return channelIDs, timestamp, SLK_CHL_ERR
		}
		log.Println("slack userID(channelID): ", channelID, " posted successfully")
	}

	return channelIDs, timestamp, SUCCESS
}

//send message to users using your token parsed from SlackNotifier
func postMsgUsers(ntf SlackNotifier, userIDs []string, msgTitle string, attachment slack.Attachment) ([]string, string, ERR) {
	return postMsgChannels(ntf, userIDs, msgTitle,
		attachment.Title, attachment.Pretext, attachment.Text)
}

//SlackNotify (to []string, subject, msg string, ntfs Notifiers)
//post a notification with subject and message provided with parameters
//to the slack userIDs(ChannelIDs) stored in(to []string)
//ChannelID and UserID are both available
func SlackNotify(to []string, subject, msg string, ntfs Notifiers) ([]string, string, ERR) {
	if len(to) == 0 {
		return []string{}, "", SLK_NOTGT
	}
	ntf := ntfs.SlackNotifier
	if ntf.Type == "slack" && (ntf.State == true) {
		attachment := slack.Attachment{Text: msg}
		return postMsgUsers(ntf, to, subject, attachment)
	}
	return []string{}, "", SLK_INVAL
}
