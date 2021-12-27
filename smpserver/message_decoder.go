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
	logger "github.com/sirupsen/logrus"
	"io"
)

type MessageDecoder struct {
	nextHandler ChannelReadHandler
}

func NewMessageDecoder(nextHandler ChannelReadHandler) *MessageDecoder {
	h := MessageDecoder{nextHandler: nextHandler}
	return &h
}

func (h MessageDecoder) Read(message []byte, channel Channel) {

	var messageSize int
	messageSize |= int(message[3])
	messageSize |= int(message[2]) << 8
	messageSize |= int(message[1]) << 16
	messageSize |= int(message[0]) << 24

	var decodedMessage []byte
	if messageSize > len(message) {
		remainingMessage := make([]byte, messageSize-len(message))
		readed, err := io.ReadFull(channel.GetConnection(), remainingMessage)
		if err != nil {
			logger.Error("Error reading client request", err)
			return
		}
		if readed != len(remainingMessage) {
			logger.Error("##READED LESS BYTES THAN ", len(remainingMessage), " (", readed, "). ", channel.GetConnection().RemoteAddr().String())
		}
		decodedMessage = append(message[MessageByteBegin:], remainingMessage...)
	} else {
		decodedMessage = message[MessageByteBegin:messageSize]
	}

	if h.nextHandler != nil {
		h.nextHandler.Read(decodedMessage, channel)
	}
}

func (h MessageDecoder) NextHandler(handler ChannelReadHandler) {
	h.nextHandler = handler
}
