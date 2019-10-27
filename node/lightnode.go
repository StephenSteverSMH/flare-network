package node

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"crypto/rand"
	"../network"
	"testing"
	"time"
	mrand "math/rand"
)

// 通道状态
type Status int
const (
	STATUS_ACTIVE Status = 0
	STATUS_CLOSED Status = 1
)
const (
	CONN_READ_TIMEOUT = 5
)

// 通道
type Channel struct {
	// n1 地址
	N1 NodeAddress
	// n2 地址
	N2 NodeAddress
	// n1 容量
	N1Cap int
	// n2 容量
	N2Cap int
	// 通道状态
	ChannelStatus Status
	// 通道创建时间撮
	CreatedTime int64
	// 通道更新时间撮
	UpdateTime int64
}
// 通道图
type Channels []Channel
// 路由表
type RouteTable struct{
	NeighborhoodMap Channels
	NeighborRadius int
	mutex sync.Mutex
}
// 闪电网络节点
type LightNode struct {
	// 网络地址
	NetAddr *net.TCPAddr
	// 节点地址
	Address NodeAddress
	// 路由表
	Route RouteTable
	// 毗邻节点
	AdjcentNodes map[NodeAddress]LightNode
	// 是否在等待毗邻节点ACK
	WaitACKMap map[NodeAddress]int
	// 毗邻节点同步状态
	SyncNBMap map[NodeAddress][]Channels
	// 事件通道
	EventChannel chan Event
	// 事件通道退出信号
	Closed chan struct{}
}

func CreateNode(radius int, event_chan_len int) LightNode{
	// 生成自己的网络地址，先模拟下
	b := make([]byte, 20)
	rand.Read(b)
	address, _ := ConvertToNodeAddress(b)
	l := LightNode{
		Route: RouteTable{
			NeighborhoodMap: []Channel{},
			NeighborRadius: radius,
		},
		AdjcentNodes:make(map[NodeAddress]LightNode),
		WaitACKMap: make(map[NodeAddress]int),
		EventChannel: make(chan Event, event_chan_len),
		SyncNBMap: make(map[NodeAddress][]Channels),
		Address: address,
	}
	return l
}

