schemas:
	$(shell protoc -I pajbot/ pajbot/pajbot.proto --go_out=plugins=grpc:pajbot)
