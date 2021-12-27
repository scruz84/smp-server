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
	"github.com/google/uuid"
	logger "github.com/sirupsen/logrus"
	"io"
	"net"
	"sync"
)

type Channel interface {
	HandleRequest(conn net.Conn)
	WriteResponse(responseType []byte, response []byte)
	write(contents []byte)
	GetConnection() net.Conn
	ChannelId() string
}

type ChannelImpl struct {
	conn                 net.Conn
	initialReaderHandler ChannelReadHandler
	initialWriterHandler ChannelWriteHandler
	channelId            string
	writeMutex           sync.Mutex
}

func NewChannelImpl(conn net.Conn, initialReaderHandler ChannelReadHandler, initialWriterHandler ChannelWriteHandler) *ChannelImpl {
	ch := ChannelImpl{conn: conn, initialReaderHandler: initialReaderHandler, initialWriterHandler: initialWriterHandler,
		channelId: uuid.New().String()}
	return &ch
}

func (ch *ChannelImpl) GetConnection() net.Conn {
	return ch.conn
}

func (ch *ChannelImpl) ChannelId() string {
	return ch.channelId
}

func (ch *ChannelImpl) HandleRequest(conn net.Conn) {
	firstMessage := make([]byte, PacketBlock)
	readed, err := io.ReadFull(conn, firstMessage)
	if err != nil {
		RemoveClient(ch)
		return
	}
	if readed != PacketBlock {
		logger.Error("##READED LESS BYTES THAN ", PacketBlock, " (", readed, "). ", conn.RemoteAddr().String())
	}

	ch.initialReaderHandler.Read(firstMessage, ch)
	go ch.HandleRequest(conn)
}

func (ch *ChannelImpl) WriteResponse(responseType []byte, response []byte) {
	if responseType != nil {
		ch.initialWriterHandler.Write(WrapMessageType(responseType, response), ch)
	} else {
		ch.initialWriterHandler.Write(response, ch)
	}
}

func (ch *ChannelImpl) write(contents []byte) {
	ch.writeMutex.Lock()
	defer ch.writeMutex.Unlock()

	var start, writed int
	var err error
	for {
		writed, err = ch.conn.Write(contents)
		if err != nil {
			logger.Error("Error writing to socket", err)
			return
		}
		start += writed
		if writed == 0 || start == len(contents) {
			break
		}
	}
}
