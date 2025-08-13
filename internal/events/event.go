package events

import (
	"encoding/json"
	"time"
)

// Event represents a domain event
type Event struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	AggregateID string                 `json:"aggregate_id"`
	Version     int                    `json:"version"`
	Data        map[string]interface{} `json:"data"`
	Metadata    map[string]interface{} `json:"metadata"`
	Timestamp   time.Time              `json:"timestamp"`
}

// EventStore interface for storing events
type EventStore interface {
	Append(aggregateID string, events ...Event) error
	GetEvents(aggregateID string) ([]Event, error)
	GetEventsByType(eventType string) ([]Event, error)
}

// EventBus interface for publishing events
type EventBus interface {
	Publish(event Event) error
	Subscribe(eventType string, handler EventHandler) error
}

// EventHandler function type for handling events
type EventHandler func(event Event) error

// NewEvent creates a new event
func NewEvent(eventType, aggregateID string, data map[string]interface{}) Event {
	println("🎯 Yeni event oluşturuluyor:", eventType, "aggregate:", aggregateID)

	event := Event{
		ID:          generateEventID(),
		Type:        eventType,
		AggregateID: aggregateID,
		Version:     1,
		Data:        data,
		Metadata:    make(map[string]interface{}),
		Timestamp:   time.Now(),
	}

	println("✅ Event oluşturuldu, ID:", event.ID)
	return event
}

// generateEventID generates a unique event ID
func generateEventID() string {
	println("🆔 Event ID oluşturuluyor...")
	id := time.Now().Format("20060102150405") + "-" + randomString(8)
	println("🆔 Event ID oluşturuldu:", id)
	return id
}

// randomString generates a random string of given length
func randomString(length int) string {
	println("🎲 Random string oluşturuluyor, uzunluk:", length)
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	result := string(b)
	println("🎲 Random string oluşturuldu:", result)
	return result
}

// Serialize converts event to JSON bytes
func (e Event) Serialize() ([]byte, error) {
	println("📦 Event serialize ediliyor:", e.ID)
	data, err := json.Marshal(e)
	if err != nil {
		println("❌ Event serialize hatası:", err.Error())
		return nil, err
	}
	println("✅ Event serialize edildi, boyut:", len(data), "bytes")
	return data, nil
}

// DeserializeEvent converts JSON bytes to Event
func DeserializeEvent(data []byte) (Event, error) {
	println("📦 Event deserialize ediliyor, boyut:", len(data), "bytes")
	var event Event
	err := json.Unmarshal(data, &event)
	if err != nil {
		println("❌ Event deserialize hatası:", err.Error())
		return Event{}, err
	}
	println("✅ Event deserialize edildi, ID:", event.ID)
	return event, nil
}
