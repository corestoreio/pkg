# Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GOTEST=go test -race -v -cover
GORUN=go run -v

DBTESTS = ./codegen ./config/... ./directory/... ./eav/... ./store/... ./storage/...

# ./net/...
NONDBTESTS = ./util/... ./locale/... ./i18n/... \
./config/... ./store/... \
./vendor/golang.org/x/text/...

test: testnodb testdb

testnodb: clean
	$(GOTEST) $(NONDBTESTS)

testdb: clean
	# setup DB
	@echo "Running tests for database schema"
	@export CS_DSN_TEST='magento2:magento2@tcp(localhost:3306)/magento2' && \
	export CS_DSN='magento2:magento2@tcp(localhost:3306)/magento2' && \
	go run codegen/tableToStruct/*.go && \
	$(GOTEST) $(DBTESTS)

clean:
	find . -name tables_generated.go -delete
	find . -name godepgraph.svg -delete

tts: clean
	@echo "Generating go source from MySQL tables"
	$(GORUN) codegen/tableToStruct/*.go

depgraph:
	# http://talks.golang.org/2015/tricks.slide#51
	find . -type d -not -iwholename '*.git*' -exec sh -c "godepgraph -horizontal {} | dot -Tsvg -o {}/godepgraph.svg" \;

cover:
	gocov test ./... | gocov report > test_coverage.txt

generate:
	# TODO add refactored version of https://github.com/mwitkow/go-proto-validators because its error messages are pretty weird and not in the sense of corestore
	find .  -not -iwholename '*.git*' -not -iwholename '*vendor*' -name *.pb.go -delete
	# config/validation
	protoc \
	--gogo_out=plugins=grpc,\
	Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
	Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,\
	Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,\
	Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api,\
	Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types:\
	./config/observer/ \
	--proto_path=../../../:../../../github.com/gogo/protobuf/protobuf/ \
	-I ./config/observer/ \
	proto.proto
