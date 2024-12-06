# Example PostgreSQL DB

Fill in variables in `variables.tf` and run:

```shell
terraform init
terraform plan
```

You can also use environment variables to set the configuration:

```shell
terraform init
TF_VAR_api_key=$DBAAS_TOKEN TF_VAR_project=$DBAAS_PROJECT TF_VAR_org=$DBAAS_ORG TF_VAR_db_password=verys3cretpassword terraform plan
```
