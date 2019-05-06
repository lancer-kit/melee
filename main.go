package melee

import (
	"sync"
	"testing"
	"time"

	"github.com/lancer-kit/melee/config"
	"github.com/lancer-kit/melee/metrics"
)

// NewFlow returns flow with default fields
func NewFlow(name string) *Flow {
	return &Flow{
		Name:  name,
		steps: []Step{},
		kv:    sync.Map{},
	}
}

func NewStep(name string, fun func(f *Flow)) Step {
	return Step{
		Name: name,
		Func: fun,
	}
}

func RunTestFlows(t *testing.T, testFlows []TestFlow) {
	m := metrics.Metrics{}.
		Init(1, 1)

	TestGroup{
		T:         t,
		Metrics:   &m,
		LoadMode:  false,
		Users:     1,
		TestFlows: testFlows,
	}.RunAll()
}

func RunLoadTestGroup(t *testing.T, testFuncs []TestFlow, await4start bool) {
	var wg sync.WaitGroup
	cfg := config.Config()
	if await4start && !wait4Signal(cfg.NATS) {
		return
	}
	m := metrics.Metrics{}.
		Init(cfg.LoadTestConfig.Threads, cfg.LoadTestConfig.Users)

	start := time.Now()
	for i := 0; i < cfg.LoadTestConfig.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			TestGroup{
				T:         t,
				Metrics:   &m,
				LoadMode:  true,
				Users:     cfg.LoadTestConfig.Users,
				TestFlows: testFuncs,
			}.RunAll()
		}()

	}

	wg.Wait()
	m.FullFlowTime = time.Since(start)
	m.ObtainMetrics()

}
