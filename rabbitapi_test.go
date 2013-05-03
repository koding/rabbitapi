package rabbitapi

import (
	"testing"
)

func TestRabbit_AlivenessTest(t *testing.T) {
	r := Auth("guest", "guest", "http://localhost:15672")
	err := r.AlivenessTest("/")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("vhost '/' is ok!")
	}
}

func TestRabbit_OverviewTest(t *testing.T) {
	r := Auth("guest", "guest", "http://localhost:15672")
	overview, err := r.GetOverview()
	if err != nil {
		t.Error(err)
	} else {
		t.Log("overview struct is", overview)
	}
}
