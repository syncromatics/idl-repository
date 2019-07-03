BUILD_PATH := ./artifacts
LINUX_BUILD_PATH = $(BUILD_PATH)/linux/idl
LINUX_ARM_BUILD_PATH = $(BUILD_PATH)/arm/idl
WINDOWS_BUILD_PATH = $(BUILD_PATH)/windows/idl.exe
MAC_BUILD_PATH = $(BUILD_PATH)/darwin/idl

export VERSION=$(shell gogitver)

build:
	mkdir -p artifacts/linux artifacts/arm artifacts/windows artifacts/darwin
	GOOS=linux GOARCH=amd64 go build -o $(LINUX_BUILD_PATH) ./cmd/idl
	GOOS=linux GOARCH=arm go build -o $(LINUX_ARM_BUILD_PATH) ./cmd/idl
	GOOS=darwin GOARCH=amd64 go build -o $(MAC_BUILD_PATH) ./cmd/idl
	GOOS=windows GOARCH=amd64 go build -o $(WINDOWS_BUILD_PATH) ./cmd/idl

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
		--tag="latest" \
		--name="${TRAVIS_BRANCH}" \
		--body="${COMMIT_LOG}" \
		"artifacts/linux.tar.gz" \
		"artifacts/windows.zip" \
		"artifacts/darwin.tar.gz"
	
version:
	@echo $(VERSION)