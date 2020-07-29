# Testing

The NFS Operator (at this time) is designed just for IBM Cloud and requires a Kubernetes cluster on IBM Cloud. This folder contain all the resources to provision an IKS cluster for development and testing, also all the Kubernetes resources required for development and testing the NFS Operator.

- [Testing](#testing)
  - [Requirements](#requirements)
  - [Build the environment](#build-the-environment)
  - [Tests](#tests)
  - [Cleanup](#cleanup)

## Requirements

Before execute the tests you need the following requirements:

1. Have an IBM Cloud account with required privileges
2. [Install IBM Cloud CLI](https://ibm.github.io/cloud-enterprise-examples/iac/setup-environment#install-ibm-cloud-cli)
3. [Install the IBM Cloud CLI Plugins](https://ibm.github.io/cloud-enterprise-examples/iac/setup-environment#ibm-cloud-cli-plugins) `infrastructure-service`, `schematics` and `container-registry`.
4. [Login to IBM Cloud with the CLI](https://ibm.github.io/cloud-enterprise-examples/iac/setup-environment#login-to-ibm-cloud)
5. [Install Terraform](https://ibm.github.io/cloud-enterprise-examples/iac/setup-environment#install-terraform)
6. [Install IBM Cloud Terraform Provider](https://ibm.github.io/cloud-enterprise-examples/iac/setup-environment#configure-access-to-ibm-cloud)
7. [Configure access to IBM Cloud](https://ibm.github.io/cloud-enterprise-examples/iac/setup-environment#configure-access-to-ibm-cloud) for Terraform and the IBM Cloud CLI setting up the `IC_API_KEY` environment variable.
8. Install the following tools:
   1. [jq](https://stedolan.github.io/jq/download/)
   2. [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

If you have an API Key but is not set neither have the JSON file when it was created, you must recreate the key. Delete the old one if won't be in use anymore. Then execute `make api-key`, set the `IC_API_KEY` and, optionally, validate the requirements with `make check`.

```bash
# Delete the old one, if won't be in use anymore
ibmcloud iam api-keys       # Identify your old API Key Name
ibmcloud iam api-key-delete OLD-NAME

# Create a new one and set it as environment variable
make api-key

export IC_API_KEY=$(grep '"apikey":' terraform_key.json | sed 's/.*: "\(.*\)".*/\1/')
# Or
export IC_API_KEY=$(jq -r .apikey terraform_key.json)

make check
```

The project also need the IBM Cloud Account ID stored in the file `terraform/.target_account`. This file is ignored by Git so it's not stored in the repository. List all the accounts you have access with the following `ibmcloud` command, then store the account GUID into the file `.target_account`, like so:

```bash
ibmcloud account list

echo 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx' > terraform/.target_account
```

The Terraform variables validation is made by terraform, before continue - at least - set them the following two variables either using the environment variables `TF_VAR_project_name` and `TF_VAR_owner` (i.e. `export TF_VAR_owner=$USER`) or the `terraform/terraform.tfvars` file, like this:

```hcl
project_name = "nfs-op-ja"
owner        = "johandry"
```

You may add more variables to customize the cluster, for example like this one, to have a larger cluster:

```hcl
region         = "us-south"
vpc_zone_names = ["us-south-1", "us-south-2", "us-south-3"]
flavors        = ["cx2.2x4", "cx2.4x8", "cx2.8x16"]
workers_count  = [3, 2, 1]
k8s_version    = "1.18"
```

## Build the environment

From the previous requirements section, make sure you have - at least - the following steps done:

1. Create the `test/terraform/.target_account` file with the IBM Cloud target account,
2. To export the `IC_API_KEY` environment variable with the API Key, and
3. Create the file `test/terraform/terraform.tfvars` with the variables `project_name` and `owner`.

Execute `make check` to verify all the requirements are all set. If the check pass then you are ready to create the environment executing the following `make` command from the `test` directory:

```bash
make environment
```

At the end of the process the `~/.kube/config` file is setup to point to the new cluster and the access is verified. If you'd like to verify this access later, execute `make test`.

If you change the environment configuration or would like to make sure everything is setup, go to the `terraform` directory and execute `make`, like this:

```bash
cd terraform
make
```

If the tests include to have a pre-created PVC, execute `make pvc` in the `test/` directory.

## Tests

After the NFS Provisioner operator is deployed to the cluster, verify everything is ready executing `make list` to list all the known resources. Or, you can use `make list-all` to list all the resources known or unknown:

```bash
make list
# Or
make list-all
```

You can use a consumer application to use or test the operator and the volume provisioned. To do this execute:

```bash
make consumer
make test
```

## Cleanup

To destroy your environment and cleanup what you have created, execute:

```bash
cd test
make clean
```

To cleanup the Kubernetes cluster of resources without destroying it, execute the `delete` rule:

```bash
make delete
```

To get the project as you cloned it, deleting all the created files, use the `purge` rule:

**IMPORTANT**: Use this command wisely, it will destroy everything you have created and setup.

```bash
make purge
```
