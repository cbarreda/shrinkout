/*
 Assumes the following:
 masterTemplates are in the ./masterTemplates folder
	masterTemplates, like masterDefault.tmpl, are used by one or more pages.
	Have things like:
		- host, db, etc
		- 	<body>{{template "defaultNavBar" .}}{{template "defaultContent" .}}</body>
 templates are in the ./templates/appname/ folder
	This templates will override the master template if they have the same name
ConfigFiles are in ./lang/english/ and each page has a template and a config file.
    templates that don't have a corresponding config file are not built.
	Note: currently build process does not differentiate build files for different apps.
		Make sure each config file and template file are unique for all pages.

*/

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"unicode"

	"shrinkout/pkg/pkgCustomTemplate"
	"shrinkout/pkg/pkgUtil"
)

const ksVersion = "V 2.0.0"
const ksDate = "01/04/2022"

//------------ Config Start

//------------ Config End

type mapTemplateFiles map[string]string
type mapTemplateDir map[string]mapTemplateFiles // [main0]TemplateFiles for main0,[main1]TemplateFiles for main1...
type mapLangConfigs map[string]mapTemplateDir   //[english]mapTemplateDir for english, [spanish]mapTemplateDir for spanish

var (
	gsConfigFile string
	gsDoTemplate string

	kConfigNode pkgUtil.ClsConfigNode
)

func SVersion()(string, string){
	return "makeTemplate Version: " + ksVersion, ", " + ksDate
}

func getLangConfigs(sPath string) (mapLangConfigs, error) {
	var tReturn mapLangConfigs

	fileDir, err := ioutil.ReadDir(sPath)

	if err == nil {
		tReturn = make(map[string]mapTemplateDir)
		for _, foldername := range fileDir {
			if foldername.IsDir() {
				tReturn[strings.ToLower(foldername.Name())], err = getFoldersAndFiles(sPath + foldername.Name() + "/")
			}
		}
	}

	if err != nil {
		fmt.Println(err)
	}
	return tReturn, err
}

func fnGetTemplateFiles(sPath string) map[string]string {
	tReturn := make(map[string]string)

	objFileDir, err := ioutil.ReadDir(sPath)

	if err == nil {
		for _, fileName := range objFileDir {
			sFilePath := sPath + "/" + fileName.Name()
			sFile := strings.ToLower(fileName.Name())
			sExtension := path.Ext(sFile)

			if (sExtension == ".config") || (sExtension == ".tmpl") {
				tReturn[sFile] = sFilePath
			}
		}
	}

	if err != nil {
		fmt.Println(err)
	}
	return tReturn
}

// sPath is path to a 'base foler'. Function will get .tmpl and .config files from folders inside the 'base folder'.
func getFoldersAndFiles(sPath string) (mapTemplateDir, error) {
	var (
		tReturn mapTemplateDir
		fileDir []os.FileInfo
		err     error
	)

	fnDirFunc := func(sPath string) bool {
		fileDir, err = ioutil.ReadDir(sPath)
		return err == nil
	}

	if fnDirFunc(sPath) {
		tReturn = make(map[string]mapTemplateFiles)
		for _, foldername := range fileDir {
			if foldername.IsDir() {
				if fnDirFunc(sPath + foldername.Name()) { // templates/default/default.tmpl, templates/default/defaultNavBar.tmpl, templates/main0/defaultNavBar.tmpl...
					templateFiles := fnGetTemplateFiles(sPath + foldername.Name())

					if len(templateFiles) > 0 {
						tReturn[strings.ToLower(foldername.Name())] = templateFiles
					}
				}
			}
		}
	}

	if err != nil {
		fmt.Println(err)
	}
	return tReturn, err
}

