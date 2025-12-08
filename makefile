.PHONY: gen
gen: genopenapi genproto

.PHONY: genopenapi
genopenapi:
	@./scripts/genopenapi.sh

.PHONY: genproto
genproto:
	@./scripts/genproto.sh

.PHONY: fmt
fmt:
	goimports -l -w internal/

.PHONY: lint
lint:
	@./scripts/lint.sh