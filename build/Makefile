prod:
	go build -ldflags "-s -w" -o bin/github-webhook .

run:
	go build -o bin/github-webhook-dev.bin
	bin/github-webhook-dev.bin --config tmp/config.yml
