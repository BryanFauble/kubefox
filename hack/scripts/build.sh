#!/bin/bash

source "$(dirname "${BASH_SOURCE[0]}")/setup.sh"

export GO111MODULE=on
export CGO_ENABLED=0

BUILD_DATE=$(TZ=UTC date --iso-8601=seconds)

if ! ${SKIP_GENERATE:-false}; then
	$SCRIPTS/generate.sh
fi

mkdir -p "${BUILD_OUT}"
rm -f "${COMPONENT_OUT}"

go build \
	-C "${COMPONENT_SRC}/" -o "${COMPONENT_OUT}" \
	-ldflags " \
		-w -s
		-X github.com/xigxog/kubefox/build.date=${BUILD_DATE}\
		-X github.com/xigxog/kubefox/build.component=${COMPONENT} \
		-X github.com/xigxog/kubefox/build.commit=${COMPONENT_COMMIT} \
		-X github.com/xigxog/kubefox/build.rootCommit=${ROOT_COMMIT} \
		-X github.com/xigxog/kubefox/build.headRef=${HEAD_REF} \
		-X github.com/xigxog/kubefox/build.tagRef=${TAG_REF}" \
	main.go
