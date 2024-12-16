package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/injoyai/conv"
	"github.com/injoyai/io/dial"
	"net"
	"time"
)

// DiscoverRequest 服务发现包
type DiscoverRequest struct {
	System string //系统,网关
	Key    string //插件key
}

type DiscoverResponse struct {
	System  string //系统,网关
	Key     string //插件key
	Address string //服务/插件地址
}

// Register 插件注册
type Register struct {
	Key       string
	Name      string
	Version   string
	Functions map[string]Function
}

type Response struct {
	Code int
	Data []byte
	Msg  string
}

type DiscoverConfig struct {
	Interval time.Duration
	Start    net.IP
	End      net.IP
	Port     []int
	Timeout  time.Duration
}

func (this *DiscoverConfig) broadcast(ctx context.Context, system, key string) (addressList []*DiscoverResponse) {
	for i := conv.Uint32([]byte(this.Start.To4())); i <= conv.Uint32([]byte(this.End.To4())); i++ {
		for _, port := range this.Port {
			ip := net.IP(conv.Bytes(i))
			addr := fmt.Sprintf("%s:%d", ip.String(), port)
			c, err := dial.NewUDPTimeout(addr, time.Millisecond*200)
			if err == nil {
				respBytes, err := c.WriteRead(conv.Bytes(DiscoverRequest{
					System: system,
					Key:    key,
				}), conv.SelectDuration(this.Timeout > 0, this.Timeout, time.Millisecond*200))
				c.Close()
				if err == nil {
					resp := new(Response)
					if err = json.Unmarshal(respBytes, resp); err != nil {
						continue
					}
					if resp.Code != 200 {
						continue
					}
					data := &DiscoverResponse{}
					if err = json.Unmarshal(resp.Data, data); err != nil {
						continue
					}
					data.Address = addr
					addressList = append(addressList, data)
				}
			}
		}
	}
	return
}
