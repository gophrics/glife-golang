package travel

import (
	"crypto/md5"
	"encoding/json"
)

func Hash(arr Trip) [16]byte {
	jsonBytes, _ := json.Marshal(arr)
	return md5.Sum(jsonBytes)
}
