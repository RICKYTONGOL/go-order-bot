// utilservice package implements functions that process the matched order list.
package utilservice

import (
	"fmt"
	"sync"
)

type UtilManager struct {
	botServer UtilInterface
	wg        *sync.WaitGroup
}

func newUtilManager() *UtilManager {
	if log == nil {
		panic(fmt.Errorf("logger not available"))
	}
	if cfg == nil {
		panic(fmt.Errorf("configuration not available"))
	}

	utilMgr := UtilManager{
		wg: new(sync.WaitGroup),
	}
	utilMgr.init()

	return &utilMgr
}

func (mgr *UtilManager) init() {
	mgr.botServer = &OrderBot{
		sigStop: make(chan bool, 1),
		mgr:     mgr,
	}

	mgr.botServer.init()
}

func (mgr *UtilManager) started(svc UtilInterface) {
	mgr.wg.Add(1)
	log.Noticef("utilservice server is running")
}

func (mgr *UtilManager) Run() {
	log.Noticef("utilservice server was started")

	go mgr.botServer.start()
}

func (mgr *UtilManager) Close() {
	log.Notice("utilservice manager received a close signal")

	// pass the signal to all the services
	log.Noticef("closing order bot server")
	mgr.botServer.close()

	// wait scanners to terminate
	log.Notice("waiting for util services to finish")
	mgr.wg.Wait()

	// we are done
	log.Notice("util service manager closed")
}
