/*
 * smp_server.
 * Copyright (C) 2021  Sergio Cruz
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 *  along with this program; if not, write to the Free Software
 *  Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
 */

package main

import (
	"flag"
	"fmt"
	"github.com/howeyc/gopass"
	logger "github.com/sirupsen/logrus"
	"net"
	"os"
	"smp/smpserver"
	"smp/smpserver/database"
)

const (
	connHost = "0.0.0.0"
	connPort = "1984"
	connType = "tcp"

	createUserFlagName = "create-user"
)

func main() {

	//input parameters
	createUserParam := flag.String(createUserFlagName, "user", "Create a new user")
	flag.Parse()

	//starts the database
	defer database.Close()

	//verify input parameter and execute them
	if isFlagSet(createUserFlagName) {
		createUser(*createUserParam)
	} else {
		runServer()
	}
}

func runServer() {
	logger.Info("Starting the server on ", connHost, ":", connPort)
	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		logger.Error("Error starting server listening:", err)
		os.Exit(1)
	}
	defer l.Close()

	database.Init(true)

	go smpserver.Broadcaster() // starts the message broadcaster

	for {
		// Listen for new connections
		conn, err := l.Accept()
		if err != nil {
			logger.Error("Error accepting new connection: ", err.Error())
		}
		// Handle new client
		messageDecoder := smpserver.NewMessageDecoder(
			smpserver.NewMessageDispatcher(nil))
		messageEncoder := smpserver.NewMessageEncoder(nil)
		channel := smpserver.NewChannelImpl(conn, messageDecoder, messageEncoder)
		go channel.HandleRequest(conn)
	}
}

func createUser(user string) {
	database.Init(false)

	fmt.Print("Password: ")
	var maskedPassword []byte
	var err error

	maskedPassword, err = gopass.GetPasswdMasked()
	if err != nil {
		logger.Error("Error accepting new connection: ", err.Error())
		return
	}

	err = database.StoreUser(user, string(maskedPassword))
	if err != nil {
		logger.Error("Error creating new user: ", err.Error())
		return
	}
}

func isFlagSet(mFlag string) bool {
	set := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == mFlag {
			set = true
		}
	})
	return set
}
