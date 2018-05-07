package main

import (
	"checkfm/parse"

	"log"
	"strconv"
	"math"
	"fmt"
)

const (
	User = "a@b.com"			// email sender user
	Password = "password"		// email sender password
	Server = "smtp.b.com"		// smtp server address
	Port = 465					// smtp server port
	ToAdd = "c@b.com"			// email receiver
	//ToCc = "d@b.com"			// 邮件抄送接收人
	Precision = 0.000001		// float64 精度
	Threshhold = 0.3			// 阀值
)

func main() {
	ec := &parse.EmailConfig{
		User:User,
		Passwd:Password,
		Server:Server,
		Port:Port,
	}
	mail := &parse.Mail{}
	mail.Email = ec
	mail.To = []string{ToAdd}

	metrics := parse.ConfigParse()
	for _,metric := range metrics {
		for _,url := range metric.Metric {
			result,err := parse.UrlParse(url)
			if err != nil {
				log.Fatal(err)
			}
			for k,v := range result {
				channelVal,err := strconv.ParseFloat(v.(string),64)
				if err != nil {
					log.Fatal(err)
				}
				if math.Dim(channelVal,Threshhold) > Precision {
					mail.Subject = fmt.Sprintf("%s of flume channel over threshhold in %s",k,metric.Host)
					mail.Body = fmt.Sprintf("The channel %s value is %f in %s,Please check it.",k,channelVal * 100,metric.Host)
					if err := mail.SendEmail();err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}
