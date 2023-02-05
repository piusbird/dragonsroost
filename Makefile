prod: clean
	go build --tags prod
dev: clean
	go build --tags dev
clean:
	rm -f dragonsroost
