package customcost

import (
	"fmt"
	"strings"
)

type CustomCostProperty string

const (
	CustomCostZoneProp           CustomCostProperty = "zone"
	CustomCostAccountNameProp    CustomCostProperty = "accountName"
	CustomCostChargeCategoryProp CustomCostProperty = "chargeCategory"
	CustomCostDescriptionProp    CustomCostProperty = "description"
	CustomCostResourceNameProp   CustomCostProperty = "resourceName"
	CustomCostResourceTypeProp   CustomCostProperty = "resourceType"
	CustomCostProviderIdProp     CustomCostProperty = "providerId"
	CustomCostUsageUnitProp      CustomCostProperty = "usageUnit"
	CustomCostDomainProp         CustomCostProperty = "domain"
	CustomCostCostSourceProp     CustomCostProperty = "costSource"
)

func ParseCustomCostProperties(props []string) ([]CustomCostProperty, error) {
	var properties []CustomCostProperty
	added := make(map[CustomCostProperty]struct{})

	for _, prop := range props {
		property, err := ParseCustomCostProperty(prop)
		if err != nil {
			return nil, fmt.Errorf("failed to parse property: %w", err)
		}

		if _, ok := added[property]; !ok {
			added[property] = struct{}{}
			properties = append(properties, property)
		}
	}

	return properties, nil
}

func ParseCustomCostProperty(text string) (CustomCostProperty, error) {
	switch strings.TrimSpace(strings.ToLower(text)) {
	case strings.TrimSpace(strings.ToLower(string(CustomCostZoneProp))):
		return CustomCostZoneProp, nil
	case strings.TrimSpace(strings.ToLower(string(CustomCostAccountNameProp))):
		return CustomCostAccountNameProp, nil
	case strings.TrimSpace(strings.ToLower(string(CustomCostChargeCategoryProp))):
		return CustomCostChargeCategoryProp, nil
	case strings.TrimSpace(strings.ToLower(string(CustomCostDescriptionProp))):
		return CustomCostDescriptionProp, nil
	case strings.TrimSpace(strings.ToLower(string(CustomCostResourceNameProp))):
		return CustomCostResourceNameProp, nil
	case strings.TrimSpace(strings.ToLower(string(CustomCostResourceTypeProp))):
		return CustomCostResourceTypeProp, nil
	case strings.TrimSpace(strings.ToLower(string(CustomCostProviderIdProp))):
		return CustomCostProviderIdProp, nil
	case strings.TrimSpace(strings.ToLower(string(CustomCostUsageUnitProp))):
		return CustomCostUsageUnitProp, nil
	case strings.TrimSpace(strings.ToLower(string(CustomCostDomainProp))):
		return CustomCostDomainProp, nil
	case strings.TrimSpace(strings.ToLower(string(CustomCostCostSourceProp))):
		return CustomCostCostSourceProp, nil
	}

	return "", fmt.Errorf("invalid custom cost property: %s", text)
}