func vfBodyDo(sMasterTemplateName, sReplace, sLangName string, objLangConfigs mapTemplateDir) {
	var arMasterConfig []string
	objMasterFolders := fnGetTemplateFiles("./masterTemplates/" + sMasterTemplateName)
	sErr := sMasterTemplateName + " No Master config File"
	if len(objMasterFolders) > 0 {
		sErr = sMasterTemplateName + " No Master Template Files"

		arMasterConfig, _ = pkgUtil.FileToArLines(objMasterFolders[sMasterTemplateName+".config"])
		if len(arMasterConfig) > 0 {
			sErr = ""
			for key, templateFolderPath := range objMasterFolders {
				if path.Ext(templateFolderPath) == ".tmpl" {
					sTemp, _ := pkgUtil.GetFile(templateFolderPath)
					objMasterFolders[key] = string(sTemp)
				}
			}
		}

	}

	if len(sErr) < 1 {

		objTemplateFolders, _ := getFoldersAndFiles("./templates/" + sMasterTemplateName + "/")
		for sFolderName, objFolders := range objTemplateFolders {
			sHeader := ""
			for sFileName, sFile := range objMasterFolders {
				if path.Ext(sFileName) == ".tmpl" {
					_, bOk := objFolders[sFileName]
					if !bOk { //local files like defaultnavbar.tmpl override master files
						sHeader += sFile
					}
				}
			}
			if len(sHeader) > 0 {
				sHeader = sfProcessCommands(sHeader, arMasterConfig)
				sHeader = strings.Replace(sHeader, "~!host@#", sReplace, -1)
				sErr = ""

			}
			if objTemplateConfigs, bOk := objLangConfigs[sFolderName]; bOk { // make sure template has a config file
				sBody := ""
				arLocalConfig, _ := pkgUtil.FileToArLines(objTemplateConfigs[sFolderName+".config"])
				for _, sFilePath := range objFolders {
					sTemp, _ := pkgUtil.GetFile(sFilePath)
					sBody = sBody + string(sTemp)
				}
				if len(sBody) > 0 {
					sHeader = sfProcessCommands(sHeader, arLocalConfig)
					sBody = sfProcessCommands(sBody, arLocalConfig)
					ioutil.WriteFile("./lang/"+sLangName+"/deploy/"+sFolderName+".tmpl", []byte(sHeader+sBody), 0666)

					fmt.Println("-----" + sFolderName)
				} else {
					fmt.Println("No BODY")
				}
			}
		}
	}
}

func vfCreateDeploy(objMasterFolderFiles mapTemplateFiles, sMasterFolderName, sReplace, sTemplatePath, sLangPath string) string {
	//	sHeader := sfHeaderDo(sMasterFolderName, sReplace)
	objLangConfigs, _ := getLangConfigs("./lang/")
	for sLangName, objLangConfigs := range objLangConfigs {
		fmt.Println("-" + sLangName)
		vfBodyDo(sMasterFolderName, sReplace, sLangName, objLangConfigs)
	}

	return ""
}

func sfTokensToLowerCase(s string) string {
	iPos := 0
	for {
		i1 := strings.Index(s[iPos:len(s)], `~!`)
		i2 := strings.Index(s[iPos:len(s)], `@#`)

		if i1 < 0 || i2 < 0 {
			break
		}

		sToken := s[iPos:len(s)][i1 : i2+2]
		s = strings.Replace(s, sToken, strings.ToLower(sToken), -1)
		iPos = iPos + i2 + 2
	}
	return s
}

func sfProcessCommands(sOutput string, arConfigFile []string) string {
	sError := ""

	sOutput = sfTokensToLowerCase(sOutput)
	for iIndex := 0; iIndex < len(arConfigFile); iIndex++ {
		v := arConfigFile[iIndex]
		arCommandLine := strings.Split(v, "sep:")
		if len(arCommandLine) > 1 {
			sCommand := strings.ToLower(strings.TrimSpace(arCommandLine[0]))
			switch sCommand {
			case "debug", "system", "replace":
				sOutput = strings.Replace(sOutput, strings.ToLower(strings.TrimSpace(arCommandLine[1])), strings.TrimSpace(arCommandLine[2]), -1)
				break
			case "multisystem", "multireplace", "multitrim":
				sReplace := ""
				sBreak := ""
				if sCommand != "multiTrim" {
					sBreak = "\n"
				}
				sTemp := strings.TrimSpace(arCommandLine[2])
				iIndex++
				for ; iIndex < len(arConfigFile); iIndex++ {
					sReplace += sTemp
					sTemp = strings.TrimSpace(arConfigFile[iIndex])
					if strings.ToLower(sTemp) == "endsep:" {
						break
					} else {
						sReplace += sBreak
					}
				}
				if strings.ToLower(sTemp) == "endsep:" {
					sOutput = strings.Replace(sOutput, strings.ToLower(strings.TrimSpace(arCommandLine[1])), sReplace, -1)
				} else {
					sError = strings.TrimSpace(arCommandLine[0]) + ", " + arCommandLine[1] + " missing endSep:"
				}
				break
			}
		}
	}
	if len(sError) > 0 {
		fmt.Println("sfProcessCommands Error " + sError)
	}
	return sOutput
}

