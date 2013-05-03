package rabbitapi

import (
	"encoding/json"
)

type Exchange struct {
	Arguments  map[string]interface{} `json:"arguments"`
	AutoDelete bool                   `json:"auto_delete"`
	Durable    bool                   `json:"durable"`
	Internal   bool                   `json:"durable"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Vhost      string                 `json:"vhost"`
}

type ExchangeSource struct {
	Arguments       map[string]interface{} `json:"arguments"`
	Destination     string                 `json:"destination"`
	DestinationType string                 `json:"destination_type"`
	PropertiesKey   string                 `json:"properties_key"`
	RoutingKey      string                 `json:"routing_key"`
	Source          string                 `json:"source"`
	Vhost           string                 `json:"vhost"`
}

// GetExchanges() returns a list of all exchanges.
func (r *Rabbit) GetExchanges() ([]Exchange, error) {
	body, err := r.doRequest("GET", "/api/exchanges", nil)
	if err != nil {
		return nil, err
	}

	exchanges := make([]Exchange, 0)
	err = json.Unmarshal(body, &exchanges)
	if err != nil {
		return nil, err
	}

	return exchanges, nil
}

// GetVhostExchanges returns a list of all exchanges in a given virtual host.
func (r *Rabbit) GetVhostExchanges(vhost string) ([]Exchange, error) {
	if vhost == "/" {
		vhost = "%2f"
	}

	body, err := r.doRequest("GET", "/api/exchanges/"+vhost, nil)
	if err != nil {
		return nil, err
	}

	exchanges := make([]Exchange, 0)
	err = json.Unmarshal(body, &exchanges)
	if err != nil {
		return nil, err
	}

	return exchanges, nil
}

// GetExchange returns an individual exchange for the given vhost and name.
func (r *Rabbit) GetExchange(vhost, name string) (Exchange, error) {
	if vhost == "/" {
		vhost = "%2f"
	}

	body, err := r.doRequest("GET", "/api/exchanges/"+vhost+"/"+name, nil)
	if err != nil {
		return Exchange{}, err
	}

	exchange := Exchange{}
	err = json.Unmarshal(body, &exchange)
	if err != nil {
		return Exchange{}, err
	}

	return exchange, nil
}

// CreateExchange creates an invididual exchange with for the given vhost and name.
func (r *Rabbit) CreateExchange(vhost, name, kind string, durable, autoDelete, internal bool, args map[string]interface{}) error {
	if vhost == "/" {
		vhost = "%2f"
	}

	if args == nil {
		args = make(map[string]interface{}, 0)
	}

	exchange := &Exchange{
		Type:       kind,
		Durable:    durable,
		AutoDelete: autoDelete,
		Internal:   internal,
		Arguments:  args,
	}

	data, err := json.Marshal(exchange)
	if err != nil {
		return err
	}

	_, err = r.doRequest("PUT", "/api/exchanges/"+vhost+"/"+name, data)
	if err != nil {
		return err
	}

	return nil
}

// DeleteExchange deletes an individual exchange for the given vhost and name.
func (r *Rabbit) DeleteExchange(vhost, name string) error {
	if vhost == "/" {
		vhost = "%2f"
	}

	_, err := r.doRequest("DELETE", "/api/exchanges/"+vhost+"/"+name, nil)
	if err != nil {
		return err
	}

	return nil
}

// GetExchangeSource returns a list of all bindings in which a given exchange
// is the source.
func (r *Rabbit) GetExchangeSource(vhost, name string) ([]ExchangeSource, error) {
	if vhost == "/" {
		vhost = "%2f"
	}

	body, err := r.doRequest("GET", "/api/exchanges/"+vhost+"/"+name+"/bindings/source", nil)

	if err != nil {
		return nil, err
	}

	exchangeSources := make([]ExchangeSource, 0)
	err = json.Unmarshal(body, &exchangeSources)
	if err != nil {
		return nil, err
	}

	return exchangeSources, nil
}
