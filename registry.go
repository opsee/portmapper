package pomapper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
  // The path to a directory where the registry will store data.
  RegistryPath string
)

func init() {
  RegistryPath = "/zuul/state/bastion"
}

// Service is a mapping between a service name and port.
type Service struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

func (s *Service) validate() error {
	if s.Name == "" {
		return fmt.Errorf("Service lacks Name field: %v", s)
	}
	if s.Port < 1 || s.Port > 65535 {
		return fmt.Errorf("Service Port is outside valid range: %v", s)
	}

	return nil
}

func (s *Service) path() string {
	return fmt.Sprintf("%s/%s.service", RegistryPath, s.Name)
}

// Marshal returns the byte array of the JSON-serialized version of the
// service.
func (s *Service) Marshal() ([]byte, error) {
	bytes, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// Given a byte array, unmarshal it.
func UnmarshalService(bytes []byte) (*Service, error) {
	s := &Service{}
	err := json.Unmarshal(bytes, s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Unregister a (service, port) tuple.
func Unregister(name string, port int) error {
	svc := &Service{name, port}
	if err := svc.validate(); err != nil {
		return err
	}

	err := os.Remove(svc.path())
	if err != nil {
		return err
	}

	return nil
}

// Register a (service, port) tuple.
func Register(name string, port int) error {
	svc := &Service{name, port}
	if err := svc.validate(); err != nil {
		return err
	}

	bytes, err := svc.Marshal()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(svc.path(), bytes, 0644)
}

// Services returns an array of Service pointers detailing the service name and
// port of each registered service.
func Services() ([]*Service, error) {
	svcPaths, err := filepath.Glob(fmt.Sprintf("%s/*.service", RegistryPath))
	if err != nil {
		return nil, err
	}

	services := make([]*Service, len(svcPaths))

	for i, svcPath := range svcPaths {
		bytes, err := ioutil.ReadFile(svcPath)
		if err != nil {
			return nil, err
		}

		svc, err := UnmarshalService(bytes)
		if err != nil {
			return nil, err
		}
		services[i] = svc
	}

	return services, nil
}
