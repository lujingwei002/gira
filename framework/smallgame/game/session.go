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
	"github.com/lujingwei002/gira/sproto"
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
	chResponse       chan *hall_grpc.StreamDataResponse // 由stream负责关闭
	chRequest        chan *hall_grpc.StreamDataRequest  // 由stream负责关闭
	chPush           chan sproto.SprotoPush             // 由session负责关闭
	streamCancelFunc context.CancelFunc
	isClosed         int32
}

func (hall *hall_server) createSession(ctx context.Context, sessionId uint64, memberId string) (session *hall_sesssion, err error) {
	var userId string
	var avatar UserAvatar
	avatar, err = hall.hallHandler.Login(ctx, memberId)
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
			log.Infow("user instread", "session_id", userId)
			var resp *hall_grpc.UserInsteadResponse
			if resp, err = rpc.Hall.UserInstead(ctx, peer, userId); err != nil {
				log.Infow("user instread fail", "session_id", sessionId, "error", err)
				return
			} else if resp.ErrorCode != 0 {
				err = gira.NewError(resp.ErrorCode, resp.ErrorMsg)
				log.Infow("user instread fail", "session_id", sessionId, "error", err)
				return
			} else {
				log.Infow("user instread success", "session_id", sessionId)
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
	var player Player
	// 加载数据
	if player, err = session.load(); err != nil {
		// 数据加载失败，释放锁
		hall.SessionDict.Delete(userId)
		return
	} else {
		session.player = player
	}
	return
}

// 主要的业务逻辑的主协程
func (self *hall_sesssion) serve() {
	self.chPush = make(chan sproto.SprotoPush, 10)
	defer func() {
		close(self.chPush)
	}()
	sessionId := self.sessionId
	for {
		select {
		case req := <-self.chRequest:
			if req != nil {
				if self.isClosed == 0 {
					if err := self.request(req); err != nil {
						log.Infow("request fail", "session_id", sessionId, "error", err)
					}
				}
			}
		case req := <-self.chPush:
			if req != nil {
				if self.isClosed == 0 {
					if err := self.push(req); err != nil {
						log.Infow("push fail", "session_id", sessionId, "error", err, "name", req.GetPushName())
					}
				}
			}
		case r := <-self.Inbox():
			r.Call()
		case <-self.ctx.Done():
			log.Infow("session exit", "session_id", sessionId)
			return
		}
	}
}

// 加载数据
func (self *hall_sesssion) load() (player Player, err error) {
	player, err = self.hall.hallHandler.NewPlayer(self.ctx, self, self.memberId, self.avatar)
	if err = player.Load(self.ctx, self.memberId, self.userId); err != nil {
		return nil, err
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

func (self *hall_sesssion) push(req sproto.SprotoPush) error {
	sessionId := self.sessionId
	var err error
	var name string = req.GetPushName()
	log.Infow("request ", "session_id", sessionId, "name", name)

	timeoutCtx, timeoutFunc := context.WithTimeout(self.ctx, 10*time.Second)
	defer timeoutFunc()
	err = self.hall.proto.PushDispatch(timeoutCtx, self.hall.playerHandler, self.player, name, req)
	if err != nil {
		log.Errorw("push dispatch fail", "session_id", sessionId, "error", err, "name", name)
		return err
	}
	return nil
}

// 处理请求
func (self *hall_sesssion) request(req *hall_grpc.StreamDataRequest) (err error) {
	sessionId := self.sessionId
	var name string
	var reqId int32
	var dataReq interface{}
	var dataResp []byte
	var dataPushArr [][]byte
	_, name, reqId, dataReq, err = self.hall.proto.RequestDecode(req.Data)
	if err != nil {
		log.Errorw("request decode fail", "session_id", sessionId, "error", err)
		return
	}
	log.Infow("request ", "session_id", sessionId, "name", name, "req_id", reqId)

	timeoutCtx, timeoutFunc := context.WithTimeout(self.ctx, 10*time.Second)
	defer timeoutFunc()
	dataResp, dataPushArr, err = self.hall.proto.RequestDispatch(timeoutCtx, self.hall.playerHandler, self.player, name, reqId, dataReq)
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
	self.chResponse <- resp
	if dataPushArr != nil {
		for _, dataPush := range dataPushArr {
			log.Debugw("data push", "session_id", sessionId, "data resp", dataPush)
			resp := &hall_grpc.StreamDataResponse{}
			resp.Type = hall_grpc.PacketType_DATA
			resp.SessionId = req.SessionId
			resp.ReqId = 0
			resp.Data = dataPush
			self.chResponse <- resp
		}
	}
	return
}

// 关闭session
// #actor(call)
func (self *hall_sesssion) close(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&self.isClosed, 0, 1) {
		return gira.ErrSessionClosed
	}
	userId := self.userId
	sessionId := self.sessionId
	// 脱离stream
	self.streamCancelFunc()
	// 开始关闭操作
	if self.player != nil {
		self.player.Logout(ctx)
	}
	self.cancelFunc()
	peer, err := facade.UnlockLocalUser(userId)
	log.Infow("unlock local user return", "session_id", sessionId, "peer", peer, "err", err)
	// 从agent dict释放
	self.hall.SessionDict.Delete(userId)
	log.Infow("session close finished", "session_id", sessionId)
	return nil
}

// 网关连接关闭
// @goroutine(hall.UpstreamServer)
func (self *hall_sesssion) OnStreamClose() {
	self.CallWithTimeout_close(5 * time.Second)
}

// 踢下线
// call by handler
func (self *hall_sesssion) Kick(ctx context.Context, reason string) (err error) {
	return self.sendPacketAndClose(ctx, hall_grpc.PacketType_KICK, reason)
}

// 顶号下线
// #actor(call)
// call by hall
func (self *hall_sesssion) instead(ctx context.Context, reason string) (err error) {
	return self.sendPacketAndClose(ctx, hall_grpc.PacketType_USER_INSTEAD, reason)
}

// 推送消息
func (self *hall_sesssion) Push(resp sproto.SprotoPush) (err error) {
	var data []byte
	if data, err = self.hall.proto.PushEncode(resp); err != nil {
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
		self.chResponse <- resp
		return
	}
}

func (self *hall_sesssion) Notify(userId string, resp sproto.SprotoPush) error {
	if v, ok := self.hall.SessionDict.Load(userId); !ok {
		return gira.ErrNoSession
	} else {
		otherSession, _ := v.(*hall_sesssion)
		return otherSession.Push(resp)
	}
}

/// 宏展开的地方，不要在文件末尾添加代码============// afafa

type SessioncloseArgument struct {
	self       *hall_sesssion
	ctx        context.Context
	r0         error
	__caller__ chan *SessioncloseArgument
}

func (__arg__ *SessioncloseArgument) Call() {
	__arg__.r0 = __arg__.self.close(__arg__.ctx)
	__arg__.__caller__ <- __arg__
}

func (self *hall_sesssion) Call_close(ctx context.Context) (r0 error) {
	__arg__ := &SessioncloseArgument{
		self:       self,
		ctx:        ctx,
		__caller__: make(chan *SessioncloseArgument),
	}
	self.Inbox() <- __arg__
	select {
	case resp := <-__arg__.__caller__:
		return resp.r0
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (self *hall_sesssion) CallWithTimeout_close(__timeout__ time.Duration) (r0 error) {
	__ctx__, __cancel__ := context.WithTimeout(context.Background(), __timeout__)
	defer __cancel__()
	return self.Call_close(__ctx__)
}

type SessioninsteadArgument struct {
	self       *hall_sesssion
	ctx        context.Context
	reason     string
	err        error
	__caller__ chan *SessioninsteadArgument
}

func (__arg__ *SessioninsteadArgument) Call() {
	__arg__.err = __arg__.self.instead(__arg__.ctx, __arg__.reason)
	__arg__.__caller__ <- __arg__
}

func (self *hall_sesssion) Call_instead(ctx context.Context, reason string) (err error) {
	__arg__ := &SessioninsteadArgument{
		self:       self,
		ctx:        ctx,
		reason:     reason,
		__caller__: make(chan *SessioninsteadArgument),
	}
	self.Inbox() <- __arg__
	select {
	case resp := <-__arg__.__caller__:
		return resp.err
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (self *hall_sesssion) CallWithTimeout_instead(__timeout__ time.Duration, reason string) (err error) {
	__ctx__, __cancel__ := context.WithTimeout(context.Background(), __timeout__)
	defer __cancel__()
	return self.Call_instead(__ctx__, reason)
}

/// =============宏展开的地方，不要在文件末尾添加代码============
