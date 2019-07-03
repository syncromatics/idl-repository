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

build:
	docker build \
		--build-arg "BUILD_FLAGS=$(BUILD_FLAGS)" \
		--target package \
		-t package:${VERSION} \
		.
	docker build \
		--build-arg "BUILD_FLAGS=$(BUILD_FLAGS)" \
		--target idl \
		-t syncromatics/idl:${VERSION} \
		.
	docker build \
		--build-arg "BUILD_FLAGS=$(BUILD_FLAGS)" \
		--target idl-repository \
		-t syncromatics/idl-repository:${VERSION} \
		.

package: build
	docker run --rm -v $$PWD:/data --entrypoint cp package:${VERSION} -R . /data

publish:
	echo "$$DOCKER_PASSWORD" | docker login -u "$$DOCKER_USERNAME" --password-stdin
	docker push syncromatics/idl:${VERSION}
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

documentation: build
	docker run --rm -v $$PWD/docs:/build/docs --workdir /build --entrypoint go package:${VERSION} run internal/cobraDocs.go
