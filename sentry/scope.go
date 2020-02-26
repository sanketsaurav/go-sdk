package sentry

type EventProcessor func(event *Event, hint *EventHint) *Event

func NewScope() *Scope {
	scope := Scope{
		breadcrumbs: make([]*Breadcrumb, 0),
		tags:        make(map[string]string),
		contexts:    make(map[string]interface{}),
		extra:       make(map[string]interface{}),
		fingerprint: make([]string, 0),
	}

	return &scope
}

type Scope struct {
	sync.RWMutex
	breadcrumbs     []*Breadcrumb
	user            User
	tags            map[string]string
	contexts        map[string]interface{}
	extra           map[string]interface{}
	fingerprint     []string
	level           Level
	transaction     string
	request         Request
	eventProcessors []EventProcessor
}

// Ad
