package yaml

import (
	"testing"
	"encoding/json"
	// "fmt"
)

func TestAverage(t *testing.T) {
	data, _ := ReadTypeIds("/vagrant/projects/eve//Proteus_1.0_109795_db/typeIds.yaml");

	keys := data.Keys();
	if (len(keys) < 10 ) {
		t.Error("Expected more.")
	}

	js, _ := json.Marshal(data.list["54"])
	if (data.list["54"].TypeId != 54 ) {
		t.Error("Incorrect Id.")
	}
	//fmt.Println(string(js))

	js, _ = json.Marshal(data.list["497"])
	//fmt.Println(string(js))

	js, _ = json.Marshal(data.list["582"])
	if (string(js) != `{"TypeId":582,"GraphicId":38,"Radius":27,"SoundId":20070,"IconId":0,"SofFactionName":"","FactionId":500001,"Masteries":{"0":[96,139,85,87,94],"1":[96,139,85,87,94],"2":[96,139,85,87,94],"3":[96,139,85,87,94],"4":[96,139,85,118,87,94]},"Traits":{"-1":{"1":{"Bonus":500,"BonusText":"bonus to \u003ca href=showinfo:3422\u003eRemote Shield Booster\u003c/a\u003e transfer range","UnitId":105}},"3330":{"1":{"Bonus":10,"BonusText":"bonus to \u003ca href=showinfo:3422\u003eRemote Shield Booster\u003c/a\u003e amount","UnitId":105},"2":{"Bonus":10,"BonusText":"reduction in \u003ca href=showinfo:3422\u003eRemote Shield Booster\u003c/a\u003e activation cost","UnitId":105}}}}`) {
		t.Error("Incorrect JSON.")
	}

	//fmt.Println(data.list["582"])
}
