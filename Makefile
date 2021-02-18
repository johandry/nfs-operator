SHELL				= /bin/bash

# Current Operator version
VERSION 			?= 0.0.1

OPERATOR_NAME	?= nfs-operator
REGISTRY 			?= johandry

# Default bundle image tag
BUNDLE_IMG 		?= controller-bundle:$(VERSION)

# Image URL to use all building/pushing image targets
IMG 	 				 = $(REGISTRY)/$(OPERATOR_NAME):$(VERSION)
MUTABLE_IMG 	 = $(REGISTRY)/$(OPERATOR_NAME):latest

# Disable WebHooks by default. To enable WebHooks generate the certificates at
# /tmp/k8s-webhook-server/serving-certs/tls.{crt,key}
ENABLE_WEBHOOKS	?= false

# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

ECHO 				= echo -e
C_STD 			= $(shell echo -e "\033[0m")
C_RED				= $(shell echo -e "\033[91m")
C_GREEN 		= $(shell echo -e "\033[92m")
P 			 		= $(shell echo -e "\033[92m> \033[0m")
OK 			 		= $(shell echo -e "\033[92m[ OK  ]\033[0m")
ERROR		 		= $(shell echo -e "\033[91m[ERROR] \033[0m")
PASS		 		= $(shell echo -e "\033[92m[PASS ]\033[0m")
FAIL		 		= $(shell echo -e "\033[91m[FAIL ] \033[0m")

# Options for 'bundle-build'
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

default: manager

all: manager install release deploy

# Run tests
test: generate fmt vet manifests
	go test ./... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	ENABLE_WEBHOOKS=$(ENABLE_WEBHOOKS) go run ./main.go

# Install CRDs into a cluster
install: manifests kustomize
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests kustomize
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests kustomize
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

# Deletes the controller from the Kubernetes cluster
delete:
	$(KUSTOMIZE) build config/default | kubectl delete -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Release the docker image with the operator
release: docker-build docker-push
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default > nfs_provisioner_operator_install.yaml

# release: docker-build docker-push manifests kustomize
# 	$(KUSTOMIZE) build config/crd > docs/install.yaml

# Build the docker image
docker-build: test
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}
	docker tag  ${IMG} ${MUTABLE_IMG}
	docker push ${MUTABLE_IMG}

clean: delete uninstall


# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.3.0 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

kustomize:
ifeq (, $(shell which kustomize))
	@{ \
	set -e ;\
	KUSTOMIZE_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$KUSTOMIZE_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/kustomize/kustomize/v3@v3.5.4 ;\
	rm -rf $$KUSTOMIZE_GEN_TMP_DIR ;\
	}
KUSTOMIZE=$(GOBIN)/kustomize
else
KUSTOMIZE=$(shell which kustomize)
endif

# export TEST_ASSET_KUBECTL=./testbin/kubectl
# export TEST_ASSET_KUBE_APISERVER=./testbin/kube-apiserver
# export TEST_ASSET_ETCD=./testbin/etcd

# Setup binaries required to run the tests
# See that it expects the Kubernetes and ETCD version
K8S_VERSION = v1.18.2
ETCD_VERSION = v3.4.3
testbin:
	curl -sSLo setup_envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/kubebuilder/master/scripts/setup_envtest_bins.sh
	chmod +x setup_envtest.sh
	./setup_envtest.sh $(K8S_VERSION) $(ETCD_VERSION)
	rm ./setup_envtest.sh

# Generate bundle manifests and metadata, then validate generated files.
bundle: manifests
	operator-sdk generate kustomize manifests -q
	kustomize build config/manifests | operator-sdk generate bundle -q --overwrite --version $(VERSION) $(BUNDLE_METADATA_OPTS)
	operator-sdk bundle validate ./bundle

# Build the bundle image.
bundle-build:
	docker build -f bundle.Dockerfile -t $(BUNDLE_IMG) .

test-environment:
	$(MAKE) -C test/terraform all

deploy-nfs:
	$(KUSTOMIZE) build config/samples | kubectl apply -f -

deploy-consumer:
	kubectl apply -f test/kubernetes/consumer/movies.yaml

deploy-test: deploy-nfs deploy-consumer

