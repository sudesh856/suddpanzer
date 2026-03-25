package scenario
import (
	"fmt"
	"gopkg.in/yaml.v3"
	"time"
	"os"

)

type Scenario struct {
	Name string `yaml:"name"`
	Stages []Stage `yaml:"stages"`
	Endpoints []Endpoint `yaml:"endpoints"`

}

type Stage struct {
	Duration string `yaml:"duration"`
	TargetVUs int `yaml:"target_vus"`
	ParsedDuration time.Duration `yaml:"-"`
}

type Endpoint struct {
	URL string `yaml:"url"`
	Method string `yaml:"method"`
	Headers map[string]string `yaml:"headers"`
	Body string `yaml:"body"`
	Weight int `yaml:"weight"`
	ExpectedStatus int `yaml:"expected_status"`
}

func LoadScenario(path string) (*Scenario, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Cannot read scenario file %q: %w", path,err)
	}

	var s Scenario
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil,fmt.Errorf("Invalid YAML in %q: %w", path, err)
	}

	if err := s.Validate(); err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *Scenario) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("Scenario is missing a 'name' field.")
	}
	if len(s.Stages) == 0 {
		return fmt.Errorf("Scenario %q  has no stages", s.Name)
	}

	for i, stage := range s.Stages {
		d, err := time.ParseDuration(stage.Duration)
		if err != nil {
			return fmt.Errorf("Stage[%d] invalid duration %q: %w", i, stage.Duration, err)
		}
		s.Stages[i].ParsedDuration = d
	}

	if len(s.Endpoints) == 0 {
		return fmt.Errorf("Scenario %q has no endpoints.", s.Name)
	}
	for i, ep := range s.Endpoints {
		if ep.URL == "" {
			return fmt.Errorf("Endpoint [%d] missing url.", i)
		}
		
		if ep.Weight <= 0 {
			return fmt.Errorf("Endpoint [%d] weight must be > 0.", i)
		}
		if ep.Method == "" {
			s.Endpoints[i].Method = "GET"
		}
	}
	return nil
}