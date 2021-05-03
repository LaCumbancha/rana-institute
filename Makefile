SHELL := /bin/bash
PWD := $(shell pwd)

prepare-deploy:
	@touch .deploy.tmp
	@echo "y" >> .deploy.tmp

deploy: prepare-deploy
	gcloud app deploy site/gcloud.yaml --version $(version) --promote < .deploy.tmp
	gcloud app deploy cache/gcloud.yaml --version $(version) --promote < .deploy.tmp
	gcloud app deploy storage/gcloud.yaml --version $(version) --promote < .deploy.tmp
	@rm .deploy.tmp

deploy-site: prepare-deploy
	gcloud app deploy site/gcloud.yaml --version $(version) --promote < .deploy.tmp
	@rm .deploy.tmp

deploy-storage: prepare-deploy
	gcloud app deploy storage/gcloud.yaml --version $(version) --promote < .deploy.tmp
	@rm .deploy.tmp

deploy-cache: prepare-deploy
	gcloud app deploy cache/gcloud.yaml --version $(version) --promote < .deploy.tmp
	@rm .deploy.tmp

run-site:
	gcloud app browse -s site

update-queues:
	gcloud tasks queues update $(queue) --routing-override=service:$(service),version:$(version)
