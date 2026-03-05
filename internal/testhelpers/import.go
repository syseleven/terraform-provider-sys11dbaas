package testhelpers

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func ComposeAggregateImportStateCheckFunc(fs ...resource.ImportStateCheckFunc) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		var result []error

		for i, f := range fs {
			if err := f(s); err != nil {
				result = append(result, fmt.Errorf("Import check %d/%d error: %w", i+1, len(fs), err))
			}
		}

		return errors.Join(result...)
	}
}

func ImportCheckResourceAttr(key, value string) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("Attribute %q expected single instance state, got %d", key, len(s))
		}

		resource := s[0]
		if resource.Attributes[key] != value {
			return fmt.Errorf("Attribute %q expected %q, got %q", key, value, resource.Attributes[key])
		}

		return nil
	}
}

func ImportCheckResourceAttrSet(key string) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("Attribute %q expected single instance state, got %d", key, len(s))
		}

		resource := s[0]
		if resource.Attributes[key] == "" {
			return fmt.Errorf("Attribute %q should be set but is not", key)
		}

		return nil
	}
}
