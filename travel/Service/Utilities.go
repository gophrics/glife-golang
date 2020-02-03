package travel

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

func Hash(arr Trip) [16]byte {
	jsonBytes, _ := json.Marshal(arr)
	fmt.Printf("\n%v\n\n", arr)
	return md5.Sum(jsonBytes)
}

func getBodyFromHttpRequest(r *http.Request, req interface{}) (interface{}, string) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err.Error()
	}
	json.Unmarshal(b, &req)
	return req, nil
}