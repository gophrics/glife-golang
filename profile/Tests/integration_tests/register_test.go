package profile_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"

	common "../../../common"
	mongodb "../../../common/mongodb"
	test_utils "../../../common/test"
	profile "../../Service"
)

var testData []map[string]interface{}

func init() {
	mongodb.Addr = "mongodb://localhost:27017"
	mongodb.Username = "testuser"
	mongodb.Password = "testpwd"

	ClearDb()
	ReadData()

	exec.Command("/bin/sh", "initialize_test.sh")
}

func ReadData() {
	file, _ := os.Open("register_test_data.json")
	bytes, _ := ioutil.ReadAll(file)
	json.Unmarshal([]byte(bytes), &testData)
}

func ClearDb() {

}

func Test_RegisterUser(t *testing.T) {
	var req []common.User
	test_utils.MapDecode(testData, &req)

	jsonData, _ := json.Marshal(req[0])
	request := httptest.NewRequest("POST", "/api/v1/profile/register", strings.NewReader(string(jsonData)))
	responseWriter := httptest.NewRecorder()
	profile.RegisterUser(responseWriter, request)
	r := responseWriter.Result()
	body, err := ioutil.ReadAll(r.Body)

	test_utils.Assert(t, err == nil)

	var token common.Token
	json.Unmarshal([]byte(body), &token)
	test_utils.Assert(t, token.Token != "")
}
