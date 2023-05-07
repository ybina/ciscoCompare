package main

import (
	"encoding/json"
)

type RuleKey struct {
	UserProfile string
	RuleName    string
}

type CsvRule struct {
	UserProfile string `json:"userProfile"`
	OrderNo     string `json:"orderNo"`
	// false: 本地rule，不需要pcrf下发；
	// true: 动态rule，需要pcrf下发，对应config中的“dynamic-only”字段
	InstalledByPcrf bool   `json:"installedByPcrf"`
	RuleName        string `json:"ruleName"`
	ChargingAction  string `json:"chargingAction"`
	MonitoringKey   string `json:"monitoringKey"`
}

func (c *CsvRule) prt() string {
	marshal, _ := json.Marshal(c)
	return string(marshal)
}

type CsvFilter struct {
	RuleName     string
	GroupName    string
	Type         string
	FilterAction string
}
type CsvFilterKey struct {
	GroupName    string
	Type         string
	FilterAction string
}
type CsvAction struct {
	ChargingName   string
	ChargingAction string
	ChargingType   string
}

//----------

type RuleBase struct {
	RuleBaseName string
	RuleArr      []ConfigRule
}

type ConfigRule struct {
	Priority        string `json:"priority"`
	IsDynamic       bool   `json:"isDynamic"`
	GroupOfRuleDefs string `json:"groupOfRuleDefs"`
	RuleDef         string `json:"ruleDef"`
	RuleName        string `json:"ruleName"`
	ChargingAction  string `json:"chargingAction"`
	MonitoringKey   string `json:"monitoringKey"`
}

type ConfigFilerAction struct {
	Actions    []string
	IsAllLines bool
}

type ConfigAction struct {
	ActionName string
	ActionPara string
}

func (c *ConfigRule) prt() string {
	marshal, _ := json.Marshal(c)
	return string(marshal)
}
