package main

import (
	"log"
	"runtime"

	"github.com/urfave/cli"
)

//MultiRoutineNotify operates all possible notifications
//with multi goroutines
//can increase CPUS with runtime.GOMAXPROCS(2)
func MultiRoutineNotify() error {
	//parse notifiers from notifyrcFile
	ntfs, err := parseNotifiers(notifyrcFile)
	if err != NIL {
		return cli.NewExitError("", int(err))
	}

	//dedicate 2 CPUs to two notifiers
	runtime.GOMAXPROCS(2)

	//ERR channel of each routine
	var EmailERR ERR
	var SlackERR ERR
	chEmailERR := make(chan ERR)
	chSlackERR := make(chan ERR)

	//do email notify
	go func() {
		err := EmailNotify(ToEmailAddrs, Subject, Message, ntfs)
		chEmailERR <- err
	}()

	//do slack notify
	go func() {
		_, _, err := SlackNotify(ToSlackUsers, Subject, Message, ntfs)
		chSlackERR <- err
	}()

	//get ERR from channels
	EmailERR = <-chEmailERR
	SlackERR = <-chSlackERR

	//check email ERR status
	if err := EmailERR; err == NIL {
		log.Println("email notification success")
	} else if err == SMTPM_INVAL {
		log.Println("email notification invalid")
	} else if err == SMTPM_NOTGT {
		log.Println("no target email address(es)")
	} else {
		cli.OsExiter(int(err))
	}

	//check slack ERR status
	if err := SlackERR; err == NIL {
		log.Println("slack notification success")
	} else if err == SLK_NOTGT {
		log.Println("no target slack users(channels)")
	} else if err == SLK_INVAL {
		log.Println("slack notification invalid")
	} else {
		cli.OsExiter(int(err))
	}
	return nil
}

//GenNotify operate all possible notifications
func GenNotify() error {
	//parse notifiers from notifiers config file
	ntfs, err := parseNotifiers(notifyrcFile)
	if err != NIL {
		return cli.NewExitError("", int(err))
	}

	//do email notify
	if err := EmailNotify(ToEmailAddrs, Subject, Message, ntfs); err == NIL {
		log.Println("email notification success")
	} else if err == SMTPM_INVAL {
		log.Println("email notification invalid")
	} else if err == SMTPM_NOTGT {
		log.Println("no target email address(es)")
	} else {
		defer cli.OsExiter(int(err))
	}

	//do slack notify
	if _, _, err := SlackNotify(ToSlackUsers, Subject, Message, ntfs); err == NIL {
		log.Println("slack notification success")
	} else if err == SLK_NOTGT {
		log.Println("no target slack users(channels)")
	} else if err == SLK_INVAL {
		log.Println("slack notification invalid")
	} else {
		defer cli.OsExiter(int(err))
	}
	return nil
}
