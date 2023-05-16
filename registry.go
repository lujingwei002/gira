package gira

import "github.com/lujingwei002/gira/registry/service/options"

type Registry interface {
	// 如果失败，则返回当前所在的节点
	LockLocalUser(userId string) (*Peer, error)
	UnlockLocalUser(userId string) (*Peer, error)
	WhereIsUser(userId string) (*Peer, error)
	RangePeers(f func(k any, v any) bool)

	RegisterService(serviceName string, opt ...options.RegisterOption) (*Peer, error)
	UnregisterService(serviceName string) (*Peer, error)
	WhereIsService(serviceName string, opt ...options.WhereOption) ([]*Peer, error)
}

// 伙伴
type Peer struct {
	Name     string            // 服务类型
	Id       int32             // 服务id
	FullName string            // 服务全名
	GrpcAddr string            // grpc地址
	Kvs      map[string]string // /server/account_1/ 下的键
}

type LocalPlayer struct {
	UserId    string
	LoginTime int
}

type Service struct {
	// <<GroupName>>/<<ShortName>>
	FullName  string
	GroupName string
	Name      string
	Peer      *Peer
	IsLocal   bool
	IsGroup   bool
}

// 伙伴节点增删改的回调
type PeerWatchHandler interface {
	OnPeerAdd(peer *Peer)
	OnPeerDelete(peer *Peer)
	OnPeerUpdate(peer *Peer)
}

// 玩家位置的回调
type LocalPlayerWatchHandler interface {
	OnLocalPlayerAdd(player *LocalPlayer)
	OnLocalPlayerDelete(player *LocalPlayer)
	OnLocalPlayerUpdate(player *LocalPlayer)
}

type ServiceWatchHandler interface {
	OnServiceAdd(service *Service)
	OnServiceDelete(service *Service)
	OnServiceUpdate(service *Service)
}
