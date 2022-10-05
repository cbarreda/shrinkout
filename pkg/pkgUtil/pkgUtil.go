package pkgUtil

import (
	"shrinkout/pkg/pkgError"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const ksVersion = "V 2.0.0"
const ksDate = "01/04/2022"


type ClsDbParam struct {
	Key   string
	DbStr string
}

type ClsConfigNode struct {
	Active     bool
	Name       string
	Url        string
	PortStr    string
	ListenPort string
	Database   ClsDbParam
}

var gPkgId = "pkUtil."

func (objConfig ClsConfigNode) SfGetNetwork() string {
	return objConfig.Url + objConfig.PortStr + "/"
}

func (objConfig ClsConfigNode) SfConfigNodeDisplay() string {
	return ("-------------\nLive:" + strconv.FormatBool(objConfig.Name == "Live") + "\nName:" + objConfig.Name +
		"\nURL:" + objConfig.Url + "\nPortStr:" + objConfig.PortStr + "\nListenPort:" + objConfig.PortStr +
		"\nDatabase:" + objConfig.Database.DbStr + "\n-------------")
}

func SVersion()(string, string){
	return "PkgUtil Version: " + ksVersion, ", " + ksDate
}

func SRandomString(iLen int) string {
	rand.Seed(time.Now().UnixNano())
	sLetters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	i := len(sLetters)
	sReturn := ""

	for z := 0; z < iLen; z++ {
		sReturn += string(sLetters[rand.Intn(i)])
	}

	return sReturn
}

func Atoi(sString string, iErrorReturn int) int {
	iCode, le := strconv.Atoi(sString)
	if le != nil {
		iCode = iErrorReturn
	}
	return iCode
}

// true, nil file exists. true, false, nil file does not exist. false, !nil file does not exist, error other than not exist
func FileExists(sPath string) (bool, error) {
	_, err := os.Stat(sPath)
	return (err == nil), err
}

// if Dir exists or successfully created it returns nil. Error if dir exists but error !- os.IsNotExist
func DirMakeIfNotExists(sDirPath string, iMode os.FileMode) error {

	bOK, localError := FileExists(sDirPath)

	if !bOK {
		localError = os.MkdirAll(sDirPath, iMode)
	}

	return localError
}

func SaveMultipartFile(file multipart.File, sPath string) error {
	sFunc := gPkgId
	sFunc += "SaveMultipartFile"
	var localError error

	sDir, _ := path.Split(sPath)
	if len(sDir) > 1 {
		localError = DirMakeIfNotExists(sDir, 0777)

		if localError == nil {
			data, err := ioutil.ReadAll(file)
			localError = err
			if localError == nil {
				localError = ioutil.WriteFile(sPath, data, 0666)
			}
		}
	} else {
		localError = errors.New("No File path supplied")
	}
	return localError
}

func GetFile(sPath string) ([]byte, error) {
	sFunc := gPkgId
	sFunc += "GetFile"

	var data []byte
	bExists, localError := FileExists(sPath)

	if bExists {
		return ioutil.ReadFile(sPath)
	}

	return data, localError
}

func FolderDel(sPath string) error {
	sFunc := gPkgId
	sFunc += "FolderDel"
	localError := os.RemoveAll(sPath)
	return localError
}

// Config File start

func sfConfigGetLine(sMatch, sStr string) string {
	sReturn := ""
	sStr = strings.TrimSpace(sStr)
	if len(sStr) > len(sMatch) {
		if sMatch == sStr[0:len(sMatch)] {
			sReturn = strings.TrimSpace(sStr[len(sMatch):])
		}

	}
	return sReturn
}

/* Reads config and returns line whos begining match sTokenString.
If sTokenString is :url it will return a line that starts with :url but not a line that starts with ;:url */
func ReadConfigFile(sFilePath, sTokenString string) ([]string, error) {
	sFunc := gPkgId
	sFunc += "ReadConfigFile"

	var lines []string
	file, err := os.Open(sFilePath)
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			sLine := sfConfigGetLine(sTokenString, scanner.Text())
			if len(sLine) > 0 {
				lines = append(lines, sLine)
			}
		}
		if len(lines) < 1 {
			err = errors.New("empty file or no matching lines")
		}
	}
	return lines, err
}

func FileToArLines(sConfigPath string) ([]string, error) {
	var lines []string
	file, err := os.Open(sConfigPath)
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

	}
	return lines, err
}

func ArLinesToFile(sConfigPath string, sLines []string) error {

	file, err := os.Create(sConfigPath)
	if err == nil {
		defer file.Close()
		fWriter := bufio.NewWriter(file)
		for _, sLine := range sLines {
			fWriter.WriteString(sLine + "\n")
		}
		err = fWriter.Flush()

	}
	return err
}

func ReplaceWhiteSpace(s string, c string) string {
	sReturn := ""
	sLastChar := ""
	for _, v := range s {
		if unicode.IsSpace(v) {
			if sLastChar != c {
				sLastChar = c
				sReturn += sLastChar
			}
		} else {
			sLastChar = string(v)
			sReturn += sLastChar
		}
	}
	return sReturn
}

func BSNeedleInSHayStack(sNeedle string, arHayStack []string) bool {
	for _, v := range arHayStack {
		if v == sNeedle {
			return true
		}
	}
	return false
}

