package test

import (
	"net/http"
	"testing"
)

func TestDemo(t *testing.T) {
	around(func(c client) {
		resp, _ := c.callRequest("POST", "/users", nil)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("status >> %s", resp.Status)
		}
		t.Log("fuck")
	})
}
