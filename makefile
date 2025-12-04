

release:
	go install github.com/mitchellh/gox@latest
	go mod tidy
	go mod vendor
	gox  -output="bin/{{.OS}}/{{.Arch}}/{{.Dir}}"  -osarch="linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64"

release-linux:
	go install github.com/mitchellh/gox@latest
	go mod tidy
	go mod vendor
	gox  -output="bin/{{.OS}}/{{.Arch}}/{{.Dir}}"  -osarch="linux/amd64"



