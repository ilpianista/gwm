package gwm_test

import (
	"strings"
	"testing"

	"github.com/ilpianista/gwm"
)

var client *gwm.GWMClient

func init() {
	client = gwm.NewClient("localhost", 9990, "admin", "password")
}

func Test_GWM_ReadAttribute(t *testing.T) {
	if strings.Compare("running", client.ReadAttribute("server-state")) != 0 {
		t.Error("Cannot read attribute")
	}
}
