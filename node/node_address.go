package node

import "errors"

const (
	ADDRESS_LEN = 20 // 节点地址长度20字节
)

type NodeAddress [ADDRESS_LEN]byte

// 将[]byte转换为NodeAddress
func ConvertToNodeAddress(b []byte) (NodeAddress,error){
	var n NodeAddress
	if len(b) != 20{
		error := errors.New("输入切片长度不为20")
		return n, error
	}
	copy(n[:], b)
	return n, nil
}