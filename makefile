JS_SRC := data/entry.ts
JS_OUTDIR := data/build
JS_OUT := $(JS_OUTDIR)/bundle.js

.PHONY: run clean build js

run: build
	go run .

build: $(JS_OUT)

$(JS_OUT): $(JS_SRC)
	npx esbuild $(JS_SRC) \
		--bundle \
		--format=iife \
		--platform=neutral \
		--outdir=$(JS_OUTDIR)

clean:
	rm -rf $(JS_OUTDIR)
