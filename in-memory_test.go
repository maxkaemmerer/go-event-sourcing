package eventsourcing_test

import (
	"fmt"
	"testing"

	es "github.com/maxkaemmerer/go-event-sourcing"

	"github.com/maxkaemmerer/go-message"
)

type testAggregate struct {
	metaData      es.AggregateMetaData
	appliedEvents []es.AggregateChanged
}

type testAggregateCreated struct {
}

func (e *testAggregateCreated) AggregateId() string {
	return "aggregate-1"
}
func (e *testAggregateCreated) ArbitraryValue() string {
	return "Hello World!"
}
func (e *testAggregateCreated) Name() string {
	return "testAggregateCreated"
}
func (e *testAggregateCreated) Context() string {
	return "testAggregate"
}

type testAggregateCreatedHandler struct {
}

func (e *testAggregateCreatedHandler) MessageName() string {
	return "testAggregateCreated"
}

func (e *testAggregateCreatedHandler) ContextName() string {
	return "testAggregate"
}

func (e *testAggregateCreatedHandler) Handle(msg message.Message) error {
	return nil
}

func NewTestAggregate() es.Aggregate {
	aggregate := &testAggregate{}

	err := aggregate.Record(&testAggregateCreated{})
	if err != nil {
		panic(err)
	}

	return aggregate

}

func (a *testAggregate) MetaData() es.AggregateMetaData {
	return a.metaData
}

func (a *testAggregate) Apply(aggregateChanged es.AggregateChanged) error {
	a.appliedEvents = append(a.appliedEvents, aggregateChanged)
	switch aggregateChanged.Name() {
	case "testAggregateCreated":
		a.metaData = es.InitMetaData(aggregateChanged.AggregateId(), aggregateChanged.Context())
	default:
		return fmt.Errorf("Unknown event %s", aggregateChanged.Name())
	}
	a.metaData.IncrementVersion()
	return nil
}
func (a *testAggregate) Record(aggregateChanged es.AggregateChanged) error {
	err := a.Apply(aggregateChanged)
	a.MetaData().Record(aggregateChanged)
	return err
}

func Test_DispatchingInMemoryEventStore_ShouldStoreAggregate(t *testing.T) {
	eventStore := &es.DispatchingInMemoryEventStore{
		Factories: map[string]es.AggregateFactory{
			"testAggregate": func(aggregateId, aggregateType string, events []es.AggregateChanged) es.Aggregate {
				aggregate := &testAggregate{}
				aggregate.metaData = es.InitMetaData(aggregateId, aggregateType)
				for _, event := range events {
					aggregate.Apply(event)
				}
				return aggregate
			},
		},
		Aggregates: map[string]map[string][]es.AggregateChanged{},
		MessageBus: message.NewSimpleMessageBus([]message.Handler{
			&testAggregateCreatedHandler{},
		}),
	}

	aggregate := NewTestAggregate()

	err := eventStore.Store(aggregate)

	if err != nil {
		t.Errorf("Should not have returned error")
	}

	if len(eventStore.Aggregates) != 1 {
		t.Errorf("Should have one aggregate type")
	}
	if len(eventStore.Aggregates[aggregate.MetaData().AggregateType()]) != 1 {
		t.Errorf("Should have one aggregate of type")
	}
	if len(eventStore.Aggregates[aggregate.MetaData().AggregateType()][aggregate.MetaData().AggregateId()]) != 1 {
		t.Errorf("Should have one aggregate with id")
	}

	fetchedAggregate, err := eventStore.Get(aggregate.MetaData().AggregateId(), aggregate.MetaData().AggregateType())
	if err != nil {
		t.Errorf("Should not have thrown error %e", err)
	}

	if !es.Equals(fetchedAggregate.MetaData(), aggregate.MetaData()) {
		t.Errorf("Fetched aggregate metadata %+v should equal stored aggregate metadata %+v", fetchedAggregate.MetaData(), aggregate.MetaData())
	}
}

func Test_DispatchingInMemoryEventStore_ShouldReturnErrorIfAggregateNotInStore(t *testing.T) {
	eventStore := &es.DispatchingInMemoryEventStore{
		Factories: map[string]es.AggregateFactory{
			"testAggregate": func(aggregateId, aggregateType string, events []es.AggregateChanged) es.Aggregate {
				aggregate := &testAggregate{}
				aggregate.metaData = es.InitMetaData(aggregateId, aggregateType)
				for _, event := range events {
					aggregate.Apply(event)
				}
				return aggregate
			},
		},
		Aggregates: map[string]map[string][]es.AggregateChanged{},
		MessageBus: message.NewSimpleMessageBus([]message.Handler{
			&testAggregateCreatedHandler{},
		}),
	}

	fetchedAggregate, err := eventStore.Get("some-id", "some-type")

	if err == nil || err.Error() != fmt.Sprintf("could not construct aggregate %s of type %s, no events found in store", "some-id", "some-type") {
		t.Errorf("should have returned error")
	}
	if fetchedAggregate != nil {
		t.Errorf("Should not have returned error")
	}

}

func Test_DispatchingInMemoryEventStore_ShouldReturnErrorIfNoFactoryFoundForAggregate(t *testing.T) {
	eventStore := &es.DispatchingInMemoryEventStore{
		Factories:  map[string]es.AggregateFactory{},
		Aggregates: map[string]map[string][]es.AggregateChanged{},
		MessageBus: message.NewSimpleMessageBus([]message.Handler{
			&testAggregateCreatedHandler{},
		}),
	}

	aggregate := NewTestAggregate()

	err := eventStore.Store(aggregate)

	fetchedAggregate, err := eventStore.Get(aggregate.MetaData().AggregateId(), aggregate.MetaData().AggregateType())

	if err == nil || err.Error() != fmt.Sprintf("could not construct aggregate %s of type %s, no factory found", aggregate.MetaData().AggregateId(), aggregate.MetaData().AggregateType()) {
		t.Errorf("should have returned error %s", err.Error())
	}
	if fetchedAggregate != nil {
		t.Errorf("Should not have returned error")
	}

}

func Test_DispatchingInMemoryEventStore_ShouldReturnErrorIfNoHandlerFound(t *testing.T) {
	eventStore := &es.DispatchingInMemoryEventStore{
		Factories:  map[string]es.AggregateFactory{},
		Aggregates: map[string]map[string][]es.AggregateChanged{},
		MessageBus: message.NewSimpleMessageBus([]message.Handler{}),
	}

	aggregate := NewTestAggregate()

	err := eventStore.Store(aggregate)

	if err == nil || err.Error() != "No handler found for message testAggregate-testAggregateCreated" {
		t.Errorf("should have returned error %s", err.Error())
	}
}
