package tcp

import (
	"testing"
	"context"
	"time"
	"bytes"
	"math/rand"
	"net"
	"errors"
)
// 注意：运行以下测试之前先启动服务端 go run examples/server.go
// 测试连续连接和关闭连接10000次，观看服务器和客户端是否正常
// go test -v -test.run TestNewClient
func TestNewClient(t *testing.T) {
	address := "127.0.0.1:7771"
	go func() {
		dial := net.Dialer{Timeout: time.Second * 3}
		conn, err := dial.Dial("tcp", address)
		if err == nil {
			for {
				// 这里发送一堆干扰数据包
				// 这里报文没有按照规范进行封包
				// 目的是为了测试服务端的解包容错性
				conn.Write([]byte(RandString(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(128))))
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()
	var (
		res  []byte
	 	data []byte
	 	client *Client
		times = 10000
		err error
	)
	for  i := 0; i < times; i++ {
		client, err = NewClient(
			context.Background(),
			address,
			SetClientConnectTimeout(time.Second * 3),
			SetWaiterTimeout(1000 * 60),
		)
		if err != nil {
			t.Errorf("NewClient error")
			return
		}
		err = nil
		for {
			data = []byte(RandString(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(128)))
			if len(data) <= 0 {
				break
			}
			wai, _, err := client.Send(data, 0)
			if err != nil {
				t.Errorf("Send fail, err=[%v]", err)
				break
			}
			res, _, err = wai.Wait(0)
			if err != nil {
				t.Errorf("Wait fail, err=[%v]", err)
				break
			}
			if !bytes.Equal(data, res) {
				t.Errorf("send != return")
				err = errors.New("send != return")
				break
			}
			break
		}
		client.Close()
		if err != nil {
			t.Errorf(err.Error())
			return
		}
	}
}

// 注意：运行以下测试之前先启动服务端 go run examples/server.go
// go test -v -test.run TestNewClient2
func TestNewClient2(t *testing.T) {
	address := "127.0.0.1:7771"
	go func() {
		dial := net.Dialer{Timeout: time.Second * 3}
		conn, err := dial.Dial("tcp", address)
		if err == nil {
			for {
				// 这里发送一堆干扰数据包
				// 这里报文没有按照规范进行封包
				// 目的是为了测试服务端的解包容错性
				conn.Write([]byte(RandString(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(128))))
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()
	client, err := NewClient(
		context.Background(),
		address,
		SetClientConnectTimeout(time.Second * 3),
		SetWaiterTimeout(1000 * 60),
	)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	defer client.Close()

	var (
		times = 10000
		res []byte
		data []byte
	)

	for  i := 0; i < times; i++ {
		data = []byte(RandString(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1024)))
		if len(data) <= 0 {
			continue
		}
		wai, _, err := client.Send(data, 0)
		if err != nil {
			t.Errorf("Send fail, err=[%v]", err)
			return
		}
		res, _, err = wai.Wait(time.Second * 3)
		if err != nil {
			t.Errorf("Wait fail, err=[%v]", err)
			return
		}
		if !bytes.Equal(data, res) {
			t.Errorf("Equal fail, send != return")
			return
		}
	}
}