package main

import (
	"encoding/json"
	"testing"
)

type cfg struct {
	Title string
	Key   string
	Value string
	Type  int32
}

func jsonStructToMap(stuObj interface{}) (map[string]interface{}, error) {
	// 结构体转json
	strRet, err := json.Marshal(stuObj)
	if err != nil {
		return nil, err
	}
	// json转map
	var mRet map[string]interface{}
	err1 := json.Unmarshal(strRet, &mRet)
	if err1 != nil {
		return nil, err1
	}
	return mRet, nil
}

func TestStructToMap(t *testing.T) {
	cfgM := &cfg{
		"sdgsdgs",
		"NewKey",
		"NewValue",
		2,
	}
	mRet, _ := jsonStructToMap(cfgM)
	t.Logf("struct convert to %v", mRet)
}
