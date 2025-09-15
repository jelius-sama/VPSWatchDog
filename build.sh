#!/bin/sh
mkdir -p ./bin

build_all() {
    # Linux
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

    # Darwin
    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_darwin_amd64 ./cmd/
    CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_darwin_arm64 ./cmd/

    # FreeBSD
    CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_freebsd_amd64 ./cmd/
    CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_freebsd_386 ./cmd/

    # OpenBSD
    CGO_ENABLED=0 GOOS=openbsd GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_openbsd_amd64 ./cmd/
    CGO_ENABLED=0 GOOS=openbsd GOARCH=386 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_openbsd_386 ./cmd/
    CGO_ENABLED=0 GOOS=openbsd GOARCH=arm64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_openbsd_arm64 ./cmd/

    # NetBSD
    CGO_ENABLED=0 GOOS=netbsd GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_netbsd_amd64 ./cmd/
    CGO_ENABLED=0 GOOS=netbsd GOARCH=386 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_netbsd_386 ./cmd/
    CGO_ENABLED=0 GOOS=netbsd GOARCH=arm go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_netbsd_arm ./cmd/

    # DragonFlyBSD
    CGO_ENABLED=0 GOOS=dragonfly GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_dragonfly_amd64 ./cmd/

    # Solaris
    CGO_ENABLED=0 GOOS=solaris GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_solaris_amd64 ./cmd/

    # Plan 9
    CGO_ENABLED=0 GOOS=plan9 GOARCH=386 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_plan9_386 ./cmd/
    CGO_ENABLED=0 GOOS=plan9 GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_plan9_amd64 ./cmd/
}

usage() {
    echo "Usage: $0 [all]"
    echo
    echo "Build VPSWatchDog in different modes:"
    echo "  ./build.sh       Build a single binary for the current system (default)."
    echo "  ./build.sh all   Cross-compile for all supported OS/architectures."
    echo
    echo "Examples:"
    echo "  ./build.sh       -> builds ./bin/VPSWatchDog"
    echo "  ./build.sh all   -> builds multiple binaries into ./bin/"
    exit 1
}

if [ $# -eq 0 ]; then
    CGO_ENABLED=0  go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog ./            
elif [ $# -eq 1 ] && [ "$1" = "all" ]; then
    build_all
else
    usage
fi
