package melee

import (
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/lancer-kit/melee/metrics"
)

type TestGroup struct {
	T         *testing.T
	Metrics   *metrics.Metrics
	LoadMode  bool
	Users     int
	TestFlows []TestFlow
}

type TestFlow struct {
	Percent  int
	TestFlow func() *Flow
}

func (lt TestGroup) RunAll() {
	//client := GetXClient()
	for i := 0; i < lt.Users; i++ {
		start := time.Now()
		//lt.RunTests(client)
		finish := time.Since(start)
		lt.Metrics.TimeExecution = append(lt.Metrics.TimeExecution, finish)
		lt.Metrics.UsersTotal++
	}
}

func (lt TestGroup) RunTests() {
	var (
		funcName string
		finish   time.Duration
	)

	//n := rand.Intn(100)

	for _, testFunc := range lt.TestFlows {
		if lt.LoadMode && testFunc.Percent == 0 {
			continue
		}
		//todo: fix percent usage
		start := time.Now()

		testFlow := testFunc.TestFlow()
		//testFlow.Client.SetClient(client)
		testFlow.Run(lt.T)

		funcName = runtime.FuncForPC(reflect.ValueOf(testFunc.TestFlow).Pointer()).Name()
		finish = time.Since(start)

		lt.Metrics.Statuses = append(lt.Metrics.Statuses, lt.T.Failed())
		lt.Metrics.TimeExecutionByScripts[funcName] = append(lt.Metrics.TimeExecutionByScripts[funcName], finish)
		lt.Metrics.QuantityScripts[funcName]++
		if lt.T.Failed() == true {
			lt.Metrics.FailsByScripts[funcName]++
		}
		//break
		//n -= percent
	}
}
