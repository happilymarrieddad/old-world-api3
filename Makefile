protoc:
	protoc -I pb/v1/ \
		--go_out=:pb \
		--go_opt=paths=source_relative \
    	--go-grpc_out=require_unimplemented_servers=false:pb \
		--go-grpc_opt=paths=source_relative \
		--gogrpcmock_out=:pb \
    	pb/v1/*.proto

protoc-js:
	mkdir -p ../frontend2/src/pb
	protoc -I pb/v1/ \
		--js_out=import_style=commonjs:../frontend2/src/pb \
		--grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:../frontend2/src/pb \
		pb/v1/*.proto

install:
	go get -u \
		google.golang.org/protobuf \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc \
		github.com/gogo/protobuf/protoc-gen-gogoslick \
		github.com/gogo/protobuf/gogoproto \
		github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-gogrpcmock
	go install github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-gogrpcmock@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	npm install protoc-gen-js protoc-gen-grpc-web protoc-gen-ts -g --save

clean:
	rm ./pb/**/*.pb.go

create.keys:
	openssl genrsa -out keys/app.rsa 4096
	openssl rsa -in keys/app.rsa -pubout -outform PEM -out keys/app.rsa.pub

install.grpcwebproxy:
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	GOPATH=~/go ; export GOPATH
	git clone https://github.com/improbable-eng/grpc-web.git $GOPATH/src/github.com/improbable-eng/grpc-web
	cd $GOPATH/src/github.com/improbable-eng/grpc-web
	dep ensure # after installing dep
	go install ./go/grpcwebproxy # installs into $GOPATH/bin/grpcwebproxy

start-proxy:
	grpcwebproxy \
		--backend_addr=localhost:50051 \
		--run_tls_server=false \
		--allow_all_origins
