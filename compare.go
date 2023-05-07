package main

import (
	"fmt"
	"log"
	"strings"
)

// startCompare
func startCompare() {
	rpt := Report{}
	upReport := compareUserProfile()
	upReport.GetCsvUpNum = len(csvAllUserProfileArr)
	upReport.GetConfigRulebaseNum = ruleBaseCount
	rpt.UpReport = upReport

	upBindingReport := compareUserProfileBindingRules()
	upBindingReport.GetCsvUpBindingReportRuleNum = csvRuleCount
	upBindingReport.GetConfigRuleNum = configRuleCount
	rpt.UpBindingReport = upBindingReport

	filterReport := compareFilters()
	log.Printf("ruleDefFilterActionCount:%v\n", ruleDefFilterActionCount)
	filterReport.CsvGroupCount = len(csvFilterMap)
	filterReport.CsvFilterActionCount = len(csvFilterActionCountMap)
	filterReport.ConfigRuleDefCount = ruleDefCount
	filterReport.ConfigActionCount = ruleDefFilterActionCount
	rpt.FilterReport = filterReport

	actionReport := compareActions()
	actionReport.CsvChargingNameCount = csvChargingNameCount
	actionReport.ConfigChargingNameCount = configChargingActionNameCount
	rpt.ActionReport = actionReport

	rpt.WriteReport()

}

// compareUserProfile
func compareUserProfile() UpReport {

	var r UpReport
	log.Printf("csv up num:%v\n, config up num:%v\n", len(csvAllUserProfileMap), len(configRuleMap))
	for k, _ := range csvAllUserProfileMap {
		//fmt.Printf("csvProfile:%v\n", k)
		if _, ok := configUps[k]; ok {
			delete(csvAllUserProfileMap, k)
			delete(configUps, k)
		}
	}
	if len(csvAllUserProfileMap) > 0 {
		for k, _ := range csvAllUserProfileMap {
			r.OnlyCsvHaveUserProfile = append(r.OnlyCsvHaveUserProfile, k)
		}
	}
	if len(configUps) > 0 {
		for k, _ := range configUps {
			r.OnlyConfigHaveUserProfile = append(r.OnlyConfigHaveUserProfile, k)
		}
	}
	return r

}

// compareUserProfileBindingRules
func compareUserProfileBindingRules() UpBindingReport {
	var rpt UpBindingReport
	rpt.Diffs = []Diff{}
	csvData := csvUPToRuleSliceMap
	configData := configRuleMap
	log.Printf("cscData len:%v, configData len:%v\n", len(csvData), len(configData))
	for k1, v1 := range csvData {
		tmpk1 := strings.TrimSpace(k1)
		tmpv1 := v1
		v2, ok := configData[tmpk1]
		tmpv2 := v2

		if !ok { // 获取在csv中有而config中没有的userProfile的冗余配置
			d := Diff{
				Msg: fmt.Sprintf("csv userProfile: %v not found in config", k1),
			}
			rpt.Diffs = append(rpt.Diffs, d)
			continue
		}
		diffs := checkRuleMatchedBySameUp(tmpv1, tmpv2)
		if len(diffs) > 0 {
			rpt.Diffs = append(rpt.Diffs, diffs...)
		}
		delete(configData, tmpk1)
	}
	// TODO:
	// 获取在config中有而csv中没有的userProfile的冗余配置
	for k, _ := range configData {
		up := k
		d := Diff{}
		d.Msg = fmt.Sprintf("config userProfile not found in csv:%v", up)
		rpt.Diffs = append(rpt.Diffs, d)
	}
	return rpt
}

// checkRuleMatchedBySameUp
func checkRuleMatchedBySameUp(v1 []CsvRule, v2 []ConfigRule) []Diff {
	var res []Diff
	csvmp := make(map[string]CsvRule)
	cfgmp := make(map[string]ConfigRule)
	for _, v11 := range v1 {
		csvmp[v11.RuleName] = v11
	}
	for _, v22 := range v2 {
		cfgmp[v22.RuleName] = v22
	}
	for k1, d1 := range csvmp {
		d2, ok := cfgmp[k1]
		if !ok { // csv中有，而config中没有的配置
			d := Diff{}
			d.CsvPara = d1
			d.Msg = "csv rule not found in config"
			res = append(res, d)
			continue
			//log.Printf("find rule diff: csv rule not found in config:%v\n", d1.prt())
		} else {
			// 检查配置参数是否一致
			if !startCheckRule(d1, d2) {
				d := Diff{}
				d.CsvPara = d1
				d.ConfigPara = d2
				d.Msg = "rule para is different"
				res = append(res, d)
				continue
			}
			// 删除cfgmp中已查找到的配置
			delete(cfgmp, k1)
		}
	}
	// config中有而csv中没有的配置
	for _, v := range cfgmp {
		r := v
		d := Diff{
			Msg:        "config rule not found in csv",
			ConfigPara: r,
		}
		res = append(res, d)
	}
	return res
}

