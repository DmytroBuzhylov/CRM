build:
	go build cmd/api/main.go

generate-notification-proto:
	protoc \
    	  --go_out=. \
    	  internal/feature/notification/interface_adapters/proto/notification.proto