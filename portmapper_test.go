package portmapper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	// IN ALPHABETICAL ORDER
	validservices = []*Service{
		&Service{Name: "serviceA", Port: 1},
		&Service{Name: "serviceB", Port: 499},
		&Service{Name: "serviceC", Port: 8888},
		&Service{Name: "serviceD", Port: 65535},
	}
	invalidservices = []*Service{
		&Service{Name: "serviceA", Port: 0},
		&Service{Name: "serviceB", Port: 65536},
		&Service{Name: "serviceC", Port: -1},
		//&Service{Name: "serviceC", Port: 23, Hostname: *(*string)(nil)},
	}
)

func Test_RegisterValidServices(t *testing.T) {
	// register these validservices
	for i := 0; i < len(validservices); i++ {
		if Register(validservices[i].Name, validservices[i].Port) != nil {
			t.Error("error registering %s on port %v", validservices[i].Name, validservices[i].Port)
		} else {
			t.Log("successfully registered %s on port %v", validservices[i].Name, validservices[i].Port)
		}
	}

	// verify that we have registered these services by getting them from etcd
	if resultantservices, err := Services(); err == nil {
		if resultantservices != nil {
			assert.Equal(t, len(validservices), len(resultantservices))

			for i := 0; i < len(resultantservices); i++ {
				t.Log("Found service %s on port %v", resultantservices[i].Name, resultantservices[i].Port)

				assert.Equal(t, resultantservices[i].Name, validservices[i].Name)
				assert.Equal(t, resultantservices[i].Port, validservices[i].Port)
			}
		} else {
			t.Log("Bad news.  Retrieved services list is nil")
		}
	} else {
		t.Error("Error retrieving services: %s", err)
	}

}

func Test_UnregisterValidServices(t *testing.T) {
	// register these validservices
	for i := 0; i < len(validservices); i++ {
		if Unregister(validservices[i].Name, validservices[i].Port) != nil {
			t.Error("error unregistering %s on port %v", validservices[i].Name, validservices[i].Port)
		} else {
			t.Log("successfully unregistered %s on port %v", validservices[i].Name, validservices[i].Port)
		}
	}
}

func Test_RegisterInvalidServices(t *testing.T) {
	// register these validservices
	for i := 0; i < len(invalidservices); i++ {
		if Register(invalidservices[i].Name, invalidservices[i].Port) != nil {
			t.Log("error registering INVALID SERVICE %s on port %v", invalidservices[i].Name, invalidservices[i].Port)
		} else {
			t.Error("successfully registered INVALID SERVICE %s on port %v", invalidservices[i].Name, invalidservices[i].Port)
		}
	}
}
