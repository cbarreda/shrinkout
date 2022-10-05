package pkgError

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type tError struct {
	SError string
	STime  string
}

type ClsError struct {
	ICode      int
	SInfo      string
	LocalError error
}

func (e *ClsError) Set(i int, s string) {
	e.ICode = i
	e.SInfo = s
}

func (obje *ClsError) Unmarshal(objInterface interface{}, sJSON string) error {
	if len(sJSON) > 0 {
		obje.SInfo = sJSON
	}
	obje.LocalError = json.Unmarshal([]byte(obje.SInfo), &objInterface)
	if obje.LocalError != nil {
		obje.ICode = -1
	}
	return obje.LocalError
}

func (obje *ClsError) JsonToMap(sJSON string) map[string]interface{} {
	m := make(map[string]interface{})
	if len(sJSON) == 0 {
		sJSON = obje.SInfo
	}
	obje.LocalError = json.Unmarshal([]byte(sJSON), &m)
	if obje.LocalError != nil {
		obje.ICode = -1
		obje.SInfo = sJSON
	}
	return m
}

func (obje *ClsError) SetDebug(iCode int, sString string) {
	obje.ICode = iCode
	obje.SInfo = sString
	obje.LocalError = errors.New(sString)
}

var (
	gPkgID    = "pkgError."
	gMapError map[tError]string
)

const (
	KLOGNONE       = -1
	KLOGDEBUG      = 0
	KLOGINFO       = 1
	KLOGWARN       = 2
	KLOGERROR      = 3
	KLOGPRINT      = 0X0f
	KCRITICALERROR = -10000
)

func ClearErrors() {
	gMapError = make(map[tError]string)
}

func VfInit() {
	ClearErrors()
}

func VfMark(sFunc string, iData ...interface{}) {
	sOutput := ""
	fmt.Println("-")
	fmt.Println("Func: " + sFunc)

	for _, v := range iData {
		sOutput = sOutput + fmt.Sprintf(", %v ", v)
	}
	fmt.Println(sOutput)
}

func VfDebug(sFunc string, iData ...interface{}) {
	sOutput := ""
	fmt.Println("-")
	fmt.Println("Func: " + sFunc)

	for _, v := range iData {
		sOutput = sOutput + fmt.Sprintf(", %v ", v)
	}
	fmt.Println(sOutput)
}

func VfPrintln(params ...interface{}) {
	sOutput := ""

	for _, v := range params {
		sOutput = sOutput + fmt.Sprintf("%v ", v)
	}
	fmt.Println(sOutput)
}

func MakeError(sError string) error {
	return errors.New(sError)
}

func BfHandleError(iLevel int, sPFunc string, localError error) bool {
	sFunc := gPkgID
	sFunc += "BfHandleError"
	bError := (iLevel > KLOGNONE) && (localError != nil)
	if bError {
		sErr := localError.Error()
		arPath := []string{"DEBUG.TXT", "INFO.TXT", "WARN.TXT", "ERROR.TXT"}
		strPath := "./logs/" + arPath[iLevel&0xf0]

		t := time.Now().Local()
		sTime := t.Format("2006/01/02 15:04:05")

		gMapError[tError{sPFunc, sTime}] = sErr

		logfile, _ := os.OpenFile(strPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		defer logfile.Close()

		logfile.WriteString(sTime + ":" + sPFunc + "-->" + sErr + "\n")

		if (iLevel & KLOGPRINT) > 0 {
			VfMark(sFunc, sTime+":"+sPFunc+"-->"+sErr)
		}
	}
	return !bError
}

func IfHandleError(iCode, iLevel int, sFunc, sError string) int {
	if iCode < 0 {
		BfHandleError(iLevel, sFunc, errors.New("Code:"+strconv.Itoa(iCode)+", Error: "+sError))
	}
	return iCode
}

func SfMaxString(iMax int, s string) string {
	if len(s) > iMax {
		return s[0:iMax]
	}
	return s
}
