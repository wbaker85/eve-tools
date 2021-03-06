package main

import (
	"github.com/wbaker85/eve-tools/pkg/models"
)

func tooExpensive(prices map[string]float64, rules []itemRule) []string {
	output := []string{}

	for _, val := range rules {
		thisName := val.ItemName
		max := val.BuyTargetPrice

		if prices[thisName] >= max {
			output = append(output, thisName)
		}
	}

	return output
}

func tooMuchInventory(hangar, escrow []models.CharacterAsset, rules []itemRule) []string {
	output := []string{}
	counts := combinedAssetCount(hangar, escrow)

	for _, v := range rules {
		thisName := v.ItemName
		max := v.MaxInventory

		if counts[thisName] >= max {
			output = append(output, thisName)
		}
	}

	return output
}

func combinedAssetCount(l1, l2 []models.CharacterAsset) map[string]int {
	output := make(map[string]int)

	for _, v := range l1 {
		output[v.Name] += v.Quantity
	}

	for _, v := range l2 {
		output[v.Name] += v.Quantity
	}

	return output
}

func shouldBeBuying(rules []itemRule, pricedOut, tooMuch []string) []string {
	output := []string{}
	m := stringSliceToMap(pricedOut, tooMuch)

	for _, val := range rules {
		if !m[val.ItemName] {
			output = append(output, val.ItemName)
		}
	}

	return output
}

func stringSliceToMap(s1, s2 []string) map[string]bool {
	output := make(map[string]bool)

	for _, val := range s1 {
		output[val] = true
	}

	for _, val := range s2 {
		output[val] = true
	}

	return output
}

func shouldBeSelling(inventory []models.CharacterAsset, rules []itemRule) []string {
	invMap := make(map[string]int)

	for _, val := range rules {
		invMap[val.ItemName] = val.MinSellLotSize
	}

	output := []string{}

	for _, val := range inventory {
		if invMap[val.Name] > 0 && val.Quantity >= invMap[val.Name] {
			output = append(output, val.Name)
		}
	}

	return output
}

func sliceDiff(base, comp []string) []string {
	compMap := make(map[string]struct{})
	for _, val := range comp {
		compMap[val] = struct{}{}
	}

	output := []string{}

	for _, val := range base {
		_, ok := compMap[val]
		if !ok {
			output = append(output, val)
		}
	}

	return output
}

func orphans(hangar, escrow []models.CharacterAsset, rules []itemRule) []string {
	counts := combinedAssetCount(hangar, escrow)
	ruleMap := make(map[string]bool)

	for _, val := range rules {
		ruleMap[val.ItemName] = true
	}

	output := []string{}

	for k, v := range counts {
		if v > 0 && !ruleMap[k] {
			output = append(output, k)
		}
	}

	return output
}

func sliceUnion(s1, s2 []string) []string {
	m1 := make(map[string]bool)
	for _, val := range s1 {
		m1[val] = true
	}

	output := []string{}

	for _, val := range s2 {
		if m1[val] {
			output = append(output, val)
		}
	}

	return output
}

func hangarPricedOut(h []models.CharacterAsset, p []string) []string {
	hangarItems := []string{}

	for _, val := range h {
		hangarItems = append(hangarItems, val.Name)
	}

	return sliceUnion(hangarItems, p)
}
