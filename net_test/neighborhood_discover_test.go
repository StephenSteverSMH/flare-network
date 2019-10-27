package net_test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)
import "../server"
func CreatePipeChannel(serverA server.Server, serverB server.Server, capA int, capB int){
	serverA.CreateChannel(serverB.GetNode().Address, *serverB.GetNode().NetAddr, capA, capB)
	serverB.CreateChannel(serverA.GetNode().Address, *serverA.GetNode().NetAddr, capB, capA)
}
func DestoryPipeChannle(serverA server.Server, serverB server.Server){
	serverA.DestoryChannel(serverB.GetNode().Address)
	serverB.DestoryChannel(serverA.GetNode().Address)
}
//测试网络01
//func TestNet01(t *testing.T){
//	serverA := server.Server{}
//	serverB := server.Server{}
//	serverC := server.Server{}
//	//serverD := server.Server{}
//	wg := sync.WaitGroup{}
//	wg.Add(3)
//	go serverA.StartNodeServer("127.0.0.1", "10005")
//	go serverB.StartNodeServer("127.0.0.1", "10006")
//	go serverC.StartNodeServer("127.0.0.1", "10007")
//	//go serverD.StartNodeServer("127.0.0.1", "10008")
//	time.Sleep(time.Second)
//	// 打印地址
//	fmt.Println("A的节点地址")
//	fmt.Println(serverA.GetNode().Address)
//	fmt.Println("B的节点地址")
//	fmt.Println(serverB.GetNode().Address)
//	fmt.Println("C的节点地址")
//	fmt.Println(serverC.GetNode().Address)
//	fmt.Println("------------------------------------------------------")
//	CreatePipeChannel(serverA, serverB, 3, 4)
//	time.Sleep(time.Second)
//	CreatePipeChannel(serverA, serverC, 3, 4)
//	time.Sleep(time.Second)
//	// 断开A和C的连接
//	DestoryPipeChannle(serverA, serverC)
//	time.Sleep(time.Second)
//	//fmt.Println(serverA.GetNode().Route.NeighborhoodMap)
//	fmt.Println("------------------------------------------------------")
//	fmt.Println("ServerA的邻居")
//	for _,node := range serverA.GetNode().AdjcentNodes{
//		fmt.Println(node.Address)
//	}
//	fmt.Println("ServerA的路由为")
//	fmt.Println(serverA.GetNode().Route.NeighborhoodMap)
//	fmt.Println("ServerB的邻居")
//	for _,node := range serverB.GetNode().AdjcentNodes{
//		fmt.Println(node.Address)
//	}
//	fmt.Println("ServerB的路由为")
//	fmt.Println(serverB.GetNode().Route.NeighborhoodMap)
//	fmt.Println("ServerC的邻居")
//	for _,node := range serverC.GetNode().AdjcentNodes{
//		fmt.Println(node.Address)
//	}
//	fmt.Println("ServerC的路由为")
//	fmt.Println(serverC.GetNode().Route.NeighborhoodMap)
//	wg.Wait()
//}
// 测试网络02
func TestNet02(t *testing.T){
	serverA := server.Server{}
	serverB := server.Server{}
	serverC := server.Server{}
	serverD := server.Server{}
	serverE := server.Server{}
	serverF := server.Server{}
	go serverA.StartNodeServer("127.0.0.1", "10005")
	go serverB.StartNodeServer("127.0.0.1", "10006")
	go serverC.StartNodeServer("127.0.0.1", "10007")
	go serverD.StartNodeServer("127.0.0.1", "10008")
	go serverE.StartNodeServer("127.0.0.1", "10009")
	go serverF.StartNodeServer("127.0.0.1", "10010")
	time.Sleep(time.Second)
	CreatePipeChannel(serverA, serverC, 3, 4)
	time.Sleep(1*time.Second)
	CreatePipeChannel(serverC, serverB, 3, 4)
	time.Sleep(1*time.Second)
	CreatePipeChannel(serverB, serverF, 3, 4)
	time.Sleep(1*time.Second)
	CreatePipeChannel(serverC, serverD, 3, 4)
	time.Sleep(1*time.Second)
	CreatePipeChannel(serverD, serverE, 3, 4)
	time.Sleep(1*time.Second)
	DestoryPipeChannle(serverB, serverF)
	time.Sleep(1*time.Second)
	wg := sync.WaitGroup{}
	wg.Add(3)
	// 打印地址
	fmt.Println("A的节点地址")
	fmt.Println(serverA.GetNode().Address)
	fmt.Println("B的节点地址")
	fmt.Println(serverB.GetNode().Address)
	fmt.Println("C的节点地址")
	fmt.Println(serverC.GetNode().Address)
	fmt.Println("D的节点地址")
	fmt.Println(serverD.GetNode().Address)
	fmt.Println("E的节点地址")
	fmt.Println(serverE.GetNode().Address)
	fmt.Println("F的节点地址")
	fmt.Println(serverF.GetNode().Address)
	fmt.Println("------------------------------------------------------")
	fmt.Println("ServerA的邻居")
	for _,node := range serverA.GetNode().AdjcentNodes{
		fmt.Println(node.Address)
	}
	fmt.Println("ServerA的路由为")
	fmt.Println(serverA.GetNode().Route.NeighborhoodMap)
	fmt.Println("ServerE的邻居")
	for _,node := range serverE.GetNode().AdjcentNodes{
		fmt.Println(node.Address)
	}
	fmt.Println("ServerE的路由为")
	fmt.Println(serverE.GetNode().Route.NeighborhoodMap)
	wg.Wait()
}
