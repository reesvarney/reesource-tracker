package sync

import (
	"io"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

var (
	syncClients   = make(map[string]map[string]chan Event)
	syncClientsMu sync.Mutex
)

func BroadcastEvent(evtType string, data interface{}) {
	syncClientsMu.Lock()
	defer syncClientsMu.Unlock()
	for _, clientMap := range syncClients {
		for _, ch := range clientMap {
			select {
			case ch <- Event{Type: evtType, Data: data}:
			default:
				// drop if client is slow
			}
		}
	}
}

func EventStream(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")
	c.Header("X-Accel-Buffering", "no")

	// For this app, we use a single global channel key
	model_id := "global"
	clientID := c.ClientIP() + uuid.New().String()
	eventChan := make(chan Event, 10)

	syncClientsMu.Lock()
	if _, ok := syncClients[model_id]; !ok {
		syncClients[model_id] = make(map[string]chan Event)
	}
	syncClients[model_id][clientID] = eventChan
	syncClientsMu.Unlock()

	defer func() {
		syncClientsMu.Lock()
		delete(syncClients[model_id], clientID)
		if len(syncClients[model_id]) == 0 {
			delete(syncClients, model_id)
		}
		syncClientsMu.Unlock()
		close(eventChan)
	}()

	c.SSEvent("info", "Connected")
	c.Writer.Flush()

	c.Stream(func(w io.Writer) bool {
		select {
		case evt := <-eventChan:
			c.SSEvent(evt.Type, evt.Data)
			c.Writer.Flush()
			return true
		case <-c.Request.Context().Done():
			return false
		}
	})
}

// RegisterSyncRoutes adds the /eventstream endpoint to the router group
func Routes(route *gin.RouterGroup) {
	route.GET("/sync", EventStream)
}
