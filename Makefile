SHELL := /bin/bash
PWD := $(shell pwd)

environment := dev

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

update-queues:
	gcloud tasks queues update $(queue) --routing-override=service:$(service),version:$(version)

run-site:
	gcloud app browse -s site

run-test:
	influxd &
	brew services start grafana &
	k6 run -e K6_INFLUXDB_USERNAME=admin -e K6_INFLUXDB_PASSWORD=admin -e ENVIRONMENT=$(environment) --out influxdb=http://localhost:8086/myk6db test/k6-performance-test.js
