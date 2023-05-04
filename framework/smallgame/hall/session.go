package hall

import (
	"context"
	"errors"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/actor"
	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/framework/smallgame/common/rpc"
	"github.com/lujingwei002/gira/framework/smallgame/gen/grpc/hall_grpc"
	"github.com/lujingwei002/gira/log"
)

// 锁是和session绑定的，因此由session来抢占会释放

// session 被动从stream接收消息
type hall_sesssion struct {
	*actor.Actor
	hall             *hall_server
	ctx              context.Context
	cancelFunc       context.CancelFunc
	sessionId        uint64
	memberId         string
	userId           string
	avatar           UserAvatar
	player           Player
	chClientResponse chan *hall_grpc.ClientMessageResponse // 由stream负责关闭
	chClientRequest  chan *hall_grpc.ClientMessageRequest  // 由stream负责关闭
	chPeerPush       chan gira.ProtoPush                   // 其他节点，或者自己节点转发来的的push消息，由session负责关闭
	clientCancelFunc context.CancelFunc
	isClosed         int32
}

func newSession(hall *hall_server, sessionId uint64, memberId string) (session *hall_sesssion, err error) {
	var userId string
	var avatar UserAvatar
	ctx := hall.ctx
	avatar, err = hall.hallHandler.NewUser(ctx, memberId)
	if err != nil {
		return
	}
	userId = avatar.GetUserId().Hex()
	// 抢占锁
	var peer *gira.Peer
	peer, err = facade.LockLocalUser(userId)
	log.Infow("lock local user return", "session_id", sessionId, "peer", peer, "err", err)
	if err != nil {
		if peer == nil {
			return nil, err
		} else {
			// 顶号下线
			log.Infow("user instead", "session_id", userId)
			var resp *hall_grpc.UserInsteadResponse
			if resp, err = rpc.Hall.UserInstead(ctx, peer, userId); err != nil {
				log.Infow("user instead fail", "session_id", sessionId, "error", err)
				return
			} else if resp.ErrorCode != 0 {
				err = gira.NewError(resp.ErrorCode, resp.ErrorMsg)
				log.Infow("user instead fail", "session_id", sessionId, "error", err)
				return
			} else {
				log.Infow("user instead success", "session_id", sessionId)
			}
		}
		// 再次抢占锁
		peer, err = facade.LockLocalUser(userId)
		log.Infow("lock local user return", "session_id", sessionId, "peer", peer, "err", err)
		if err != nil {
			return
		}
	}
	defer func() {
		if err != nil {
			peer, err := facade.UnlockLocalUser(userId)
			log.Infow("unlock local user return", "session_id", sessionId, "peer", peer, "err", err)
		}
	}()

	// 如果还没释放完成，则失败
	if _, ok := hall.SessionDict.Load(userId); ok {
		log.Infow("unexpect session", "session_id", sessionId)
		err = gira.ErrUserInstead
		return
	}
	// 创建会话
	session = &hall_sesssion{
		hall:      hall,
		sessionId: sessionId,
		userId:    userId,
		memberId:  memberId,
		avatar:    avatar,
		Actor:     actor.NewActor(hall.config.Framework.Hall.SessionActorBuffSize),
	}
	session.ctx, session.cancelFunc = context.WithCancel(ctx)

	if _, loaded := hall.SessionDict.LoadOrStore(userId, session); loaded {
		log.Infow("session store fail", "session_id", sessionId)
		err = gira.ErrUserInstead
		return
	}
	atomic.AddInt64(&hall.sessionCount, 1)
	var player Player
	// 加载数据
	if player, err = session.load(); err != nil {
		// 数据加载失败，释放锁
		hall.SessionDict.Delete(userId)
		atomic.AddInt64(&hall.sessionCount, -1)
		return
	} else {
		session.player = player
	}
	return
}

