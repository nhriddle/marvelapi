install:
	go get github.com/go-redis/redis/v8
	go get -u github.com/ilyakaznacheev/cleanenv
	go get -u github.com/davecgh/go-spew/spew
	go get -u github.com/gorilla/mux

swagger:
	GOMODULE=off swagger generate spec -o ./swagger.yml --scan-models
