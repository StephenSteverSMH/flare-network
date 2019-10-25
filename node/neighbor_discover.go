package node

import (
	"math/rand"
	. "../utils"
)
const (
	NEIGHBOR_HELLO = iota
	NEIGHBOR_UPD
	NEIGHBOR_RST
	NEIGHBOR_ACK
)

// 数据发现包
type DiscoverMsg struct {
	Type int
	From NodeAddress
	To NodeAddress
	NewRoute Channels
	CreatedTime int64
	Identify int
}

// 消息发现数据包格式
// discover_packet := [type(1B)|from(20B)|created_time(8B)|indetify(4B)|chan_route_num(4B)|channel_payload(?B)]
// channel_payload := [n1(20B)|n2(20B)|n1_cap(4B)|n2_cap(4B)|status(1B)] * chan_route_num


func (msg *DiscoverMsg)ConvertToRaw() []byte{
	var raw []byte
	raw_type := byte(msg.Type)
	raw_from := msg.From[:]
	raw_created_time := Int64ToBytes(msg.CreatedTime)
	var raw_identify []byte
	if msg.Identify == 0{
		// 如果没有标识符则随机生成一个
		raw_identify = Int32ToBytes(rand.Int31())
	}else{
		raw_identify = Int32ToBytes(int32(msg.Identify))
	}
	raw_chan_route_num := Int32ToBytes(int32(len(msg.NewRoute)))
	// 可以优化速度
	raw = append(raw, raw_type)
	raw = append(raw, raw_from...)
	raw = append(raw, raw_created_time...)
	raw = append(raw, raw_identify...)
	raw = append(raw, raw_chan_route_num...)
	// channel_payload
	if len(msg.NewRoute)!=0{
		var raw_n1, raw_n2, raw_n1_cap, raw_n2_cap []byte
		var raw_status byte
		for _, channel := range msg.NewRoute{
			raw_n1 = channel.N1[:]
			raw_n2 = channel.N2[:]
			raw_n1_cap = Int32ToBytes(int32(channel.N1Cap))
			raw_n2_cap = Int32ToBytes(int32(channel.N2Cap))
			raw_status = byte(channel.ChannelStatus)
			raw = append(raw, raw_n1...)
			raw = append(raw, raw_n2...)
			raw = append(raw, raw_n1_cap...)
			raw = append(raw, raw_n2_cap...)
			raw = append(raw, raw_status)
		}
	}
	return raw
}

func (msg *DiscoverMsg)ConvertFromRaw(raw []byte){
	base := 0
	raw_type := raw[base]
	base += 1
	raw_from := raw[base:base+20]
	base += 20
	raw_created_time := raw[base:base+8]
	base += 8
	raw_identify := raw[base:base+4]
	base += 4
	raw_chan_route_num := raw[base:base+4]
	base += 4
	chan_route_num:=BytesToInt32(raw_chan_route_num)
	// type
	msg.Type = int(raw_type)
	// from
	copy(msg.From[:], raw_from)
	// created_time
	msg.CreatedTime = BytesToInt64(raw_created_time)
	// indentify
	msg.Identify = int(BytesToInt32(raw_identify))
	if chan_route_num!=0{
		msg.NewRoute = make([]Channel, chan_route_num)
		// 开始处理channel
		for i:=0;i<int(chan_route_num);i++{
			raw_n1 := raw[base:base+20]
			base+=20
			raw_n2 := raw[base:base+20]
			base+=20
			raw_n1_cap := raw[base:base+4]
			base+=4
			raw_n2_cap := raw[base:base+4]
			base+=4
			raw_status := raw[base]
			// N1
			copy(msg.NewRoute[i].N1[:], raw_n1)
			// N2
			copy(msg.NewRoute[i].N2[:], raw_n2)
			// N1_Cap
			msg.NewRoute[i].N1Cap = int(BytesToInt32(raw_n1_cap))
			// N2_Cap
			msg.NewRoute[i].N2Cap = int(BytesToInt32(raw_n2_cap))
			// status
			msg.NewRoute[i].ChannelStatus = Status(int(raw_status))
		}
	}
}

