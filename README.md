# Notifier

Notifier is a simple command line tool written in GO and can be used to send notifications through email and slack.

(This project is ready to be updated from now on. Incoming Webhook)

## Overview

Notifier is a command line tool that can send emails and/or slack notifications. More notification methods are to be added. Currently the supported methods are:

- e-mails
- slack message (users have to obtain a slack token before using Notifier)

## Prerequisites

Any environments on which GOLang supports are required. More specifically,
if you use the binary file `notifier`, then the requirements are as follows:

- darwin (Windows may need a small change in the source code.)
- amd64

If you have GOLang on your system, there are no extra requirements. `go get` will handle everything.

Notifier is tested only on macOS and linux.

## Installation

You can install Notifier either by downloading the binary file `notifier` or by using `go get`.

### Download directly

Download the binary file:

- notifier

and put it in `/usr/local/bin` (or any other directory which is included in `$PATH`), then you can use it as a command.

Download the following config files:

- .notifdef.yml
- .notifyrc.yml

and put them direcly under `$HOME`.

Optionally, you can download the following files. You can ignore these optional files.

- error.log
- slackListFile
- emailListFile

### Using `go get`

If you have installed GOLang, then you can easily install Notifier with:

```
go get github.com/charleshenryhugo/Notifier
```

which will download all files to `$GOPATH/src/github.com/` and build a binary file `Notifier` to `$GOPATH/bin/`

Then put binary file in `/usr/local/bin` (or anywhere you like) and the config files just under `$HOME` as described above.

``` shell
cp $GOPATH/bin/Notifier /usr/local/bin/
cp $GOPATH/src/github.com/charleshenryhugo/Notifier/.notifdef.yml $HOME
cp $GOPATH/src/github.com/charleshenryhugo/Notifier/.notifyrc.yml $HOME
```

The second method (`go get`) is recommended because `go get` builds a binary file from GO code optimized to your OS settings.

You can refer to <https://github.com/golang/go> for GO installation.

## Usage

### Options and Commands

Just type `notifier --help` or `notifier -h`, and you see the usage for options and commands:

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

### Configuration files

- $HOME/.notifdef.yml

This file is used for configuring default settings such as a default notification message and a subject.
You can find more details in the file in the repository.

- $HOME/.notifyrc.yml

This file is used for configuring the notification methods such as a slack token and an email account

There is a key `state` in .notifyrc.yml. When its value is `off` (or `false`), any operations associated with that notifier will not be executed. So set the `state` as `on` (or `true`) to make sure that that notifier is valid.
You can find more details in the file in the repository.

### Option Usage

#### Example 1

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

If `U7BL3HC86` is the ID of channel `general`, then `U7BL3HC86` is totally equal to `#general`. Similarly, if `U7BL3IC87` is the ID of user `hugo`, then `U7BL3HC87` is equal to `@hugo`.

Thus, the `slackListFile` above can be modified to the following form:

```
#general
@hugo
#random
@Peter
U7BL3IC90
```

Notifications are still able to be delivered to these channels/users.

Don't forget to add `-x` or `-exe` to explicitly confirm the sending operation

#### Example 2

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

In addition, subject and message will be set according to`$HOME/.notifdef.yml`, because no message or subject option is specified.

#### Example 3

```
notifier -x
```

There are no command line options specified, so all the parameters will be set according to `$HOME/.notifdef.yml`
So the trick is, write all necessary default settings in advance and things become easy.

That is:

- create a file(e.g emailListFile) and write all the target email accounts.
- create a file(e.g slackListFile) and write all the target slack IDs(they must be in your group).
- create a file(e.g error.log) and write the default message you want to send in the next minute or in the future.
- configure the files you have just created (or downloaded) in `$HOME/.notifdef.yml`.
- do some other default settings(please refer to `$HOME/.notifdef.yml`)

### Command Usage

For the usage of each command, just type `notifier [COMMAND] --help`.

#### Example 1

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

#### Example 2

```
notifier default msg "this is a new notification message"
```

This will rewrite the current default notification message to `"this is a new notification message"`.
You can use the command `default` to overwrite any default settings.

However, modifying config files manually is highly recommended.

## Exit Codes

You might want to know if `notifier` did a job or an error occurred. An exit code will tell you the case (e.g. code `130` for `CTRL-C` termination, and use `echo $?` to see it).

`notifier` will exit with an exit code ranging from 1~127 (not all values are used) if any error happened during sending notification (e.g. code `30` for invalid slack token).

For general UNIX/LINUX exit codes, please refer to <http://www.tldp.org/LDP/abs/html/exitcodes.html>

Exit Code |   Temporary or Permanent   |  Meaning | What to Do |
---     |   --- |   --- | --- |
0       |   -   |   notification success | -
1       |   P   |   general error | restart
55      |   P   |   error during parsing .notityrc.yml | check config files
56      |   P   |   error during parsing .defaults.yml | check config files
12 | P | lose internet connection or get refused by remote host | check network, host and port (in config file)
13 | P | error occurs while building a smtp email client | check network, host and port
14 | P | error occurs while authenticating mail account | check your account(address, pasword) and network
15 | P | error occurs while applying email sender  | check your email address
16 | P | error occurs while adding email receivers | check receivers' email address
17 | P | error occurs while initializing or close a iostream for email client | restart
18 | P | error occurs while writing message to email client | restart
19 | P | error occurs while closing an email client | restart
30 | P | slack token is invalid | check your slack token (in config file)
31 | P | target slack user ID or channel ID invalid | check target slack IDs
32 | T | lose internet connection or get refused by slack host | wait for seconds and try again

## Uninstallation

You can also easily uninstall `Notifier` just by removing all the related files and directory, which are:

- Notifier
- .notifdef.yml
- .notifyrc.yml
- $GOPATH/src/github.com/charleshenryhugo/Notifier/

Remove them with:

``` shell
rm $GOPATH/bin/Notifier
rm -rf $GOPATH/src/github.com/charleshenryhugo/Notifier/
rm /usr/local/bin/Notifier
rm $HOME/.notifyrc.yml
rm $HOME/.notifdef.yml
```

## Author

ZHU YUE
