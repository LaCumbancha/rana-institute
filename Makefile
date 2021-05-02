SHELL := /bin/bash
PWD := $(shell pwd)

deploy:
	gcloud app deploy site/gcloud.yaml --version $(version) --promote
	gcloud app deploy infra/gcloud.yaml --version $(version) --promote

deploy-site:
	gcloud app deploy site/gcloud.yaml --version $(version) --promote

deploy-infra:
	gcloud app deploy infra/gcloud.yaml --version $(version) --promote

run-site:
	gcloud app browse -s site

run-infra:
	gcloud app browse -s infra

update-queues:
	gcloud tasks queues update $(queue) --routing-override=service:$(service),version:$(version)
