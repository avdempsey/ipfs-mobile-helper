build-prod:
	env GOOS=linux GOARCH=amd64 go build -o ipfs-mobile-helper