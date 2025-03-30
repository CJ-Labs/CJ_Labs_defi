package crypto

import (
	"log"
	"testing"

	"github.com/jamie20241210/CJ_Labs_crypto_go/json"
)

func TestJson(t *testing.T) {
	jsonString := json.ToJSONString("dads")
	log.Println(jsonString)
	println("d11111", jsonString)
}