// startCheckRule
func startCheckRule(v1 CsvRule, v2 ConfigRule) bool {
	if v1.OrderNo != v2.Priority {
		return false
	}
	if v1.InstalledByPcrf != v2.IsDynamic {
		return false
	}
	if v1.RuleName != v2.RuleName {
		return false
	}
	if v1.ChargingAction != v2.ChargingAction {
		return false
	}
	if v1.MonitoringKey != v2.MonitoringKey {
		return false
	}
	return true
}

// compareFilters
func compareFilters() FilterReport {
	var rpt FilterReport
	rpt.FilterDiffs = []FilterDiff{}
	for kcsv, vcsv := range csvFilterMap {
		vconf, ok := configRuleDefMap[kcsv]
		if !ok {
			d := FilterDiff{
				Msg: fmt.Sprintf("csv group:%v not found in config", kcsv),
			}
			rpt.FilterDiffs = append(rpt.FilterDiffs, d)
			continue
		} else {
			diffs := compareFilterByGroupName(vcsv, vconf, kcsv)
			if len(diffs) > 0 {
				rpt.FilterDiffs = append(rpt.FilterDiffs, diffs...)
			}
			delete(configRuleDefMap, kcsv)
		}
	}
	for k, _ := range configRuleDefMap {
		rk := k
		d := FilterDiff{
			Msg: fmt.Sprintf("config group:%v not found in csv", rk),
		}
		rpt.FilterDiffs = append(rpt.FilterDiffs, d)
	}
	return rpt
}

// compareFilterByGroupName
func compareFilterByGroupName(csvFilter []CsvFilter, confFilter []string, groupName string) []FilterDiff {

	var res []FilterDiff
	csvTmp := make(map[string]CsvFilter)
	confTmp := make(map[string]interface{})
	//var t1 string
	//var t2 string

	for _, v1 := range csvFilter {

		tmp := v1

		action := strings.TrimLeft(tmp.FilterAction, "	 ")
		//if action == "http url starts-with http://www.google.com/gen_204" {
		//	t1 = tmp.FilterAction
		//	log.Printf("t1 match:%v\n", t1)
		//	log.Printf("confFilter:%v\n", confFilter)
		//}
		csvTmp[action] = tmp
	}
	for _, v2 := range confFilter {
		tmp := strings.TrimLeft(v2, "	 ")
		//if tmp == "http url starts-with http://www.google.com/gen_204" {
		//	log.Printf("t2 match:%v\n", tmp)
		//	t2 = tmp
		//	log.Printf("t1:%v, t2:%v\n", t1, t2)
		//}

		if tmp == "multi-line-or all-lines" {
			for _, v := range csvFilter {
				if v.Type == "AND" {
					d := FilterDiff{
						Msg:     fmt.Sprintf("csv filter type AND is not match to config multi-line-or all-lines"),
						CsvPara: v,
					}
					res = append(res, d)
				}
			}
		} else {
			confTmp[tmp] = nil
		}
	}
	//fmt.Printf("t1:%v, t2:%v\n", t1, t2)
	//log.Printf("group:%v, \ncsv:%v, \nconf:%v\n, csvTmp:%v\n, confTmp:%v\n-------\n", groupName, len(csvFilter), len(confFilter), len(csvTmp), len(confTmp))
	for k, v := range csvTmp {
		if _, ok := confTmp[v.FilterAction]; !ok {
			d := FilterDiff{
				Msg:     "filterAction not found in config",
				CsvPara: v,
			}
			res = append(res, d)
		} else {
			delete(csvTmp, k)
			delete(confTmp, k)
		}
	}

	for _, v := range csvTmp {
		d := FilterDiff{
			Msg:     "Only csv have filter action",
			CsvPara: v,
		}
		res = append(res, d)
	}
	for k, _ := range confTmp {
		d := FilterDiff{
			Msg:        "Only config have filter action",
			ConfigPara: k,
		}
		res = append(res, d)
	}
	return res
}

// compareActions
func compareActions() ActionReport {
	var rpt ActionReport
	rpt.ActionDiffs = []ActionDiff{}
	for k, _ := range csvChargingNameMap {
		if _, ok := configChargingActionMap[k]; ok {
			delete(csvChargingNameMap, k)
			delete(configChargingActionMap, k)
		}
	}
	for k, _ := range csvChargingNameMap {
		d := ActionDiff{
			Msg:             "Charging name only csv have",
			CsvChargingName: k,
		}
		rpt.ActionDiffs = append(rpt.ActionDiffs, d)

	}
	for k, _ := range configChargingActionMap {
		d := ActionDiff{
			Msg:                "Charging name only config have",
			ConfigChargingName: k,
		}
		rpt.ActionDiffs = append(rpt.ActionDiffs, d)
		// writeNewChargings(k, c)
	}
	return rpt
}

// writeNewChargings
func writeNewChargings(chargingName string, chargingAction []ConfigAction) {
	var s string
	s += "    charging-action " + chargingName + "\n"
	for _, v := range chargingAction {
		s += v.ActionPara
		s += "\n"
	}
	s += "    #exit" + "\n"
	_, err := chargingActionFile.WriteString(s)
	if err != nil {
		log.Printf("ERROR:%v\n", err.Error())
	}
}
