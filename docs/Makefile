.PHONY: clean dep build serve

DOC_SRC = source/v1
DOC_FILE = fortune_app_openapi.yml

clean:
	rm -rf node_modules
	rm -f $(DOC_FILE)

dep: clean
	npm install

build: dep
	echo "bla"
	@cat $(DOC_SRC)/main.yml >> $(DOC_FILE)
	@cat $(DOC_SRC)/fortune.yml >> $(DOC_FILE)
	@cat $(DOC_SRC)/models.yml >> $(DOC_FILE)
	
	# @node_modules/redoc-cli/index.js validate $(DOC_FILE)
	@node_modules/redoc-cli/index.js bundle	-t static/template.hbs --options option.hideLoading $(DOC_FILE) --output ../assets/docs.html

serve: build
	@node_modules/redoc-cli/index.js serve $(DOC_FILE)