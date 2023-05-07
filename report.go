package main

import (
	"encoding/json"
	"log"
	"os"
)

type Report struct {
	UpReport        UpReport        `json:"upReport"`
	UpBindingReport UpBindingReport `json:"upBindingReport"`
	FilterReport    FilterReport    `json:"filterReport"`
	ActionReport    ActionReport    `json:"actionReport"`
}
type UpReport struct {
	GetCsvUpNum               int      `json:"getCsvUpNum"`
	GetConfigRulebaseNum      int      `json:"getConfigRulebaseNum"`
	OnlyCsvHaveUserProfile    []string `json:"onlyCsvHaveUserProfile"`
	OnlyConfigHaveUserProfile []string `json:"onlyConfigHaveUserProfile"`
}

type UpBindingReport struct {
	GetCsvUpBindingReportRuleNum int    `json:"getCsvUpBindingReportRuleNum"`
	GetConfigRuleNum             int    `json:"getConfigRuleNum"`
	Diffs                        []Diff `json:"diffs"`
}

type Diff struct {
	Line       int        `json:"line"`
	CsvPara    CsvRule    `json:"CsvRule"`
	ConfigPara ConfigRule `json:"ConfigRule"`
	Msg        string     `json:"msg"`
}
type FilterReport struct {
	CsvGroupCount        int          `json:"csvGroupCount"`
	CsvFilterActionCount int          `json:"csvFilterActionCount"`
	ConfigRuleDefCount   int          `json:"configRuleDefCount"`
	ConfigActionCount    int          `json:"configActionCount"`
	FilterDiffs          []FilterDiff `json:"filterDiffs"`
}

type FilterDiff struct {
	Msg        string    `json:"msg"`
	CsvPara    CsvFilter `json:"csvPara"`
	ConfigPara string    `json:"configPara"`
}
type ActionReport struct {
	CsvChargingNameCount    int          `json:"csvChargingCount"`
	ConfigChargingNameCount int          `json:"configChargingCount"`
	ActionDiffs             []ActionDiff `json:"actionDiffs"`
}

type ActionDiff struct {
	Msg                string `json:"msg"`
	CsvChargingName    string `json:"CsvChargingName"`
	ConfigChargingName string `json:"configChargingName"`
}

func (r *Report) WriteReport() {
	file, err := os.OpenFile("./data/report.txt", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	marshal, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(string(marshal))
	if err != nil {
		panic(err)
	}
	log.Printf("work finished success")
}
