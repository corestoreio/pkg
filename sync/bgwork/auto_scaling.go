package bgwork

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type stateID uint

const (
	eventWorkerStateNone stateID = iota
	eventWorkerStateInit
	eventWorkerStateWIP
	eventWorkerStateIdle
	eventWorkerStateTerminate
	eventWorkerStateTerminated
	eventWorkerStateMax
)

func (s stateID) String() string {
	ret := "unknown"
	switch s {
	case eventWorkerStateNone:
		ret = "none"
	case eventWorkerStateInit:
		ret = "init"
	case eventWorkerStateWIP:
		ret = "wip"
	case eventWorkerStateIdle:
		ret = "idle"
	case eventWorkerStateTerminate:
		ret = "terminate"
	case eventWorkerStateTerminated:
		ret = "terminated"
	}
	return ret
}

type internalState struct {
	die   chan struct{} // signals termination of a worker goroutine
	state stateID
}

type statePool struct {
	minWorkers uint16

	mu         sync.RWMutex
	state      map[uint]internalState
	statistics [eventWorkerStateMax]int
}

func newWrkrState(minWorkers, maxWorkers uint16) *statePool {
	return &statePool{
		minWorkers: minWorkers,
		state:      make(map[uint]internalState, maxWorkers),
	}
}

func (ws *statePool) init(wrkr uint) chan struct{} {
	is := internalState{
		die:   make(chan struct{}),
		state: eventWorkerStateInit,
	}
	ws.mu.Lock()
	ws.state[wrkr] = is
	ws.statistics[eventWorkerStateInit]++
	ws.mu.Unlock()
	return is.die
}

func (ws *statePool) terminateWaiting() {
	_, allCount := ws.countState(eventWorkerStateTerminate)
	ws.mu.Lock()
	for wrkrID, is := range ws.state {
		if is.state == eventWorkerStateTerminate {
			if allCount > ws.minWorkers {
				close(is.die)
				ws.statistics[eventWorkerStateTerminated]++
				// remove entry from map to free space and dont let the map grow
				// unlimited.
				delete(ws.state, wrkrID)
				allCount--
			} else {
				is.state = eventWorkerStateIdle
				ws.state[wrkrID] = is
			}
		}
	}
	ws.mu.Unlock()
}

func (ws *statePool) set(wrkr uint, state stateID) {
	ws.mu.Lock()
	is := ws.state[wrkr]
	is.state = state
	ws.statistics[state]++
	ws.state[wrkr] = is
	ws.mu.Unlock()
}

func (ws *statePool) transitionState(oldState, newState stateID) {
	ws.mu.Lock()
	for i, is := range ws.state {
		if is.state == oldState {
			is.state = newState
			ws.statistics[newState]++
			ws.state[i] = is
		}
	}
	ws.mu.Unlock()
}

func (ws *statePool) countState(state stateID) (stateCount uint16, allCount uint16) {
	ws.mu.RLock()
	for _, s := range ws.state {
		if s.state == state {
			stateCount++
		}
		allCount++
	}
	ws.mu.RUnlock()
	return stateCount, allCount
}

func (ws *statePool) printStat() string {
	var buf strings.Builder
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	fmt.Fprintf(&buf, "Total Init %d\n", ws.statistics[eventWorkerStateInit])
	fmt.Fprintf(&buf, "Total WIP %d\n", ws.statistics[eventWorkerStateWIP])
	fmt.Fprintf(&buf, "Total Waiting %d\n", ws.statistics[eventWorkerStateIdle])
	fmt.Fprintf(&buf, "Total Terminate %d\n", ws.statistics[eventWorkerStateTerminate])
	fmt.Fprintf(&buf, "Total Terminated %d\n", ws.statistics[eventWorkerStateTerminated])
	workerIDs := make([]int, 0, len(ws.state))
	for wrkrID := range ws.state {
		workerIDs = append(workerIDs, int(wrkrID))
	}
	sort.Ints(workerIDs)
	for _, wrkerID := range workerIDs {
		s := ws.state[uint(wrkerID)].state
		if s != eventWorkerStateNone {
			fmt.Fprintf(&buf, "WrkrID:%d State: %s\n", wrkerID, s)
		}
	}
	return buf.String()
}

// ScalingOptions sets various configurations to AutoScaling function.
type ScalingOptions struct {
	// MinWorkers defaults to 2
	MinWorkers uint16
	// MaxWorkers defaults to 16
	MaxWorkers uint16
	// WorkerCheckInterval defaults to one second. Duration when to check for
	// termination and/or creation of workers. Each interval the GetStatistics
	// function will be called.
	WorkerCheckInterval time.Duration
	// GetStatistics gets called each WorkerCheckInterval
	GetStatistics func(string)
}

// AutoScaling runs the same task across min to max background workers. Idle
// workers are getting terminated until the min amount of workers get reached.
// There won't be more workers than max amount. AutoScaling blocks once called
// so it must start in its own goroutine.
func AutoScaling(ctx context.Context, chanEvents <-chan interface{}, jobFn func(event interface{}), opt ScalingOptions) {
	if opt.MinWorkers == 0 {
		opt.MinWorkers = 2
	}
	if opt.MaxWorkers == 0 {
		opt.MaxWorkers = 16
	}
	if opt.WorkerCheckInterval < 1 {
		opt.WorkerCheckInterval = time.Second
	}

	wrkrMngr := newWrkrState(opt.MinWorkers, opt.MaxWorkers)

	// wrkrFn runs in a goroutine
	wrkrFn := func(wrkrID uint, jobFn func(interface{})) {
		chanDie := wrkrMngr.init(wrkrID)
		wrkrMngr.set(wrkrID, eventWorkerStateIdle)
		for {
			select {
			case <-chanDie:
				return
			case event := <-chanEvents:
				wrkrMngr.set(wrkrID, eventWorkerStateWIP)
				jobFn(event)
				wrkrMngr.set(wrkrID, eventWorkerStateIdle)
			case <-ctx.Done(): // stops all of them
				wrkrMngr.set(wrkrID, eventWorkerStateTerminated)
				return
			}
		}
	}

	tkr := time.NewTicker(opt.WorkerCheckInterval)
	defer tkr.Stop()
	var i uint
	for {
		wip, _ := wrkrMngr.countState(eventWorkerStateWIP)
		wait, allWrkrs := wrkrMngr.countState(eventWorkerStateIdle)
		switch {
		case allWrkrs < opt.MinWorkers:
			go wrkrFn(i, jobFn)
			i++
		case wait > opt.MinWorkers:
			wrkrMngr.transitionState(eventWorkerStateIdle, eventWorkerStateTerminate)
			wrkrMngr.terminateWaiting()
		case wip >= opt.MinWorkers && wip < opt.MaxWorkers:
			go wrkrFn(i, jobFn)
			i++
		}
		select {
		case <-ctx.Done():
			return
		case <-tkr.C: // sleeps the for loop
			if opt.GetStatistics != nil {
				opt.GetStatistics(wrkrMngr.printStat())
			}
		}
	}
}
