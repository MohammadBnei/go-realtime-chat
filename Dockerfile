FROM golang:1.19

RUN apt-get update && apt-get install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

RUN export PATH="$PATH:$(go env GOPATH)/bin"

CMD protoc -I /gen/proto --proto_path=/gen/proto message.proto \
    --go-grpc_out /out/go/ --go_out /out/go/