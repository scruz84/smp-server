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

const MESSAGE_TYPE_BYTES = 4

type MessageDispatcher struct {
	nextHandler ChannelReadHandler
}

func NewMessageDispatcher(nextHandler ChannelReadHandler) *MessageDispatcher {
	h := MessageDispatcher{nextHandler: nextHandler}
	return &h
}

func (h MessageDispatcher) Read(message []byte, channel Channel) {

	if len(message) > 0 {
		switch message[0] {
		case LOGIN_REQUEST:
			doLogin(message[MESSAGE_TYPE_BYTES:], channel)
		case TOPIC_SUBSCRIPTION_REQUEST:
			doTopicSubscription(message[MESSAGE_TYPE_BYTES:], channel)
		case SEND_MESSAGE:
			go doSendMessage(message[MESSAGE_TYPE_BYTES:], channel)
		}
	}
}

func (h MessageDispatcher) NextHandler(handler ChannelReadHandler) {
	h.nextHandler = handler
}
