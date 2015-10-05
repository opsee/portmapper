package portmapper

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"os"
	"sync"
	"time"
)

var (
	// RegistryPath sets the location in etcd where pomapper will store data.
	// Default: /opsee.co/portmapper
	RegistryPath string
	// EtcdHost is the IP_ADDRESS:PORT location of Etcd.
	// Default: http://127.0.0.1:2379
	EtcdHost   string
	etcdClient *etcd.Client
	ServiceMap map[string]*ServiceRegistration
	once       sync.Once
)

func init() {
	RegistryPath = "/opsee.co/portmapper"
	EtcdHost = "http://127.0.0.1:2379"
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

type ServiceRegistration struct {
	Service   *Service
	Timestamp int64
	err       error
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
	etcdClient = etcd.NewClient([]string{EtcdHost})

	svc := &Service{name, port, os.Getenv("HOSTNAME")}
	if err := svc.validate(); err != nil {
		return err
	}

	if _, err := etcdClient.Delete(svc.path(), false); err != nil {
		return err
	}

	return nil
}

// Loops over each service type in the map and attempt to register
func RegisterServices() {
	for {
		wg := &sync.WaitGroup{}
		for service_name := range ServiceMap {
			wg.Add(1)
			go func() {
				defer wg.Done()
				println("registering: " + service_name)
				etcdClient = etcd.NewClient([]string{EtcdHost})
				svc := ServiceMap[service_name].Service
				ServiceMap[service_name].Timestamp = time.Now().Unix()
				ServiceMap[service_name].err = nil

				if err := svc.validate(); err != nil {
					ServiceMap[service_name].err = err
				}

				bytes, err := svc.Marshal()
				if err != nil {
					ServiceMap[service_name].err = err
				}

				if _, err = etcdClient.Set(svc.path(), string(bytes), 0); err != nil {
					ServiceMap[service_name].err = err
				}
			}()
		}

		// Wait for all services to be registered
		wg.Wait()
		time.Sleep(60 * time.Second)
	}
}

// Register a (service, port) tuple.
func Register(name string, port int) error {
	svc := &Service{name, port, os.Getenv("HOSTNAME")}
	ServiceMap[name] = &ServiceRegistration{Service: svc, Timestamp: 0, err: nil}

	//XXX I assume that this will not exit ever
	once.Do(RegisterServices)
	return nil
}

// Services returns an array of Service pointers detailing the service name and
// port of each registered service. (from etcd)
func Services() ([]*Service, error) {
	etcdClient = etcd.NewClient([]string{EtcdHost})

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
