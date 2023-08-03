package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strings"
	"unicode"
)

const (
	equalName   = "equal"
	greaterName = "greater"
	lowerName   = "lower"
)

func dataSourceVersionValidator() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Version validator datasource.",

		ReadContext: dataSourceVersionValidatorRead,

		Schema: map[string]*schema.Schema{
			"current_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Currently deployed version.",
			},
			"required_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required version.",
			},
		},
	}
}

func dataSourceVersionValidatorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	currentVersionString := d.Get("current_version").(string)
	requiredVersionString := d.Get("required_version").(string)
	comparison := compare(currentVersionString, DigitPostfix(requiredVersionString))
	digitPrefix := DigitPrefix(requiredVersionString)

	switch digitPrefix {
	case ">":
		if comparison == greaterName {
			log.Printf("current_version: %s\nrequired_version:%s\nCurrent version is greater then required. Check succeed", currentVersionString, requiredVersionString)
			return diags
		}
		diags = diag.Errorf("current_version: %s\nrequired_version:%s\nCurrent version lower then required. Check failed", currentVersionString, requiredVersionString)
		return diags
	case "<":
		if comparison == lowerName {
			log.Printf("current_version: %s\nrequired_version:%s\nCurrent version is lower then required. Check succeed", currentVersionString, requiredVersionString)
			return diags
		}
		diags = diag.Errorf("current_version: %s\nrequired_version:%s\nCurrent version greater then required. Check failed", currentVersionString, requiredVersionString)
		return diags
	case "":
		if comparison == equalName {
			log.Printf("current_version: %s\nrequired_version:%s\nCurrent version is equal required. Check succeed", currentVersionString, requiredVersionString)
			return diags
		}
		diags = diag.Errorf("current_version: %s\nrequired_version:%s\nCurrent version is NOT equal required. Check failed", currentVersionString, requiredVersionString)
		return diags
	case ">=":
		if comparison == equalName || comparison == greaterName {
			log.Printf("current_version: %s\nrequired_version:%s\nCurrent version is greater or equal required. Check succeed", currentVersionString, requiredVersionString)
			return diags
		}
		diags = diag.Errorf("current_version: %s\nrequired_version:%s\nCurrent version is lower then required. Check failed", currentVersionString, requiredVersionString)
		return diags
	case "<=":
		if comparison == equalName || comparison == lowerName {
			log.Printf("current_version: %s\nrequired_version:%s\nCurrent version is lower or equal required. Check succeed", currentVersionString, requiredVersionString)
			return diags
		}
		diags = diag.Errorf("current_version: %s\nrequired_version:%s\nCurrent version is greater then required. Check failed", currentVersionString, requiredVersionString)
		return diags
	default:
		diags = diag.Errorf(fmt.Sprintf("required_version:%s\nWrong symbols in `required_version`", requiredVersionString))
		return diags
	}
}

func compare(currentVersion string, requiredVersion string) string {
	currentVersionArray := strings.Split(currentVersion, ".")
	requiredVersionArray := strings.Split(requiredVersion, ".")
	if currentVersionArray[0] > requiredVersionArray[0] {
		return greaterName
	} else if currentVersionArray[0] == requiredVersionArray[0] {
		if currentVersionArray[1] > requiredVersionArray[1] {
			return greaterName
		} else if currentVersionArray[1] == requiredVersionArray[1] {
			if currentVersionArray[2] > requiredVersionArray[2] {
				return greaterName
			} else if currentVersionArray[2] == requiredVersionArray[2] {
				return equalName
			}
		} else {
			return lowerName
		}
	} else {
		return lowerName
	}
	return ""
}

func DigitPrefix(s string) string {
	for i, r := range s {
		if unicode.IsDigit(r) {
			return s[:i]
		}
	}
	return s
}

func DigitPostfix(s string) string {
	for i, r := range s {
		if unicode.IsDigit(r) {
			return s[i:]
		}
	}
	return s
}
