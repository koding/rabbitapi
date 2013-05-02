# rabbitapi

Implementation of [RabbitMq Management HTTP
Api](http://hg.rabbitmq.com/rabbitmq-management/raw-file/rabbitmq_v3_1_0/priv/www/api/index.html)
in Go. Alpha status.

For more information and documenation please read [Godoc RabbitApi
page](http://godoc.org/github.com/koding/rabbitapi)

# setup

```
go get github.com/koding/rabbitapi
```

# example usage

First create a rabbitapi instance with your api credentials

```
r := rabbitapi.Auth("guest", "guest", "http://localhost:15672")
```

To get a list of all vhosts

```
vhosts, err := r.GetVhosts()
if err != nil {
	fmt.Println(err)
} else {
	fmt.Println("vhosts:", vhosts)
}
```

Create an exchange on vhost `/` with the name `rabbitapi`, `durable=false`,
`autoDelete=true`, `internal=false,` and `arguments=nil`

```
err = r.CreateExchange("/", "rabbitapi", "topic", false, true, false, nil)
if err != nil {
	fmt.Println(err)
}
```

Get an exchange we created previously on the vhost `/`

```
exchange, err := r.GetExchange("/", "rabbitapi")
if err != nil {
	fmt.Println(err)
}
fmt.Println(exchange) // exchange.Type is 'topic'
```

for more examples look into `*_test.go` files.

