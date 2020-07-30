# ui

## Project setup

This project requires Go and yarn to be installed.

### Development

Terminal 1
``` shell
go run . serve -p 8081 testdata/example1.forensicstore
```

Terminal 2
```
yarn serve
```

Open http://localhost:8080/

### Compile for production
``` shell
go generate # includes yarn build
go build .
```
