S2I_BUILD_IMAGE := docker.io/dimssss/golang-s2i:0.4
VERSION := $(shell git rev-parse --short HEAD)


.PHONY: build-docker check-env



check-docker-registry:
ifndef DOCKER_REGISTRY
	$(error DOCKER_REGISTRY is undefined, plese export DOCKER_REGISTRY, example: export DOCKER_REGISTRY=docker.io/dimssss/uac)
endif

check-uac-webhook-service-name:
ifndef UAC_WEBHOOK_SERVICE_NAME
	$(error UAC_WEBHOOK_SERVICE_NAME is undefined, plese export UAC_WEBHOOK_SERVICE_NAME, example: export UAC_WEBHOOK_SERVICE_NAME=uac.bnhp-system.svc.cluster.local)
endif


build-docker: check-docker-registry
	@echo VERIONS: $(VERSION)
	@echo DOCKER IMAGE: ${DOCKER_REGISTRY}:$(VERSION)
	@echo STARTING S2I BUILD
	s2i build . $(S2I_BUILD_IMAGE) $(DOCKER_REGISTRY):$(VERSION)

generate-tls: check-uac-webhook-service-name
	rm -fr /tmp/webhook_deployment
	mkdir /tmp/webhook_deployment
	@deploy/setuptls/create-certs.sh $(UAC_WEBHOOK_SERVICE_NAME)





