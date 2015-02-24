package main

import (
	"encoding/xml"
	"encoding/json"
	"fmt"
	"os/exec"
	"io"
	"io/ioutil"
	"net/http"
)

type ByteWriter []byte

func (writer *ByteWriter) Write(p []byte) (n int, err error) {
	*writer = append(*writer, p...)
	n = len(*writer)
	return
}

// http://golang.org/pkg/encoding/xml/#example_Unmarshal
type VNStat struct {
	XMLName xml.Name `xml:"vnstat" json:"-"`
	Interfaces []VNStatInterface `xml:"interface"`
}

func (vnstat *VNStat) MarshalJSON() ([]byte, error) {
	for idxi, in := range vnstat.Interfaces {
		for idxd, ds := range in.Traffic.Days {
			formattedDate := fmt.Sprintf("%s-%s-%s", ds.VNSDate.Year, ds.VNSDate.Month, ds.VNSDate.Day)
			vnstat.Interfaces[idxi].Traffic.Days[idxd].Date = formattedDate
		}
	}

	data, err := json.MarshalIndent(*vnstat, "", "    ")
	return data, err
}

type VNStatInterface struct {
	XMLName xml.Name `xml:"interface" json:"-"`
	Id string `xml:"id,attr"`
	Traffic VNStatTraffic `xml:"traffic"`
}

type VNStatTraffic struct {
	XMLName xml.Name `xml:"traffic" json:"-"`
	Days []VNStatDay `xml:"days>day"`
}

type VNStatDay struct {
	XMLName xml.Name `xml:"day" json:"-"`
	VNSDate VNStatDate `xml:"date" json:"-"`
	Id int64 `xml:"id,attr"`
	Date string
	Rx int64 `xml:"rx"`
	Tx int64 `xml:"tx"`
}

type VNStatDate struct {
	XMLName xml.Name `xml:"date"`
	Year string `xml:"year"`
	Month string `xml:"month"`
	Day string `xml:"day"`
}

func main() {
	http.HandleFunc("/", servePage)
	http.HandleFunc("/stats", serveStat)
	http.ListenAndServe(":4545", nil)
}

func servePage(w http.ResponseWriter, _ *http.Request) {
	fileData, err := ioutil.ReadFile("vnstat.html")
	if err != nil {
		http.Error(w, "SERVER_ERROR", http.StatusInternalServerError)
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	w.Write(fileData)
}

func serveStat(w http.ResponseWriter, _ *http.Request) {
	cmd := exec.Command("vnstat", "-i", "eth0", "-d", "--xml")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, "COMMAND_ERROR", http.StatusInternalServerError)
		return
	}

	err = cmd.Start()
	if err != nil {
		http.Error(w, "COMMAND_ERROR", http.StatusInternalServerError)
		return
	}
	defer cmd.Wait()

	ios := new(ByteWriter)
	io.Copy(ios, stdout)

	v := new(VNStat)
	err = xml.Unmarshal(*ios, &v)
	if err != nil {
		http.Error(w, "XML_ERROR", http.StatusInternalServerError)
		return
	}

	data, err := v.MarshalJSON()
	if err != nil {
		http.Error(w, "JSON_ERROR", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
