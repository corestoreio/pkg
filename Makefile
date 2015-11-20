VPKGS = ./eav ./store

test: test1 test2

testnodb:
	go test ./utils/...

test1: clean
	@echo "Running tests for Mage1 database schema"
	@export CS_DSN_TEST='magento-1-9:magento-1-9@tcp(localhost:3306)/magento-1-9' && \
	export CS_DSN='magento-1-9:magento-1-9@tcp(localhost:3306)/magento-1-9' && \
	go run codegen/tableToStruct/*.go && \
	go test -tags mage1 $(VPKGS)

test2: clean
	@echo "Running tests for Mage2 database schema"
	@export CS_DSN_TEST='magento2:magento2@tcp(localhost:3306)/magento2' && \
	export CS_DSN='magento2:magento2@tcp(localhost:3306)/magento2' && \
	go run codegen/tableToStruct/*.go && \
	go test -tags mage2 $(VPKGS)

clean:
	@find . -name generated_tables.go -delete

tts: clean
	go run codegen/tableToStruct/*.go
