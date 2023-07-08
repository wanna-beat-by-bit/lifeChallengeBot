package events

type Fethcer interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(event Event) error
}

type Event struct {
	Type Type
	Text string
	Meta interface{}
}

type Type int

const (
	Message Type = iota
	Unknown
)