// 主要的业务逻辑的主协程
func (session *hall_sesssion) serve() {
	ticker := time.NewTicker(1 * time.Second)
	hall := session.hall
	saveTicker := time.NewTicker(time.Duration(session.hall.config.Framework.Hall.BgSaveInterval) * time.Second)
	session.chPeerPush = make(chan gira.ProtoPush, hall.config.Framework.Hall.PushBufferSize)
	defer func() {
		ticker.Stop()
		saveTicker.Stop()
		close(session.chPeerPush)
	}()
	sessionId := session.sessionId
	for {
		select {
		// 定时保存数据
		case <-saveTicker.C:
			if session.isClosed == 0 {
				session.save()
			}
		// 逻辑定时器
		case <-ticker.C:
			if session.isClosed == 0 {
				session.update()
			}
		// 处理客户端消息
		case req := <-session.chClientRequest:
			if req == nil {
				continue
			}
			if session.isClosed == 0 {
				if err := session.processClientMessage(req); err != nil {
					log.Infow("request fail", "session_id", sessionId, "error", err, "name")
				}
			} else {
				log.Infow("request fail", "session_id", sessionId, "error", "already closed")
			}
		case req := <-session.chPeerPush:
			if req == nil {
				continue
			}
			if session.isClosed == 0 {
				if err := session.processPeerPush(req); err != nil {
					log.Infow("push fail", "session_id", sessionId, "error", err, "name", req.GetPushName())
				} else {
					log.Infow("push fail", "session_id", sessionId, "error", "already closed", "name", req.GetPushName())
				}
			}
		case r := <-session.Inbox():
			r.Call()
		case <-session.ctx.Done():
			log.Infow("session exit", "session_id", sessionId)
			session.close(context.TODO())
			return
		}
	}
}

func (session *hall_sesssion) update() {
	defer func() {
		if e := recover(); e != nil {
			log.Error(e)
			debug.PrintStack()
		}
	}()
	session.player.Update()
}

func (session *hall_sesssion) save() {
	defer func() {
		if e := recover(); e != nil {
			log.Error(e)
			debug.PrintStack()
		}
	}()
	session.player.Save(session.ctx)
}

// 加载数据
func (session *hall_sesssion) load() (player Player, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error(e)
			debug.PrintStack()
			err = e.(error)
		}
	}()
	// player的ctx和session的ctx平级，player并不和session绑定，player可以将自己缓存起来，下次相同玩家登录的时候再复用
	player, err = session.hall.hallHandler.NewPlayer(session.hall.ctx, session, session.memberId, session.avatar)
	if err != nil {
		return
	}
	err = player.Load(session.hall.ctx, session.memberId, session.userId)
	if err != nil {
		return
	}
	return player, nil
}

// 发送控制消息，并关闭会话
func (self *hall_sesssion) sendPacketAndClose(ctx context.Context, typ hall_grpc.PacketType, reason string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Errorw("sendPacketAndClose", "type", typ, "reason", reason)
			log.Error(e)
			debug.PrintStack()
			err = e.(error)
		}
	}()
	resp := &hall_grpc.ClientMessageResponse{}
	resp.Type = typ
	resp.SessionId = self.sessionId
	resp.ReqId = 0
	resp.Data = []byte(reason)
	self.chClientResponse <- resp
	err = self.close(ctx)
	return nil
}

// 处理peer的push消息
func (session *hall_sesssion) processPeerPush(req gira.ProtoPush) (err error) {
	sessionId := session.sessionId
	var name string = req.GetPushName()
	log.Infow("peer push", "session_id", sessionId, "name", name)
	timeoutCtx, timeoutFunc := context.WithTimeout(session.ctx, 10*time.Second)
	defer func() {
		timeoutFunc()
		if e := recover(); e != nil {
			log.Error(e)
			debug.PrintStack()
			err = e.(error)
		}
	}()
	err = session.hall.proto.PushDispatch(timeoutCtx, session.hall.playerHandler, session.player, name, req)
	if err != nil {
		log.Errorw("push dispatch fail", "session_id", sessionId, "error", err, "name", name)
		return
	}
	return
}

