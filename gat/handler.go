package gat

type MiddleWareInterface interface {
	ServeMessage(r *Request)
	OnSessionOpen(s *Session)
	OnSessionClose(s *Session)
}