//// 处理NEIGHBOR_HELLO和NEIGHBOR_UPD消息，并改动路由表
func (node *LightNode)processRoute(msg DiscoverMsg){
	if msg.Type != NEIGHBOR_HELLO && msg.Type != NEIGHBOR_UPD{
		// 返回错误
		fmt.Println(node.Address[:4], "processRoute中不该处理的消息，Type为", msg.Type)
		return
	}
	for _, channel := range msg.NewRoute{
		// 如果该通道还未出现
		//fmt.Println(node.Address[:4], "此时的路由",node.Route.NeighborhoodMap)
		//fmt.Println(node.Address[:4], channel, IsChannelsInclude(node.Route.NeighborhoodMap, channel))
		if !IsChannelsInclude(node.Route.NeighborhoodMap, channel){
			// 如果该channel距离node的跳数小于node.Route.NeighborRadius
			// 先用true表示
			if true{
				// 确认该通道真实存在于区块链
				// VerifyChannel()
				// 将该通道添加到路由表中
				node.Route.mutex.Lock()
				node.Route.NeighborhoodMap = append([]Channel(node.Route.NeighborhoodMap), channel)
				node.Route.mutex.Unlock()
			}
		} else{
			// 如果已经出现, 判断是否是通道关闭消息
			if channel.ChannelStatus == STATUS_CLOSED{
				// 如果是通道关闭消息，则将该通道从路由表中移除
				node.Route.mutex.Lock()
				node.Route.NeighborhoodMap = RemoveChannel(node.Route.NeighborhoodMap, channel)
				node.Route.mutex.Unlock()
			}
		}
	}

}
// 处理节点状态变化
func (node *LightNode)ProcessEvent(){
	LOOP:
	for{
		select {
			case event := <-node.EventChannel:
				// 根据不同事件类型操作
				switch event.Type {
				case INIT:
					// 通道创建
					// 如果通道已经存在，通知创建失败
					// hook
					fmt.Println(node.Address[:4],"节点收到INIT事件")
					msg := event.Data.(InitNotify)
					to := msg.To
					channel := Channel{
						N1: node.Address,
						N2: to,
						N1Cap: msg.N1Cap,
						N2Cap: msg.N2Cap,
						ChannelStatus:STATUS_ACTIVE,
						CreatedTime: time.Now().Unix(),
						UpdateTime: time.Now().Unix(),
					}
					if IsChannelsInclude(node.Route.NeighborhoodMap, channel){
						// 如果路由已经存在
						fmt.Println(node.Address[:4], "该路由已经存在，不能创建")
						continue
					}
					if _, ok := node.AdjcentNodes[to];ok{
						// 如果该邻居节点已经存在
					}else{
						// 如果该邻居节点在节点表中不存在
						l := LightNode{
							NetAddr:&msg.ToNetAddr,
							Address: to,
						}
						node.AdjcentNodes[to] = l
					}
					// 将通道加入到路由表
					node.Route.mutex.Lock()
					node.Route.NeighborhoodMap = append(node.Route.NeighborhoodMap, channel)
					node.Route.mutex.Unlock()
					// 设置该节点ACK同步状态
					node.SyncNBMap[to] = []Channels{[]Channel{},[]Channel{}}
					// 设置该节点ACK标识符
					node.WaitACKMap[to] = -1
					// 产生UPD事件
					node.EventChannel <- Event{Type:RT, Data:msg}
					break
				case ACT:
					// 邻居节点打开
					break
				case RT:
					// 自己路由表发生变化
					// 向所有邻居节点发送NEIGHBOR_UPD包
					switch msg:= event.Data.(type) {
					case DiscoverMsg:
						// 因为收到HELLO或UPD导致路由发生变化的情况
						from := msg.From
						node.sendNB_UPD(from)
					case InitNotify:
						// 建立通道的通知
						// 不向建立通道的一方发送UPD
						fmt.Println(node.Address[:4],"产生INIT_HELLO事件")
						to := msg.To
						//none_arr := [20]byte{}
						node.sendNB_HELLO(to)
						break
					case DestoryChannelNotify:
						// 因为解除邻居通道的通知
						// to := msg.To
						none_arr := [20]byte{}
						node.sendNB_UPD(none_arr)
						break
					}
					break
				case ACK:
					// 收到确认
					msg := event.Data.(DiscoverMsg)
					// 检验是否有个这个邻居
					if _, ok := node.AdjcentNodes[msg.From];!ok{
						fmt.Println(node.Address[:4],"没有这个邻居，不需要它的ACK")
						break
					}
					// 检验消息是否在等待ack
					if node.WaitACKMap[msg.From] == -1{
						fmt.Println(node.Address[:4],"突如其来的ACK")
						break
					}
					fmt.Println(node.Address[:4],"收到",msg.From[:4],"的ACK包")
					// 检验消息是否是需要的那个
					if node.WaitACKMap[msg.From]!=msg.Identify{
						fmt.Println(node.Address[:4],"并非需要的ACK")
					}else{
						// 完成ACK
						node.WaitACKMap[msg.From] = -1
						// 更新镜像
						node.SyncNBMap[msg.From][0] = node.SyncNBMap[msg.From][1]
					}
					break
				case RST:
					break
				case TO:
					break
				case NB_UP:
					// 接收道NEIGHBOR_HELLO和UPD
					//fmt.Println("产生NB_UP事件")
					discoverMsg := event.Data.(DiscoverMsg)
					if discoverMsg.Type==NEIGHBOR_UPD{
						fmt.Println(node.Address[:4],"收到UPD包，其upds为",discoverMsg.NewRoute)
					}
					if discoverMsg.Type==NEIGHBOR_HELLO{
						fmt.Println(node.Address[:4],"收到HELLO包，其route为",discoverMsg.NewRoute)
					}
					if _, ok := node.AdjcentNodes[discoverMsg.From];!ok{
						fmt.Println(node.Address[:4], "收到非邻居发送的HELLO或UPD包")
					}
					if len(discoverMsg.NewRoute)!=0{
						// 更新路由
						node.processRoute(discoverMsg)
						// 发起路由更新事件
						node.EventChannel <- Event{
							Type: RT,
							Data: discoverMsg,
						}
					}
					// 向事件源发回ACK包
					ack_msg := DiscoverMsg{Identify:event.Data.(DiscoverMsg).Identify, Type:NEIGHBOR_ACK, From:node.Address, To:discoverMsg.From, CreatedTime:time.Now().Unix(), NewRoute:[]Channel{}}
					if _, ok :=node.AdjcentNodes[discoverMsg.From];!ok{
						fmt.Println(node.Address[:4],"收到非邻居的UDP包")
						continue
					}
					conn, err :=net.Dial("tcp", node.AdjcentNodes[discoverMsg.From].NetAddr.String())
					if err!=nil{
						// 打印连接邻居节点失败
						fmt.Println(node.Address[:4],"发送ACK失败", err)
						return
					}
					raw_ack_msg := ack_msg.ConvertToRaw()
					packet := network.Packet{Type:network.PACKET_NEIGHBOR, Size:len(raw_ack_msg), Payload:raw_ack_msg}
					raw_packet := packet.ConvertToRaw()
					conn.Write(raw_packet)
					conn.Close()
					break
				case CHAN_DOWN:
					// 邻居节点通道关闭
					destroy_msg := event.Data.(DestoryChannelNotify)
					// 清除邻居节点
					if _, ok := node.AdjcentNodes[destroy_msg.To];!ok{
						// 不存在该邻居
						fmt.Println(node.Address[:4],"通道和邻居不存在，无法关闭")
						continue
					}else{
						temp_channel := Channel{N1: node.Address, N2: destroy_msg.To, N1Cap: 1, N2Cap: 1}
						// 清理邻居通道
						node.Route.mutex.Lock()
						if IsChannelsInclude(node.Route.NeighborhoodMap, temp_channel){
							node.Route.NeighborhoodMap = RemoveChannel(node.Route.NeighborhoodMap, temp_channel)
						}
						node.Route.mutex.Unlock()
						// 清除邻居节点
						delete(node.AdjcentNodes, destroy_msg.To)
						delete(node.SyncNBMap, destroy_msg.To)
						delete(node.WaitACKMap, destroy_msg.To)
						// 产生RT事件
						fmt.Println(node.Address[:4], "结束通道，产生RT事件")
						msg:=DiscoverMsg{
							From: destroy_msg.To,
						}
						event := Event{
							Type: RT,
							Data: msg,
						}
						node.EventChannel <- event
					}

					break
				default:
				}
			case <-node.Closed:
				// 通知关闭
				break LOOP
		}
	}
}

