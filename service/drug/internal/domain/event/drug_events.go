package event

import "time"

type DrugEvent struct {
	eventType string
	drugID    int64
	operatorID int64
	data      map[string]interface{}
	occurredAt time.Time
}

func NewDrugEvent(eventType string, drugID, operatorID int64, data map[string]interface{}) *DrugEvent {
	return &DrugEvent{
		eventType:  eventType,
		drugID:     drugID,
		operatorID: operatorID,
		data:       data,
		occurredAt: time.Now(),
	}
}

func (e *DrugEvent) EventType() string { return e.eventType }
func (e *DrugEvent) DrugID() int64 { return e.drugID }
func (e *DrugEvent) OperatorID() int64 { return e.operatorID }
func (e *DrugEvent) Data() map[string]interface{} { return e.data }
func (e *DrugEvent) OccurredAt() time.Time { return e.occurredAt }

const (
	EventTypeDrugCreated    = "drug.created"
	EventTypeStockUpdated   = "drug.stock_updated"
	EventTypeStockAdded     = "drug.stock_added"
	EventTypeStockReduced   = "drug.stock_reduced"
	EventTypeStatusChanged  = "drug.status_changed"
)

func DrugCreatedEvent(drugID int64, name string) *DrugEvent {
	return NewDrugEvent(EventTypeDrugCreated, drugID, 0, map[string]interface{}{
		"name": name,
	})
}

func StockUpdatedEvent(drugID, operatorID int64, quantity int) *DrugEvent {
	return NewDrugEvent(EventTypeStockUpdated, drugID, operatorID, map[string]interface{}{
		"quantity": quantity,
	})
}

func StockAddedEvent(drugID, operatorID int64, quantity int) *DrugEvent {
	return NewDrugEvent(EventTypeStockAdded, drugID, operatorID, map[string]interface{}{
		"quantity": quantity,
	})
}

func StockReducedEvent(drugID, operatorID int64, quantity int) *DrugEvent {
	return NewDrugEvent(EventTypeStockReduced, drugID, operatorID, map[string]interface{}{
		"quantity": quantity,
	})
}