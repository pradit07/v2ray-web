package web

import (
	"bytes"
	"sync"
	"text/template"

	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"v2ray.com/core"
	"v2ray.com/ext/tools/conf/serial"
)

var srv = _Service{
	listeners: make(map[int64]ListenerFunc),
	status:    &ListenerStatus{},
}

// ListenerStatus .
type ListenerStatus struct {
	Run          bool   `json:"run,omitempty"`
	ID           uint64 `json:"id,omitempty"`
	Subscription uint64 `json:"subscription,omitempty"`
}

// ListenerFunc .
type ListenerFunc func(*ListenerStatus)
type _Service struct {
	server    *core.Instance
	listeners map[int64]ListenerFunc
	id        int64
	status    *ListenerStatus
	sync.Mutex
}

func (s *_Service) AddListener(listener ListenerFunc) (id int64) {
	s.Lock()
	id = s.id
	s.id++
	s.listeners[id] = listener
	listener(s.status)
	s.Unlock()
	return
}
func (s *_Service) RemoveListener(id int64) {
	s.Lock()
	delete(s.listeners, id)
	s.Unlock()
	return
}
func (s *_Service) Start(element *data.Element) (e error) {
	s.Lock()
	defer s.Unlock()
	var mSettings manipulator.Settings
	str, e := mSettings.GetV2ray()
	if e != nil {
		return
	}
	ctx, e := element.Outbound.ToContext()
	if e != nil {
		return
	}
	t := template.New("v2ray")
	t, e = t.Parse(str)
	if e != nil {
		return
	}
	var buffer bytes.Buffer
	e = t.Execute(&buffer, ctx)
	if e != nil {
		return
	}
	// v2ray
	cnf, e := serial.LoadJSONConfig(&buffer)
	if e != nil {
		return
	}
	server, e := core.New(cnf)
	if e != nil {
		return
	}
	var closed bool
	if s.server != nil {
		s.server.Close()
		s.server = nil
		closed = true
	}

	e = server.Start()
	if e == nil {
		s.server = server
		s.notify(&ListenerStatus{
			Run:          true,
			ID:           element.ID,
			Subscription: element.Subscription,
		})
	} else {
		if closed {
			s.notify(&ListenerStatus{})
		}
	}
	return
}
func (s *_Service) Stop() {
	s.Lock()
	defer s.Unlock()
	if s.server == nil {
		return
	}
	s.server.Close()
	s.notify(&ListenerStatus{})
}
func (s *_Service) notify(status *ListenerStatus) {
	if s.status.Run {
		if status.Run && s.status.ID == status.ID && s.status.Subscription == status.Subscription {
			return
		}
	} else if !status.Run {
		return
	}

	s.status = status
	for _, f := range s.listeners {
		f(status)
	}
}