single:
	env GOOS=linux GOARCH=amd64 go build -o main single.go benchmark.capnp.go
	rsync -avz main root@$(ENDPOINT):/root

concurrent:
	env GOOS=linux GOARCH=amd64 go build -o main concurrent.go benchmark.capnp.go
	rsync -avz main root@$(ENDPOINT):/root

capnp:
	capnp compile -I$(GOPATH)/pkg/mod/zombiezen.com/go/capnproto2@v2.17.0+incompatible/std -ogo *.capnp