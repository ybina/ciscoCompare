package main

import "os"

var (
	// userprofile
	csvAllUserProfileArr []string
	csvAllUserProfileMap = make(map[string]interface{})
	userProfileCount     = 0

	//csv userprofile binding
	csvUserProfileBindingArr []CsvRule
	//key: RuleKey; val: CsvRule
	csvUpbMap = make(map[RuleKey]CsvRule)
	//key:USER_PROFILE; val:
	csvUPToRuleSliceMap = make(map[string][]CsvRule)
	csvRuleCount        = 0

	// csv filters
	csvFilterArr []CsvFilter
	// key: GroupName; val: []csvFilter
	csvFilterMap            = make(map[string][]CsvFilter)
	csvFilterActionCountMap = make(map[CsvFilterKey]interface{})
	csvFilterCount          = 0

	// csv actions
	csvActionArr []CsvAction
	// key: chargingName; val: chargingStruct
	csvChargingNameMap     = make(map[string][]CsvAction)
	csvChargingNameCount   = 0
	csvChargingActionCount = 0
	actionCount            = 0
)

// --------

var (
	//rulebase
	configUps     = make(map[string]interface{})
	ruleBaseCount = 0

	// rules
	configRuleCount = 0
	// key：userProfile
	configRuleMap = make(map[string][]ConfigRule)

	// ruledef
	ruleDefCount             = 0
	ruleDefFilterActionCount = 0
	configRuleDefMap         = make(map[string][]string)

	// chargingAction
	configChargingActionMap       = make(map[string][]ConfigAction)
	configChargingActionNameCount = 0
	configChargingActionCount     = 0
)

// -------------
var (
	chargingActionFile *os.File
)

func initFile() {
	f, err := os.OpenFile("./res/actionRes.txt", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	chargingActionFile = f
}
