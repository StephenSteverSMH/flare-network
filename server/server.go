package server

import (
	"../network"
	"../node"
	"../utils"
	"fmt"
	"io"
	"net"
)
type Server struct {
	node node.LightNode
}
func (server *Server)StartNodeServer(ip string, port string) error{
	server.node = node.CreateNode(2, 5)
	tcpAddr, err := net.ResolveTCPAddr("tcp", ip+":"+port)
	server.node.NetAddr = tcpAddr
	if err!=nil{
		// 打印本地节点地址解析失败
		return err
	}
	listen, err := net.ListenTCP("tcp", tcpAddr)
	if err!=nil{
		// 打印监听地址失败
		return err
	}
	go server.node.ProcessEvent()
	for {
		conn, err := listen.Accept()
		if err != nil {
			// 打印获取连接失败
			break
		}
		go server.HandleConn(conn)
	}
	return nil
}
func (server *Server)HandleConn(conn net.Conn){
	packet, err :=getPacket(conn)
	if err!=nil{
		// 数据包格式错误
		fmt.Println("数据包格式错误")
		return
	}
	switch packet.Type {
	case network.PACKET_NEIGHBOR:
		msg := node.DiscoverMsg{}
		msg.ConvertFromRaw(packet.Payload)
		//msg.To
		msg.To = server.node.Address
		switch msg.Type {
		case node.NEIGHBOR_HELLO:
			server.node.EventChannel <- node.Event{
				Type:node.NB_UP,
				Data: msg,
			}
			break
		case node.NEIGHBOR_UPD:
			server.node.EventChannel <- node.Event{
				Type:node.NB_UP,
				Data: msg,
			}
			break
		case node.NEIGHBOR_ACK:
			server.node.EventChannel <- node.Event{
				Type:node.ACK,
				Data: msg,
			}
			break
		case node.NEIGHBOR_RST:
			server.node.EventChannel <- node.Event{
				Type:node.RST,
				Data: msg,
			}
			break
		}
	}
	conn.Close()
}

// 获得数据包
func getPacket(conn net.Conn) (network.Packet, error){
	packet := network.Packet{}
	raw_type := make([]byte, 1)
	raw_size := make([]byte, 4)
	_, err :=conn.Read(raw_type)
	if err!=nil{
		if err==io.EOF{
			// 客户端断开连接
			return packet, err
		}
		// 发生其他错误
		return packet, err
	}
	_, err = conn.Read(raw_size)
	if err!=nil{
		if err==io.EOF{
			// 客户端断开连接
			return packet, err
		}
		// 发生其他错误
		return packet, err
	}
	packet.Type = int(utils.BytesToInt32(raw_type))
	packet.Size = int(utils.BytesToInt32(raw_size))
	packet.Payload = make([]byte, packet.Size)
	_, err = conn.Read(packet.Payload)
	if err!=nil{
		if err==io.EOF{
			// 客户端断开连接
			return packet, err
		}
		// 发生其他错误
		return packet, err
	}
	return packet, nil
}
