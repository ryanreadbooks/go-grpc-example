proto:
	protoc -Iproto proto/*.proto \
		--go_out=pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=pb \
		--go-grpc_opt=paths=source_relative

clean:
	rm pb/*.pb.go

.PHONY: proto clean