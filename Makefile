generate-docs:
	cd tools; go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-dir .. -provider-name sys11dbaas
