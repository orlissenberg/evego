package main

import (
	"encoding/xml"
	"encoding/json"
	"fmt"
	"os/exec"
	"io"
	"io/ioutil"
	"net/http"
	"flag"
	"errors"
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
		for idxh, h := range in.Traffic.Hours {			
			vnstat.Interfaces[idxi].Traffic.Hours[idxh].Date = h.VNSDate.Iso8601String()
		}

		for idxd, ds := range in.Traffic.Days {			
			vnstat.Interfaces[idxi].Traffic.Days[idxd].Date = ds.VNSDate.Iso8601String()
		}

		for idxm, m := range in.Traffic.Months {			
			vnstat.Interfaces[idxi].Traffic.Months[idxm].Date = m.VNSDate.Iso8601String()
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
	Months []VNStatMonth `xml:"months>month"`
	Hours []VNStatHour `xml:"hours>hour"`
}

type VNStatHour struct {
	XMLName xml.Name `xml:"hour" json:"-"`	
	VNSDate VNStatDate `xml:"date" json:"-"`
	Id int64 `xml:"id,attr"`
	Date string
	Rx int64 `xml:"rx"`
	Tx int64 `xml:"tx"`
}

type VNStatDay struct {
	XMLName xml.Name `xml:"day" json:"-"`
	VNSDate VNStatDate `xml:"date" json:"-"`
	Id int64 `xml:"id,attr"`
	Date string
	Rx int64 `xml:"rx"`
	Tx int64 `xml:"tx"`
}

type VNStatMonth struct {
	XMLName xml.Name `xml:"month" json:"-"`	
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

func (d *VNStatDate) Iso8601String() string {
	formattedDate := fmt.Sprintf("%s-%s-%s", d.Year, d.Month, d.Day)	

	if d.Day == "" {
		formattedDate = fmt.Sprintf("%s-%s", d.Year, d.Month)
	}

	return formattedDate
}

func main() {
	port := flag.Int("port", 4545, "Server port.")
	debug := flag.Bool("debug", false, "Debug.")
	flag.Parse()

	if *debug {
		data, _ := getVnStatData()
		fmt.Println(string(data))
		return
	}

	http.HandleFunc("/", servePage)
	http.HandleFunc("/stats", serveStat)
	
	p := fmt.Sprintf(":%d", *port)
	fmt.Printf("Server at port %d.\n", *port)
	err := http.ListenAndServe(p, nil)
	if err != nil {
		fmt.Println(err)
	}
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

func getVnStatData() (data []byte, err error) {
	cmd := exec.Command("vnstat", "-i", "eth0", "-d", "--xml")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		err = errors.New("COMMAND_ERROR")		
		return
	}

	err = cmd.Start()
	if err != nil {
		err = errors.New("COMMAND_ERROR")
		return
	}
	defer cmd.Wait()

	ios := new(ByteWriter)
	io.Copy(ios, stdout)

	v := new(VNStat)
	err = xml.Unmarshal(*ios, &v)
	if err != nil {
		err = errors.New("XML_ERROR")		
		return
	}

	data, err = v.MarshalJSON()
	if err != nil {
		err = errors.New("JSON_ERROR")
		return
	}

	return
}

func serveStat(w http.ResponseWriter, _ *http.Request) {
	data, err := getVnStatData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