// 处理客户端的消息
func (session *hall_sesssion) processClientMessage(message *hall_grpc.ClientMessageRequest) (err error) {
	sessionId := session.sessionId
	var name string
	var reqId int32
	var req interface{}
	var resp []byte
	var pushArr []gira.ProtoPush
	name, reqId, req, err = session.hall.proto.RequestDecode(message.Data)
	if err != nil {
		log.Errorw("request decode fail", "session_id", sessionId, "error", err)
		return
	}
	log.Infow("request ", "session_id", sessionId, "name", name, "req_id", reqId, "data", message.Data)

	timeoutCtx, timeoutFunc := context.WithTimeout(session.ctx, 10*time.Second)
	defer func() {
		timeoutFunc()
		if e := recover(); e != nil {
			// err = gira.ErrBrokenChannel
			log.Error(e)
			debug.PrintStack()
			err = e.(error)
		}
	}()
	resp, pushArr, err = session.hall.proto.RequestDispatch(timeoutCtx, session.hall.playerHandler, session.player, name, reqId, req)
	if err != nil {
		log.Errorw("request dispatch fail", "session_id", sessionId, "error", err)
		return
	}
	if session.isClosed != 0 {
		return nil
	}
	// log.Debugw("data response", "session_id", sessionId, "data response", len(dataResp))
	response := &hall_grpc.ClientMessageResponse{}
	response.Type = hall_grpc.PacketType_DATA
	response.SessionId = message.SessionId
	response.ReqId = message.ReqId
	response.Data = resp
	response.Route = name
	session.chClientResponse <- response
	if pushArr != nil {
		for _, push := range pushArr {
			if dataPush, err := session.hall.proto.PushEncode(push); err != nil {
				log.Errorw("push fail", "session_id", sessionId, "error", err)
			} else {
				log.Infow("push to client", "session_id", sessionId, "name", push.GetPushName(), "data", len(dataPush))
				response := &hall_grpc.ClientMessageResponse{}
				response.Type = hall_grpc.PacketType_DATA
				response.SessionId = message.SessionId
				response.ReqId = 0
				response.Data = dataPush
				response.Route = push.GetPushName()
				session.chClientResponse <- response
			}
		}
	}
	return
}

// 关闭session
// #actor(call)
func (session *hall_sesssion) close(ctx context.Context) (err error) {
	if !atomic.CompareAndSwapInt32(&session.isClosed, 0, 1) {
		err = gira.ErrSessionClosed
		return
	}
	userId := session.userId
	sessionId := session.sessionId
	log.Infow("session close", "session_id", sessionId)
	// 脱离stream
	session.clientCancelFunc()
	defer func() {
		if e := recover(); e != nil {
			log.Error(err)
			debug.PrintStack()
			var ok bool
			var s string
			if err, ok = e.(error); ok {
			} else if s, ok = e.(string); ok {
				err = errors.New(s)
			} else {
				err = errors.New("panic")
			}
		}
		session.cancelFunc()
		peer, err := facade.UnlockLocalUser(userId)
		log.Infow("unlock local user return", "session_id", sessionId, "peer", peer, "err", err)
		// 从agent dict释放
		session.hall.SessionDict.Delete(userId)
		atomic.AddInt64(&session.hall.sessionCount, -1)
		log.Infow("session close finished", "session_id", sessionId)
	}()
	// 开始关闭操作
	if session.player != nil {
		session.player.Logout(ctx)
	}
	return
}

// 踢下线
// call by handler
func (self *hall_sesssion) Kick(ctx context.Context, reason string) (err error) {
	return self.sendPacketAndClose(ctx, hall_grpc.PacketType_KICK, reason)
}

// 顶号下线
// call by hall
// #actor(call)
func (self *hall_sesssion) instead(ctx context.Context, reason string) (err error) {
	return self.sendPacketAndClose(ctx, hall_grpc.PacketType_USER_INSTEAD, reason)
}

