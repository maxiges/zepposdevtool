

release:
	go install github.com/mitchellh/gox@latest
	go mod tidy
	go mod vendor
	gox  -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" -os="linux darwin windows"