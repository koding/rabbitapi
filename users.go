package rabbitapi

import (
	"encoding/json"
)

type User struct {
	Name         string `json:"name`
	PasswordHash string `json:"password_hash"`
	Password     string `json:"password"`
	Tags         string `json:"tags"`
}

// GetUsers() returns a list of all users.
func (r *Rabbit) GetUsers() ([]User, error) {
	body, err := r.doRequest("GET", "/api/users", nil)
	if err != nil {
		return nil, err
	}

	users := make([]User, 0)
	err = json.Unmarshal(body, &users)
	if err != nil {
		return nil, err
	}

	return users, nil

}

// GetUser returns an individual user.
func (r *Rabbit) GetUser(name string) (User, error) {
	body, err := r.doRequest("GET", "/api/users/"+name, nil)
	if err != nil {
		return User{}, err
	}

	user := User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// CreateUser creates a new user with the given password and tags.
// The tags key is mandatory (means you can give an empty string). tags is a
// comma-separated list of tags for the user. Currently recognised tags are
// "administrator", "monitoring" and "management" (please aware that tags
// should be in the form of "foo, bar").
func (r *Rabbit) CreateUser(name, password, tags string) error {
	user := &User{
		Password: password,
		Tags:     tags,
	}

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	_, err = r.doRequest("PUT", "/api/users/"+name, data)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes an individual user.
func (r *Rabbit) DeleteUser(name string) error {
	_, err := r.doRequest("DELETE", "/api/users/"+name, nil)
	if err != nil {
		return err
	}

	return nil
}

// GetUserPermissions returns a list of all permissions for a given user.
func (r *Rabbit) GetUserPermissions(name string) ([]Permission, error) {
	body, err := r.doRequest("GET", "/api/users/"+name+"/permissions", nil)
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
