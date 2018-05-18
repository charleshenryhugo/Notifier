package main

import (
	"io/ioutil"
	"log"
	"notifier/consts"
	"notifier/parsers"
	"strings"

	"github.com/urfave/cli"
)

//confirm to send notifications, default by false
var (
	//SendConfirm(bool) is to confirm the notif-sending operation(set by boolflag -x or -exe)
	SendConfirm = false
)

//global input parameters
//to be added for more notifiers
var (
	Subject          string
	Message          string
	MessageFile      string
	ToEmailAddrs     []string
	ToEmailAddrsFile string
	ToSlackUsers     []string
	ToSlackUsersFile string
)

//usage of global input parameters
//to be added for more notifiers

const (
	subjectFlgUsg          = "Specify the title/subject of your notification (UTF-8, maximum 256 bytes for email notification)"
	messageFlgUsg          = "Specify the message of your notification (UTF-8)"
	msgFileFlgUsg          = "Specify the file that stores your notification message (UTF-8)"
	toEmailAddrsFlgUsg     = "Specify the target email address(es). Do nothing if the email state is off"
	toSlackUsersFlgUsg     = "Specify the target slack userID(s). Do nothing if the slack state is off"
	toEmailAddrsFileFlgUsg = "Specify the file that stores target email address list (one address per line). Do nothing if the email state is off"
	toSlackUsersFileFlgUsg = "Specify the file that stores target slack userID list (one address per line). Do nothing if the email state is off"
)

func appInit() *cli.App {
	app := cli.NewApp()

	app.Name = consts.AppName
	app.Usage = consts.AppUsage
	app.HelpName = consts.AppHelpName
	app.Version = consts.AppVersion
	app.Author = consts.AppAuthor

	return app
}

func appAction(ctx *cli.Context) error {

	//if user didn't specify any arguments
	if !(ctx.IsSet("execute-send") && SendConfirm) {
		log.Println("\nPlease confirm execution using -x or --exe.\nUse -h or --help for more help.")
		return nil
	}

	//parse target IDs from flag arguments
	ToEmailAddrs = ctx.StringSlice("email-addrs")
	ToSlackUsers = ctx.StringSlice("slack-ids")
	//append those email addrs stored in the file, only if the file is available
	//and user didn't specify any email addrs
	if fileBytes, err := ioutil.ReadFile(ToEmailAddrsFile); err == nil && len(ToEmailAddrs) == 0 {
		//ToEmailAddrs = append(strings.Fields(string(fileBytes)), ToEmailAddrs...)
		ToEmailAddrs = strings.Fields(string(fileBytes))
	}
	//append those slack user IDs stored in the file, only if the file is available
	//and user didn't specify any target slack IDs
	if fileBytes, err := ioutil.ReadFile(ToSlackUsersFile); err == nil && len(ToSlackUsers) == 0 {
		//ToSlackUsers = append(strings.Fields(string(fileBytes)), ToSlackUsers...)
		ToSlackUsers = strings.Fields(string(fileBytes))
	}
	//get message from the file(usually error.log), only if the file is available
	//and user didn't specify any message
	if fileBytes, err := ioutil.ReadFile(MessageFile); err == nil && Message == "" {
		Message = string(fileBytes)
	}
	//apply the default settings to message, subject, emails or slacks
	//if any of them is empty
	dflt, err := parsers.ParseDefaults(consts.DefaultsFile)
	if err == consts.NIL {
		//Apply default settings for any empty CLI flags
		if Message == "" {
			Message = dflt.GetDfltmsg()
		}
		if Subject == "" {
			Subject = dflt.GetDfltSbjt()
		}
		if len(ToEmailAddrs) == 0 {
			ToEmailAddrs = dflt.GetDfltEmailList()
		}
		if len(ToSlackUsers) == 0 {
			ToSlackUsers = dflt.GetDfltSlackList()
		}
	}

	//operate all possible notifications
	//using global variables
	return MultiRoutineNotify()
}

func appFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:        "execute-send, exe, x",
			Usage:       "explicitly confirm to send notifications",
			Destination: &SendConfirm,
		},
		cli.StringFlag{
			Name:        "subject, s",
			Usage:       subjectFlgUsg,
			Destination: &Subject,
		},
		cli.StringFlag{
			Name:        "msg, m",
			Usage:       messageFlgUsg,
			Destination: &Message,
		},
		cli.StringFlag{
			Name:        "msgfile, mf",
			Usage:       msgFileFlgUsg,
			Destination: &MessageFile,
		},
		cli.StringFlag{
			Name:        "emails-file, ef",
			Usage:       toEmailAddrsFileFlgUsg,
			Destination: &ToEmailAddrsFile,
		},
		cli.StringFlag{
			Name:        "slacks-file, kf",
			Usage:       toSlackUsersFileFlgUsg,
			Destination: &ToSlackUsersFile,
		},
		cli.StringSliceFlag{
			Name:  "email-addrs, e",
			Usage: toEmailAddrsFlgUsg,
		},
		cli.StringSliceFlag{
			Name:  "slack-ids, k",
			Usage: toSlackUsersFlgUsg,
		},
	}
}

func appCommands() []cli.Command {
	return []cli.Command{
		//change default settings(modify config file)
		{
			Name:    "setdefault",
			Aliases: []string{"default", "def"},
			Usage:   "Change(set) default settings (with some subcommands)",
			Subcommands: []cli.Command{
				{
					Name:    "message",
					Aliases: []string{"msg"},
					Usage:   "Change(set) default message to be sent",
					Action: func(c *cli.Context) error {
						newMsg := c.Args().First()
						return parsers.CfgDfltMsg(newMsg)
					},
				},
				{
					Name:        "subject",
					Aliases:     []string{"title", "sbjt"},
					Usage:       "Change(set) default subject/title to be sent",
					Description: "hahahah",
					Action: func(c *cli.Context) error {
						newSbjt := c.Args().First()
						return parsers.CfgDfltSbjt(newSbjt)
					},
				},
				{
					Name:    "messageFile",
					Aliases: []string{"msgFile", "msgf"},
					Usage:   "Change(set) default file name which stores message",
					Action: func(c *cli.Context) error {
						newMsgFile := c.Args().First()
						return parsers.CfgDfltMsgFile(newMsgFile)
					},
				},
				{
					Name:    "slackListFile",
					Aliases: []string{"kfile", "kf"},
					Usage:   "Change(set) default file name which stores target slack userID(s)",
					Action: func(c *cli.Context) error {
						newSlackListFile := c.Args().First()
						return parsers.CfgDfltSlackListFile(newSlackListFile)
					},
				},
				{
					Name:    "emailListFile",
					Aliases: []string{"efile", "ef"},
					Usage:   "Change(set) default file name which stores target email address(es)",
					Action: func(c *cli.Context) error {
						newEmailListFile := c.Args().First()
						return parsers.CfgDfltEmailListFile(newEmailListFile)
					},
				},
			},
		},
		//change notifiers settings(modify config file)
		{
			Name:        "setnotif",
			Aliases:     []string{"notif"},
			Usage:       "Change(set) notifiers settings, (e.g. slack token, email account)",
			Subcommands: []cli.Command{},
		},
		{
			Name:    "toggle",
			Aliases: []string{"tog"},
			Usage:   "toggle notifier state between 'on' and 'off' ",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "email",
					Usage: "toggle email notifier state",
				},
				cli.BoolFlag{
					Name:  "slack",
					Usage: "toggle slack notifier state",
				},
			},
			Action: func(ctx *cli.Context) error {
				if ctx.Bool("email") {
					parsers.CfgToggStat(consts.EmailNotifier)
				}
				if ctx.Bool("slack") {
					parsers.CfgToggStat(consts.SlackNotifier)
				}
				return nil
			},
		},
	}
}
