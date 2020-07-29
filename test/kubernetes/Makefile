SHELL				= /bin/bash

REGISTRY 			?= johandry
APP_NAME			?= movies
VERSION				?= 1.0
IMAGE 				:= $(REGISTRY)/$(APP_NAME):$(VERSION)
MUTABLE_IMAGE  = $(REGISTRY)/$(APP_NAME):latest


ECHO 				= echo -e
C_STD 			= $(shell echo -e "\033[0m")
C_RED				= $(shell echo -e "\033[91m")
C_GREEN 		= $(shell echo -e "\033[92m")
P 			 		= $(shell echo -e "\033[92m> \033[0m")
OK 			 		= $(shell echo -e "\033[92m[ OK  ]\033[0m")
ERROR		 		= $(shell echo -e "\033[91m[ERROR] \033[0m")
PASS		 		= $(shell echo -e "\033[92m[PASS ]\033[0m")
FAIL		 		= $(shell echo -e "\033[91m[FAIL ] \033[0m")
NONE		 		= $(shell echo -e "\033[91m[NONE ] \033[0m")

default: nfs consumer

all: pvc nfs consumer

get-endpoint:
	@-kubectl get svc movies -o=jsonpath='{.status.loadBalancer.ingress[0].hostname}'

## PVC: validation, deploy, list, wait & delete

check-pvc:
	@if -kubectl get pvc nfs-block-custom | grep -q 'Bound'; then $(ECHO) "$(P) $(OK) IBM Cloud Block PVC"; else $(ECHO) "$(P) $(ERROR) IBM Cloud Block PVC"; exit 1; fi

deploy-pvc:
	-kubectl apply -f pvc.yaml

list-pvc:
	-@$(ECHO) "$(P) Persistent Volume Claim"; kubectl get pvc 2>/dev/null
	-@kubectl get pvc | awk '{print $$3}' | grep -v VOLUME | \
		while read pv; do $(ECHO) "$(P) Persistent Volume: $$pv"; kubectl get pv $$pv 2>/dev/null; done

wait-pvc:
	@printf "Waiting for IBM Cloud Block PVC to be Bound"
	@while -kubectl get pvc nfs-block-custom | grep -q 'Pending'; do printf .; sleep 3; done; echo

delete-pvc:
	-kubectl delete -f pvc.yaml 2>/dev/null || true

pvc: deploy-pvc wait-pvc

## Consumer: validation, build, deploy, list, wait & delete

check-consumer:
	@if -kubectl get pods | grep movies | grep -q Running;\
		then $(ECHO) "$(P) $(OK) Consumer Pod Running";\
		else $(ECHO) "$(P) $(ERROR) Consumer Pod Not Running"; exit 1; fi
	@if [[ $$($(MAKE) get-endpoint) == *.cloud ]];\
		then $(ECHO) "$(P) $(OK) Consumer Service Ready: http://"$$($(MAKE) get-endpoint)"/movies";\
		else $(ECHO) "$(P) $(ERROR) Consumer Service Not Ready"; exit 1; fi

build-consumer:
	docker build -t $(MUTABLE_IMAGE) ./consumer
	docker tag  $(MUTABLE_IMAGE) $(IMAGE)
	docker push $(IMAGE)
	docker push $(MUTABLE_IMAGE)

deploy-consumer:
	-kubectl apply -f consumer/cm.yaml
	-kubectl apply -f consumer/deployment.yaml

list-consumer:
	-@$(ECHO) "$(P) ConfigMap"; kubectl get cm movies-db 2>/dev/null || true
	-@$(ECHO) "$(P) Service"; kubectl get svc movies 2>/dev/null || true
	-@$(ECHO) "$(P) Deployment"; kubectl get deploy movies 2>/dev/null || true

wait-consumer:
	@printf "Waiting for consumer service"
	@while [[ $$($(MAKE) get-endpoint) != *.cloud ]]; do printf .; sleep 3; done; echo
	@printf "Waiting for consumer pod"
	@while -kubectl get pods | grep movies | grep -q -v Running; do printf .; sleep 3; done; echo

delete-consumer:
	-kubectl delete -f consumer/deployment.yaml 2>/dev/null || true
	-kubectl delete -f consumer/cm.yaml 2>/dev/null || true

consumer: build-consumer deploy-consumer wait-consumer

# wait-nfs:
# 	@printf "Waiting for NFS PVC to be Bound"
# 	@while -kubectl get pvc nfs | grep -q 'Pending'; do printf .; sleep 3; done; echo

