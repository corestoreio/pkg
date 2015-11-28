# Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

# all packages tested here must pass the tests. packages not listed here
# are under development will break things.

DBTESTS = ./codegen ./config/... ./directory/... ./eav/... ./store/... ./storage/...
NONDBTESTS = ./utils/... ./net/... ./locale/... ./i18n/...

test: testnodb test1 test2

testnodb: clean
	@go test -race -v $(NONDBTESTS)

test1: clean
	@echo "Running tests for Mage1 database schema"
	@export CS_DSN_TEST='magento-1-9:magento-1-9@tcp(localhost:3306)/magento-1-9' && \
	export CS_DSN='magento-1-9:magento-1-9@tcp(localhost:3306)/magento-1-9' && \
	go run codegen/tableToStruct/*.go && \
	go test -v -race -tags mage1 $(DBTESTS)

test2: clean
	@echo "Running tests for Mage2 database schema"
	@export CS_DSN_TEST='magento2:magento2@tcp(localhost:3306)/magento2' && \
	export CS_DSN='magento2:magento2@tcp(localhost:3306)/magento2' && \
	go run codegen/tableToStruct/*.go && \
	go test -v -race -tags mage2 $(DBTESTS)

clean:
	@find . -name tables_generated.go -delete

tts: clean
	go run -v codegen/tableToStruct/*.go

depgraph:
	# http://talks.golang.org/2015/tricks.slide#51
	find . -type d -not -iwholename '*.git*' -exec sh -c "godepgraph -horizontal {} | dot -Tsvg -o {}/godepgraph.svg" \;
