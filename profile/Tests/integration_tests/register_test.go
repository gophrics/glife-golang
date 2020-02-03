package profile_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	common "../../../common"
	mongodb "../../../common/mongodb"
	test_utils "../../../common/test"
	profile "../../Service"
)

var testData map[string]interface{}

// mongo --eval 'db.createUser({user: "testuser", pwd:"testpwd", roles:[{role:"userAdminAnyDatabase", db: "admin"}]})' admin

func init() {
	mongodb.Addr = "mongodb://localhost:27017"
	mongodb.Username = "testuser"
	mongodb.Password = "testpwd"

	ClearDb()
	ReadData()
}

func ReadData() {
	file, _ := os.Open("data.json")
	bytes, _ := ioutil.ReadAll(file)
	json.Unmarshal([]byte(bytes), &testData)
}

func ClearDb() {

}

func Test_RegisterUser(t *testing.T) {
	var req []common.User
	test_utils.MapDecode(testData["register"], &req)

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
