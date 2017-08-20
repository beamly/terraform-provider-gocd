package gocd

import (
	"encoding/json"
	"errors"
)

// StagesService exposes calls for interacting with Stage objects in the GoCD API.
type StagesService service

// Stage represents a GoCD Stage object.
type Stage struct {
	Name                  string    `json:"name"`
	FetchMaterials        bool      `json:"fetch_materials"`
	CleanWorkingDirectory bool      `json:"clean_working_directory"`
	NeverCleanupArtifacts bool      `json:"never_cleanup_artifacts"`
	Approval              *Approval `json:"approval,omitempty"`
	EnvironmentVariables  []string  `json:"environment_variables,omitempty"`
	Resources             []string  `json:"resource,omitempty"`
	Jobs                  []*Job    `json:"jobs,omitempty"`
}

// JSONString returns a string of this stage as a JSON object.
func (s *Stage) JSONString() (string, error) {
	err := s.Validate()
	if err != nil {
		return "", err
	}
	s.Clean()
	bdy, err := json.MarshalIndent(s, "", "  ")
	return string(bdy), err
}

// Validate ensures the attributes attached to this structure are ready for submission to the GoCD API.
func (s *Stage) Validate() error {
	if s.Name == "" {
		return errors.New("`gocd.Stage.Name` is empty")
	}

	if len(s.Jobs) == 0 {
		return errors.New("At least one `Job` must be specified")
	}

	for _, job := range s.Jobs {
		if err := job.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Clean the approvel step.
func (s *Stage) Clean() {
	s.Approval.Clean()
}
