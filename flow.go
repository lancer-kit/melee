package melee

import (
	"sync"
	"testing"
)

type CtxKey string

const (
	KeyClient CtxKey = "client"
)

type Step struct {
	Name string
	Func func(f *Flow)
}

func (s *Step) run(f *Flow) {
	s.Func(f)
}

type Flow struct {
	*testing.T
	kv    sync.Map
	steps []Step
	Name  string
}

// AddTests adds test collection to flow
func (f *Flow) AddSteps(steps ...Step) *Flow {
	for _, step := range steps {
		f.steps = append(f.steps, step)
	}
	return f
}

func (f *Flow) WithValue(key CtxKey, value interface{}) *Flow {
	f.kv.Store(key, value)
	return f
}

func (f *Flow) GetValue(key CtxKey) (interface{}, bool) {
	return f.kv.Load(key)
}

func (f *Flow) Run(t *testing.T) {
	f.T = t
	f.Logf("Starting flow: %s\n", f.Name)
	for _, step := range f.steps {
		f.Logf("Starting step: %s\n", step.Name)
		step.run(f)
	}
}
