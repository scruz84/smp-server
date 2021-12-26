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

type MessageEncoder struct {
	nextHandler ChannelWriteHandler
}

func NewMessageEncoder(nextHandler ChannelWriteHandler) *MessageEncoder {
	h := MessageEncoder{nextHandler: nextHandler}
	return &h
}

func (h MessageEncoder) Write(message []byte, channel Channel) {

	messageSize := len(message) + 4
	messageSizeBytes := make([]byte, 4)
	messageSizeBytes[0] = byte(messageSize >> 24 & 0xff)
	messageSizeBytes[1] = byte(messageSize >> 16 & 0xff)
	messageSizeBytes[2] = byte(messageSize >> 8 & 0xff)
	//messageSizeBytes[3] = byte(messageSize >> 0 & 0xff)
	messageSizeBytes[3] = byte(messageSize >> 0)

	var encodedMessage []byte
	if messageSize < PacketBlock {
		encodedMessage = append(messageSizeBytes, message...)
		encodedMessage = append(encodedMessage, make([]byte, PacketBlock-messageSize)...)

	} else {
		encodedMessage = append(messageSizeBytes, message...)
	}

	if h.nextHandler != nil {
		h.nextHandler.Write(encodedMessage, channel)
	} else {
		channel.write(encodedMessage)
	}
}

func (h MessageEncoder) NextHandler(handler ChannelWriteHandler) {
	h.nextHandler = handler
}
