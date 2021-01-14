// Copyright 2014 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.
//
//
//
// Copyright 2017-2019 Lei Ni (nilei81@gmail.com) and other Dragonboat authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package syncutil

import (
	"log"
	"sync"

	"github.com/lni/goutils/envutil"
	"github.com/lni/goutils/lang"
)

// Stopper is a manager struct for managing worker goroutines. It is modified
// from an early version of the stopper struct found in CockroachDB's codebase.
type Stopper struct {
	mu          sync.Mutex
	shouldStopC chan struct{}
	wg          sync.WaitGroup
	debug       bool
}

// NewStopper return a new Stopper instance.
func NewStopper() *Stopper {
	s := &Stopper{
		shouldStopC: make(chan struct{}),
		debug:       envutil.GetBoolEnvVarOrDefault("LEAKTEST", false),
	}

	return s
}

// RunWorker creates a new goroutine and invoke the f func in that new
// worker goroutine.
func (s *Stopper) RunWorker(f func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.stopped() {
		return
	}
	s.runWorker(f, "")
}

func (s *Stopper) runWorker(f func(), name string) {
	s.wg.Add(1)
	var gid uint64
	go func() {
		if s.debug {
			gid = lang.GetGIDForDebugOnly()
			log.Printf("goroutine %d started, name %s", gid, name)
		}
		f()
		s.wg.Done()
		if s.debug {
			log.Printf("goroutine %d stopped, name %s", gid, name)
		}
	}()
}

// ShouldStop returns a chan struct{} used for indicating whether the
// Stop() function has been called on Stopper.
func (s *Stopper) ShouldStop() chan struct{} {
	return s.shouldStopC
}

// Stop signals all managed worker goroutines to stop and wait for them
// to actually stop.
func (s *Stopper) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	close(s.shouldStopC)
	s.wg.Wait()
}

// Close closes the internal shouldStopc chan struct{} to signal all
// worker goroutines that they should stop.
func (s *Stopper) Close() {
	close(s.shouldStopC)
}

func (s *Stopper) stopped() bool {
	select {
	case <-s.shouldStopC:
		return true
	default:
	}
	return false
}
