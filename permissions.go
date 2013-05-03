package rabbitapi

import (
	"encoding/json"
)

type Permission struct {
	Configure string `json:"configure"`
	Read      string `json:"read"`
	User      string `json:"user"`
	Vhost     string `json:"vhost"`
	Write     string `json:"write"`
}

// GetPermissions returns a list of all permissions for all users.
func (r *Rabbit) GetPermissions() ([]Permission, error) {
	body, err := r.doRequest("GET", "/api/permissions", nil)
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

// GetPermissions returns an individual permission of a user and virtual host
func (r *Rabbit) GetPermission(vhost, user string) (Permission, error) {
	if vhost == "/" {
		vhost = "%2f"
	}

	body, err := r.doRequest("GET", "/api/permissions/"+vhost+"/"+user, nil)
	if err != nil {
		return Permission{}, err
	}

	permission := Permission{}
	err = json.Unmarshal(body, &permission)
	if err != nil {
		return Permission{}, err
	}

	return permission, nil

}

// CreatePermission creates the necessery configure, write and read permissions
// for the the given vhost and user. For more info please look at:
// http://www.rabbitmq.com/access-control.html
func (r *Rabbit) CreatePermission(vhost, user, configure, write, read string) error {
	if vhost == "/" {
		vhost = "%2f"
	}

	permission := &Permission{
		Configure: configure,
		Write:     write,
		Read:      read,
	}

	data, err := json.Marshal(permission)
	if err != nil {
		return err
	}

	_, err = r.doRequest("PUT", "/api/permissions/"+vhost+"/"+user, data)
	if err != nil {
		return err
	}

	return nil
}

// DeletePermission deletes the permission for the given vhost and user
func (r *Rabbit) DeletePermission(vhost, user string) error {
	if vhost == "/" {
		vhost = "%2f"
	}

	_, err := r.doRequest("DELETE", "/api/permissions/"+vhost+"/"+user, nil)
	if err != nil {
		return err
	}

	return nil

}
