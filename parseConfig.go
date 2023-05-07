package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	lineCount = 1
	// 用于匹配config.rulebase中的规则
	ruleBaseStart = regexp.MustCompile("rulebase (.*)")
	ruleBaseEnd   = regexp.MustCompile("#exit")
	endCount      = 0
	ruleReg       = regexp.MustCompile("action priority .*")

	// 用于匹配ruledef中的filterAction
	ruleDefStart = regexp.MustCompile("ruledef (.*)")

	ruleDefEnd = regexp.MustCompile("#exit")

	// 用于匹配charingAction
	s1                  = `charging-action (.*)`
	chargingActionStart = regexp.MustCompile(s1)
	chargingActionEnd   = regexp.MustCompile("#exit")

	// 用于匹配group-of-ruledefs
	groupOfRuledefsStart = regexp.MustCompile(`group-of-ruledefs (.*)`)
)

func parseConfig() {
	//testParseConfig.txt
	//20221109_SPGW_HER1.txt
	file, err := os.OpenFile("./data/config.txt", os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		//log.Printf("l:%v\n", string(line))
		//if strings.Contains(string(line), "CLARO-CENAM-RG11") {
		//	log.Printf("---%v;%v\n", ruleBaseStart.MatchString(strings.TrimLeft(string(line), " 	")), string(line))
		//
		//}
		//log.Printf("l:%v\n", string(line))
		if strings.Contains(string(line), "IP-FLAT-CA_Pospago") {
			log.Printf("------%v\n", string(line))
		}
		if ruleBaseStart.MatchString(strings.TrimLeft(string(line), " ")) {
			if strings.Contains(string(line), "active-charging") {
				continue
			}
			//fmt.Printf("ruleBase start:%v\n", string(line))
			if strings.Contains(string(line), "schema") {
				continue
			}
			up := strings.SplitN(strings.TrimLeft(string(line), " "), " ", -1)
			//for i, v := range up {
			//
			//}
			startParseConfigRules(reader, up[1])
			continue
		}
		//if strings.Contains(string(line), "F-IP-INT-PENTCLOUD-DNS-SNOOPING-RD120") {
		//
		//	log.Printf("---- %v; %v\n", ruleDefStart.MatchString(strings.Trim(string(line), " ")), string(line))
		//	fmt.Printf("%v\n", strings.SplitN(strings.TrimLeft(string(line), ""), " ", -1))
		//}
		if ruleDefStart.MatchString(strings.TrimLeft(string(line), "	 ")) {
			ruleDef := strings.SplitN(strings.TrimLeft(string(line), " "), " ", -1)
			actions := startParseRuleDefines(reader)
			configRuleDefMap[ruleDef[1]] = actions
			continue

		}

		if chargingActionStart.MatchString(strings.TrimLeft(string(line), "	 ")) {
			chargingName := strings.SplitN(strings.TrimLeft(string(line), " "), " ", -1)
			if chargingName[0] != "charging-action" {
				continue
			}
			actions := startParseChargingActions(reader, chargingName[1])
			configChargingActionMap[chargingName[1]] = actions
			continue
		}

	}
	fmt.Printf("configUps len:%v\n", len(configUps))
	fmt.Printf("ruleBaseCount:%v, endCount:%v, configRuleCount:%v\n", ruleBaseCount, endCount, configRuleCount)
}

func startParseConfigRules(reader *bufio.Reader, up string) {

	//fmt.Printf("config up:%v\n", up[1])
	//log.Printf("rulebase:%v\n", string(line))

	if _, ok := configUps[up]; ok {
		log.Printf("config up is duplicate:%v\n", up[1])
	} else {
		configUps[up] = nil
		ruleBaseCount++
	}
	parseConfigRule(reader, up)
}

var parseConfigRuleInCount = 0

func parseConfigRule(f *bufio.Reader, up string) {
	parseConfigRuleInCount++
	//fmt.Printf("%v parseConfigRule:%v\n", parseConfigRuleInCount, up)
	for {
		line, _, err := f.ReadLine()
		if err == io.EOF {
			return
		}
		if ruleBaseEnd.MatchString(strings.TrimLeft(string(line), "	 ")) {
			endCount++
			return
		}

		if ruleReg.MatchString(strings.TrimLeft(string(line), "	 ")) {
			configRuleCount++
			//fmt.Printf("action:%v\n", string(line))
			addConfigRules(string(line), up)
		}
	}
}

// CLARO-COSTA_RICA-INTERNET
// CLARO-COSTA_RICA-INTERNET
func addConfigRules(s string, up string) {
	ruleArr, ok := configRuleMap[up]
	if !ok {
		arr := []ConfigRule{unmarshalRule(s)}
		configRuleMap[up] = arr
	} else {
		ruleArr = append(ruleArr, unmarshalRule(s))
		configRuleMap[up] = ruleArr
	}
}

func unmarshalRule(s string) ConfigRule {
	s = strings.TrimLeft(s, " ")
	arr := strings.SplitN(s, " ", -1)
	//fmt.Printf("reg arr:%v\n", arr)
	c := ConfigRule{}
	for i := 0; i < len(arr); i++ {
		switch arr[i] {
		case "dynamic-only":
			c.IsDynamic = true
		case "group-of-ruledefs":
			c.GroupOfRuleDefs = arr[i+1]
			c.RuleName = arr[i+1]
		case "ruledef":
			c.RuleDef = arr[i+1]
			c.RuleName = arr[i+1]
		case "charging-action":
			c.ChargingAction = arr[i+1]
		case "monitoring-key":
			//log.Printf("---- key:%v\n", arr[i+1])
			c.MonitoringKey = arr[i+1]
		case "priority":

			//log.Printf("priority:%v\n", arr[i+1])
			c.Priority = arr[i+1]
			//if c.Priority == "804" {
			//	fmt.Printf("804 arr:%v\n", arr[i:])
			//}
		default:
		}
	}
	//if arr[len(arr)-2] != "monitoring-key" {
	//	c.MonitoringKey = "0"
	//}
	//fmt.Printf("configRule:%v\n", c)
	if c.MonitoringKey == "" {
		c.MonitoringKey = "0"
	}
	return c
}

func startParseRuleDefines(reader *bufio.Reader) []string {
	//log.Printf("startParseRuleDefines in\n")
	ruleDefCount++
	var actions []string
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			panic(err)
		}
		if ruleDefEnd.MatchString(strings.TrimLeft(string(line), "	 ")) {
			break
		}
		actions = append(actions, strings.TrimLeft(string(line), " "))
		ruleDefFilterActionCount++
		//log.Printf("line:%v\n", string(line))
	}
	return actions
}

func startParseChargingActions(reader *bufio.Reader, chargingName string) []ConfigAction {
	var res []ConfigAction
	configChargingActionNameCount++

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			panic(err)
		}
		if chargingActionEnd.MatchString(strings.TrimLeft(string(line), " 	")) {
			break
		}
		c := ConfigAction{
			ActionName: chargingName,
			ActionPara: string(line),
		}
		configChargingActionCount++
		res = append(res, c)
	}
	return res
}
