package eventsourcing

import (
	m "github.com/maxkaemmerer/go-message"
)

type AggregateChanged interface {
	m.Message
	AggregateId() string
}

type Aggregate interface {
	MetaData() AggregateMetaData
	Apply(aggregateChanged AggregateChanged) error
}

type AggregateFactory func(aggregateId, aggregateType string, events []AggregateChanged) Aggregate

type EventStore interface {
	Get(aggregateId string, aggregateType string) (Aggregate, error)
	Store(aggregate *Aggregate) error
}
