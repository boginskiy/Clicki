package audit

import "time"

// Event is message
type Event struct {
	TimeS  int64  `json:"ts"`
	Action string `json:"action"`
	UserID int    `json:"user_id"`
	URL    string `json:"url"`
}

func NewEvent(action string, userID int, url string) *Event {
	return &Event{
		TimeS:  time.Now().Unix(),
		Action: action,
		UserID: userID,
		URL:    url,
	}
}

// Publish is publisher
type Publish struct {
	Subscribers map[int]Subscriber
	Count       int
}

func NewPublish(subs ...Subscriber) *Publish {
	tmpMap := make(map[int]Subscriber, len(subs))
	cnt := len(subs)

	for _, sub := range subs {
		if sub.GetID() == 0 {
			cnt--
		} else {
			tmpMap[sub.GetID()] = sub
		}
	}

	return &Publish{Subscribers: tmpMap, Count: cnt}
}

func (p *Publish) CheckSubscribers() bool {
	return p.Count != 0
}

func (p *Publish) Register(sub Subscriber) {
	if sub == nil {
		return
	}
	p.Subscribers[sub.GetID()] = sub
	p.Count++
}

func (p *Publish) Deregister(sub Subscriber) {
	delete(p.Subscribers, sub.GetID())
	p.Count--
}

func (p *Publish) Send(event any) {
	for _, subscriber := range p.Subscribers {
		if subscriber != nil {
			subscriber.Update(event)
		}
	}
}
