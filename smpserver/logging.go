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

package smpserver

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"smp/smpserver/database"
)

type LoginInfo struct {
	User     string
	Password string
}

type LoginReponse struct {
	Status bool
	Error  string
}

//message:
// XC
//  X: auth type (0: user/password)
//  C: message
//response:
// SC
//  S: status (0: error, 1: ok)
//  C: reponse string status
func doLogin(message []byte, channel Channel) {

	var authError error
	if message[0] == 0 {
		authError = doLoginPasswordAuthentication(message[1:])
	}
	var response []byte
	response = make([]byte, 1)
	if authError != nil {
		response[0] = 0
		response = append(response, []byte(authError.Error())...)
	} else {
		response[0] = 1
		AddAuthenticatedClient(channel)
	}
	go channel.WriteResponse([]byte{LOGIN_RESPONSE}, response)
}

//message:
// SUP
// S: user size
// U: user name
// P: password
func doLoginPasswordAuthentication(message []byte) error {
	userSize := int(message[0])
	userName := string(message[1 : userSize+1])
	password := string(message[userSize+1:])

	var errorResult error
	err := database.VerifyUser(userName, password)
	if err != nil {
		errorResult = errors.New("authentication error")
		logger.Error("authentication error ", err.Error())
	}
	return errorResult
}
