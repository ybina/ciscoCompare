package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func parseCsv() {
	profile, err := os.OpenFile("./data/allUserProfile.csv", os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	allUserProfile := csv.NewReader(profile)
	getCsvAllUserProfile(allUserProfile)
	fmt.Printf("userProfileCount:%v\n", userProfileCount)

	upb, err := os.OpenFile("./data/userProfileBinding.csv", os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	userProfileb := csv.NewReader(upb)
	getCsvUserProfileBinding(userProfileb)
	fmt.Printf("rule count:%v\n", csvRuleCount)

	fil, err := os.OpenFile("./data/filter.csv", os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	filter := csv.NewReader(fil)
	getCsvFilters(filter)
	fmt.Printf("filterCount:%v\n", csvFilterCount)

	act, err := os.OpenFile("./data/action.csv", os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	action := csv.NewReader(act)
	GetCsvActions(action)
	fmt.Printf("actionCount:%v\n", actionCount)
}

func getCsvAllUserProfile(r *csv.Reader) {
	var readCount = 0
	for {
		read, err := r.Read()
		if readCount == 0 {
			readCount++
			continue
		}
		if err == io.EOF {
			break
		}
		readCount++
		csvAllUserProfileArr = append(csvAllUserProfileArr, read[0])
		csvAllUserProfileMap[read[0]] = nil
		userProfileCount++
	}
}

func getCsvUserProfileBinding(r *csv.Reader) {
	var readCount = 0
	upErrCount := 0
	for {
		read, err := r.Read()
		if readCount == 0 {
			readCount++
			continue
		}
		if err == io.EOF {
			break
		}
		readCount++
		rk := RuleKey{
			UserProfile: read[0],
			RuleName:    read[3],
		}
		var isInstalled bool
		if strings.TrimSpace(read[2]) == "TRUE" {
			isInstalled = true
		}
		up := read[0]
		rule := CsvRule{
			UserProfile:     strings.TrimSpace(read[0]),
			OrderNo:         strings.TrimSpace(read[1]),
			InstalledByPcrf: isInstalled,
			RuleName:        strings.TrimSpace(read[3]),
			ChargingAction:  strings.TrimSpace(read[4]),
			MonitoringKey:   strings.TrimSpace(read[5]),
		}
		csvRuleCount++
		csvUserProfileBindingArr = append(csvUserProfileBindingArr, rule)
		csvUpbMap[rk] = rule
		addCsvUPRuleSliceMap(rule, up)

	}
	fmt.Printf(".....upCount:%v\n", upErrCount)
}
func addCsvUPRuleSliceMap(r CsvRule, up string) {
	upArr, ok := csvUPToRuleSliceMap[up]
	if !ok {
		arr := []CsvRule{r}
		csvUPToRuleSliceMap[r.UserProfile] = arr
	} else {
		upArr = append(upArr, r)
		csvUPToRuleSliceMap[r.UserProfile] = upArr
	}
}

// group-of-ruledefs 下绑定多个filter
//
// group-of-ruledefs DNS-QUERY-PROHIBITED-01
//
//	  add-ruledef priority 10 ruledef DNS-QUERY-PROHIBITED-RD01
//	  add-ruledef priority 20 ruledef DNS-QUERY-PROHIBITED-RD02
//	#exit
//
// TODO: BLOCK-TETHERING; WAZE-HTTP-RD01

func getCsvFilters(r *csv.Reader) {
	readCount := 0
	for {
		read, err := r.Read()
		if readCount == 0 {
			readCount++
			continue
		}
		if err == io.EOF {
			break
		}
		readCount++
		f := CsvFilter{
			RuleName:     read[0],
			GroupName:    read[1],
			Type:         read[2],
			FilterAction: read[3],
		}
		if f.GroupName == "" {
			f.GroupName = f.RuleName
			f.RuleName = ""
		}
		csvFilterArr = append(csvFilterArr, f)
		m, ok := csvFilterMap[f.GroupName]
		if !ok {
			arr := []CsvFilter{f}
			csvFilterMap[f.GroupName] = arr
		} else {
			m = append(m, f)
			csvFilterMap[f.GroupName] = m
		}
		csvFilterCount++
		k := CsvFilterKey{
			GroupName:    f.GroupName,
			Type:         f.Type,
			FilterAction: f.FilterAction,
		}
		csvFilterActionCountMap[k] = nil
	}
}

func GetCsvActions(r *csv.Reader) {
	readCount := 0
	for {
		read, err := r.Read()
		if readCount == 0 {
			readCount++
			continue
		}
		if err == io.EOF {
			break
		}
		k := read[0]
		a := CsvAction{
			ChargingName:   read[0],
			ChargingAction: read[1],
			ChargingType:   read[2],
		}
		csvActionArr = append(csvActionArr, a)
		if arr, ok := csvChargingNameMap[k]; !ok {
			csvChargingNameCount++
			csvChargingActionCount++
			ar := []CsvAction{a}
			csvChargingNameMap[k] = ar
		} else {
			csvChargingActionCount++
			arr = append(arr, a)
			csvChargingNameMap[k] = arr
		}
		actionCount++
	}
}
