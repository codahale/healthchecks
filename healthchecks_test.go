package healthchecks

import (
	"encoding/json"
	"expvar"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestHealthchecks(t *testing.T) {
	start := time.Now()
	v := expvar.Get("healthchecks").String()
	elapsed := time.Now().Sub(start)

	if elapsed.Seconds() > 1.2 {
		t.Errorf("Took %v, not running in parallel", elapsed)
	}

	var actual map[string]string
	if err := json.Unmarshal([]byte(v), &actual); err != nil {
		t.Fatal(err)
	}

	expected := map[string]string{
		"good":  "OK",
		"bad":   "bad thing happen",
		"ugly":  "well dang",
		"slow1": "OK",
		"slow2": "OK",
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Was %#v, but expected %#v", actual, expected)
	}
}

func init() {
	Add("good", func() error {
		return nil
	})

	Add("bad", func() error {
		return fmt.Errorf("bad thing happen")
	})

	Add("ugly", func() error {
		panic("well dang")
	})

	Add("slow1", func() error {
		time.Sleep(1 * time.Second)
		return nil
	})

	Add("slow2", func() error {
		time.Sleep(1 * time.Second)
		return nil
	})
}
