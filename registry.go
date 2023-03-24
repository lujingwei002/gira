package gira

// 伙伴
type Peer struct {
	Name     string            // 服务类型
	Id       int32             // 服务id
	FullName string            // 服务全名
	GrpcAddr string            // grpc地址
	Kvs      map[string]string // /server/account_1/ 下的键
}

type LocalPlayer struct {
	MemberId  string
	LoginTime int
}

// 伙伴节点增删改的回调
type PeerHandler interface {
	OnPeerAdd(peer *Peer)
	OnPeerDelete(peer *Peer)
	OnPeerUpdate(peer *Peer)
}

type LocalPlayerHandler interface {
	OnLocalPlayerAdd(player *LocalPlayer)
	OnLocalPlayerDelete(player *LocalPlayer)
	OnLocalPlayerUpdate(player *LocalPlayer)
}
