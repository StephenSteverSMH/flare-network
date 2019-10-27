package node

import (
	. "../../node"
	"crypto/rand"
	"fmt"
	"testing"
	"time"
)
func createRandChannel() Channel{
	b1 := make([]byte, 20)
	rand.Read(b1)
	n1, _ := ConvertToNodeAddress(b1)

	b2 := make([]byte, 20)
	rand.Read(b2)
	n2, _ := ConvertToNodeAddress(b1)
	return Channel{
		N1:n1,
		N2:n2,
		N1Cap: 3,
		N2Cap: 3,
		CreatedTime: time.Now().Unix(),
		UpdateTime:time.Now().Unix(),
		ChannelStatus: STATUS_ACTIVE,
	}
}

// 判断Channels是否包含这个Channel
func TestIsChannelsInclude(t *testing.T){
	n1 := createRandChannel()
	n2 := createRandChannel()
	n3 := createRandChannel()
	n4 := createRandChannel()
	n5 := createRandChannel()
	chns := Channels{n1, n2, n3, n4, n5}
	// 测试是否包含
	if IsChannelsInclude(chns, n3){
		fmt.Println("IsChannelsInclude测试成功")
	}else{
		fmt.Println("IsChannelsInclude测试失败")
	}
	n3.N1, n3.N2 = n3.N2, n3.N1
	if IsChannelsInclude(chns, n3){
		fmt.Println("IsChannelsInclude测试成功")
	}else{
		fmt.Println("IsChannelsInclude测试失败")
	}
}

// 从Channels中移除这个Channel
func TestRemoveChannel(t *testing.T){
	n1 := createRandChannel()
	n2 := createRandChannel()
	n3 := createRandChannel()
	n4 := createRandChannel()
	n5 := createRandChannel()
	chns := Channels{n1, n2, n3, n4, n5}
	chns = RemoveChannel(chns, n3)
	// 测试是否包含
	if len(chns)==4{
		fmt.Println("RemoveChannel测试成功")
		if IsChannelsInclude(chns, n3){
			fmt.Println("RemoveChannel测试成功")
		}
		fmt.Println("")
	}else{
		fmt.Println("RemoveChannel测试失败")
	}
}

// 获取两个channels的UPD
func TestChannelsUPDs(t *testing.T){
	n1 := createRandChannel()
	n2 := createRandChannel()
	n3 := createRandChannel()
	n4 := createRandChannel()
	n5 := createRandChannel()
	chns_old := Channels{n1, n2, n3, n4, n5}
	chns_new := Channels{n1, n2, n3}
	upd := ChannelsUPD(chns_old, chns_new)
	fmt.Printf("%+v\n", upd)
	chns_new = Channels{n1, n2, n3, n4, n5}
	chns_old = Channels{n1, n2, n3}
	upd = ChannelsUPD(chns_old, chns_new)
	fmt.Printf("%+v\n", upd)
}