// 推送消息
func (self *hall_sesssion) Push(push gira.ProtoPush) (err error) {
	var data []byte
	if data, err = self.hall.proto.PushEncode(push); err != nil {
		return
	} else {
		defer func() {
			if e := recover(); e != nil {
				log.Error(e)
				debug.PrintStack()
				err = e.(error)
			}
		}()
		resp := &hall_grpc.ClientMessageResponse{}
		resp.Type = hall_grpc.PacketType_DATA
		resp.SessionId = self.sessionId
		resp.ReqId = 0
		resp.Data = data
		resp.Route = push.GetPushName()
		// WARN: chResponse有可能已经关闭
		self.chClientResponse <- resp
		return
	}
}

// func (self *hall_sesssion) Notify(userId string, resp gira.ProtoPush) error {
// 	if v, ok := self.hall.SessionDict.Load(userId); !ok {
// 		return gira.ErrNoSession
// 	} else {
// 		otherSession, _ := v.(*hall_sesssion)
// 		return otherSession.Push(resp)
// 	}
// }
//
/// 宏展开的地方，不要在文件末尾添加代码============// afafa

type hall_sesssioncloseArgument struct {
	session *hall_sesssion
	ctx context.Context
	err error
	__caller__ chan*hall_sesssioncloseArgument
}

func (__arg__ *hall_sesssioncloseArgument) Call() {
	__arg__.err = __arg__.session.close(__arg__.ctx)
	__arg__.__caller__ <- __arg__
}

func (session *hall_sesssion) Call_close (ctx context.Context, opts ...actor.CallOption) (err error){
	var __options__ actor.CallOptions
	for _, v := range opts {
		v.Config(&__options__)	
	}
	__arg__ := &hall_sesssioncloseArgument {
		session: session,
		ctx: ctx,
		__caller__: make(chan*hall_sesssioncloseArgument),
	}
	session.Inbox() <- __arg__
	if __options__.TimeOut != 0 {
		__timer__ := time.NewTimer(__options__.TimeOut)
		defer __timer__.Stop()
		select {
		case resp :=<-__arg__.__caller__:
			return resp.err
		case <-ctx.Done():
			return ctx.Err()
		case <-__timer__.C:
	        log.Errorw("actor call time out", "func", "func (session *hall_sesssion) close(ctx context.Context) (err error) ")
			return actor.ErrCallTimeOut
		}
	} else {
		select {
		case resp :=<-__arg__.__caller__:
			return resp.err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}



	
type hall_sesssioninsteadArgument struct {
	self *hall_sesssion
	ctx context.Context
	reason string
	err error
	__caller__ chan*hall_sesssioninsteadArgument
}

func (__arg__ *hall_sesssioninsteadArgument) Call() {
	__arg__.err = __arg__.self.instead(__arg__.ctx, __arg__.reason)
	__arg__.__caller__ <- __arg__
}

func (self *hall_sesssion) Call_instead (ctx context.Context, reason string, opts ...actor.CallOption) (err error){
	var __options__ actor.CallOptions
	for _, v := range opts {
		v.Config(&__options__)	
	}
	__arg__ := &hall_sesssioninsteadArgument {
		self: self,
		ctx: ctx,
		reason: reason,
		__caller__: make(chan*hall_sesssioninsteadArgument),
	}
	self.Inbox() <- __arg__
	if __options__.TimeOut != 0 {
		__timer__ := time.NewTimer(__options__.TimeOut)
		defer __timer__.Stop()
		select {
		case resp :=<-__arg__.__caller__:
			return resp.err
		case <-ctx.Done():
			return ctx.Err()
		case <-__timer__.C:
	        log.Errorw("actor call time out", "func", "func (self *hall_sesssion) instead(ctx context.Context, reason string) (err error) ")
			return actor.ErrCallTimeOut
		}
	} else {
		select {
		case resp :=<-__arg__.__caller__:
			return resp.err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}



	
/// =============宏展开的地方，不要在文件末尾添加代码============