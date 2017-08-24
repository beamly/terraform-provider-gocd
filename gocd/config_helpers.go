package gocd

import (
	"encoding/json"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"reflect"
	"sort"
	"strconv"
)

// Give an abstract list of strings cast as []interface{}, convert them back to []string{}.
func decodeConfigStringList(lI []interface{}) []string {

	if len(lI) == 1 {
		return []string{lI[0].(string)}
	}
	ret := make([]string, len(lI))
	for i, vI := range lI {
		ret[i] = vI.(string)
	}
	sort.Sort(sort.StringSlice(ret))
	return ret
}

// Take our object we parsed from the TF resource, and encode it in JSON.
func definitionDocFinish(d *schema.ResourceData, r interface{}) error {
	jsonDoc, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	jsonString := string(jsonDoc)
	d.Set("json", jsonString)
	d.SetId(strconv.Itoa(hashcode.String(jsonString)))

	return nil
}

func supressJSONDiffs(k, old, new string, d *schema.ResourceData) bool {
	var j1, j2 interface{}
	if (old == "" || new == "") && old != new {
		return false
	}
	if err := json.Unmarshal([]byte(old), &j1); err != nil {
		panic(err)
	}
	if err := json.Unmarshal([]byte(new), &j2); err != nil {
		panic(err)
	}
	return reflect.DeepEqual(j2, j1)
}
