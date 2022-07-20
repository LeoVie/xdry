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
	GOOS=${goos} GOARCH=${goarch} cd src/cmd && go build -o ${binaryFile}

.PHONY: build_for_all_platforms
build_for_all_platforms:
	chmod +x ./multiplatform_build.sh
	./multiplatform_build.sh