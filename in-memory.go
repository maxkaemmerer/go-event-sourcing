package eventsourcing

import (
	"fmt"

	m "github.com/maxkaemmerer/go-message"
)

type DispatchingInMemoryEventStore struct {
	Factories  map[string]AggregateFactory
	Aggregates map[string]map[string][]AggregateChanged
	MessageBus m.MessageBus
}

func (es *DispatchingInMemoryEventStore) Store(aggregate *Aggregate) error {
	metaData := (*aggregate).MetaData()

	aggregatesOfType, ok := es.Aggregates[metaData.AggregateType()]
	if !ok {
		es.Aggregates[metaData.AggregateType()] = map[string][]AggregateChanged{}
	}

	existingEvents, ok := aggregatesOfType[metaData.AggregateId()]
	if !ok {
		es.Aggregates[metaData.AggregateType()][metaData.AggregateId()] = []AggregateChanged{}
	}

	es.Aggregates[metaData.AggregateType()][metaData.AggregateId()] = append(existingEvents, metaData.PendingEvents()...)

	for _, pendingEvent := range metaData.PendingEvents() {
		err := es.MessageBus.Dispatch(pendingEvent)
		if err != nil {
			return err
		}
	}

	metaData.ClearPendingEvents()

	return nil
}

func (es *DispatchingInMemoryEventStore) Get(aggregateId string, aggregateType string) (Aggregate, error) {
	existingEvents, ok := es.Aggregates[aggregateType][aggregateId]
	if !ok {
		return nil, fmt.Errorf("could not construct aggregate %s of type %s, no events found in store", aggregateId, aggregateType)
	}

	factory, ok := es.Factories[aggregateType]
	if !ok {
		return nil, fmt.Errorf("could not construct aggregate %s of type %s, no factory found", aggregateId, aggregateType)
	}

	return factory(aggregateId, aggregateType, existingEvents), nil
}
