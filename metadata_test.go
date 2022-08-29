package eventsourcing_test

import (
	"testing"

	es "github.com/maxkaemmerer/go-event-sourcing"
)

func Test_InitMetaData_ShouldInitialize(t *testing.T) {
	metaData := es.InitMetaData("id", "type")

	if metaData.AggregateId() != "id" {
		t.Errorf("Should have correct id but has %s", metaData.AggregateId())
	}

	if metaData.AggregateType() != "type" {
		t.Errorf("Should have correct type but has %s", metaData.AggregateType())
	}

	if metaData.AggregateVersion() != 0 {
		t.Errorf("Should have initialized with version 0 but got %d", metaData.AggregateVersion())
	}

	if len(metaData.PendingEvents()) != 0 {
		t.Errorf("Should have no pending events but got %d", metaData.AggregateVersion())
	}
}

func Test_InitMetaData_ShouldIncrementVersion(t *testing.T) {
	metaData := es.InitMetaData("id", "type")

	if metaData.AggregateVersion() != 0 {
		t.Errorf("Should have initialized with version 0 but got %d", metaData.AggregateVersion())
	}
	metaData.IncrementVersion()

	if metaData.AggregateVersion() != 1 {
		t.Errorf("Should have incremented version to 1 %d", metaData.AggregateVersion())
	}
}
func Test_InitMetaData_ShouldRecordEvent(t *testing.T) {
	metaData := es.InitMetaData("id", "type")

	if len(metaData.PendingEvents()) != 0 {
		t.Errorf("Should have no pending events but got %d", metaData.AggregateVersion())
	}

	metaData.Record(&testAggregateCreated{})

	if len(metaData.PendingEvents()) != 1 {
		t.Errorf("Should have recorded 1 event but has %d", metaData.AggregateVersion())
	}
}

func Test_InitMetaData_ShouldClearPendingEvents(t *testing.T) {
	metaData := es.InitMetaData("id", "type")

	if len(metaData.PendingEvents()) != 0 {
		t.Errorf("Should have no pending events but got %d", metaData.AggregateVersion())
	}

	metaData.Record(&testAggregateCreated{})

	if len(metaData.PendingEvents()) != 1 {
		t.Errorf("Should have recorded 1 event but has %d", metaData.AggregateVersion())
	}

	metaData.ClearPendingEvents()

	if len(metaData.PendingEvents()) != 0 {
		t.Errorf("Should have cleared events but has %d", metaData.AggregateVersion())
	}
}
