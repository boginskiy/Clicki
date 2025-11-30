package audit

import "net/http"

// Auditer - interface. Doesn't use.
type Auditer interface {
	NoticeCreateLink(req *http.Request)
	NoticeFollowLink(req *http.Request)
	NeedAudit(req *http.Request) bool
}

// Publisher - interface.
type Publisher interface {
	CheckSubscribers() bool
	Deregister(Subscriber)
	Register(Subscriber)
	Send(event any)
}

// Subscriber - interface.
type Subscriber interface {
	Update(event any)
	GetID() int
	Close()
}
