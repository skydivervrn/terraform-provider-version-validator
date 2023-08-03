package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
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
	currentVersionArray := fillZeroes(strings.Split(currentVersionString, "."))
	requiredVersionArray := fillZeroes(strings.Split(DigitPostfix(requiredVersionString), "."))
	currentVersionArrayInt := stringArrToNumbers(currentVersionArray)
	requiredVersionArrayInt := stringArrToNumbers(requiredVersionArray)
	comparison := compare(currentVersionArrayInt, requiredVersionArrayInt)
	digitPrefix := DigitPrefix(requiredVersionString)
	inRangeMessageFailed := fmt.Sprintf("current_version: %s\nrequired_version:%s\nCurrent version NOT in required range. Check failed", currentVersionString, requiredVersionString)
	inRangeMessageSucceed := fmt.Sprintf("current_version: %s\nrequired_version:%s\nCurrent version IN required range. Check succeed", currentVersionString, requiredVersionString)
	lowerMessageFailed := fmt.Sprintf("current_version: %s\nrequired_version:%s\nCurrent version lower then required. Check failed", currentVersionString, requiredVersionString)
	lowerMessageSucceed := fmt.Sprintf("current_version: %s\nrequired_version:%s\nCurrent version is lower then required. Check succeed", currentVersionString, requiredVersionString)
	lowerEqualMessageSucceed := fmt.Sprintf("current_version: %s\nrequired_version:%s\nCurrent version is lower or equal required. Check succeed", currentVersionString, requiredVersionString)
	equalMessageFailed := fmt.Sprintf("current_version: %s\nrequired_version:%s\nCurrent version is NOT equal required. Check failed", currentVersionString, requiredVersionString)
	equalMessageSucceed := fmt.Sprintf("current_version: %s\nrequired_version:%s\nCurrent version is equal required. Check succeed", currentVersionString, requiredVersionString)
	greaterMessageFailed := fmt.Sprintf("current_version: %s\nrequired_version:%s\nCurrent version greater then required. Check failed", currentVersionString, requiredVersionString)
	greaterEqualMessageSucceed := fmt.Sprintf("current_version: %s\nrequired_version:%s\nCurrent version is greater or equal required. Check succeed", currentVersionString, requiredVersionString)
	switch digitPrefix {
	case ">":
		if comparison == greaterName {
			log.Printf("current_version: %s\nrequired_version:%s\nCurrent version is greater then required. Check succeed", currentVersionString, requiredVersionString)
			return diags
		}
		diags = diag.Errorf(lowerMessageFailed)
		return diags
	case "<":
		if comparison == lowerName {
			log.Printf(lowerMessageSucceed)
			return diags
		}
		diags = diag.Errorf(greaterMessageFailed)
		return diags
	case "":
		if comparison == equalName {
			log.Printf(equalMessageSucceed)
			return diags
		}
		diags = diag.Errorf(equalMessageFailed)
		return diags
	case ">=":
		if comparison == equalName || comparison == greaterName {
			log.Printf(greaterEqualMessageSucceed)
			return diags
		}
		diags = diag.Errorf(lowerMessageFailed)
		return diags
	case "<=":
		if comparison == equalName || comparison == lowerName {
			log.Printf(lowerEqualMessageSucceed)
			return diags
		}
		diags = diag.Errorf(greaterMessageFailed)
		return diags
	case "^":
		if comparison == lowerName {
			diags = diag.Errorf(lowerMessageFailed)
			return diags
		}
		if comparison == equalName {
			log.Printf(equalMessageSucceed)
			return diags
		}
		if comparison == greaterName {
			if requiredVersionArrayInt[0] != 0 {
				tempArr := []int{requiredVersionArrayInt[0] + 1, 0, 0}
				comp := compare(currentVersionArrayInt, tempArr)
				if comp == lowerName {
					log.Printf(inRangeMessageSucceed)
					return diags
				} else {
					diags = diag.Errorf(inRangeMessageFailed)
					return diags
				}
			} else if requiredVersionArrayInt[1] != 0 {
				tempArr := []int{0, requiredVersionArrayInt[1] + 1, 0}
				comp := compare(currentVersionArrayInt, tempArr)
				if comp == lowerName {
					log.Printf(inRangeMessageSucceed)
					return diags
				} else {
					diags = diag.Errorf(inRangeMessageFailed)
					return diags
				}
			}
			diags = diag.Errorf(inRangeMessageFailed)
			return diags
		}
		diags = diag.Errorf(inRangeMessageFailed)
		return diags
	case "~":
		//if comparison == equalName || comparison == lowerName {
		//	log.Printf("current_version: %s\nrequired_version:%s\nCurrent version is lower or equal required. Check succeed", currentVersionString, requiredVersionString)
		//	return diags
		//}
		//diags = diag.Errorf("current_version: %s\nrequired_version:%s\nCurrent version is greater then required. Check failed", currentVersionString, requiredVersionString)
		return diags
	default:
		diags = diag.Errorf(fmt.Sprintf("required_version:%s\nWrong symbols in `required_version`", requiredVersionString))
		return diags
	}
}

func compare(currentVersionArrayInt []int, requiredVersionArrayInt []int) string {
	if currentVersionArrayInt[0] > requiredVersionArrayInt[0] {
		return greaterName
	} else if currentVersionArrayInt[0] == requiredVersionArrayInt[0] {
		if currentVersionArrayInt[1] > requiredVersionArrayInt[1] {
			return greaterName
		} else if currentVersionArrayInt[1] == requiredVersionArrayInt[1] {
			if currentVersionArrayInt[2] > requiredVersionArrayInt[2] {
				return greaterName
			} else if currentVersionArrayInt[2] == requiredVersionArrayInt[2] {
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

func fillZeroes(arr []string) (resultArr []string) {
	resultArr = []string{"0", "0", "0"}
	for index, element := range arr {
		resultArr[index] = element
	}
	return
}

func stringArrToNumbers(stringArr []string) (intArr []int) {
	intArr = []int{0, 0, 0}
	for i, v := range stringArr {
		err := errors.New("")
		intArr[i], err = strconv.Atoi(v)
		if err != nil {
			diag.Errorf(fmt.Sprintf("Error: %s", err))
		}
	}
	return
}
