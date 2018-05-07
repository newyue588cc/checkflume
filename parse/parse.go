package parse

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Unknwon/goconfig"
	"net/http"
	"io/ioutil"
	"strings"
)

var cfg *goconfig.ConfigFile

type Metrics struct {
	Host    string
	Metric 	map[string]string
}

func init() {
	config, err := goconfig.LoadConfigFile("config/flume_metrics.conf")
	if err != nil {
		log.Println("load config flume_metrics.conf faild!")
		os.Exit(-1)
	}
	cfg = config
}

// json translate to map
func JsontoMap(s []byte) (map[string]interface{},error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		return nil, err
	}
	return result, nil
}

// parse url and get channelFilll percentage value.
func UrlParse(url string) (map[string]interface{},error) {
	result := make(map[string]interface{})
	resp,err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	data,err := JsontoMap(body)
	if err != nil {
		log.Fatal(err)
	}

	for key,val := range data {
		if strings.HasPrefix(key,"CHANNEL") {
			temp := val.(map[string]interface{})
			result[key] = temp["ChannelFillPercentage"]
		}
	}
	return result,nil
}

// parse config file and return sec,key,value.
func ConfigParse() []Metrics {
	m := Metrics{}
	ms := []Metrics{}
	m.Metric = make(map[string]string)
	secs := cfg.GetSectionList()
	for _, sec := range secs {
		//fmt.Printf("host %v: \n",sec)
		m.Host = sec
		for _, key := range cfg.GetKeyList(sec) {
			url, err := cfg.GetValue(sec, key)
			//fmt.Printf("keys %v value %v \n",key,url)
			if err != nil {
				log.Fatal(err)
			}
			m.Metric[key] = url
		}
		ms = append(ms,m )
	}
	//fmt.Println("host is:",m.Host)
	return ms
}
