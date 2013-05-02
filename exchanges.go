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

// GetExchanges() returns a list of all exchanges.
func (r *Rabbit) GetExchanges() ([]Exchange, error) {
	body, err := r.getRequest("/api/exchanges")
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

	body, err := r.getRequest("/api/exchanges/" + vhost)
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

	body, err := r.getRequest("/api/exchanges/" + vhost + "/" + name)
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

	err = r.putRequest("/api/exchanges/"+vhost+"/"+name, data)
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

	err := r.deleteRequest("/api/exchanges/" + vhost + "/" + name)
	if err != nil {
		return err
	}

	return nil
}
