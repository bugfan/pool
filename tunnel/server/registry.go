package server

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"pool/log"
)

const (
	cacheSaveInterval time.Duration = 10 * time.Minute
)

type cacheUrl string

func (url cacheUrl) Size() int {
	return len(url)
}

// ControlRegistry maps a client ID to Control structures
type ControlRegistry struct {
	controls map[string]*Control
	log.Logger
	sync.RWMutex
}

func NewControlRegistry() *ControlRegistry {
	return &ControlRegistry{
		controls: make(map[string]*Control),
		Logger:   log.NewPrefixLogger("registry", "ctl"),
	}
}

func (r *ControlRegistry) Get(clientId string) *Control {
	r.RLock()
	defer r.RUnlock()
	return r.controls[clientId]
}

func (r *ControlRegistry) Add(clientId string, ctl *Control) (oldCtl *Control) {
	r.Lock()
	defer r.Unlock()

	oldCtl = r.controls[clientId]
	if oldCtl != nil {
		oldCtl.Replaced(ctl)
	}

	r.controls[clientId] = ctl
	r.Debug("Registered control with id %s", clientId)
	return
}

func (r *ControlRegistry) Del(clientId string) error {
	r.Lock()
	defer r.Unlock()
	if r.controls[clientId] == nil {
		return fmt.Errorf("No control found for client id: %s", clientId)
	} else {
		r.Debug("Removed control registry id %s", clientId)
		delete(r.controls, clientId)
		return nil
	}
}

func (r *ControlRegistry) IDs() (ids map[string]string) {
	ids = make(map[string]string)
	for _, c := range r.controls {
		ip := strings.Split(c.conn.RemoteAddr().String(), ":")[0]
		ids[ip] = c.id
	}
	return
}