/*  arFlags[x] get replaced with arReplace[x]
	line with flag must have form -->x//flag or just //flag
  	function leaves //flag but replaces the line.
 	gs_post = //UrlStringFlag becomes gs_post = "http://localhost/";//UrlStringFlag
*/
func modifyConfig(arLines, arFlags, arReplace []string) error {
	sError := ""
	for iFlag, sFlag := range arFlags {
		sFlagError := sFlag + ","
		for i, sLine := range arLines {
			sLine = strings.TrimSpace(sLine)
			if len(sLine) > len(sFlag) {
				if sLine[0:2] != "//" {
					if sFlag == (sLine[len(sLine)-len(sFlag):]) {
						arLines[i] = arReplace[iFlag] + sFlag
						sFlagError = ""
						break
					}
				}
			}
		}
		sError = sError + sFlagError
	}

	if len(sError) > 0 {
		return errors.New(sError)
	}

	return nil
}

func ModifyConfig(sConfigPath string, arFlags, arReplace []string) error {

	arLines, err := pkgUtil.FileToArLines(sConfigPath)
	if err == nil {
		err = modifyConfig(arLines, arFlags, arReplace)
		if err == nil {
			err = pkgUtil.ArLinesToFile(sConfigPath, arLines)
		}
	}

	if err != nil {
		return errors.New(sConfigPath + ", " + err.Error())
	}

	return err
}

func vfSetMasterConfig(sFileName, sPath, sReplace string) {
	sFinalPath := sPath + sFileName + ".config"
	arConfig, err := pkgUtil.FileToArLines(sFinalPath)
	for i := 0; i < len(arConfig); i++ {
		iPos := strings.Index(arConfig[i], "~!host@#")
		if iPos > 0 {
			arConfig[i] = "system sep:~!host@# sep: " + sReplace
			break
		}
	}
	if err == nil {
		pkgUtil.ArLinesToFile(sFinalPath, arConfig)
	}
}

func vfProcessTemplates(sPath string) {
	var err error
	sLangPath := sPath
	sLangPath += "lang"

	sUrl := kConfigNode.SfGetNetwork()
	sReplace := `<base href="` + strings.TrimSpace(sUrl) + `">`

	objMasterFolders, _ := getFoldersAndFiles("./masterTemplates/")
	if len(gsDoTemplate) > 0 { // do only one ie luis
		objMasterFolderFiles, bOk := objMasterFolders[gsDoTemplate]
		if bOk {
			vfSetMasterConfig(gsDoTemplate, "./masterTemplates/"+gsDoTemplate+"/", sReplace)
			vfCreateDeploy(objMasterFolderFiles, gsDoTemplate, sReplace, "./templates/", "./lang")
		} else {
			fmt.Println("VfProcessTemplates :" + gsDoTemplate + " Not Found in ./masterTemplates")
		}
	} else {		
		for sMasterFolderName, objMasterFolderFiles := range objMasterFolders {
			vfSetMasterConfig(sMasterFolderName, "./masterTemplates/"+sMasterFolderName+"/", sReplace)
			vfCreateDeploy(objMasterFolderFiles, sMasterFolderName, sReplace, "./templates/", "./lang")
		}
	}

	if err != nil {
		fmt.Println(err)
	}

}

func bfInit() error {
	var(
		localError error
		iMode int
	)
	localError = pkgCustomTemplate.FnPreInit()
	if localError == nil{
		kConfigNode, iMode, localError = pkgUtil.Init()
		if localError == nil {
			localError = pkgCustomTemplate.FnAfterInit(iMode)
			sUrl := kConfigNode.SfGetNetwork()
			arFlags, arReplace := pkgCustomTemplate.ConfigArs(sUrl)

			localError = ModifyConfig("./static/js/globals.js", arFlags, arReplace)
			fmt.Println(kConfigNode.SfConfigNodeDisplay())
		} else {
			fmt.Println(localError)
		}
	}

	if localError == nil{
		localError = pkgCustomTemplate.FnPostInit(iMode)
	}
	return localError
}

