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

default: test

all: init test

## Validations:

check-kubectl:
	@if kubectl version | grep -q 'Client Version'; then $(ECHO) "$(P) $(OK) kubectl"; else $(ECHO) "$(P) $(ERROR) kubectl"; exit 1; fi

check-kubernetes:
	$(MAKE) -C ./terraform test

check: check-kubectl check-kubernetes

## Provisioning and Deployments:

environment:
	$(MAKE) -C terraform all

pvc:
	$(MAKE) -C kubernetes pvc

consumer:
	$(MAKE) -C kubernetes consumer

## List resources

list:
	$(MAKE) -C kubernetes list

list-all:
	$(MAKE) -C kubernetes list-all

## Test

test:
	$(MAKE) -C kubernetes test

## Cleanup

destroy:
	$(MAKE) -C terraform clean

delete:
	$(MAKE) -C kubernetes clean
	$(MAKE) -C kubernetes clean-pvc

remove:
	$(RM) -rf terraform/.terraform
	$(RM) -rf terraform/.terraform.tfstate*
	$(RM) -rf terraform/terraform.tfstate*
	$(RM) -rf terraform/terraform.tfvars
	$(RM) -rf .terraform_key.json
	$(RM) -rf .target_account

clean: delete destroy

purge: clean remove
