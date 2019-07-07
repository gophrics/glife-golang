package travel

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
)

func Hash(arr Trip) [16]byte {
	jsonBytes, _ := json.Marshal(arr)
	fmt.Printf("\n%v\n\n", arr)
	return md5.Sum(jsonBytes)
}
