package entity

type Steps = []*Step

type Path struct {
	steps Steps
}

func (p *Path) AddStep(s *Step) {
	if len(p.steps) > 0 {
		lastStep := p.steps[len(p.steps)-1]
		lastStep.next = s
	}
	p.steps = append(p.steps, s)
}

func (p *Path) Execute() {
	for _, s := range p.steps {
		response, err := s.ExecuteRequest()
	}
}
