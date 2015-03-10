package yaml

import (
	"strconv"
	"io/ioutil"
	yaml "gopkg.in/yaml.v2"
)

func ReadData(yamlPath string) (result map[string]map[string]interface{}, err error) {
	var dat []byte;
	dat, err = ioutil.ReadFile(yamlPath)
	result = make(map[string]map[string]interface{})
	yaml.Unmarshal(dat, &result)

	return;
}

func ReadTypeIds(path string) (result EveTypeIdList, err error) {
	var dat []byte;
	dat, err = ioutil.ReadFile(path)
	ids := make(map[string]EveTypeId)
	yaml.Unmarshal(dat, &ids)

	for key, _ := range ids {
		var val = ids[key]
		val.TypeId, _ = strconv.ParseInt(key, 10, 32)
		ids[key] = val
	}

	result = EveTypeIdList{list:ids};
	return;
}

func Save(path string, writer EveTypeIdWriter) {
	types, _ := ReadTypeIds(path)
	types.WriteTypeIds(writer)
}

func (types *EveTypeIdList) WriteTypeIds(writer EveTypeIdWriter) (err error) {
	for _, val := range types.list {
		err = writer.Write(val)
		if err != nil {
			return
		}
	}

	return
}


