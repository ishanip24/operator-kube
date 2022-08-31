package health

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/go-logr/logr"
)

type HealthCheck interface {
	// Check implements healthz.Checker
	Check(*http.Request) error
	// Trigger update last run time
	Trigger()
}

type healthCheck struct {
	// clock clock.PassiveClock
	clock     time.Time
	tolerance time.Duration
	log       logr.Logger
	lastTime  atomic.Value
}

var _ HealthCheck = (*healthCheck)(nil)

// NewHealthCheck creates a new HealthCheck that will calculate the time inactive
// based on the provided clock and configuration.
func NewHealthCheck(clock time.Time,
	tolerance time.Duration,
	log logr.Logger,
) HealthCheck {
	answer := &healthCheck{
		clock:     clock,
		tolerance: tolerance,
		log:       log.WithName("HealthCheck"),
	}
	answer.Trigger()
	return answer
}

func (a *healthCheck) Trigger() {
	// a.lastTime.Store(a.clock.Now())
	a.lastTime.Store(a.clock)
}

func (a *healthCheck) Check(_ *http.Request) error {
	lastActivity := a.lastTime.Load().(time.Time)
	if a.clock.After(lastActivity.Add(a.tolerance)) {
		err := fmt.Errorf("last activity more than %s ago (%s)",
			a.tolerance, lastActivity.Format(time.RFC3339))
		a.log.Error(err, "Failing activity health check")
		return err
	}
	return nil
}
