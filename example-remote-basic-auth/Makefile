.PHONY: build run
build:
	docker run --rm -v `pwd`:/go/src/go-filter -w /go/src/go-filter \
		-e GOPROXY=https://goproxy.cn \
		golang:1.19 \
		go build -v -o libgolang.so -buildmode=c-shared .

run:
	docker run --rm -v `pwd`/envoy.yaml:/etc/envoy/envoy.yaml \
		-v `pwd`/libgolang.so:/etc/envoy/libgolang.so \
		-p 10000:10000 \
		envoyproxy/envoy:contrib-dev \
		envoy -c /etc/envoy/envoy.yaml

run-server:
	go run remote-server/main.go

test:
	curl -s -I 'http://localhost:10000/'
	curl -s -I 'http://localhost:10000/' -H 'Authorization: basic Zm9vOmJhcg=='
