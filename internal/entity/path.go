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
		r, err := s.ExecuteRequest()
		if err != nil {
			p.failedSteps = p.steps[i:len(p.steps)]
			return errors.New(fmt.Sprintf("Could not execute the step #%d. Error: %s", i, err.Error()))
		}

		if err2 := p.resolvePlaceholders(s, r, i); err2 != nil {
			return err2
		}

		for placeholder, value := range s.Placeholders {
			if s.next != nil {
				s.next.Placeholders[placeholder] = value
			}
		}

	}
	return nil
}

func (p *Path) resolvePlaceholders(s *Step, r []byte, i int) error {
	if len(s.PlaceholderNameToPath) > 0 {
		for placeholder, path := range s.PlaceholderNameToPath {
			value, err := AccessJsonByPath(r, path.Split())
			if err != nil {
				return errors.New(fmt.Sprintf("Could not access json by bath in the step #%d. Error: %s", i, err.Error()))
			}

			s.Placeholders[placeholder] = value
		}
	}
	return nil
}
