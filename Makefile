.PHONY: build
build:
ifndef goos
	$(error goos is not set)
endif
ifndef goarch
	$(error goarch is not set)
endif
ifndef binaryFile
	$(error binaryFile is not set)
endif
ifndef version
	$(error version is not set)
endif
	GOOS=${goos} GOARCH=${goarch} cd src/cmd && go build -o ${binaryFile} -ldflags "-X main.version=${version}"

.PHONY: build_for_all_platforms
build_for_all_platforms:
ifndef version
	$(error version is not set)
endif
	@echo "Building for all platforms (version ${version})"
	chmod +x ./multiplatform_build.sh
	./multiplatform_build.sh ${version}

.PHONY: test
test:
	go test -v ./src...