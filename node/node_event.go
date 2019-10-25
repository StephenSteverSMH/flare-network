package node

import "net"

const (
	INIT int = iota // 打开通道事件
	ACT	// 某个邻居节点从Close -> ACTIVE
	RT  // 自己的路由表发生变化
	ACK // 接受了从其他节点发来的ACK
	RST // 接收从其他节点发来的NEIGHBOR_RST
	TO // 准备发送NEIGHBOR_UPD或NEIGHTBOR_HELLO给对端节点
	NB_UP // 接受了从NB节点发来的路由表
	CHAN_DOWN // 关闭通道
)

type Event struct{
	Type int
	Data interface{}
}

// 建立通道通知
type InitNotify struct {
	To NodeAddress
	N1Cap int
	N2Cap int
	ToNetAddr net.TCPAddr
}

// 摧毁邻居通道通知
type DestoryChannelNotify struct {
	To NodeAddress
	ToNetAddr net.TCPAddr
}