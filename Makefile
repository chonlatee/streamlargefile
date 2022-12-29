server:
	go run main.go

client1:
	go run client/main.go -addr=:3333 -dirName=3333

client2:
	go run client/main.go -addr=:3334 -dirName=3334


clean:
	rm -rf ./client/3333
	rm -rf ./client/3334