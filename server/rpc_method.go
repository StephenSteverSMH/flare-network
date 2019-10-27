package server

import (
	"../node"
	"net"
)

func (server *Server)CreateChannel(to node.NodeAddress, netAddr net.TCPAddr, capFrom int, capTo int){
	init_msg := node.InitNotify{
		To:to,
		N1Cap:capFrom,
		N2Cap:capTo,
		ToNetAddr: netAddr,
	}
	event := node.Event{
		Type:node.INIT,
		Data:init_msg,

	}
	// hook
	server.node.EventChannel <- event
}

func (server *Server)DestoryChannel(to node.NodeAddress){
	destroy_msg := node.DestoryChannelNotify{
		To:to,
	}
	event := node.Event{
		Type: node.CHAN_DOWN,
		Data: destroy_msg,
	}
	server.GetNode().EventChannel <- event
}