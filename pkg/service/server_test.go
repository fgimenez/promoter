package service_test

import (
	"net"
	"testing"

	"github.com/fgimenez/promoter/pkg/service"
)

func TestServerListens(t *testing.T) {
	port := "3000"

	server := service.NewServer("localhost:" + port)
	if err := server.Start(); err != nil {
		t.Errorf("not able to start server, error %v", err)
		return
	}
	defer server.Stop()

	conn, err := net.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		t.Errorf("not able to connect to 127.0.0.1:%s, error %v", port, err)
		return
	}
	conn.Close()
}
