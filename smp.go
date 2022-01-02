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
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/howeyc/gopass"
	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
	"os"
	"smp/smpserver"
	"smp/smpserver/database"
	"strconv"
	"sync"
)

const (
	connType           = "tcp"
	createUserFlagName = "create-user"
)

type Tls struct {
	Enabled    bool   `yaml:"enabled"`
	Port       int    `yaml:"tls_port"`
	ServerKey  string `yaml:"server_key"`
	ServerCert string `yaml:"server_cert"`
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Tls  Tls    `yaml:"tls"`
}

func main() {

	serverConfig := loadConfig()

	//input parameters
	createUserParam := flag.String(createUserFlagName, "user", "Create a new user")
	flag.Parse()

	//starts the database
	defer database.Close()

	//verify input parameter and execute them
	if isFlagSet(createUserFlagName) {
		createUser(*createUserParam)
	} else {
		runServer(serverConfig)
	}
}

func runServer(serverConfiguration map[string]Server) {
	database.Init(true)

	var waitingGroup sync.WaitGroup

	waitingGroup.Add(1)
	go startListener(serverConfiguration["server"], &waitingGroup)

	if serverConfiguration["server"].Tls.Enabled {
		waitingGroup.Add(1)
		go startTLSListener(serverConfiguration["server"], &waitingGroup)
	}

	waitingGroup.Wait()
}

func startListener(serverConfiguration Server, waitingGroup *sync.WaitGroup) {
	defer waitingGroup.Done()
	l, err := net.Listen(connType, serverConfiguration.Host+":"+strconv.Itoa(serverConfiguration.Port))
	if err != nil {
		logger.Error("Error starting listening:", err)
		os.Exit(1)
	}
	logger.Info("Starting listening on ", serverConfiguration.Host, ":", serverConfiguration.Port)
	defer l.Close()

	go smpserver.Broadcaster() // starts the message broadcaster

	listenConnections(l, false)
}

func startTLSListener(serverConfiguration Server, waitingGroup *sync.WaitGroup) {
	defer waitingGroup.Done()

	cert, err := tls.LoadX509KeyPair(serverConfiguration.Tls.ServerCert, serverConfiguration.Tls.ServerKey)
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.Rand = rand.Reader
	l, err := tls.Listen(connType, serverConfiguration.Host+":"+strconv.Itoa(serverConfiguration.Tls.Port),
		&config)
	if err != nil {
		logger.Error("Error starting TLS listening:", err)
		os.Exit(1)
	}
	defer l.Close()
	logger.Info("Starting listening TLS on ", serverConfiguration.Host, ":", serverConfiguration.Tls.Port)

	go smpserver.Broadcaster() // starts the message broadcaster

	listenConnections(l, true)
}

func listenConnections(l net.Listener, isTLS bool) {
	for {
		// Listen for new connections
		conn, err := l.Accept()
		if err != nil {
			logger.Error("Error accepting new connection: ", err.Error())
		}

		//TLS if required
		if isTLS {
			tlsCon, ok := conn.(*tls.Conn)
			if ok {
				state := tlsCon.ConnectionState()
				for _, v := range state.PeerCertificates {
					logger.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
				}
			}
			err = tlsCon.Handshake()
			if err != nil {
				logger.Error("error on handshake. ", err)
			}
		}

		// Handle new client
		messageDecoder := smpserver.NewMessageDecoder(smpserver.NewMessageDispatcher(nil))
		messageEncoder := smpserver.NewMessageEncoder(
			smpserver.NewFixedLengthFragmentEncoder(nil, smpserver.PacketBlock))
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

func loadConfig() map[string]Server {
	yamlFile, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		logger.Error("Error opening configuration file. ", err)
		panic(err)
	}
	config := make(map[string]Server)
	err2 := yaml.Unmarshal(yamlFile, &config)
	if err2 != nil {
		logger.Error("Error opening configuration file. ", err2)
		panic(err2)
	}
	return config
}
