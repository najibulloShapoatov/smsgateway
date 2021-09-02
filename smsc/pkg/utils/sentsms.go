package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"smsc/pkg/config"
	"smsc/pkg/log"
)

//SentSms ....
func SentSms(From, PhoneNo, SmsText string) int {

	defer func() {
		if err := recover(); err != nil {
			//log.Error("panic SentSms(PhoneNo, SmsText string) occurred:", err)
			SentSms(From, PhoneNo, SmsText)
		}
	}()

	cfg := config.GetKannelCfg()

	var link = cfg.Link + "?username=" + cfg.UserName + "&password=" + cfg.Password + "&charset=utf-8&smsc=" + cfg.SmsC + "&to=" + PhoneNo + "&from=" + From + "&text=" + url.QueryEscape(SmsText)

	log.Info("LINK ->", link)

	client := &http.Client{}

	req, err := http.NewRequest("GET", cfg.Link, nil)
	if err != nil {
		log.Error("sending SMS error", err)
	}
	q := url.Values{}
	q.Add("username", cfg.UserName)
	q.Add("password", cfg.Password)
	q.Add("charset", "utf-8")
	q.Add("coding", "2")
	q.Add("from", From)
	q.Add("smsc", cfg.SmsC)
	q.Add("to", PhoneNo)
	q.Add("text", SmsText)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Warn("client.Do(req) ->\n", err)
	}
	log.Info("requested to kannel ", resp.StatusCode, resp.Status)
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		log.Info("\t>>>sms sent>>> status = " + fmt.Sprint(resp.StatusCode) + "\t ")
		return 1
	}
	log.Error("sms don`t sent to -> ", PhoneNo)

	return 0

}
