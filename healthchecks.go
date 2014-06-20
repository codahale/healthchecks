// Package healthchecks provides functionality to add application-level
// healthchecks which evaluate the health of components at runtime and expose
// the results as expvars.
package healthchecks

import (
	"expvar"
	"fmt"
	"sync"
)

// A Healthcheck checks the functionality of a component of your application and
// returns and error if anything is wrong.
type Healthcheck func() error

// Add adds the given healthcheck to the set of known healthchecks.
func Add(name string, healthcheck Healthcheck) {
	healthchecksMutex.Lock()
	defer healthchecksMutex.Unlock()

	healthchecks[name] = healthcheck
}

func execAll() map[string]string {
	healthchecksMutex.Lock()
	defer healthchecksMutex.Unlock()

	var m sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(healthchecks))

	results := make(map[string]string, len(healthchecks))
	for name, healthcheck := range healthchecks {
		go func(name string, healthcheck Healthcheck) {
			defer wg.Done()

			res := exec(healthcheck)

			m.Lock()
			defer m.Unlock()

			results[name] = res
		}(name, healthcheck)
	}
	wg.Wait()
	return results
}

func exec(hc Healthcheck) (s string) {
	defer func() {
		err := recover()
		if err != nil {
			s = fmt.Sprintf("%v", err)
		}
	}()

	err := hc()
	if err != nil {
		return err.Error()
	}
	return "OK"
}

var (
	healthchecks      = make(map[string]Healthcheck)
	healthchecksMutex sync.Mutex
)

func init() {
	expvar.Publish("healthchecks", expvar.Func(func() interface{} {
		return execAll()
	}))
}
