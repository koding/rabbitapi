package rabbitapi

import (
	"testing"
)

func TestRabbit_GetExchanges(t *testing.T) {
	r := Auth("guest", "guest", "http://localhost:15672")
	exchanges, err := r.GetExchanges()
	if err != nil {
		t.Error(err)
	} else {
		t.Log("exchanges:", exchanges)
	}
}

func TestRabbit_GetVhostExchanges(t *testing.T) {
	r := Auth("guest", "guest", "http://localhost:15672")
	exchanges, err := r.GetVhostExchanges("/")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("exchanges:", exchanges)
	}
}

func TestRabbit_CreateExchange(t *testing.T) {
	r := Auth("guest", "guest", "http://localhost:15672")
	err := r.CreateExchange("/", "rabbitapi", "topic", false, true, false, nil)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("exchange with name 'rabbitapi', type 'topic' is created successfull")
	}
}

func TestRabbit_GetExchange(t *testing.T) {
	r := Auth("guest", "guest", "http://localhost:15672")
	exchange, err := r.GetExchange("/", "rabbitapi")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("exchange 'rabbitapi':", exchange)
	}
}

func TestRabbit_DeleteExchange(t *testing.T) {
	r := Auth("guest", "guest", "http://localhost:15672")
	err := r.DeleteExchange("/", "rabbitapi")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("exchange 'rabbitapi' is deleted successfully")
	}

	exchange, err := r.GetExchange("/", "rabbitapi")
	if err != nil {
		t.Log("retriving exchange 'rabbitapi'")
		t.Log(err)
	} else {
		t.Error("exchange 'rabbitapi':", exchange)
	}
}