delete-test:
	kubectl delete -f test/kubernetes/consumer/movies.yaml
	$(KUSTOMIZE) build config/samples | kubectl delete -f -

list-nfs:
	@$(ECHO) "$(P)$(C_GREEN) NFS:$(C_STD)"; kubectl get nfs.nfs.storage.ibmcloud.ibm.com/nfs-sample

list-consumer:
	@$(ECHO) "$(P)$(C_GREEN) Consumer Namespace:$(C_STD)"; kubectl get namespace/nfs-consumer-app
	@$(ECHO) "$(P)$(C_GREEN) Consumer PVC:$(C_STD)"; kubectl get persistentvolumeclaim/nfs-consumer-movies --namespace nfs-consumer-app
	@$(ECHO) "$(P)$(C_GREEN) Consumer Service:$(C_STD)"; kubectl get service/nfs-consumer-movies --namespace nfs-consumer-app
	@$(ECHO) "$(P)$(C_GREEN) Consumer ConfigMap:$(C_STD)"; kubectl get configmap/nfs-consumer-movies-db --namespace nfs-consumer-app
	@$(ECHO) "$(P)$(C_GREEN) Consumer Deployment:$(C_STD)"; kubectl get deployment.apps/nfs-consumer-movies --namespace nfs-consumer-app
	@$(ECHO) "$(P)$(C_GREEN) Consumer Pods:$(C_STD)"; kubectl get pods --namespace nfs-consumer-app

list-test: list-nfs list-consumer

list-operator:
	@$(ECHO) "$(P)$(C_GREEN) Namespace:$(C_STD)"; kubectl get namespace/nfs-operator-system
	@$(ECHO) "$(P)$(C_GREEN) CRD:$(C_STD)"; kubectl get customresourcedefinition.apiextensions.k8s.io/nfs.nfs.storage.ibmcloud.ibm.com --namespace nfs-operator-system
	@$(ECHO) "$(P)$(C_GREEN) RBAC:$(C_STD)"; \
		kubectl get role.rbac.authorization.k8s.io/nfs-operator-leader-election-role --namespace nfs-operator-system; \
		kubectl get clusterrole.rbac.authorization.k8s.io/nfs-operator-manager-role --namespace nfs-operator-system; \
		kubectl get clusterrole.rbac.authorization.k8s.io/nfs-operator-proxy-role --namespace nfs-operator-system; \
		kubectl get clusterrole.rbac.authorization.k8s.io/nfs-operator-metrics-reader --namespace nfs-operator-system; \
		kubectl get rolebinding.rbac.authorization.k8s.io/nfs-operator-leader-election-rolebinding --namespace nfs-operator-system; \
		kubectl get clusterrolebinding.rbac.authorization.k8s.io/nfs-operator-manager-rolebinding --namespace nfs-operator-system; \
		kubectl get clusterrolebinding.rbac.authorization.k8s.io/nfs-operator-proxy-rolebinding --namespace nfs-operator-system
	@$(ECHO) "$(P)$(C_GREEN) Cert-Manager & WebHook:$(C_STD)"; \
		kubectl get mutatingwebhookconfiguration.admissionregistration.k8s.io/nfs-operator-mutating-webhook-configuration --namespace nfs-operator-system; \
		kubectl get service/nfs-operator-webhook-service --namespace nfs-operator-system; \
		kubectl get certificate.cert-manager.io/nfs-operator-serving-cert --namespace nfs-operator-system; \
		kubectl get issuer.cert-manager.io/nfs-operator-selfsigned-issuer --namespace nfs-operator-system; \
		kubectl get validatingwebhookconfiguration.admissionregistration.k8s.io/nfs-operator-validating-webhook-configuration
	@$(ECHO) "$(P)$(C_GREEN) Service:$(C_STD)"; kubectl get service/nfs-operator-controller-manager-metrics-service --namespace nfs-operator-system
	@$(ECHO) "$(P)$(C_GREEN) Deployment:$(C_STD)"; kubectl get deployment.apps/nfs-operator-controller-manager --namespace nfs-operator-system
	@$(ECHO) "$(P)$(C_GREEN) Pods:$(C_STD)"; kubectl get pod --namespace nfs-operator-system
