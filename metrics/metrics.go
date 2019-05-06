package metrics

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gonum.org/v1/gonum/floats"
)

type Metrics struct {
	Threads    int `json:"threads"`
	Users      int `json:"users"`
	UsersTotal int `json:"users_total"`

	SuccessStatus  bool           `json:"success"`
	Statuses       []bool         `json:"-"`
	Fails          int            `json:"failed"`
	FailsPerc      float32        `json:"fails_percent"`
	FailsByScripts map[string]int `json:"fails_by_scripts"`

	TimeOfAll            float64 `json:"time"`
	Timeout              int     `json:"timeout"`
	MaxTimeExecution     float64 `json:"max_execution"`
	MinTimeExecution     float64 `json:"min_execution"`
	AverageTimeExecution float64 `json:"avg_execution"`

	FullFlowTime  time.Duration   `json:"-"`
	TimeExecution []time.Duration `json:"-"`

	TimeExecutionByRequest   map[string][]time.Duration `json:"-"`
	TimeExecutionByScripts   map[string][]time.Duration `json:"-"`
	MaxTimeExecutionByScript map[string]float64         `json:"max_time_execution_by_script"`
	MinTimeExecutionByScript map[string]float64         `json:"min_time_execution_by_script"`
	AvgTimeExecutionByScript map[string]float64         `json:"avg_time_execution_by_script"`
	QuantityScripts          map[string]int             `json:"quantity_scripts"`
}

const SEC = float64(time.Second)

func (m Metrics) Init(threads, users int) Metrics {
	m.Threads = threads
	m.Users = users
	m.TimeExecutionByScripts = make(map[string][]time.Duration)
	m.QuantityScripts = make(map[string]int)
	m.FailsByScripts = make(map[string]int)
	return m

}
func (m *Metrics) ObtainMetrics() {
	m.obtainTime()
	m.obtainStatus()
	m.prepareFile()
}

func (m *Metrics) obtainStatus() {
	for _, status := range m.Statuses {
		if status { //t.Failed == true
			m.Fails++
			m.SuccessStatus = false
			continue
		}
		m.SuccessStatus = true
	}
	if m.UsersTotal > 0 {
		m.FailsPerc = float32(m.Fails) / float32(m.UsersTotal) * 100
	}
}

func (m *Metrics) obtainTime() {
	m.MaxTimeExecutionByScript = make(map[string]float64)
	m.MinTimeExecutionByScript = make(map[string]float64)
	m.AvgTimeExecutionByScript = make(map[string]float64)

	max, min, sum := m.checkScopes(m.TimeExecution)

	for key, value := range m.TimeExecutionByScripts {
		maxS, minS, sumS := m.checkScopes(value)
		m.MaxTimeExecutionByScript[key] = floats.Round(float64(maxS)/SEC, 4)
		m.MinTimeExecutionByScript[key] = floats.Round(float64(minS)/SEC, 4)
		m.AvgTimeExecutionByScript[key] = floats.Round(float64(sumS)/float64(len(value))/SEC, 4)
	}

	m.MaxTimeExecution = floats.Round(float64(max)/SEC, 4)
	m.MinTimeExecution = floats.Round(float64(min)/SEC, 4)
	m.AverageTimeExecution = floats.Round(float64(sum)/float64(len(m.TimeExecution))/SEC, 4)
}

func (m *Metrics) checkScopes(f []time.Duration) (max, min, sum int64) {
	min = 1000
	for _, t := range f {
		if t.Nanoseconds() > max {
			max = t.Nanoseconds()
		}
		if t.Nanoseconds() < min {
			min = t.Nanoseconds()
		}
		sum += t.Nanoseconds()
	}
	return
}

func (m *Metrics) prepareFile() {
	m.TimeOfAll = floats.Round(m.FullFlowTime.Seconds(), 3)
	newFile, _ := os.Create(fmt.Sprintf("metrics.%d.json", time.Now().Unix()))
	metricsJson, _ := json.MarshalIndent(m, "", "  ")
	newFile.Write(metricsJson)
}
