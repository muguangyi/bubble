// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package action

import (
	"bubble/env"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	log "github.com/cihub/seelog"
)

// --- Factory ---

// EmailFactory struct.
type EmailFactory struct {
	host     string
	port     int
	auth     bool
	username string
	password string
}

// Validate whether the configure is correct.
func (f *EmailFactory) Validate(conf env.IAny) error {
	if !conf.IsMap() {
		return errors.New("email configure format is incorrect")
	}

	m := conf.Map()

	host, ok := m["smtp"]
	if !ok {
		return errors.New("not setting \"smtp\" for email")
	}
	f.host = host.String()

	port, ok := m["port"]
	if !ok {
		return errors.New("not setting \"port\" for email")
	}
	f.port = port.Int()

	auth, ok := m["auth"]
	if !ok {
		return errors.New("not setting \"auth\" for email")
	}
	f.auth = auth.Bool()

	if f.auth {
		username, ok := m["username"]
		if !ok {
			return errors.New("not setting \"username\" for email")
		}
		f.username = username.String()

		password, ok := m["password"]
		if !ok {
			return errors.New("not setting \"password\" for email")
		}
		f.password = password.String()
	}

	c, err := dial(fmt.Sprintf("%s:%d", f.host, f.port))
	if err != nil {
		return err
	}
	c.Close()

	return nil
}

// Create email action.
func (f *EmailFactory) Create() IAction {
	return &email{f: f}
}

// --- Action ---

type email struct {
	Action
	f *EmailFactory
}

func (e *email) Execute(script env.IAny, target string, env env.IEnv, log ILog) chan bool {
	success := make(chan bool, 1)

	if !script.IsMap() {
		log.Error("email command format is incorrect!")
		e.error = errors.New("email command format is incorrect")
		success <- false
		return success
	}

	m := script.Map()

	to, ok := m["to"]
	if !ok {
		log.Error("Not setting \"to\" in email command!")
		e.error = errors.New("not setting \"to\" in email command")
		success <- false
		return success
	}

	subject, ok := m["subject"]
	if !ok {
		log.Error("Not setting \"subject\" in email command!")
		e.error = errors.New("not setting \"subject\" in email command")
		success <- false
		return success
	}

	body, ok := m["body"]
	if !ok {
		log.Error("Not setting \"body\" in email command!")
		e.error = errors.New("not setting \"body\" in email command")
		success <- false
		return success
	}

	var a smtp.Auth
	if e.f.auth {
		a = smtp.PlainAuth("", e.f.username, e.f.password, e.f.host)
	}

	recipients := env.Format(to)

	header := make(map[string]string)
	header["From"] = e.f.username
	header["To"] = recipients
	header["Subject"] = env.Format(subject)
	header["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + env.Format(body)

	e.error = sendMail(fmt.Sprintf("%s:%d", e.f.host, e.f.port), a, e.f.username, strings.Split(recipients, ";"), []byte(message))
	success <- (e.error == nil)

	return success
}

func dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Split host and port.
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

func sendMail(addr string, auth smtp.Auth, from string, to []string, msg []byte) (err error) {
	//create smtp client
	c, err := dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Error(err)
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		log.Error(err)
		return err
	}

	for _, rcpt := range to {
		if err = c.Rcpt(rcpt); err != nil {
			log.Error(err)
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		log.Error(err)
		return err
	}

	err = w.Close()
	if err != nil {
		log.Error(err)
		return err
	}

	return c.Quit()
}
