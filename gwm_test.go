package gwm_test

import (
	"testing"

	"github.com/ilpianista/gwm"
)

var client *gwm.GWMClient

func init() {
	client = gwm.NewClient("localhost", 9990, "admin", "password")
}

func Test_GWM_ReadAttribute(t *testing.T) {
	client.ReadAttribute("server-name")
}
