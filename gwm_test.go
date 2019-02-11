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
	attribute, err := client.ReadAttribute("", "server-state")

	if err != nil {
		t.Error(err)
	} else if strings.Compare("running", attribute) != 0 {
		t.Error("Cannot read attribute")
	}

	attribute, err = client.ReadAttribute("/subsystem=undertow/server=default-server/host=default-host", "alias")

	if err != nil {
		t.Error(err)
	} else if strings.Compare("localhost", attribute) != 0 {
		t.Error("Cannot read attribute of default-host")
	}

	_, err = client.ReadAttribute("", "servr-stte")

	if err == nil {
		t.Error("This should fail")
	}
}
