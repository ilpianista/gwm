package gwm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	digestRequest "github.com/ilpianista/go-digest-request"
)

const HTTP = "http"
const APPLICATION_JSON = "application/json"

type GWMClient struct {
	host       string
	port       int
	username   string
	password   string
	HttpClient *http.Client
}

func NewClient(controller_host string, controller_port int, management_user string, management_user_pwd string) *GWMClient {
	client := &GWMClient{
		host:       controller_host,
		port:       controller_port,
		username:   management_user,
		password:   management_user_pwd,
		HttpClient: http.DefaultClient,
	}

	return client
}

func (c GWMClient) ReadAttribute(address string, attribute string) (string, error) {
	sep := func(r rune) bool {
		return r == '/' || r == '='
	}

	command := map[string]interface{}{
		"operation":   "read-attribute",
		"name":        attribute,
		"address":     strings.FieldsFunc(address, sep),
		"json.pretty": 1,
	}

	bytesRepresentation, err := json.Marshal(command)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	var response *http.Response
	if c.username == "" || c.password == "" {
		response, err = http.Post(HTTP+"://"+c.host+":"+strconv.Itoa(c.port)+"/management", APPLICATION_JSON, bytes.NewBuffer(bytesRepresentation))
	} else {
		r := digestRequest.New(context.Background(), c.username, c.password)
		request, err := http.NewRequest(http.MethodPost, HTTP+"://"+c.host+":"+strconv.Itoa(c.port)+"/management", bytes.NewBuffer(bytesRepresentation))
		request.Header.Add("Content-Type", APPLICATION_JSON)

		if err != nil {
			log.Fatal(err)
			return "", err
		}

		response, err = r.Do(request)
	}

	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode == 401 {
		return "", errors.New("Unauthorized")
	} else {
		var result map[string]interface{}

		if err := json.Unmarshal(body, &result); err != nil {
			log.Fatal(err)
			return "", err
		} else {
			if strings.Compare("success", result["outcome"].(string)) == 0 {
				switch result["result"].(type) {
				case string:
					return result["result"].(string), nil
				default:
					return result["result"].([]interface{})[0].(string), nil
				}
			} else {
				return "", errors.New(result["failure-description"].(string))
			}
		}
	}
}
