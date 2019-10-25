package server

import (
	"testing"
)

func TestServer(t *testing.T){
	server := Server{}
	server.StartNodeServer("127.0.0.1", "8001")
}
