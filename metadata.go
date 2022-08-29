package eventsourcing

type AggregateMetaData interface {
	AggregateId() string
	AggregateType() string
	AggregateVersion() uint
	PendingEvents() []AggregateChanged
	Record(aggregateChanged AggregateChanged)
	IncrementVersion()
	ClearPendingEvents()
}
type aggregateMetaData struct {
	id            string
	aggregateType string
	version       uint
	pendingEvents []AggregateChanged
}

func (md *aggregateMetaData) AggregateId() string {
	return md.id
}

func (md *aggregateMetaData) AggregateType() string {
	return md.aggregateType
}

func (md *aggregateMetaData) AggregateVersion() uint {
	return md.version
}

func (md *aggregateMetaData) PendingEvents() []AggregateChanged {
	return md.pendingEvents
}

func (md *aggregateMetaData) ClearPendingEvents() {
	md.pendingEvents = []AggregateChanged{}
}

func (md *aggregateMetaData) Record(aggregateChanged AggregateChanged) {
	md.pendingEvents = append(md.pendingEvents, aggregateChanged)
}
func (md *aggregateMetaData) IncrementVersion() {
	md.version++
}

func InitMetaData(aggregateId, aggregateType string) AggregateMetaData {
	return &aggregateMetaData{
		id:            aggregateId,
		aggregateType: aggregateType,
		version:       0,
		pendingEvents: []AggregateChanged{},
	}
}

func Equals(md AggregateMetaData, mdTwo AggregateMetaData) bool {
	return md.AggregateId() == mdTwo.AggregateId() &&
		md.AggregateType() == mdTwo.AggregateType() &&
		md.AggregateVersion() == mdTwo.AggregateVersion() &&
		len(md.PendingEvents()) == len(mdTwo.PendingEvents())
}
