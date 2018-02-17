# Notifier

Notifier is a simple command line tool written in GO and can be used to send notifications through email and slack.

## Overview

Notifier is a command line tool, which can send email and slack notifications.(More notification methods to be added)

Users will have to obtain a valid email account and slack token before using Notifier.

Users can only send notifications to their own slack group specified by the slack token.

Get slack token from:

<https://api.slack.com/custom-integrations/legacy-tokens>

## Requirements

- darwin (UNIX-like, Mach, BSD)
- amd64

## Installation

### download directly

download the binary file:

- Notifier

and put it in `/usr/local/bin` (or any other directory which is included by $PATH)

download the following config files:

- .defaults.yml
- .notifyrc.yml

and put them in $HOME

Optional: (you can ignore these optional files below)

- error.log
- slackListFile
- emailListFile

### use homebrew

To be added.

## Usage

### Options and Commands

Just type `notifier --help` or `notifier -h` , out comes the usage for options and commands:

```
COMMANDS:
     setdefault, default, def  Change(set) default settings (with some subcommands)
     setnotif, notif           Change(set) notifiers settings, (e.g. slack token, email account)
     toggle, tog               toggle notifier state between 'on' and 'off'
     help, h                   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --email-addrs value, -e value    Specify the target email address(es). Do nothing if the email state is off
   --emails-file value, --ef value  Specify the file that stores target email address list (one address per line). Do nothing if the email state is off
   --execute-send, --exe, -x        explicitly confirm to send notifications
   --msg value, -m value            Specify the message of your notification (UTF-8)
   --msgfile value, --mf value      Specify the file that stores your notification message (UTF-8)
   --slack-ids value, -k value      Specify the target slack userID(s). Do nothing if the slack state is off
   --slacks-file value, --kf value  Specify the file that stores target slack userID list (one address per line). Do nothing if the email state is off
   --subject value, -s value        Specify the title/subject of your notification (UTF-8, maximum 256 bytes for emailnotification)
   --help, -h                       show help
   --version, -v                    print the version
```

### files

- $HOME/.defaults.yml

This file is used for configuring default settings such as default notification message and subject.

(You can find more details in the file itself.)

- $HOME/.notifyrc.yml

This file is used for configuring the notification methods such as slack token and email account

There is an key `state` in .notifyrc.yml. If it's value is `off` or `false`, any operations associated with that notifier will not be executed. So set the `state` as `on` or `true` to make sure that notifier is valid.

(You can find more details in the file itself.)

### Options

e.g.1

```
notifier -x -s "new notif" -m "some error happened!" -e "google@gmail.com" -e "yahoo@gmail.com" -kf "somedir/slackListFile"
```

This will send a notification with subject:"new notif!" and message:"some error happened!"
to google@gmail.com and yahoo@gmail.com as well as slack users(or channels) that have IDs stored in "somedir/slackListFile",
which looks like this:

somedir/slackListFile

```
U7BL3HC86
U7BL3IC87
U7BL3IC88
U7BL3IC89
U7BL3IC90
```

One ID in a line and no blank line.

Don't forget to add `-x` or `-exe` to explicitly confirm the sending operation

e.g.2

```
notifier -x -ef "somedir/emailListFile" -k U7BL3HC86 -k U7BL3HC87 -k U7BL3HC88
```

This will send a notification to slack ID U7BL3HC86, U7BL3HC87, U7BL3HC88 and the email addresses stored in somedir/emailListFile
which looks like this:
somedir/emailListFile

```
google@gmail.com
yahoo@gmail.com
```

One email address in a line and no blank line.

In addition, subject and message will be set according to`$HOME/.defaults.yml`, because no message or subject option is specified.

e.g.3

```
notifier -x
```

There are no command line options specified, so all the parameters w ill be set according to `$HOME/.defaults.yml`
So the trick is, write all necessary default settings in advance and things become easy.

That is:

- create a file(e.g emailListFile) and write all the target email accounts.
- create a file(e.g slackListFile) and write all the target slack IDs(they must be in your group).
- create a file(e.g error.log) and write the default message you want to send in the next minute or in the future.
- configure the files you have just created (or downloaded) in `$HOME/.defaults.yml`.
- do some other default settings(please refer to `$HOME/.defaults.yml`)

### Commands

For the usage of each command, just type `notifier [COMMAND] --help`.

e.g.1

```
notifier default --help
```

out comes usage for command `setdefault`(or `default`, `def`) and it's subcommands:

```
NAME:
   Notifier setdefault - Change(set) default settings (with some subcommands)

USAGE:
   Notifier setdefault command [command options] [arguments...]

COMMANDS:
     message, msg                Change(set) default message to be sent
     subject, title, sbjt        Change(set) default subject/title to be sent
     messageFile, msgFile, msgf  Change(set) default file name which stores message
     slackListFile, kfile, kf    Change(set) default file name which stores target slack userID(s)
     emailListFile, efile, ef    Change(set) default file name which stores target email address(es)

OPTIONS:
   --help, -h  show help
```

e.g.2

```
notifier default msg "this is a new notification message"
```

This will rewrite the current default notification message to `"this is a new notification message"`.
You can use the command `default` to overwrite any default settings.

However, modifying config files manually is highly recommended.

## Error Codes

## Demo

## Author

ZHU YUE
