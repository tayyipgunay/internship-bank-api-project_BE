package events

import (
	"fmt"
	"sync"
)

// InMemoryEventBus implements EventBus using in-memory storage
type InMemoryEventBus struct {
	handlers map[string][]EventHandler
	mutex    sync.RWMutex
}

// NewInMemoryEventBus creates a new in-memory event bus
func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Publish publishes an event to all subscribers
func (eb *InMemoryEventBus) Publish(event Event) error {
	println("📢 Event yayınlanıyor, tip:", event.Type, "ID:", event.ID)

	eb.mutex.RLock()
	handlers, exists := eb.handlers[event.Type]
	eb.mutex.RUnlock()

	if !exists {
		println("ℹ️ Bu event tipi için handler bulunamadı:", event.Type)
		return nil // No handlers for this event type
	}

	println("📋", len(handlers), "handler bulundu, event işleniyor...")

	for i, handler := range handlers {
		println("🔧 Handler", i+1, "çalıştırılıyor...")
		if err := handler(event); err != nil {
			println("❌ Handler", i+1, "hatası:", err.Error())
			return fmt.Errorf("handler failed for event %s: %w", event.Type, err)
		}
		println("✅ Handler", i+1, "başarıyla tamamlandı")
	}

	println("✅ Event başarıyla yayınlandı, tip:", event.Type)
	return nil
}

// Subscribe adds a handler for a specific event type
func (eb *InMemoryEventBus) Subscribe(eventType string, handler EventHandler) error {
	println("📝 Event handler kaydediliyor, tip:", eventType)

	if eventType == "" {
		println("❌ Event tipi boş")
		return fmt.Errorf("event type cannot be empty")
	}

	if handler == nil {
		println("❌ Handler nil")
		return fmt.Errorf("handler cannot be nil")
	}

	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	if eb.handlers[eventType] == nil {
		eb.handlers[eventType] = make([]EventHandler, 0)
	}

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
	println("✅ Event handler kaydedildi, tip:", eventType, "toplam handler:", len(eb.handlers[eventType]))
	return nil
}

// Unsubscribe removes a handler for a specific event type
func (eb *InMemoryEventBus) Unsubscribe(eventType string, handler EventHandler) error {
	println("🗑️ Event handler kaldırılıyor, tip:", eventType)

	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	handlers, exists := eb.handlers[eventType]
	if !exists {
		println("ℹ️ Bu event tipi için handler bulunamadı:", eventType)
		return nil
	}

	for i, h := range handlers {
		if &h == &handler {
			eb.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			println("✅ Event handler kaldırıldı, tip:", eventType)
			break
		}
	}

	return nil
}