func sfGetKey(sFileName, sData string) (int, string) {
	sAdd := ""
	i := strings.Index(sData, "~!")
	if i > -1 {
		e := strings.Index(sData, "@#")
		if e < i {
			panic("====================" + sFileName + " Has ~! but no closing @#")
		}

		sAdd = sData[i : e+2]
		i = e + 2
	}
	return i, sAdd
}

func arGetKeys(sFileName, sTemplateFile string) []string {
	var arKeys []string
	if len(sTemplateFile) > len("~!@#") {
		for {
			i, sAdd := sfGetKey(sFileName, sTemplateFile)
			if i > -1 {
				if len(sAdd) >= len("Key@#") {
					if strings.ToLower(sAdd[len(sAdd)-len("Key@#"):len(sAdd)]) == "key@#" {
						sAdd = ""
					}
				}
				if len(sAdd) > 0 {
					arKeys = append(arKeys, sAdd)
				}
				sTemplateFile = sTemplateFile[i:len(sTemplateFile)]
			} else {
				break
			}

		}
	}
	return arKeys
}

func mapGetKeys(sFileName, sConfigFile string) map[string]string {
	arLines := strings.Split(sConfigFile, "\n")
	m := make(map[string]string)

	for _, v := range arLines {
		_, sKey := sfGetKey(sFileName, v)
		if len(sKey) > 0 {
			arCommand := strings.Split(v, " ")
			m[strings.ToLower(sKey)] = arCommand[0]
		}
	}
	return m
}

func ifFindKeyInArray(sKey string, arAr []string) int {
	var i int
	sKey = strings.ToLower(sKey)
	for i = len(arAr) - 1; i > -1; i-- {
		if strings.ToLower(arAr[i]) == sKey {
			break
		}
	}

	return i
}

func vfProcessKeys(sConfigFilePath, sTemplateFilePath string, bAdd bool) {
	sConfigFile, err := pkgUtil.GetFile(sConfigFilePath)
	if len(sConfigFile) < 1 || err != nil {
		fmt.Println(sConfigFilePath+" empty: Error: ", err)
	} else {
		arTemplatePaths := strings.Split(sTemplateFilePath, ";")
		sTemplateFile := ""
		for _, filePath := range arTemplatePaths {
			st, err := pkgUtil.GetFile(filePath)
			if err == nil {
				sTemplateFile = sTemplateFile + string(st)
			}
		}

		if len(sTemplateFile) < 1 {
			fmt.Println("Empty: ", sTemplateFilePath)
		} else {
			mapConfigKeys := mapGetKeys(sConfigFilePath, string(sConfigFile))
			arTemplateKeys := arGetKeys(sTemplateFilePath, sTemplateFile)

			var arAddToConfig, arDelFromConfig []string
			if bAdd {
				for _, v := range arTemplateKeys {
					if _, bOk := mapConfigKeys[strings.ToLower(v)]; !bOk {
//						arAddToConfig = append(arAddToConfig, "replace sep:" + v + "\tsep:<!-- -->")
						arAddToConfig = append(arAddToConfig, v)
						delete(mapConfigKeys, v)
					}
				}
			}

			for k, _ := range mapConfigKeys {
				if ifFindKeyInArray(k, arTemplateKeys) < 0 {
					arDelFromConfig = append(arDelFromConfig, k)
				}
			}
			vfProcessFile(sConfigFilePath, string(sConfigFile), arAddToConfig, arDelFromConfig)
		}
	}
}

func vfExtractMulti(arLines []string) []string {
	var arTemp []string
	iz := -1

	for {
		iz++
		if iz >= len(arLines) {
			break
		}

		if len(arLines[iz]) >= len("multi") {
			if strings.ToLower(arLines[iz][0:len("multi")]) == "multi" {
				arTemp = append(arTemp, arLines[iz])
				arLines[iz] = ""
				for {
					iz++
					if iz >= len(arLines) {
						panic("No endsep:")
					}
					sLine := arLines[iz]
					arTemp = append(arTemp, sLine)
					arLines[iz] = ""
					if strings.ToLower(sLine) == "endsep:" || iz >= len(arLines) {
						break
					}

				}
			}

		}
	}
	return arTemp
}

