mkdir -p ./bin

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_linux_amd64 ./cmd/
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_linux_arm ./cmd/
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_linux_arm64 ./cmd/
CGO_ENABLED=0 GOOS=linux GOARCH=ppc64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_linux_ppc64 ./cmd/
CGO_ENABLED=0 GOOS=linux GOARCH=ppc64le go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_linux_ppc64le ./cmd/
CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_linux_mips ./cmd/
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_linux_mipsle ./cmd/
CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_linux_mips64 ./cmd/
CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_linux_mips64le ./cmd/
CGO_ENABLED=0 GOOS=linux GOARCH=s390x go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_linux_s390x ./cmd/
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_darwin_amd64 ./cmd/
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_darwin_arm64 ./cmd/
