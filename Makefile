SHELL := /bin/bash
PWD := $(shell pwd)

deploy:
	gcloud app deploy site/gcloud.yaml --version $(version) --promote
	gcloud app deploy cache/gcloud.yaml --version $(version) --promote
	gcloud app deploy storage/gcloud.yaml --version $(version) --promote

deploy-site:
	gcloud app deploy site/gcloud.yaml --version $(version) --promote

deploy-storage:
	gcloud app deploy storage/gcloud.yaml --version $(version) --promote

deploy-cache:
	gcloud app deploy cache/gcloud.yaml --version $(version) --promote

run-site:
	gcloud app browse -s site

update-queues:
	gcloud tasks queues update $(queue) --routing-override=service:$(service),version:$(version)
