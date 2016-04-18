package helpers

import (
	"fmt"
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	s "strings"
)

func DoRequest(sType string, sURL string, aParams [][2]string, aHeader [][2]string, sBody string) *http.Response {
	return DoRequestDebug(sType, sURL, aParams, aHeader, sBody, false)
}
func DoRequestDebug(sType string, sURL string, aParams [][2]string, aHeader [][2]string, sBody string, bDebug bool) *http.Response {
	if bDebug {
		fmt.Printf("header=%+v\n", aHeader)
	}
	cClient := &http.Client{}
	sParams := ""
	sAmp := ""
	if aParams != nil && len(aParams) > 0 {
		sURL += "?"
		for aP := range aParams {
			sParams += sAmp + aParams[aP][0] + "=" + aParams[aP][1]
			sAmp = "&"
		}
	}
	var cBody io.Reader
	if "" != sBody {
		if bDebug {
			fmt.Printf("BODY=%s\n", sBody)
		}
		aBody := []byte(sBody)
		cBody = bytes.NewBuffer(aBody)
	}
	req, err := http.NewRequest(sType, sURL+sParams, cBody)
	if bDebug {
		fmt.Printf("URL=%s\n", sURL+sParams)
	}
	for nI := range aHeader {
		req.Header.Add(aHeader[nI][0], aHeader[nI][1])
	}
	resp, err := cClient.Do(req)
	if err != nil {
		if bDebug {
			fmt.Printf("ERROR %v\n", err)
		}
		return nil
	}
	return resp
}
func ReaderToString(cSrc io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(cSrc)
	sSrc := buf.String()
	return sSrc
}
func GoFolderGet() string {
	var sGO = ""
	var err error
	if sGO, err = filepath.Abs("go"); err == nil {
		nIndx := s.Index(sGO, string(os.PathSeparator)+"go"+string(os.PathSeparator))
		if nIndx == -1 {
			fmt.Printf("======== PATH TO GO FOLDER = error detecting2\n")
			sGO = ""
		} else {
			sGO = sGO[0 : nIndx+4]
		}
	} else {
		fmt.Printf("======== PATH TO GO FOLDER = error detecting1\n")
	}
	return sGO
}
func BinFolderGet() string {
	var sBin = ""
	if sBin = GoFolderGet(); sBin != "" {
		sBin += "bin" + string(os.PathSeparator)
	}
	return sBin
}
