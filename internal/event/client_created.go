package event

import "time"

type ClientCreated struct {
	Name    string
	Payload interface{}
}

func NewClientCreated() *ClientCreated {
	return &ClientCreated{
		Name: "ClientCreated",
	}
}

func (e *ClientCreated) GetName() string {
	return e.Name
}

func (e *ClientCreated) GetPayload() interface{} {
	return e.Payload
}

func (e *ClientCreated) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (e *ClientCreated) GetDateTime() time.Time {
	return time.Now()
}
