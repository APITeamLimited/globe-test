#!/bin/bash
# Copyright 2020 gRPC authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -eu -o pipefail

WORKDIR=$(mktemp -d)

function finish ***REMOVED***
  rm -rf "$WORKDIR"
***REMOVED***
trap finish EXIT

export GOBIN=$***REMOVED***WORKDIR***REMOVED***/bin
export PATH=$***REMOVED***GOBIN***REMOVED***:$***REMOVED***PATH***REMOVED***
mkdir -p $***REMOVED***GOBIN***REMOVED***

echo "remove existing generated files"
# grpc_testing_not_regenerate/*.pb.go is not re-generated,
# see grpc_testing_not_regenerate/README.md for details.
rm -f $(find . -name '*.pb.go' | grep -v 'grpc_testing_not_regenerate')

echo "go install google.golang.org/protobuf/cmd/protoc-gen-go"
(cd test/tools && go install google.golang.org/protobuf/cmd/protoc-gen-go)

echo "go install cmd/protoc-gen-go-grpc"
(cd cmd/protoc-gen-go-grpc && go install .)

echo "git clone https://github.com/grpc/grpc-proto"
git clone --quiet https://github.com/grpc/grpc-proto $***REMOVED***WORKDIR***REMOVED***/grpc-proto

echo "git clone https://github.com/protocolbuffers/protobuf"
git clone --quiet https://github.com/protocolbuffers/protobuf $***REMOVED***WORKDIR***REMOVED***/protobuf

# Pull in code.proto as a proto dependency
mkdir -p $***REMOVED***WORKDIR***REMOVED***/googleapis/google/rpc
echo "curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/rpc/code.proto"
curl --silent https://raw.githubusercontent.com/googleapis/googleapis/master/google/rpc/code.proto > $***REMOVED***WORKDIR***REMOVED***/googleapis/google/rpc/code.proto

mkdir -p $***REMOVED***WORKDIR***REMOVED***/out

# Generates sources without the embed requirement
LEGACY_SOURCES=(
  $***REMOVED***WORKDIR***REMOVED***/grpc-proto/grpc/binlog/v1/binarylog.proto
  $***REMOVED***WORKDIR***REMOVED***/grpc-proto/grpc/channelz/v1/channelz.proto
  $***REMOVED***WORKDIR***REMOVED***/grpc-proto/grpc/health/v1/health.proto
  $***REMOVED***WORKDIR***REMOVED***/grpc-proto/grpc/lb/v1/load_balancer.proto
  profiling/proto/service.proto
  reflection/grpc_reflection_v1alpha/reflection.proto
)

# Generates only the new gRPC Service symbols
SOURCES=(
  $(git ls-files --exclude-standard --cached --others "*.proto" | grep -v '^\(profiling/proto/service.proto\|reflection/grpc_reflection_v1alpha/reflection.proto\)$')
  $***REMOVED***WORKDIR***REMOVED***/grpc-proto/grpc/gcp/altscontext.proto
  $***REMOVED***WORKDIR***REMOVED***/grpc-proto/grpc/gcp/handshaker.proto
  $***REMOVED***WORKDIR***REMOVED***/grpc-proto/grpc/gcp/transport_security_common.proto
  $***REMOVED***WORKDIR***REMOVED***/grpc-proto/grpc/lookup/v1/rls.proto
  $***REMOVED***WORKDIR***REMOVED***/grpc-proto/grpc/lookup/v1/rls_config.proto
  $***REMOVED***WORKDIR***REMOVED***/grpc-proto/grpc/testing/*.proto
  $***REMOVED***WORKDIR***REMOVED***/grpc-proto/grpc/core/*.proto
)

# These options of the form 'Mfoo.proto=bar' instruct the codegen to use an
# import path of 'bar' in the generated code when 'foo.proto' is imported in
# one of the sources.
#
# Note that the protos listed here are all for testing purposes. All protos to
# be used externally should have a go_package option (and they don't need to be
# listed here).
OPTS=Mgrpc/core/stats.proto=google.golang.org/grpc/interop/grpc_testing/core,\
Mgrpc/testing/benchmark_service.proto=google.golang.org/grpc/interop/grpc_testing,\
Mgrpc/testing/stats.proto=google.golang.org/grpc/interop/grpc_testing,\
Mgrpc/testing/report_qps_scenario_service.proto=google.golang.org/grpc/interop/grpc_testing,\
Mgrpc/testing/messages.proto=google.golang.org/grpc/interop/grpc_testing,\
Mgrpc/testing/worker_service.proto=google.golang.org/grpc/interop/grpc_testing,\
Mgrpc/testing/control.proto=google.golang.org/grpc/interop/grpc_testing,\
Mgrpc/testing/test.proto=google.golang.org/grpc/interop/grpc_testing,\
Mgrpc/testing/payloads.proto=google.golang.org/grpc/interop/grpc_testing,\
Mgrpc/testing/empty.proto=google.golang.org/grpc/interop/grpc_testing

for src in $***REMOVED***SOURCES[@]***REMOVED***; do
  echo "protoc $***REMOVED***src***REMOVED***"
  protoc --go_out=$***REMOVED***OPTS***REMOVED***:$***REMOVED***WORKDIR***REMOVED***/out --go-grpc_out=$***REMOVED***OPTS***REMOVED***:$***REMOVED***WORKDIR***REMOVED***/out \
    -I"." \
    -I$***REMOVED***WORKDIR***REMOVED***/grpc-proto \
    -I$***REMOVED***WORKDIR***REMOVED***/googleapis \
    -I$***REMOVED***WORKDIR***REMOVED***/protobuf/src \
    $***REMOVED***src***REMOVED***
done

for src in $***REMOVED***LEGACY_SOURCES[@]***REMOVED***; do
  echo "protoc $***REMOVED***src***REMOVED***"
  protoc --go_out=$***REMOVED***OPTS***REMOVED***:$***REMOVED***WORKDIR***REMOVED***/out --go-grpc_out=$***REMOVED***OPTS***REMOVED***,require_unimplemented_servers=false:$***REMOVED***WORKDIR***REMOVED***/out \
    -I"." \
    -I$***REMOVED***WORKDIR***REMOVED***/grpc-proto \
    -I$***REMOVED***WORKDIR***REMOVED***/googleapis \
    -I$***REMOVED***WORKDIR***REMOVED***/protobuf/src \
    $***REMOVED***src***REMOVED***
done

# The go_package option in grpc/lookup/v1/rls.proto doesn't match the
# current location. Move it into the right place.
mkdir -p $***REMOVED***WORKDIR***REMOVED***/out/google.golang.org/grpc/internal/proto/grpc_lookup_v1
mv $***REMOVED***WORKDIR***REMOVED***/out/google.golang.org/grpc/lookup/grpc_lookup_v1/* $***REMOVED***WORKDIR***REMOVED***/out/google.golang.org/grpc/internal/proto/grpc_lookup_v1

# grpc_testing_not_regenerate/*.pb.go are not re-generated,
# see grpc_testing_not_regenerate/README.md for details.
rm $***REMOVED***WORKDIR***REMOVED***/out/google.golang.org/grpc/reflection/grpc_testing_not_regenerate/*.pb.go

# grpc/testing does not have a go_package option.
mv $***REMOVED***WORKDIR***REMOVED***/out/grpc/testing/*.pb.go interop/grpc_testing/
mv $***REMOVED***WORKDIR***REMOVED***/out/grpc/core/*.pb.go interop/grpc_testing/core/

cp -R $***REMOVED***WORKDIR***REMOVED***/out/google.golang.org/grpc/* .
