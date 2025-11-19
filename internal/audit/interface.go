package audit

import "net/http"

type Auditer interface {
	NoticeCreateLink(req *http.Request)
	NoticeFollowLink(req *http.Request)
	NeedAudit(req *http.Request) bool
}

type Publisher interface {
	CheckSubscribers() bool
	Deregister(Subscriber)
	Register(Subscriber)
	Send(event any)
}

type Subscriber interface {
	Update(event any)
	GetID() int
	Clouse()
}