func BfIsImageExtensionValid(sExtension string) bool {
	sExtension = strings.ToLower(sExtension)
	ar := []string{"jpg", "jpeg", "png"}
	for _, v := range ar {
		if sExtension == v {
			return true
		}
	}
	return false
}

func SfFileExtension(sFile string) string {
	arAr := strings.Split(sFile, ".")

	if len(arAr) > 0 {
		return arAr[len(arAr)-1]
	}

	return ""
}

func WildCardDel(sPath string) (int, error) {

	iCode := 0
	files, localError := filepath.Glob(sPath)
	if localError == nil {
		for _, f := range files {
			iCode++
			localError = os.Remove(f)
			if localError != nil {
				break
			}
		}
	}
	if localError != nil {
		iCode = -1
	}
	return iCode, localError

}

func IOFileCopy(sSource, sDest string) (int, error) {

	objSource, localError := os.Open(sSource)
	if localError != nil {
		return -1, localError
	}
	defer objSource.Close()

	objDest, localError := os.OpenFile(sDest, os.O_RDWR|os.O_CREATE, 0666)
	if localError != nil {
		return -2, localError
	}
	defer objDest.Close()

	_, localError = io.Copy(objDest, objSource)
	if localError != nil {
		return -3, localError
	}
	localError = objDest.Sync()
	if localError != nil {
		return -4, localError
	}
	return 0, nil

}

func DataAsStringGet(sJSON, sLeftDelim, sRightDelim, sStripChar string) string {
	sJSON = strings.TrimSpace(sJSON)[strings.Index(sJSON, `"Data":`)+len(`"Data":`) : len(sJSON)-1]

	if len(sStripChar) > 0 {
		sJSON = strings.Trim(sJSON, sStripChar)
	}
	return sLeftDelim + sJSON + sRightDelim
}

func FileToBrowserDownload(w http.ResponseWriter, r *http.Request, sFilePath string, sFile string) {
	downloadBytes, localError := ioutil.ReadFile(sFilePath)

	if localError == nil {
		// set the default MIME type to send
		mime := http.DetectContentType(downloadBytes)

		fileSize := len(string(downloadBytes))

		// Generate the server headers
		w.Header().Set("Content-Type", mime)
		w.Header().Set("Content-Disposition", "attachment; filename="+sFile+"")
		w.Header().Set("Expires", "0")
		w.Header().Set("Content-Transfer-Encoding", "binary")
		w.Header().Set("Content-Length", strconv.Itoa(fileSize))
		w.Header().Set("Content-Control", "private, no-transform, no-store, must-revalidate")

		//b := bytes.NewBuffer(downloadBytes)
		//if _, err := b.WriteTo(w); err != nil {
		//              fmt.Fprintf(w, "%s", err)
		//      }

		// force it down the client's.....
		http.ServeContent(w, r, sFilePath, time.Now(), bytes.NewReader(downloadBytes))
	}

}

func SfQueryGet(r *http.Request, sQueryFor string) string {
	dat, bOk := r.URL.Query()[sQueryFor]

	if bOk {
		return dat[0]
	}

	return ""
}

func ArQueryUrl(r *http.Request, queryFor []string) []string {
	var arResult []string

	for _, v := range queryFor {
		arResult = append(arResult, SfQueryGet(r, v))
	}
	return arResult
}

func Init() (ClsConfigNode, int, error) {
	var (
		arConfigNodes []ClsConfigNode
		objConfigNode ClsConfigNode
	)

	sFunc := gPkgId
	sFunc += "Init"

	iMode := 0

	configLines, localError := GetFile("./go.config")
	if localError == nil {
		localError = json.Unmarshal([]byte(configLines), &arConfigNodes)
		if localError == nil {
			for _, v := range arConfigNodes {
				if v.Active == true {
					objConfigNode = v
					break
				}
			}
			if len(objConfigNode.Name) < 1 {
				panic("No Active node in go.config")
			} else {
				switch objConfigNode.Name {
				case "Live":
					iMode = 0
					break
				case "TestOnServer":
					iMode = 1
					break
				default:
					iMode = 2 //localhost
					break
				}
			}
		}
	}
	if localError != nil {
		fmt.Println("--------------------------- Initialize Error/n", localError)
	}
	return objConfigNode, iMode, localError
}

func SendMail(sTo, sSubject, sBody string) pkgError.ClsError {
	var objError pkgError.ClsError
	sTo = strings.ToLower(strings.TrimSpace(sTo))

	if len(sTo) > len("@test.com") {
		if sTo[len(sTo)-len("@test.com"):] == "@test.com" {
			objError.ICode = -1
		}
	}

	if objError.ICode > -1 {
		sPass := "ZelfPublish!Pa"
		sFrom := "cbarreda@zelfpublish.com"

		sMsg := "From: " + sFrom + "\n" +
			"To: " + sTo + "\n" +
			"Subject:" + sSubject + "\n\n" +
			sBody

		objError.LocalError = smtp.SendMail("smtpout.secureserver.net:25",
			smtp.PlainAuth("", sFrom, sPass, "smtpout.secureserver.net"),
			sFrom, []string{sTo}, []byte(sMsg))
		fmt.Println(objError)
	} else {
		fmt.Println("SendMail: sTo: " + sTo + ", sSubject: " + sSubject + ", sBody: " + sBody)
	}
	return objError
}
