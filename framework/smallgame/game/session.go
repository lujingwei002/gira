package game

import (
	"context"
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
	hall             *hall
	ctx              context.Context
	cancelFunc       context.CancelFunc
	sessionId        uint64
	memberId         string
	userId           string
	avatar           UserAvatar
	player           Player
	chResponse       chan *hall_grpc.StreamDataResponse // 由stream负责关闭
	chRequest        chan *hall_grpc.StreamDataRequest  // 由stream负责关闭
	chPush           chan gira.ProtoPush                // 由session负责关闭
	streamCancelFunc context.CancelFunc
	isClosed         int32
}

func newSession(hall *hall, sessionId uint64, memberId string) (session *hall_sesssion, err error) {
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
		Actor:     actor.NewActor(16),
	}
	session.ctx, session.cancelFunc = context.WithCancel(ctx)

	if _, loaded := hall.SessionDict.LoadOrStore(userId, session); loaded {
		log.Infow("session store fail", "session_id", sessionId)
		err = gira.ErrUserInstead
		return
	}
	atomic.AddInt64(&hall.SessionCount, 1)
	var player Player
	// 加载数据
	if player, err = session.load(); err != nil {
		// 数据加载失败，释放锁
		hall.SessionDict.Delete(userId)
		atomic.AddInt64(&hall.SessionCount, -1)
		return
	} else {
		session.player = player
	}
	return
}

