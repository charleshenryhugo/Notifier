package main

import (
	"log"
	"notifier/consts"
	eml "notifier/emailNotify"
	"notifier/parsers"
	slk "notifier/slackNotify"
	"runtime"

	"github.com/urfave/cli"
)

//MultiRoutineNotify operates all possible notifications
//with multi goroutines
//can increase CPUS with runtime.GOMAXPROCS(2)
func MultiRoutineNotify() error {
	//parse notifiers from notifyrcFile
	ntfs, err := parsers.ParseNotifiers(consts.NotifyrcFile)
	if err != consts.NIL {
		return cli.NewExitError("", int(err))
	}

	//dedicate 2 CPUs to two notifiers
	runtime.GOMAXPROCS(2)

	//ERR channel of each routine
	var EmailERR consts.ERR
	var SlackERR consts.ERR
	chEmailERR := make(chan consts.ERR)
	chSlackERR := make(chan consts.ERR)

	//do email notify
	go func() {
		err := eml.EmailNotify(ToEmailAddrs, Subject, Message, ntfs)
		chEmailERR <- err
	}()

	//do slack notify
	go func() {
		_, _, err := slk.SlackNotify(ToSlackUsers, Subject, Message, ntfs)
		chSlackERR <- err
	}()

	//get ERR from channels
	EmailERR = <-chEmailERR
	SlackERR = <-chSlackERR

	//check email ERR status
	if err := EmailERR; err == consts.NIL {
		log.Println("email notification success")
	} else if err == consts.SMTPM_INVAL {
		log.Println("email notification invalid")
	} else if err == consts.SMTPM_NOTGT {
		log.Println("no target email address(es)")
	} else {
		cli.OsExiter(int(err))
	}

	//check slack ERR status
	if err := SlackERR; err == consts.NIL {
		log.Println("slack notification success")
	} else if err == consts.SLK_NOTGT {
		log.Println("no target slack users(channels)")
	} else if err == consts.SLK_INVAL {
		log.Println("slack notification invalid")
	} else {
		cli.OsExiter(int(err))
	}
	return nil
}

//GenNotify operate all possible notifications
func GenNotify() error {
	//parse notifiers from notifiers config file
	ntfs, err := parsers.ParseNotifiers(consts.NotifyrcFile)
	if err != consts.NIL {
		return cli.NewExitError("", int(err))
	}

	//do email notify
	if err := eml.EmailNotify(ToEmailAddrs, Subject, Message, ntfs); err == consts.NIL {
		log.Println("email notification success")
	} else if err == consts.SMTPM_INVAL {
		log.Println("email notification invalid")
	} else if err == consts.SMTPM_NOTGT {
		log.Println("no target email address(es)")
	} else {
		defer cli.OsExiter(int(err))
	}

	//do slack notify
	if _, _, err := slk.SlackNotify(ToSlackUsers, Subject, Message, ntfs); err == consts.NIL {
		log.Println("slack notification success")
	} else if err == consts.SLK_NOTGT {
		log.Println("no target slack users(channels)")
	} else if err == consts.SLK_INVAL {
		log.Println("slack notification invalid")
	} else {
		defer cli.OsExiter(int(err))
	}
	return nil
}
