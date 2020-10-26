package entity

import (
	"errors"
	"fmt"
)

type Steps = []*Step

type Path struct {
	steps       Steps
	failedSteps Steps
}

func (p *Path) AddStep(s *Step) {
	if len(p.steps) > 0 {
		lastStep := p.steps[len(p.steps)-1]
		lastStep.next = s
	}
	p.steps = append(p.steps, s)
}

func (p *Path) Execute() error {
	for i, s := range p.steps {
		_, err := s.ExecuteRequest()
		if err != nil {
			p.failedSteps = p.steps[i:len(p.steps)]
			return errors.New(fmt.Sprintf("Could not execute the step #%d. Error: %s", i, err.Error()))
		}
	}
	return nil
}