func (node *LightNode)sendNB_UPD(exclude NodeAddress) {
	for _, adjcent := range node.AdjcentNodes {
		// 不向事件源发回更新包
		if bytes.Equal(adjcent.Address[:], exclude[:]) {
			fmt.Println(node.Address[:4],"阻拦UPD的的发送地址",adjcent.Address[:4])
			continue
		}
		// 如果该邻居还在等待该邻居的ACK，跳过
		if node.WaitACKMap[adjcent.Address] != -1 {
			// 设置ACK超时时间，优化
			continue
		}
		if bytes.Equal(adjcent.Address[:], exclude[:]) {
			fmt.Println(node.Address[:4],"进不来这里")
			continue
		}
		// 发起连接，并发送NEIGHBOR_UPD包
		go func(adjcent LightNode) {
			// 计算UPD
			upds := ChannelsUPD(node.SyncNBMap[adjcent.Address][0], node.Route.NeighborhoodMap)
			fmt.Println(node.Address[:4],"准备发往UPD的地址为", adjcent.Address[:4], "，upds为",upds)
			identify := int(mrand.Int31())
			msg := DiscoverMsg{
				Type:        NEIGHBOR_UPD,
				From:        node.Address,
				To:          adjcent.Address,
				NewRoute:    upds,
				CreatedTime: time.Now().Unix(),
				Identify:    identify,
			}
			raw_msg := msg.ConvertToRaw()
			conn, err := net.Dial("tcp", adjcent.NetAddr.String())
			if err != nil {
				// 打印连接邻居节点失败
				fmt.Println(node.NetAddr, adjcent.NetAddr)
				fmt.Println(node.Address[:4],"连接邻居节点失败", err)
				return
			}
			// 默认读超时5s
			conn.SetReadDeadline(time.Now().Add(CONN_READ_TIMEOUT * time.Second))
			// 发送PACKET_NEIGHBOR包
			packet := network.Packet{Type: network.PACKET_NEIGHBOR, Size: len(raw_msg), Payload: raw_msg}
			raw_packet := packet.ConvertToRaw()
			conn.Write(raw_packet)
			// 存储ACK标识符
			node.WaitACKMap[adjcent.Address] = identify
			// 更新镜像
			node.SyncNBMap[adjcent.Address][1] = []Channel{}
			node.SyncNBMap[adjcent.Address][1] = append(node.SyncNBMap[adjcent.Address][1], node.Route.NeighborhoodMap...)
		}(adjcent)
	}
}
func (node *LightNode)sendNB_HELLO(to NodeAddress) {

	fmt.Println(node.Address[:4],"准备发送HELLO的地址为", to)
	// 发起连接，并发送NEIGHBOR_UPD包
	go func() {
		identify := int(mrand.Int31())
		msg := DiscoverMsg{
			Type:        NEIGHBOR_HELLO,
			From:        node.Address,
			To:          to,
			NewRoute:    node.Route.NeighborhoodMap,
			CreatedTime: time.Now().Unix(),
			Identify:    identify,
		}
		raw_msg := msg.ConvertToRaw()
		conn, err := net.Dial("tcp", node.AdjcentNodes[to].NetAddr.String())
		if err != nil {
			// 打印连接邻居节点失败
			fmt.Println(node.Address[:4],"连接邻居节点失败", err)
			return
		}
		// 默认读超时5s
		conn.SetReadDeadline(time.Now().Add(CONN_READ_TIMEOUT * time.Second))
		// 发送PACKET_NEIGHBOR包
		packet := network.Packet{Type: network.PACKET_NEIGHBOR, Size: len(raw_msg), Payload: raw_msg}
		raw_packet := packet.ConvertToRaw()
		conn.Write(raw_packet)
		// 存储ACK标识符
		node.WaitACKMap[node.AdjcentNodes[to].Address] = identify
		// 更新镜像
		node.SyncNBMap[node.AdjcentNodes[to].Address][1] = []Channel{}
		node.SyncNBMap[node.AdjcentNodes[to].Address][1] = append(node.SyncNBMap[to][1], node.Route.NeighborhoodMap...)
	}()

}

// 获取两个channels的UPD
func TestChannelsUPD(t *testing.T){

}