package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/cache"
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/goutil/oss/linux/systemctl"
	"github.com/injoyai/io"
	"github.com/injoyai/io/listen"
	"net/url"
	"path/filepath"
	"time"
)

type Server interface {

	// Discover 发现插件,通过UDP广播,插件获取到广播信息后,开始注册插件
	Discover(ctx context.Context, req *DiscoverConfig) ([]*DiscoverResponse, error)

	// All 全部插件
	All(ctx context.Context) ([]*Plugin, error)

	// Install 根据插件地址安装插件
	Install(ctx context.Context, url string) error

	// Restart 重新启动插件
	// 让插件重新启动,重新注册,例如升级版本
	Restart(ctx context.Context, key string) error

	// Disable 禁用插件,停止插件
	Disable(ctx context.Context, key string) error

	// Enable 启用插件
	Enable(ctx context.Context, key string) error

	// Remove 移除插件,
	// todo 是否允许重新注册? 如果可以相当于重新加载,如果不可以之后怎么重新加载
	Remove(ctx context.Context, key string) error

	// Execute 执行插件函数
	// 需要参数 插件key,插件的函数名称,和函数需要的参数
	// 返回响应的业务数据,和通讯错误信息
	Execute(ctx context.Context, key string, method string, args map[string]interface{}) (*Response, error)
}

func NewServer(dir, system string, port int) (Server, error) {
	s, err := listen.NewTCPServer(port)
	if err != nil {
		return nil, err
	}
	return &server{
		Server: s,
		dir:    dir,
		system: system,
		cache:  cache.NewFile("plugin", system),
	}, nil
}

type server struct {
	*io.Server             //插件服务
	dir        string      //插件目录,例 /root/xxx/data/plugin
	system     string      //插件系统名称
	cache      *cache.File //插件缓存
}

func (this *server) Discover(ctx context.Context, cfg *DiscoverConfig) ([]*DiscoverResponse, error) {
	return cfg.broadcast(ctx, this.system, ""), nil
}

func (this *server) All(ctx context.Context) (list []*Plugin, err error) {
	this.cache.Range(func(key interface{}, value interface{}) bool {
		if val, ok := value.(*Plugin); ok {
			list = append(list, val)
		} else {
			l := &Plugin{}
			if err = json.Unmarshal(conv.Bytes(value), l); err != nil {
				return false
			}
			list = append(list, l)
		}
		return true
	})
	return
}

func (this *server) Install(ctx context.Context, u string) error {
	u2, err := url.Parse(u)
	if err != nil {
		return err
	}
	key := filepath.Base(u2.Path)
	dir := filepath.Join(this.dir, this.system)
	filename := filepath.Join(dir, key)
	if _, err := http.GetToFile(u, filename); err != nil {
		return err
	}
	if err := systemctl.Install(key, dir); err != nil {
		return err
	}
	this.cache.Set(key, &Plugin{Key: key})
	return this.cache.Cover()
}

func (this *server) Restart(ctx context.Context, key string) error {
	return systemctl.Restart(key)
}

func (this *server) Disable(ctx context.Context, key string) error {
	if err := systemctl.Disable(key); err != nil {
		return err
	}
	return systemctl.Stop(key)
}

func (this *server) Enable(ctx context.Context, key string) error {
	if err := systemctl.Enable(key); err != nil {
		return err
	}
	return systemctl.Start(key)
}

func (this *server) Remove(ctx context.Context, key string) error {
	this.cache.Del(key)
	return this.cache.Cover()
}

func (this *server) Execute(ctx context.Context, key string, function string, args map[string]interface{}) (*Response, error) {
	c := this.GetClient(key)
	if c == nil {
		return nil, fmt.Errorf("插件(%s)不存在", key)
	}
	respBytes, err := c.WriteRead(conv.Bytes(&ExecuteRequest{
		Function: function,
		Args:     args,
	}), time.Second*2)
	if err != nil {
		return nil, err
	}
	resp := new(Response)
	err = json.Unmarshal(respBytes, resp)
	return resp, err
}

type ExecuteRequest struct {
	Function string
	Args     map[string]interface{}
}
