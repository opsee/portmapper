package portmapper

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/coreos/go-etcd/etcd"
)

var (
	// RegistryPath sets the location in etcd where pomapper will store data.
	// Default: /opsee.co/portmapper
	RegistryPath string
	// EtcdHost is the IP_ADDRESS:PORT location of Etcd.
	// Default: 127.0.0.1:2379
	EtcdHost   string
	etcdClient *etcd.Client
)

func init() {
	RegistryPath = "/opsee.co/portmapper"
	EtcdHost = "127.0.0.1:2379"
}

// Service is a mapping between a service name and port. It may also contain
// the hostname where the service is running or the container ID in the
// Hostname field. It will attempt to get this from the HOSTNAME environment
// variable.
type Service struct {
	Name     string `json:"name"`
	Port     int    `json:"port"`
	Hostname string `json:"hostname,omitempty"`
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
	return fmt.Sprintf("%s/%s:%d", RegistryPath, s.Name, s.Port)
}

// Marshal a service object to a byte array.
func (s *Service) Marshal() ([]byte, error) {
	bytes, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// UnmarshalService deserializes a Service object from a byte array.
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
	svc := &Service{name, port, os.Getenv("HOSTNAME")}
	if err := svc.validate(); err != nil {
		return err
	}

	if _, err := etcdClient.Delete(svc.path(), false); err != nil {
		return err
	}

	return nil
}

// Register a (service, port) tuple.
func Register(name string, port int) error {
	svc := &Service{name, port, os.Getenv("HOSTNAME")}
	if err := svc.validate(); err != nil {
		return err
	}

	bytes, err := svc.Marshal()
	if err != nil {
		return err
	}

	if _, err = etcdClient.Set(svc.path(), string(bytes), 0); err != nil {
		return err
	}

	return nil
}

// Services returns an array of Service pointers detailing the service name and
// port of each registered service.
func Services() ([]*Service, error) {
	resp, err := etcdClient.Get(RegistryPath, false, false)
	if err != nil {
		return nil, err
	}

	svcNodes := resp.Node.Nodes
	services := make([]*Service, len(svcNodes))

	for i, node := range svcNodes {
		svcStr := node.Value
		svc, err := UnmarshalService([]byte(svcStr))
		if err != nil {
			return nil, err
		}
		services[i] = svc
	}

	return services, nil
}
