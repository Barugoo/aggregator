port = $(shell echo $(SERVER_ADDR) | sed -e 's,^.*:,:,g' -e 's,.*:\([0-9]*\).*,\1,g' -e 's,[^0-9],,g')

run: .check-env-vars
	docker build  --build-arg PORT=${port} -t server .
	docker run  -p ${port}:${port} -e SERVER_ADDR=${SERVER_ADDR} server

test: 
	go test ./...

.check-env-vars:
	@test $${SERVER_ADDR?Please set environment variable SERVER_ADDR}