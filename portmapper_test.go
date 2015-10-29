package portmapper

import (
	"testing"
)

var (
	services = []*Service{
		&Service{Name: "serviceA", Port: 40, Hostname: "serviceA"},
		&Service{Name: "serviceB", Port: 41, Hostname: "serviceB"},
		&Service{Name: "serviceC", Port: 42, Hostname: "serviceC"},
		&Service{Name: "serviceD", Port: 43, Hostname: "serviceD"},
	}
)

func Test_RegisterServices(t *testing.T) {

	// register these services
	for i := 0; i < len(services); i++ {
		if Register(services[i].Name, services[i].Port) != nil {
			t.Error("error registering %s on port %v", services[i].Name, services[i].Port)
		} else {
			t.Log("successfully registered %s on port %v", services[i].Name, services[i].Port)
		}
	}
}

func Test_UnregisterServices(t *testing.T) {

	// register these services
	for i := 0; i < len(services); i++ {
		if Unregister(services[i].Name, services[i].Port) != nil {
			t.Error("error unregistering %s on port %v", services[i].Name, services[i].Port)
		} else {
			t.Log("successfully unregistered %s on port %v", services[i].Name, services[i].Port)
		}
	}
}
