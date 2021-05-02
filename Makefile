SHELL := /bin/bash
PWD := $(shell pwd)

deploy-site:
	gcloud app deploy site/gcloud.yaml --version $(version)

deploy-infra:
	gcloud app deploy infra/gcloud.yaml --version $(version)

run-site:
	gcloud app browse -s site

run-infra:
	gcloud app browse -s infra