func extractSort(arLines []string, sCommand string) ([]string, int) {
	var arTemp []string

	objMap := make(map[string]int)
	iReturn := 1

	iz := -1
	for {
		iz++

		if iz >= len(arLines) {
			break
		}

		if len(arLines[iz]) > len(sCommand) {

			if strings.ToLower(arLines[iz][0:len(sCommand)]) == strings.ToLower(sCommand) {

				arInsert := strings.Split(arLines[iz], "sep:")
				if len(arInsert) > 0 {
					objMap[arInsert[1]]++
					arTemp = append(arTemp,
						strings.ToLower(strings.TrimSpace(arInsert[1]))+" "+arLines[iz])
					arLines[iz] = ""
				}
			}
		}
	}

	if len(arTemp) > 0 {
		sort.Strings(arTemp)

		iz = -1
		for {
			iz++
			if iz >= len(arTemp) {
				break
			}

			arTimp := strings.Split(arTemp[iz], " ")
			arTemp[iz] = strings.TrimSpace(arTemp[iz][len(arTimp[0]):])
		}
	}
	for _, v := range objMap {
		if v > 1 {
			iReturn = v
			break
		}
	}
	return arTemp, iReturn

}

func arSortArLines(arLines []string) ([]string, int) {
	arCommands := []string{"system", "replace"}
	iRepeat := 0
	arTomp := vfExtractMulti(arLines)
	for _, v := range arCommands {
		arTimp, iR := extractSort(arLines, v)
		if iR > iRepeat {
			iRepeat = iR
		}
		for _, v := range arTimp {
			arTomp = append(arTomp, v)
		}
	}

	for _, v := range arLines {
		if len(strings.TrimSpace(v)) > 0 {
			arTomp = append(arTomp, strings.TrimSpace(v))
		}
	}

	return arTomp, iRepeat
}

func vfProcessFile(sPath, sFileData string, arAdd, arDel []string) {
	bModified := false
	if len(arDel)+len(arAdd) > 0 {
		arLines := strings.Split(sFileData, "\n")
		if len(arDel) > 0 {
			varFDelKey :=
				func(sKey string) {
					sKey = strings.ToLower(sKey)
					iPos := -1
					for {
						iPos++
						if iPos >= len(arLines) {
							break
						}
						sLine := arLines[iPos]
						arCommand := strings.Split(sLine, " ")
						if len(arCommand) > 0 {
							vfProcessMulti := func(bDel bool) {
								for {
									iPos = iPos + 1
									if iPos > len(arLines)-1 {
										panic("====================" + sPath + " Missing endsep: on key " + sKey)
									} else {
										arCommand = strings.Split(arLines[iPos], " ")
										if bDel {
											arLines[iPos] = ""
										}
										if len(arCommand[0]) > len("endsep") {
											if strings.ToLower(arCommand[0][0:len("endsep")]) == "endsep" {
												break
											}
										}
									}
								}
							}
							if strings.Index(strings.ToLower(arCommand[0]), "system") < 0 {
								_, s := sfGetKey(sPath, sLine)
								if strings.ToLower(s) == sKey {
									bModified = true
									arLines[iPos] = ""
									if arCommand[0][0:len("multi")] == "multi" {
										vfProcessMulti(true)
									}
								}

							} else if strings.Index(strings.ToLower(arCommand[0]), "multi") > -1 {
								vfProcessMulti(false)
							}
						}
					}
				}

			for _, v := range arDel {
				varFDelKey(v)
			}
		}

		arTimp, iRepeat := arSortArLines(arLines)
		for _, v := range arAdd {
			bModified = true
			arTimp = append(arTimp, v)
		}

		pkgUtil.ArLinesToFile(sPath, arTimp)
		//carlos pkgUtil.ArLinesToFile("./temp/temp.txt", arTimp)

		if iRepeat > 1 || bModified {
			fmt.Println("----------", sPath, "-----------")
			fmt.Println("Modified:", bModified, ", Repeat:", iRepeat)
		}
	}

}

