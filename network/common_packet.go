package network

import . "../utils"

const (
	PACKET_NEIGHBOR = iota
)

type Packet struct{
	Type int
	Size int
	Payload []byte
}
// 数据包格式 type(1B)|size(4B)|payload

func (packet *Packet)ConvertToRaw() []byte{
	raw := []byte{}
	raw_type := byte(packet.Type)
	raw_size := Int32ToBytes(int32(packet.Size))
	raw  = append(raw, raw_type)
	raw  = append(raw, raw_size...)
	raw  = append(raw, packet.Payload...)
	return raw
}

func (packet *Packet)ConvertFromRaw(raw []byte){
	base := 0
	raw_type := raw[base]
	base = base + 1
	raw_size := raw[base:base + 4]
	base = base + 4
	size := int(BytesToInt32(raw_size))

	packet.Type = int(raw_type)
	packet.Size = size
	packet.Payload = raw[base: base + size]
}