// 主要的业务逻辑的主协程
func (session *hall_sesssion) serve() {
	ticker := time.NewTicker(1 * time.Second)
	saveTicker := time.NewTicker(time.Duration(session.hall.config.BgSaveInterval) * time.Second)
	session.chPush = make(chan gira.ProtoPush, 10)
	defer func() {
		ticker.Stop()
		saveTicker.Stop()
		close(session.chPush)
	}()
	sessionId := session.sessionId
	for {
		select {
		case <-saveTicker.C:
			if session.player != nil && session.isClosed == 0 {
				session.player.Save(session.ctx)
			}
		case <-ticker.C:
			if session.player != nil && session.isClosed == 0 {
				session.player.Update()
			}
		case req := <-session.chRequest:
			if req != nil {
				if session.isClosed == 0 {
					if err := session.request(req); err != nil {
						log.Infow("request fail", "session_id", sessionId, "error", err)
					}
				}
			}
		case req := <-session.chPush:
			if req != nil {
				if session.isClosed == 0 {
					if err := session.processPush(req); err != nil {
						log.Infow("push fail", "session_id", sessionId, "error", err, "name", req.GetPushName())
					}
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

// 加载数据
func (session *hall_sesssion) load() (player Player, err error) {
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
			err = gira.ErrBrokenChannel
		}
	}()
	resp := &hall_grpc.StreamDataResponse{}
	resp.Type = typ
	resp.SessionId = self.sessionId
	resp.ReqId = 0
	resp.Data = []byte(reason)
	self.chResponse <- resp
	err = self.close(ctx)
	return nil
}

// 处理processPush请求
func (session *hall_sesssion) processPush(req gira.ProtoPush) error {
	sessionId := session.sessionId
	var err error
	var name string = req.GetPushName()
	log.Infow("request ", "session_id", sessionId, "name", name)

	timeoutCtx, timeoutFunc := context.WithTimeout(session.ctx, 10*time.Second)
	defer timeoutFunc()
	err = session.hall.proto.PushDispatch(timeoutCtx, session.hall.playerHandler, session.player, name, req)
	if err != nil {
		log.Errorw("push dispatch fail", "session_id", sessionId, "error", err, "name", name)
		return err
	}
	return nil
}

// 处理请求
func (session *hall_sesssion) request(req *hall_grpc.StreamDataRequest) (err error) {
	sessionId := session.sessionId
	var name string
	var reqId int32
	var dataReq interface{}
	var dataResp []byte
	var pushArr []gira.ProtoPush
	name, reqId, dataReq, err = session.hall.proto.RequestDecode(req.Data)
	if err != nil {
		log.Errorw("request decode fail", "session_id", sessionId, "error", err)
		return
	}
	log.Infow("request ", "session_id", sessionId, "name", name, "req_id", reqId)

	timeoutCtx, timeoutFunc := context.WithTimeout(session.ctx, 10*time.Second)
	defer timeoutFunc()
	dataResp, pushArr, err = session.hall.proto.RequestDispatch(timeoutCtx, session.hall.playerHandler, session.player, name, reqId, dataReq)
	if err != nil {
		log.Errorw("request dispatch fail", "session_id", sessionId, "error", err)
		return
	}
	defer func() {
		if e := recover(); e != nil {
			err = gira.ErrBrokenChannel
		}
	}()

	// log.Debugw("data resp", "session_id", sessionId, "data resp", len(dataResp))
	resp := &hall_grpc.StreamDataResponse{}
	resp.Type = hall_grpc.PacketType_DATA
	resp.SessionId = req.SessionId
	resp.ReqId = req.ReqId
	resp.Data = dataResp
	resp.Route = name
	session.chResponse <- resp
	if pushArr != nil {
		for _, push := range pushArr {
			if dataPush, err := session.hall.proto.PushEncode(push); err != nil {
			} else {
				log.Debugw("data push", "session_id", sessionId, "data resp", dataPush)
				resp := &hall_grpc.StreamDataResponse{}
				resp.Type = hall_grpc.PacketType_DATA
				resp.SessionId = req.SessionId
				resp.ReqId = 0
				resp.Data = dataPush
				resp.Route = push.GetPushName()
				session.chResponse <- resp
			}
		}
	}
	return
}

// 关闭session
// #actor(call)
func (session *hall_sesssion) close(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&session.isClosed, 0, 1) {
		return gira.ErrSessionClosed
	}
	log.Infow("session close")
	userId := session.userId
	sessionId := session.sessionId
	// 脱离stream
	session.streamCancelFunc()
	// 开始关闭操作
	if session.player != nil {
		session.player.Logout(ctx)
	}
	session.cancelFunc()
	peer, err := facade.UnlockLocalUser(userId)
	log.Infow("unlock local user return", "session_id", sessionId, "peer", peer, "err", err)
	// 从agent dict释放
	session.hall.SessionDict.Delete(userId)
	atomic.AddInt64(&session.hall.SessionCount, -1)
	log.Infow("session close finished", "session_id", sessionId)
	return nil
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

func (self *hall_sesssion) serverDown(ctx context.Context, reason string) (err error) {
	return self.sendPacketAndClose(ctx, hall_grpc.PacketType_SERVER_DOWN, reason)
}

// 推送消息
func (self *hall_sesssion) Push(push gira.ProtoPush) (err error) {
	var data []byte
	if data, err = self.hall.proto.PushEncode(push); err != nil {
		return
	} else {
		defer func() {
			if e := recover(); e != nil {
				err = gira.ErrBrokenChannel
			}
		}()
		resp := &hall_grpc.StreamDataResponse{}
		resp.Type = hall_grpc.PacketType_DATA
		resp.SessionId = self.sessionId
		resp.ReqId = 0
		resp.Data = data
		resp.Route = push.GetPushName()
		self.chResponse <- resp
		return
	}
}

func (self *hall_sesssion) Notify(userId string, resp gira.ProtoPush) error {
	if v, ok := self.hall.SessionDict.Load(userId); !ok {
		return gira.ErrNoSession
	} else {
		otherSession, _ := v.(*hall_sesssion)
		return otherSession.Push(resp)
	}
}

/// 宏展开的地方，不要在文件末尾添加代码============// afafa

type hall_sesssioncloseArgument struct {
	self       *hall_sesssion
	ctx        context.Context
	r0         error
	__caller__ chan *hall_sesssioncloseArgument
}

func (__arg__ *hall_sesssioncloseArgument) Call() {
	__arg__.r0 = __arg__.self.close(__arg__.ctx)
	__arg__.__caller__ <- __arg__
}

func (self *hall_sesssion) Call_close(ctx context.Context, opts ...actor.CallOption) (r0 error) {
	var __options__ actor.CallOptions
	for _, v := range opts {
		v.Config(&__options__)
	}
	__arg__ := &hall_sesssioncloseArgument{
		self:       self,
		ctx:        ctx,
		__caller__: make(chan *hall_sesssioncloseArgument),
	}
	self.Inbox() <- __arg__
	if __options__.TimeOut != 0 {
		__timer__ := time.NewTimer(__options__.TimeOut)
		defer __timer__.Stop()
		select {
		case resp := <-__arg__.__caller__:
			return resp.r0
		case <-ctx.Done():
			return ctx.Err()
		case <-__timer__.C:
			log.Errorw("actor call time out", "func", "func (self *hall_sesssion) close(ctx context.Context) error ")
			return actor.ErrCallTimeOut
		}
	} else {
		select {
		case resp := <-__arg__.__caller__:
			return resp.r0
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

type hall_sesssioninsteadArgument struct {
	self       *hall_sesssion
	ctx        context.Context
	reason     string
	err        error
	__caller__ chan *hall_sesssioninsteadArgument
}

func (__arg__ *hall_sesssioninsteadArgument) Call() {
	__arg__.err = __arg__.self.instead(__arg__.ctx, __arg__.reason)
	__arg__.__caller__ <- __arg__
}

func (self *hall_sesssion) Call_instead(ctx context.Context, reason string, opts ...actor.CallOption) (err error) {
	var __options__ actor.CallOptions
	for _, v := range opts {
		v.Config(&__options__)
	}
	__arg__ := &hall_sesssioninsteadArgument{
		self:       self,
		ctx:        ctx,
		reason:     reason,
		__caller__: make(chan *hall_sesssioninsteadArgument),
	}
	self.Inbox() <- __arg__
	if __options__.TimeOut != 0 {
		__timer__ := time.NewTimer(__options__.TimeOut)
		defer __timer__.Stop()
		select {
		case resp := <-__arg__.__caller__:
			return resp.err
		case <-ctx.Done():
			return ctx.Err()
		case <-__timer__.C:
			log.Errorw("actor call time out", "func", "func (self *hall_sesssion) instead(ctx context.Context, reason string) (err error) ")
			return actor.ErrCallTimeOut
		}
	} else {
		select {
		case resp := <-__arg__.__caller__:
			return resp.err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

/// =============宏展开的地方，不要在文件末尾添加代码============
