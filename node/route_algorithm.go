package node

import "bytes"

// 判断Channels是否包含这个Channel
func IsChannelsInclude(channels Channels, newChannel Channel) bool{
	for _, channel :=range channels{
		if bytes.Equal(channel.N1[:], newChannel.N1[:])&&bytes.Equal(channel.N2[:], newChannel.N2[:]){
			return false
		}
		if bytes.Equal(channel.N1[:], newChannel.N2[:])&&bytes.Equal(channel.N2[:], newChannel.N1[:]){
			return false
		}
	}
	return true
}
// 从Channels中移除这个Channel
func RemoveChannel(channels Channels, newChannel Channel) Channels{
	// 如果存在，则获得该channel的下标
	var getIndex = func() int{
		for index, channel :=range channels{
			if bytes.Equal(channel.N1[:], newChannel.N1[:])&&bytes.Equal(channel.N2[:], newChannel.N2[:]){
				return index
			}
			if bytes.Equal(channel.N1[:], newChannel.N2[:])&&bytes.Equal(channel.N2[:], newChannel.N1[:]){
				return index
			}
		}
		return -1
	}
	// 获得当前索引
	index := getIndex()
	if index == -1{
		return channels
	}
	// 剔除当前channel
	channels = append(channels[:index], channels[index+1:]...)
	return channels
}
// 获取两个channels的UPD
func ChannelsUPD(oldChannels Channels, newChannels Channels) Channels{
	upds := []Channel{}
	for _, channel := range newChannels{
		upd_channel := channel
		upds = append(upds, upd_channel)
	}
	for _, oldChannel := range oldChannels{
		if IsChannelsInclude(upds, oldChannel){
			// 如果包含，则无需更新
			upds = RemoveChannel(upds, oldChannel)
		}else{
			temp := oldChannel
			temp.ChannelStatus = STATUS_CLOSED
			upds = append(upds, temp)
		}
	}
	return upds
}