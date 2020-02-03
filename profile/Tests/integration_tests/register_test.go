package profile_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"

	common "../../../common"
	profile "../../Service"
	"github.com/mitchellh/mapstructure"
	"gotest.tools/assert"
)

var testData map[string]interface{}

func ReadData() {
	file, _ := os.Open("data.json")
	bytes, _ := ioutil.ReadAll(file)
	json.Unmarshal([]byte(bytes), &testData)
}

func Test_RegisterUser(t *testing.T) {

	ReadData()

	var req []common.User
	mapstructure.Decode(testData["register"], &req)

	request := httptest.NewRequest("POST", "/api/v1/profile/register", bytes.NewReader([]byte(fmt.Sprintf("%v", req[0]))))
	responseWriter := httptest.NewRecorder()
	profile.RegisterUser(responseWriter, request)
	r := responseWriter.Result()
	body, err := ioutil.ReadAll(r.Body)

	assert.Assert(t, err != nil)
	assert.Equal(t, body, "abc")
}
