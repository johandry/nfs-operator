SHELL				= /bin/bash

ECHO 				= echo -e
C_STD 			= $(shell echo -e "\033[0m")
C_RED				= $(shell echo -e "\033[91m")
C_GREEN 		= $(shell echo -e "\033[92m")
P 			 		= $(shell echo -e "\033[92m> \033[0m")
OK 			 		= $(shell echo -e "\033[92m[ OK  ]\033[0m")
ERROR		 		= $(shell echo -e "\033[91m[ERROR] \033[0m")
PASS		 		= $(shell echo -e "\033[92m[PASS ]\033[0m")
FAIL		 		= $(shell echo -e "\033[91m[FAIL ] \033[0m")

default: apply test

all: check init apply test

environment: all

## Validations:

check-ibmcloud:
	@if ibmcloud --version | grep -q 'ibmcloud version'; then $(ECHO) "$(P) $(OK) ibmcloud"; else $(ECHO) "$(P) $(ERROR) ibmcloud"; exit 1; fi
	@if ibmcloud plugin show infrastructure-service | grep -q 'Plugin Version'; then $(ECHO) "$(P) $(OK) Plugin infrastructure-service"; else $(ECHO) "$(P) $(ERROR) Plugin infrastructure-service"; exit 1; fi
	@if ibmcloud plugin show kubernetes-service | grep -q 'Plugin Version'; then $(ECHO) "$(P) $(OK) Plugin kubernetes-service"; else $(ECHO) "$(P) $(ERROR) Plugin kubernetes-service"; exit 1; fi

check-api-key:
	@if [[ -n $$IC_API_KEY ]];\
		then $(ECHO) "$(P) $(OK) IC_API_KEY";\
		else $(ECHO) '$(P) $(ERROR) IC_API_KEY. Execute: make api-key; export IC_API_KEY=$$(jq -r .apikey test/terraform/.terraform_key.json)';  exit 1; fi

check-terraform:
	@if terraform version | grep -q 'Terraform v'; then $(ECHO) "$(P) $(OK) terraform"; else $(ECHO) "$(P) $(ERROR) terraform"; exit 1; fi
	@if ls ~/.terraform.d/plugins/terraform-provider-ibm_* | grep -q 'provider'; then $(ECHO) "$(P) $(OK) IBM Cloud Provider"; else $(ECHO) "$(P) $(ERROR) IBM Cloud Provider"; exit 1; fi

print-target:
	@$(ECHO) "$(P) ibmcloud target: "; ibmcloud target

check: check-ibmcloud check-terraform check-api-key

## Initializations:

api-key:
	ibmcloud iam api-key-create TerraformKey -d "API Key for Terraform" --file .terraform_key.json

init:
	@if [[ -e .target_account ]]; then $(ECHO) "$(P) $(OK) Target IBM Cloud Account"; else $(ECHO) "$(P) $(ERROR) Missing File with Target IBM Cloud Account: .target_account"; exit 1; fi
	ibmcloud target -c $$(cat .target_account) -r us-south -g Default
	terraform init && \
	terraform validate && \
	terraform plan

## Apply

apply: check-terraform check-ibmcloud check-api-key
	terraform apply -auto-approve
	ibmcloud target -g $$(terraform output resource_group)
	ibmcloud ks cluster config -cluster $$(terraform output cluster_id)

## Test

test-kubernetes:
	@if TERM=dumb kubectl cluster-info | grep -q 'Kubernetes master is running at';\
		then $(ECHO) "$(P) $(OK) Kubernetes Cluster Running";\
		else $(ECHO) "$(P) $(ERROR) Kubernetes Cluster Running"; exit 1; fi

test-entrypoint:
	@if terraform output entrypoint | grep -q 'https://';\
		then $(ECHO) "$(P) $(OK) Kubernetes Cluster Created by Terraform";\
		else $(ECHO) "$(P) $(ERROR) Kubernetes Cluster by Terraform"; exit 1; fi

test: test-entrypoint test-kubernetes

## Cleanup

clean:
	terraform destroy -auto-approve


