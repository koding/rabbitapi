/*
Implementation of RabbitMq Management HTTP Api in Go

Currently supported api calls are:

    GET     /api/vhosts
    GET     /api/vhost/name
    PUT     /api/vhost/name
    DELETE  /api/vhost/name
    GET     /api/vhost/name/permissions

    GET     /api/users
    GET     /api/users/name
    PUT     /api/users/name
    DELETE  /api/users/name
    GET     /api/users/name/permissions

    GET     /api/permissions
    GET     /api/permissions/vhost/user
    PUT     /api/permissions/vhost/user
    DELETE  /api/permissions/vhost/user

Example code:

	r := rabbitapi.Auth("guest", "guest", "http://localhost:15672")
	vhosts, err := r.GetVhosts()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("vhosts:", vhosts)
	}


*/
package rabbitapi
