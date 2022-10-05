/*
	bpsetup sAppToReplace sAppReplacement 
	ie bpsetup boilerplate ZelfPublish will:
		- replace all boilerplate/ with /ZelfPublish/ in all import section of .go files in current and nested folders
		- rename ./masterTemplates/sAppToReplace to ./masterTemplates/sAppReplacement  t
		- rename ./masterTemplates/sAppToReplace/sAppToReplace.config  to ./masterTemplates/sAppReplacement/sAppToReplace/sAppToReplace.config
		- rename ./templates/sAppToReplace to ./templates/sAppReplacement  
	ie bpsetup boilerplate ZelfPublish
	
	Note: 
	1) import format must be: (Opening parenthesis must be in the same line as import, closing parenthesis on a different line)
		import(
			files to import
		)
	2) Do not have empy imports.
		Not:
			import()
		nor:
			import(
			)
	

*/

package main

import(
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode"

)

const ksVersion = "V 1.0.0"
const ksDate = "01/03/2022"


func main() {	
	var sError string
	fmt.Println("Version:", ksVersion)
	fmt.Println("Written:", ksDate)
	
	if len(os.Args) < 3{
		sError = "Not enough parameters. formatImports sImportToReplace sReplacement. Ie: formatImports boilerplate/ ZelfPublish/"
	}else{
		fnFormatImports(os.Args[1], os.Args[2])
		fnSetupFolders(os.Args[1], os.Args[2])
	}
	if len(sError) > 0{
		fmt.Println(sError)
	}
	fmt.Println("Done")
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

func fnSetupFolders(sStringToReplace, sReplacement string){
	sStringToReplace = strings.TrimSpace(strings.ToLower(sStringToReplace))
	sReplacement = strings.TrimSpace(strings.ToLower(sReplacement))
	os.Rename("./templates/" + sStringToReplace,"./templates/" + sReplacement)
	os.Rename("./masterTemplates/" + sStringToReplace,"./masterTemplates/" + sReplacement)
	os.Rename("./masterTemplates/" + sReplacement + "/" + sStringToReplace + ".config",
		"./masterTemplates/" + sReplacement + "/" + sReplacement + ".config")
		

}

func fnFormatImports(sStringToReplace, sReplacement string){
	var (
		fileDir []os.FileInfo
		err     error
	)
	fnDirFunc := func(sPath string) bool {
		fileDir, err = ioutil.ReadDir(sPath)
		return err == nil
	}

	fnReplace := func(sPath string){
		for _, foldername := range fileDir {
			if ! foldername.IsDir() {
				sExtension :=  SfFileExtension(foldername.Name())
				if sExtension == "go"{
					arLines, localError :=  FileToArLines(sPath + foldername.Name())
					bModified := false
					if localError == nil{
						var(
							arImport, arCopy []string							
							iStatus int
						)
						for iLineCount := 0;iLineCount < len(arLines);iLineCount++{
							sLine := arLines[iLineCount]
							switch(iStatus){
								case 0:
									if fnStripAllBlanks(sLine) == "import("{										
										arImport = append(arImport,sLine)
										iStatus = 1
									}else{
										arCopy = append(arCopy,sLine)
									}
									break
								case 1:
									if strings.TrimSpace(sLine) == ")"{
										arCopy = append(arCopy,arImport...)
										arCopy = append(arCopy,sLine)										
										iStatus = 2
									}else{
										sNewLine := strings.Replace(sLine,sStringToReplace + "/", sReplacement + "/",1)
										arImport = append(arImport,sNewLine)
										if !bModified{
											bModified = strings.TrimSpace(sNewLine) != strings.TrimSpace(sLine)
										}										
									}
									break
								default:
									arCopy = append(arCopy,sLine)
								break;
							}
						}
						if iStatus > 0 && bModified{
							sFileName := sPath +  foldername.Name()
							ArLinesToFile(sFileName,arCopy)
							fmt.Println(sPath + foldername.Name())
							iStatus = 0
						}
						
								
					}	
				}
			}
		}
	}
	
	if fnDirFunc("./") {
		fnReplace("./")
		if fnDirFunc("./pkg"){
			fnReplace("./pkg")
			for _, foldername := range fileDir {
				if foldername.IsDir() {
					if fnDirFunc("./pkg/" + foldername.Name()){
						fnReplace("./pkg/" + foldername.Name() + "/")
					}
				}
			}
		
		}
	}

	if err != nil {
		fmt.Println(err)
	}
	
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

func SfFileExtension(sFile string) string {
	arAr := strings.Split(sFile, ".")

	if len(arAr) > 0 {
		return arAr[len(arAr)-1]
	}

	return ""
}

func GetFile(sPath string) ([]byte, error) {

	var data []byte

	_, localError := os.Stat(sPath)

	if localError == nil {
		return ioutil.ReadFile(sPath)
	}

	return data, localError
}


