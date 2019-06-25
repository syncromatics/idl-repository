BUILD_PATH := ./artifacts
LINUX_BUILD_PATH = $(BUILD_PATH)/linux/idl
LINUX_ARM_BUILD_PATH = $(BUILD_PATH)/arm/idl
WINDOWS_BUILD_PATH = $(BUILD_PATH)/windows/idl.exe
MAC_BUILD_PATH = $(BUILD_PATH)/darwin/idl

export VERSION=$(shell gogitver)
export COMMIT_HASH=$(shell git rev-parse HEAD)
export BUILD_DATE=$(shell date +%Y-%m-%dT%T%z)
export BUILD_FLAGS=\
	-X github.com/syncromatics/idl-repository/cmd/idl/cmd.version=$(VERSION) \
	-X github.com/syncromatics/idl-repository/cmd/idl/cmd.commit=$(COMMIT_HASH) \
	-X github.com/syncromatics/idl-repository/cmd/idl/cmd.date=$(BUILD_DATE)

documentation:
	go run internal/cobraDocs.go
	git add docs/idl docs/idl-repository
	git commit -m "(+semver: patch) Updated documentation" docs/idl docs/idl-repository || echo 'Ignoring "no changes" error'

build:
	mkdir -p artifacts/linux artifacts/arm artifacts/windows artifacts/darwin
	GOOS=linux GOARCH=amd64 go build -ldflags "$(BUILD_FLAGS)" -o $(LINUX_BUILD_PATH) ./cmd/idl
	GOOS=linux GOARCH=arm go build -ldflags "$(BUILD_FLAGS)" -o $(LINUX_ARM_BUILD_PATH) ./cmd/idl
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(BUILD_FLAGS)" -o $(MAC_BUILD_PATH) ./cmd/idl
	GOOS=windows GOARCH=amd64 go build -ldflags "$(BUILD_FLAGS)" -o $(WINDOWS_BUILD_PATH) ./cmd/idl

	docker build -f ./cmd/idl-repository/Dockerfile -t syncromatics/idl-repository:${VERSION} .

test: build
	go test ./... -race -coverprofile=coverage.txt -covermode=atomic

package: test
	cd $(BUILD_PATH)/darwin && tar -zcvf ../darwin.tar.gz *
	cd $(BUILD_PATH)/linux && tar -zcvf ../linux.tar.gz *
	cd $(BUILD_PATH)/arm && tar -zcvf ../arm.tar.gz *
	cd $(BUILD_PATH)/windows && zip -r ../windows.zip *
	rm -R $(BUILD_PATH)/darwin $(BUILD_PATH)/linux $(BUILD_PATH)/arm $(BUILD_PATH)/windows

publish:
	echo "$$DOCKER_PASSWORD" | docker login -u "$$DOCKER_USERNAME" --password-stdin
	docker push syncromatics/idl-repository:${VERSION}

	COMMIT_LOG=`git log -1 --format='%ci %H %s'`
	github-release upload \
		--owner=syncromatics \
		--repo=idl-repository \
		--tag="${VERSION}" \
		--name="${TRAVIS_BRANCH}" \
		--body="${COMMIT_LOG}" \
		"artifacts/linux.tar.gz" \
		"artifacts/arm.tar.gz" \
		"artifacts/windows.zip" \
		"artifacts/darwin.tar.gz"	
version:
	@echo $(VERSION)