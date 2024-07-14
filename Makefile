build:
	@go build -o bin/gitlab-mr-reminder .

run: build
	@./bin/gitlab-mr-reminder