# check-nfs:
# 	@if -kubectl get pvc nfs | grep -q 'Bound'; then $(ECHO) "$(P) $(OK) NFS PVC"; else $(ECHO) "$(P) $(ERROR) NFS PVC"; exit 1; fi

## Test

test-movie:
	@ADDRESS=$$($(MAKE) get-endpoint); \
	if [[ $$(curl -s $$ADDRESS/movies/1067 | grep title | sed 's/.*: "\(.*\)",/\1/') == "Seven Samurai" ]]; \
		then $(ECHO) "$(P) $(PASS) Test: Check Get Movie API"; \
		else $(ECHO) "$(P) $(FAIL) Test: Check Get Movie API"; exit 1; fi

test-file:
	@if -kubectl exec --stdin --tty $(-kubectl get pods | grep movies | awk '{print $1}') -- cat /data/db.min.json | grep -q '"movies": '; \
		then $(ECHO) "$(P) $(PASS) Test: Check JSON File"; \
		else $(ECHO) "$(P) $(FAIL) Test: Check JSON File"; exit 1; fi

test: check-consumer test-movie test-file

list-nfs-provisioner:
	-@$(ECHO) "$(P) Service Account"; kubectl get sa nfs-provisioner 2>/dev/null || true
	-@$(ECHO) "$(P) Service"; kubectl get svc nfs-provisioner 2>/dev/null || true
	-@$(ECHO) "$(P) Deployment"; kubectl get deploy nfs-provisioner 2>/dev/null || true
	-@$(ECHO) "$(P) Role"; kubectl get clusterrole nfs-provisioner-runner 2>/dev/null || true
	-@$(ECHO) "$(P) RoleBinding"; kubectl get clusterrolebinding run-nfs-provisioner 2>/dev/null || true
	-@$(ECHO) "$(P) Role"; kubectl get role leader-locking-nfs-provisioner 2>/dev/null || true
	-@$(ECHO) "$(P) RoleBinding"; kubectl get rolebinding leader-locking-nfs-provisioner 2>/dev/null || true
	-@$(ECHO) "$(P) Storage Class"; kubectl get storageclass ibmcloud-nfs 2>/dev/null || true
	-@$(ECHO) "$(P) NFS Persistent Volume Claim"; kubectl get pvc nfs 2>/dev/null || true
	-@kubectl get pvc nfs 2>/dev/null | awk '{print $$3}' | grep -v VOLUME | \
		while read pv; do $(ECHO) "$(P) NFS Persistent Volume: $$pv"; kubectl get pv $$pv 2>/dev/null; done || true

list-nfs-operator:
	-@$(ECHO) "$(P) Custom Resource NFS"; kubectl get nfs 2>/dev/null || true
	-@$(ECHO) "$(P) Custom Resource Definition"; kubectl get crd nfs 2>/dev/null || true
	-@$(ECHO) "$(P) Deployment"; kubectl get deploy nfs-operator 2>/dev/null || true
	-@$(ECHO) "$(P) Role"; kubectl get role nfs-operator 2>/dev/null || true
	-@$(ECHO) "$(P) RoleBinding"; kubectl get rolebinding nfs-operator 2>/dev/null || true
	-@$(ECHO) "$(P) Service Account"; kubectl get sa nfs-operator 2>/dev/null || true

list: list-pvc list-nfs-provisioner list-nfs-operator list-consumer

list-all:
	-@for r in pv pvc deploy cm po role rolebinding nfs; do\
		$(ECHO) "$(P) $$r"; \
		kubectl get $$r 2>/dev/null; \
	done || true
	-@$(ECHO) "$(P) svc"; kubectl get svc 2>/dev/null | grep -v kubernetes || true
	-@$(ECHO) "$(P) sa"; kubectl get sa 2>/dev/null | grep -v default || true
	-@$(ECHO) "$(P) storageclass"; kubectl get storageclass 2>/dev/null | grep -v ibmc-vpc-block || true
	-@$(ECHO) "$(P) clusterrole"; kubectl get clusterrole 2>/dev/null | grep nfs || true
	-@$(ECHO) "$(P) clusterrolebinding"; kubectl get clusterrolebinding 2>/dev/null | grep nfs || true
	-@$(ECHO) "$(P) crd"; kubectl get crd 2>/dev/null | grep nfs || true

## Clean

delete: delete-consumer delete-nfs

clean-pvc: delete-pvc

clean: delete
