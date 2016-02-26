BEATNAME=openstackbeat
BEAT_DIR=github.com/fassisrosa/beats
SYSTEM_TESTS=false
TEST_ENVIRONMENT=false
ES_BEATS=./vendor/github.com/elastic/beats
GOPACKAGES=$(shell glide novendor)
PREFIX?=.

# Path to the libbeat Makefile
-include $(ES_BEATS)/libbeat/scripts/Makefile

.PHONY: init
init:
	glide update  --no-recursive
	make update
	git init
	git add .

.PHONY: install-cfg
install-cfg:
	mkdir -p $(PREFIX)
	cp etc/openstackbeat.template.json     $(PREFIX)/openstackbeat.template.json
	cp etc/openstackbeat.yml               $(PREFIX)/openstackbeat.yml
	cp etc/openstackbeat.yml               $(PREFIX)/openstackbeat-linux.yml
	cp etc/openstackbeat.yml               $(PREFIX)/openstackbeat-binary.yml
	cp etc/openstackbeat.yml               $(PREFIX)/openstackbeat-darwin.yml
	cp etc/openstackbeat.yml               $(PREFIX)/openstackbeat-win.yml

.PHONY: update-deps
update-deps:
	glide update  --no-recursive
