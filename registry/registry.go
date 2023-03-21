package registry

/// 参考:
/// - https://pkg.go.dev/go.etcd.io/etcd/client/v3#section-readme
/// - https://www.cnblogs.com/hsyw/p/16026627.html

/// watch的版本如果已经compact,则watch会返回一个compacted error,且管道会被关闭

/// 如果etcd server关闭后再恢复，watch goroutine可以正常运行

/// key设置超时的话，如果etcd变成不可用，所有的key就会超时被删除，等etcd变得可用时，需要重新注册全部的key.

/// key不设置超时的话, 非正常关闭, key会一直保留着，但事实上服务已经变得不可用。玩家也会被锁着，无法登录到其他服。如果挂了的服有能力恢复玩家的数据的话，
/// 锁着玩家登录其他服，就可以防止玩家的数据丢失。这种情况下可以增加一个心跳功能，使服务被选中的优先级降低。

///

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lujingwei/gira"
	"github.com/lujingwei/gira/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

const GRPC_KEY string = "grpc"

type Registry struct {
	Config         gira.EtcdConfig
	FullName       string // 节点全名
	Name           string // 节点名
	Client         *clientv3.Client
	cancelCtx      context.Context
	cancelFunc     context.CancelFunc
	facade         gira.ApplicationFacade
	PeerRegistry   *PeerRegistry
	PlayerRegistry *PlayerRegistry
}

func (r *Registry) Notify() error {
	if err := r.PeerRegistry.notify(r); err != nil {
		return err
	}
	if err := r.PlayerRegistry.notify(r); err != nil {
		return err
	}
	return nil
}

func (r *Registry) GetPeer(fullName string) *gira.Peer {
	return r.PeerRegistry.getPeer(r, fullName)
}

func NewConfigRegistry(config *gira.EtcdConfig, facade gira.ApplicationFacade) (*Registry, error) {
	r := &Registry{
		Config: *config,
	}
	r.cancelCtx, r.cancelFunc = context.WithCancel(context.Background())
	r.FullName = facade.GetFullName()
	r.Name = facade.GetName()
	r.facade = facade

	// 配置endpoints
	endpoints := make([]string, 0)
	for _, v := range r.Config.Endpoints {
		endpoints = append(endpoints, fmt.Sprintf("http://%s:%d", v.Host, v.Port))
	}
	c := clientv3.Config{
		Endpoints:   endpoints,                                         // 节点信息
		DialTimeout: time.Duration(r.Config.DialTimeout) * time.Second, // 超时时间
		DialOptions: []grpc.DialOption{grpc.WithBlock()},               // 使用阻塞模式，确认启动时etcd是可用的
		Username:    r.Config.Username,
		Password:    r.Config.Password,
		Context:     r.cancelCtx,
	}
	var client *clientv3.Client
	var err error

	// 建立连接
	if client, err = clientv3.New(c); err != nil {
		log.Infof("connect to etcd failed, err:%v\n", err)
		return nil, err
	}
	r.Client = client
	log.Info("connect to etcd success")
	if v, err := newConfigPeerRegistry(r); err != nil {
		return nil, err
	} else {
		r.PeerRegistry = v
	}
	if v, err := newConfigPlayerRegistry(r); err != nil {
		return nil, err
	} else {
		r.PlayerRegistry = v
	}
	return r, nil
}

func (r *Registry) LockLocalMember(memberId string) (*gira.Peer, error) {
	return r.PlayerRegistry.LockLocalMember(r, memberId)
}
func (r *Registry) UnlockLocalMember(memberId string) (*gira.Peer, error) {
	return r.PlayerRegistry.UnlockLocalMember(r, memberId)
}
func explodeServerFullName(fullName string) (name string, id int, err error) {
	pats := strings.Split(string(fullName), "_")
	if len(pats) != 4 {
		err = gira.ErrInvalidPeer
		return
	}
	name = pats[0]
	id, err = strconv.Atoi(pats[3])
	return
}
