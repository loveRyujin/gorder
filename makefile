.PHONY: gen
gen: genopenapi genproto

.PHONY: genopenapi
genopenapi:
	@./scripts/genopenapi.sh

.PHONY: genproto
genproto:
	@./scripts/genproto.sh
