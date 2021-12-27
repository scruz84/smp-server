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
	"sync"
)

type TopicMessage struct {
	topic   string  //the topic name
	message []byte  //the message content
	channel Channel //the channel were the message was received
}

var (
	authClients          = make(map[string]*Channel)
	topicSubscribers     = make(map[string]map[string]bool)
	topicMessages        = make(chan TopicMessage)
	topicSubscriversLock sync.RWMutex
)

func Broadcaster() {
	//subscribers := make(map[subscriber]bool)
	for {
		select {
		case topicMessage := <-topicMessages:
			{
				topicSubscriversLock.RLock()
				channelsSubscribedToTopic := topicSubscribers[topicMessage.topic]
				srcChannel := topicMessage.channel.ChannelId()
				for channelId := range channelsSubscribedToTopic {
					if channelId != srcChannel {
						channel := authClients[channelId]
						if channel != nil {
							go (*channel).WriteResponse([]byte{TOPIC_MESSAGE}, topicMessage.message)
						}
					}
				}
				topicSubscriversLock.RUnlock()
			}
		}
	}
}

func AddAuthenticatedClient(channel Channel) {
	authClients[channel.ChannelId()] = &channel
}

func RemoveClient(channel Channel) {
	topicSubscriversLock.Lock()
	defer topicSubscriversLock.Unlock()
	//delete from authenticated clients
	delete(authClients, channel.ChannelId())
	//delete from subscriptions
	for _, channels := range topicSubscribers {
		delete(channels, channel.ChannelId())
	}
}

func SubscribeToTopic(topic string, channel Channel, subscribe bool) {
	topicSubscriversLock.Lock()
	defer topicSubscriversLock.Unlock()
	var subscribers = topicSubscribers[topic]
	if subscribers == nil {
		subscribers = make(map[string]bool)
		topicSubscribers[topic] = subscribers
	}
	if subscribe {
		subscribers[channel.ChannelId()] = true
	} else {
		delete(subscribers, channel.ChannelId())
	}
}

func AddTopicMessage(topicMessage TopicMessage) {
	topicMessages <- topicMessage
}
