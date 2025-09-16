#!/bin/sh
# POSIX sh build script — runs all targets, records failures, prints summary.

mkdir -p ./bin

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

FAILED_BUILDS=""

build() {
    label="$1"
    shift
    logfile="./bin/${label}.err"

    # Run the command using env so VAR=val tokens are treated as environment assignments
    # and the final token(s) form the command to execute.
    if env "$@" 2> "$logfile"; then
        printf '%b\n' "${GREEN}Build succeeded: ${label}${NC}"
        rm -f "$logfile"
    else
        printf '%b\n' "${RED}Build failed: ${label}${NC}"
        printf '   (see %s for error output)\n' "$logfile"
        FAILED_BUILDS="${FAILED_BUILDS} ${label}"
    fi
}

build_all() {
    # Linux
    build "linux_amd64" CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_linux_amd64 ./cmd/
    build "linux_arm" CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_linux_arm ./cmd/
    build "linux_arm64" CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_linux_arm64 ./cmd/
    build "linux_ppc64" CGO_ENABLED=0 GOOS=linux GOARCH=ppc64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_linux_ppc64 ./cmd/
    build "linux_ppc64le" CGO_ENABLED=0 GOOS=linux GOARCH=ppc64le go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_linux_ppc64le ./cmd/
    build "linux_mips" CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_linux_mips ./cmd/
    build "linux_mipsle" CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_linux_mipsle ./cmd/
    build "linux_mips64" CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_linux_mips64 ./cmd/
    build "linux_mips64le" CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_linux_mips64le ./cmd/
    build "linux_s390x" CGO_ENABLED=0 GOOS=linux GOARCH=s390x go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_linux_s390x ./cmd/

    # Darwin
    build "darwin_amd64" CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_darwin_amd64 ./cmd/
    build "darwin_arm64" CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_darwin_arm64 ./cmd/

    # FreeBSD
    build "freebsd_amd64" CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_freebsd_amd64 ./cmd/
    build "freebsd_386" CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_freebsd_386 ./cmd/

    # OpenBSD
    build "openbsd_amd64" CGO_ENABLED=0 GOOS=openbsd GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_openbsd_amd64 ./cmd/
    build "openbsd_386" CGO_ENABLED=0 GOOS=openbsd GOARCH=386 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_openbsd_386 ./cmd/
    build "openbsd_arm64" CGO_ENABLED=0 GOOS=openbsd GOARCH=arm64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_openbsd_arm64 ./cmd/

    # NetBSD
    build "netbsd_amd64" CGO_ENABLED=0 GOOS=netbsd GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_netbsd_amd64 ./cmd/
    build "netbsd_386" CGO_ENABLED=0 GOOS=netbsd GOARCH=386 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_netbsd_386 ./cmd/
    build "netbsd_arm" CGO_ENABLED=0 GOOS=netbsd GOARCH=arm go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_netbsd_arm ./cmd/

    # DragonFlyBSD
    build "dragonfly_amd64" CGO_ENABLED=0 GOOS=dragonfly GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_dragonfly_amd64 ./cmd/

    # Solaris
    build "solaris_amd64" CGO_ENABLED=0 GOOS=solaris GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_solaris_amd64 ./cmd/

    # Plan 9
    build "plan9_386" CGO_ENABLED=0 GOOS=plan9 GOARCH=386 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_plan9_386 ./cmd/
    build "plan9_amd64" CGO_ENABLED=0 GOOS=plan9 GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog_plan9_amd64 ./cmd/
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
    build "host_default" CGO_ENABLED=0  go build -trimpath -buildvcs=false -o ./bin/VPSWatchDog ./cmd/            
elif [ $# -eq 1 ] && [ "$1" = "all" ]; then
    build_all
else
    usage
fi

# Final summary
if [ -n "$FAILED_BUILDS" ]; then
    echo
    printf '%b\n' "${RED}Some builds failed:${NC}"
    for target in $FAILED_BUILDS; do
        printf '   - %s\n' "$target"
    done
    exit 1
else
    echo
    printf '%b\n' "${GREEN}All builds completed successfully.${NC}"
    exit 0
fi
