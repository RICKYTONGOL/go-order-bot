// Package svc implements trading services.
package svc

import (
	"fmt"
	"sync"

	"github.com/dongle/go-order-bot/internal/repository"
)

// ServiceManager implements service manager.
type ServiceManager struct {
	wg *sync.WaitGroup
}

// newServiceManager creates a new instance of service manager.
func newServiceManager() *ServiceManager {
	// make sure we have what we need
	if log == nil {
		panic(fmt.Errorf("logger not available"))
	}
	if cfg == nil {
		panic(fmt.Errorf("configuration not available"))
	}

	sm := ServiceManager{
		wg: new(sync.WaitGroup),
	}

	sm.init()
	return &sm
}

// init the svc manager.
func (mgr *ServiceManager) init() {

}

// Run starts all the services prepared to be run.
func (mgr *ServiceManager) Run() {
	repo = repository.R()
}

// Close signals orchestrator to terminate all orchestrated services.
func (mgr *ServiceManager) Close() {

	// wait scanners to terminate
	log.Notice("waiting for services to finish")
	mgr.wg.Wait()

	log.Notice("svc manager closed")
}
