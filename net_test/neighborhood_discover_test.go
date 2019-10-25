package net_test

//func TestConvertRaw(t *testing.T){
//	b := make([]byte, 20)
//	n1 := make([]byte, 20)
//	n2 := make([]byte, 20)
//	rand.Read(b)
//	rand.Read(n1)
//	rand.Read(n2)
//	msg := node.DiscoverMsg{
//		Type:node.NEIGHBOR_HELLO,
//		CreatedTime: time.Now().Unix(),
//		Identify: 12,
//		NewRoute: []node.Channel{},
//	}
//	copy(msg.From[:], b)
//	channel1 := node.Channel{
//		N1Cap: 2,
//		N2Cap: 3,
//		ChannelStatus: node.STATUS_ACTIVE,
//	}
//	copy(channel1.N1[:], n1)
//	copy(channel1.N2[:], n2)
//	msg.NewRoute = append(msg.NewRoute, channel1)
//	raw := msg.ConvertToRaw()
//	fmt.Println(len(raw))
//	fmt.Println(raw)
//	new_msg := node.DiscoverMsg{}
//	new_msg.ConvertFromRaw(raw)
//	fmt.Println(fmt.Sprintf("%+v", msg))
//	fmt.Println(fmt.Sprintf("%+v", new_msg))
//}
