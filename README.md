> go build -o weatherman cmd/main.go # go 1.22 required

> ./weatherman &

> curl http://localhost:8080/conditions?lat=<Latitude>&lon=<Longitude>&appid=<API KEY>

You may also run the application with debug logging enabled

> LOG=Debug ./weatherman
