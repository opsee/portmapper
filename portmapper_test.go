package portmapper

import (
	"testing"
)

func Test_RegisterService(t *testing.T) {
	if Register("ServiceA", 4000) != nil {
		t.Error("error registering service")
	} else {
		t.Log("successfully registered service")
	}
}
