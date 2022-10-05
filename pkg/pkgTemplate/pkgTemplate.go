package pkgTemplate

import (
	"io/ioutil"
	"os"
	"strings"
	

	"shrinkout/pkg/pkgError"
)

type templateId struct {
	sLanguage string
	sName     string
}


var (  
    gMapTokens	map[string]string
    gMapTemplates map[templateId]string
    gSLayout string
)

func SfGetTemplate(sLang, sFileName string) string {
	return gMapTemplates[templateId{sLang, sFileName}]
}

// path is path to /lang folder
func IfReadTemplate(sPath string) error {
	pkgError.VfDebug("pkgTemplate.IfReadTemplate","******************",sPath)
	var (
		fileDir []os.FileInfo
		err     error
		sData   []byte
	)

	gMapTemplates = make(map[templateId]string)

	fileDirFunc := func(sPath string) bool {
		fileDir, err = ioutil.ReadDir(sPath)

		return err == nil
	}

	if fileDirFunc(sPath) {
		langDir := fileDir
		for _, langFolder := range langDir {
			if fileDirFunc(sPath + "/" + langFolder.Name() + "/deploy") {
				for _, file := range fileDir {
					sFilePath := sPath + "/" + langFolder.Name() + "/deploy" + "/" + file.Name()
					sData, err = ioutil.ReadFile(sFilePath)
					if len(sData) > 0 {
						pkgError.VfPrintln("Read ", sFilePath)
						gMapTemplates[templateId{langFolder.Name(), file.Name()}] = string(sData)
					}
				}
			}
		}
	}

	return err
}

/* Added from original shrinkout */
func PstrTemplate(sLang,sName string)string{    	
	strReturn :=  gMapTemplates[templateId{sLang,sName}]
    if len(strReturn) > 0{
		for sKey,sValue := range(gMapTokens){
			strReturn = strings.Replace(strReturn,sKey,sValue,-1)
		}		
		strReturn = gSLayout + strReturn // gSLayout has the tokens inside the commented text so m_layout added here not at the start
	}
    return strReturn
}
