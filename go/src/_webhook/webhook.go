package main

import (
	"encoding/json"
	"errors"
	"flag"
	. "helpers"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
	. "types"

	"github.com/ratsil/iniflags"
)

const (
	ResponseStatusSuccess = "success"
	ResponseStatusError   = "error"
	DepthNew              = 50
	DepthRevise           = 20
)

var (
	_sLogFolder      = flag.String("log_folder", "/var/log", "Log files folder")
	_sWebServerPort  = flag.String("wh_port", "1649", "Web Server Port")
	_sAsanaUserToken = flag.String("wh_asana_token", "", "access token for asana")
)

func init() {
	flag.Set("config", "preferences.ini")
	iniflags.Parse()
}

func main() {
	sLog := filepath.Join(*_sLogFolder, "webhook_"+time.Now().Format("2006_01_02")+".log")
	f, err := os.OpenFile(sLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error opening log file: %v", err)
		return
	}
	defer f.Close()
	log.SetOutput(f)

	log.Printf("======== WEB PORT = %v\n", ":"+*_sWebServerPort)
	server := &http.Server{
		Addr:           ":" + *_sWebServerPort,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	http.HandleFunc("/", Handler)

	if err = server.ListenAndServe(); nil != err {
		log.Printf("main(): %s\n", err)
	}
	log.Printf("======== THE END\n")
}

func Handler(oResponseWriter http.ResponseWriter, pRequest *http.Request) {
	log.Println("request detected")
	oResponseWriter.Header().Add("Content-Type", "application/json")

	log.Println("-----------------------request start")
	err := Request(pRequest)
	log.Println("-----------------------request end")

	if nil != err {
		log.Println("ERROR during request handling: " + err.Error())
		return
	}
}

func Request(pRequest *http.Request) (err error) {
	sHeader := pRequest.Header["User-Agent"][0]
	log.Printf("request HEADER.USER_AGENT is: [%s]", sHeader)
	if bMatch, _ := regexp.MatchString("(?i)Bitbucket(?-i)", sHeader); bMatch == false {
		err = errors.New("ERROR  - This isn't a BitBucket request! The end.")
		return
	}

	//	log.Printf("request BODY is: \n%+v\n", pRequest.Body)

	pBBI := &BitBucketRequest{}
	if err = json.NewDecoder(pRequest.Body).Decode(pBBI); err != nil {
		log.Printf("ERROR-1 = %+v\n", err)
		return
	}

	for _, oChange := range pBBI.Push.Changes {
		oTarget := oChange.New.Target
		if "commit" != oTarget.Type {
			continue
		}
		log.Print(regexp.MustCompile("(?:#(\\d+))").FindAllStringSubmatch(oTarget.Message, -1))
		for _, aTaskID := range regexp.MustCompile("(?:#(\\d+))").FindAllStringSubmatch(oTarget.Message, -1) {
			log.Print("task:" + aTaskID[1])
			PostCommentToAsana(aTaskID[1], oTarget.Author.User.Username+" just committed with message:\n"+oTarget.Message+"\nto repository "+pBBI.Actor.Links.HTML.Href+pBBI.Repository.Name)
		}
	}

	//log.Printf("BitBucketInfo = %+v\n", pBBI)
	return
}
func PostCommentToAsana(sTaskID, sText string) error {
	log.Print("--------------Request to Asana--------------\n")
	var aParams [][2]string = nil //[]aPair{{"tagged", sTag}, {"__a", "1"},}
	sRequestData := "text=" + sText
	aHeader := [][2]string{
		{"Accept", "*/*"},
		{"Accept-Encoding", ""},
		{"User-Agent", "Mozilla/5.0"},
		{"Cache-Control", "no-cache"},
		{"Connection", "Keep-Alive"},
		{"Host", "app.asana.com"},
		{"Content-Type", "application/x-www-form-urlencoded"},
		{"Connection", "Keep-Alive"},
		{"Content-Length", strconv.Itoa(len(sRequestData))},
		{"Authorization", "Bearer " + *_sAsanaUserToken},
	}
	pResp := DoRequestDebug("POST", "https://app.asana.com/api/1.0/tasks/"+sTaskID+"/stories", aParams, aHeader, sRequestData, false)
	sBody := ReaderToString(pResp.Body)
	if pResp.StatusCode != 200 && pResp.StatusCode != 201 {
		log.Printf("\t\tERROR  - STATUS IS NOT 200 or 201 = %v\n", pResp.Status)
		log.Printf("\t\tBODY is\n%+v\n\n", sBody)
		return errors.New("ERROR  - STATUS IS NOT 200 or 201 = " + string(pResp.StatusCode))
	}
	log.Printf("BODY is [%+v]\n", sBody)
	log.Print("--------------Request to Asana was OK--------------\n\n")
	return nil
}
