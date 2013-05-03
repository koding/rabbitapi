package rabbitapi

import (
	"encoding/json"
)

type Vhost struct {
	Name    string
	Tracing bool
}

// GetVhosts returns a list of all vhosts.
func (r *Rabbit) GetVhosts() ([]Vhost, error) {
	body, err := r.doRequest("GET", "/api/vhosts", nil)
	if err != nil {
		return nil, err
	}

	vhosts := make([]Vhost, 0)
	err = json.Unmarshal(body, &vhosts)
	if err != nil {
		return nil, err
	}

	return vhosts, nil
}

// GetVhost returns an individual vhost.
func (r *Rabbit) GetVhost(name string) (Vhost, error) {
	if name == "/" {
		name = "%2f"
	}

	body, err := r.doRequest("GET", "/api/vhosts/"+name, nil)
	if err != nil {
		return Vhost{}, err
	}

	vhost := Vhost{}
	err = json.Unmarshal(body, &vhost)
	if err != nil {
		return Vhost{}, err
	}

	return vhost, nil

}

// CreateVhost creates an invididual vhost.
func (r *Rabbit) CreateVhost(name string) error {
	if name == "/" {
		name = "%2f"
	}

	_, err := r.doRequest("PUT", "/api/vhosts/"+name, nil)
	if err != nil {
		return err
	}

	return nil
}

// DeleteVhost deletes an individual vhost.
func (r *Rabbit) DeleteVhost(name string) error {
	if name == "/" {
		name = "%2f"
	}

	_, err := r.doRequest("DELETE", "/api/vhosts/"+name, nil)
	if err != nil {
		return err
	}

	return nil
}

// GetVhostPermissions returns a list of all permissions for a given virtual
// host.
func (r *Rabbit) GetVhostPermissions(vhost string) ([]Permission, error) {
	if vhost == "/" {
		vhost = "%2f"
	}

	body, err := r.doRequest("GET", "/api/vhosts/"+vhost+"/permissions", nil)
	if err != nil {
		return nil, err
	}

	list := make([]Permission, 0)
	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}