func cleanConfig() {
	type clsConfigs struct {
		SConfigPath   string
		STemplatePath string
		IUsed         int
	}

	objMasters, err := getFoldersAndFiles("./masterTemplates/")
	if err == nil {
		for sFolderName, objFiles := range objMasters {
			if sConfigPath, bOk := objFiles[sFolderName+".config"]; bOk {
				sFilePath := ""
				for _, sPath := range objFiles {
					if path.Ext(sPath) == ".tmpl" {
						if len(sFilePath) > 0 {
							sFilePath = sFilePath + ";"
						}
						sFilePath += sPath
					}
				}
				if len(sFilePath) > 0 {
					vfProcessKeys(sConfigPath, sFilePath, true)
				}
			}
		}

		var objConfig clsConfigs
		objConfigs := make(map[string]clsConfigs)

		objLangConfigs, err := getFoldersAndFiles("./lang/english/") // need to do other languages
		delete(objLangConfigs, "deploy")
		if err == nil {
			for sLangName, objLangFolders := range objLangConfigs { // ie ./lang/engish
				objConfig.SConfigPath = objLangFolders[sLangName+".config"]
				objConfigs[sLangName] = objConfig
			}

			fileDir, err := ioutil.ReadDir("./templates/")

			if err == nil {
				for _, v := range fileDir {
					objTemplateFolders, err := getFoldersAndFiles("./templates/" + v.Name() + "/")
					if err == nil {
						for sFolderName, sFolderPath := range objTemplateFolders {
							obj, bOk := objConfigs[sFolderName]
							if bOk {
								obj.IUsed = obj.IUsed + 1
								if obj.IUsed > 1 {
									panic("====================More than one template with name: " + sFolderName + ", only one config")
								} else {
									for _, sFilePath := range sFolderPath {
										obj.STemplatePath += sFilePath + ";"
									}
									objConfigs[sFolderName] = obj
								}
							}
						}
					}
				}

				for _, v := range objConfigs {
					if v.IUsed < 1 {
						fmt.Println("No Template for: " + v.SConfigPath)
					} else {
						vfProcessKeys(v.SConfigPath, v.STemplatePath, true)
					}
				}
			}
		} else {
			fmt.Println(err)
		}

	}
}

func sfCommand()(string,string){
	var(
		sCommand, sError string
	)
	
	_, localError := os.Stat("./go.config")

	if localError == nil{
		if len(os.Args) < 2{
			sError = "Incomplete number of parameters. makeTemplate command options. ie makeTemplate makethis zelfpublish\n options are: \ncleanconfig \nmakeall \nmakethis name"		
		}else{
			sCommand = strings.ToLower(os.Args[1])
		}
	}else{
		sError = "missing go.config file missing"
	}

	return sCommand, sError
}

func fnStripAllBlanks(sLine string)string{
	var sReturn []rune
	sLine = strings.ToLower(sLine)
	for _,v := range(sLine){
		if ! unicode.IsSpace(v){
			sReturn = append(sReturn, v)
		}
	}
	return string(sReturn)
}

func fnTest(){
	fmt.Println("Test")
}

func main() {
	fmt.Println(pkgUtil.SVersion())
	fmt.Println(SVersion())
	sCommand, sError := sfCommand()
	if len(sError) < 1{
		switch(sCommand){
			case "cleanconfig":
				cleanConfig()
				break;
			case "makethis":
				gsDoTemplate = os.Args[2]
				fmt.Println(gsDoTemplate)
				localError := bfInit()
				if localError == nil {
					vfProcessTemplates("./")
				} else {
					sError = localError.Error()
				}
				break
			case "makeall":
				fmt.Println("MakeAll")			
				localError := bfInit()
				if localError == nil {
					vfProcessTemplates("./")
				} else {
					sError = localError.Error()
				}
				break
			case "test":
				fnTest()
				break;
			default:
				fnPrintError(sCommand + " is not a valid command")
		}
	}
	if len(sError) > 0{
		fmt.Println("Error: " + sError)
	}else{
		fmt.Println("Done")
	}
}

func fnPrintError(sError string){
	fmt.Println("*************** " + sError )
}
