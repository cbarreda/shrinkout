package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"


	"shrinkout/pkg/pkgError"
	"shrinkout/pkg/pkgDatabase"
	"shrinkout/pkg/pkgTemplate"	
	"shrinkout/pkg/pkgUser"
		
	"shrinkout/pkg/pkgUtil"
	"shrinkout/pkg/pkgSession"

    "shrinkout/pkg/pkgAdmin"
    "bytes"
    "time"
)

type TcustUrl struct {
	Url  	string
	PortStr	string
	Cust	string
	CustId	int	
	DbStr	string
}



/*----------------  GLOBAL VARIABLES START -------------------------------*/

var (
	kConfigNode   	pkgUtil.ClsConfigNode
	kCustUrl 		TcustUrl
	ksLangPath    = "./lang"
	ksDefaultLang = "english"


	/*
	gKey  = []byte("cb!boi!erpl@te$ecret%Key")
	
	gstore = sessions.NewCookieStore(pkgSession.GetKey())
	*/

	router = mux.NewRouter()


	gPkgID = "shrinkout."
	/*
        gSessionKey = []byte("cb!App$ecret%Key")
        gSessionName = "cb$e$$ionName"
        store = sessions.NewCookieStore(gSessionKey)
	 
	 */
	 
)


/*----------------  GLOBAL VARIABLES END -------------------------------*/

/*----------------  INIT START -------------------------------*/

func getLanguage(w http.ResponseWriter, r *http.Request) string {
	objError := pkgSession.GetLang(w, r)
	if objError.ICode < 0 {
		objError.SInfo = ksDefaultLang
	}

	return objError.SInfo

}

/*----------------  INIT END -------------------------------*/


func main() {

	var(
		localError error
		iMode int
	)		
	
	
	kConfigNode, iMode, localError  = pkgUtil.Init()

	if localError == nil {
		pkgError.VfInit()

		pkgDatabase.SetDatabase(kConfigNode.Database.DbStr)
		fmt.Println("main Initialized------- Mode",iMode)
		fmt.Println("Network ", kConfigNode.SfGetNetwork())
		fmt.Println(kConfigNode.SfConfigNodeDisplay())
		fmt.Println("Template Read, ", pkgTemplate.IfReadTemplate(ksLangPath))

		fs := http.FileServer(http.Dir("static"))
		http.Handle("/static/", http.StripPrefix("/static/", fs))

		router.NotFoundHandler = http.HandlerFunc(fnNotFound)		

		router.HandleFunc("/fnRouter/", fnRouter)

		router.HandleFunc("/test/", fnTest)
		router.HandleFunc("/setupcolors/", fnSetupColors)
		router.HandleFunc("/readTemplate", fnTemplatesRead)
		router.HandleFunc("/", fnRoot)
		router.HandleFunc("/english/", fnRootE)
		router.HandleFunc("/spanish/", fnRootS)

        router.HandleFunc("/fnRouter/",fnRouter)
        router.HandleFunc("/client",fnClient)
        router.HandleFunc("/report",fnReport)
        router.HandleFunc("/fnDoShrink",fnDoShrink)
        router.HandleFunc("/download", fnDownload)
        router.HandleFunc("/readTemplate/",fnReadTemplates)
        router.HandleFunc("/admin",fnAdmin)
        router.HandleFunc("/device",fnGetDevice)
        router.HandleFunc("/fnSpc",fnSpc)        
        router.HandleFunc("/fnAdminLogin",fnAdminLogin)
        router.HandleFunc("/fnDeviceSetup",fnDeviceSetup)

		http.Handle("/", router)
		
		localError = http.ListenAndServe(kConfigNode.ListenPort, nil)

		if localError == nil{
			fmt.Println("PORT", kCustUrl.PortStr)
			
		}else{
			pkgError.VfMark("func main", localError)
		}

	} else {
		pkgError.VfMark("func main", localError)
		fmt.Println("Template Read, ", pkgTemplate.IfReadTemplate(ksLangPath))
	}
}

func fnRouter(w http.ResponseWriter, r *http.Request) {
	sFunc := gPkgID
	sFunc += "fnRouter"


	type clsReturn struct {
		IntCode    int
		StrCommand string
		SJSON      string
	}

	var objReturn clsReturn

	objReturn.IntCode = -1000
	r.ParseForm()
	if len(r.Form["sCommand"]) > 0 {
		objReturn.IntCode = -1001
		objReturn.StrCommand = r.Form["sCommand"][0]
		sJSON := r.Form["sData"][0]
		// Carlos Opening all CORS
		(w).Header().Set("Access-Control-Allow-Origin", "*") 
		if len(sJSON) > 0 {
			objReturn.IntCode = -9999
			switch objReturn.StrCommand {
				case "PstTest":
					PstTest(w, r)
					break
				case "PstSaveColors":
					fnColorsGetPut(w,r,sJSON)
					break
				case "andDoShrink":
					objReturn.IntCode, objReturn.SJSON = andDoShrink(w,r,sJSON)
				case "andFloorLogin":
					objReturn.IntCode, objReturn.SJSON = andFloorLogin(w,r,sJSON)
				case "cmdFloorLogin":
					objReturn.IntCode, objReturn.SJSON = cmdFloorLogin(w,r,sJSON)
				case "cmdFloorLogout":
					objReturn.IntCode, objReturn.SJSON = cmdFloorLogout(w,r)
				case "andSetStore":
					objReturn.IntCode, objReturn.SJSON = andSetStore(w,r,sJSON)
				case "cmdSetStore":
					objReturn.IntCode, objReturn.SJSON = cmdSetStore(w,r,sJSON)
				case "cmdShrinkLogToHtml":
					objReturn.IntCode, objReturn.SJSON = cmdShrinkLogToHtml(w,r)
							// Admin start
				case "cmdAddPlu":
					objReturn.IntCode, objReturn.SJSON  = pkgAdmin.FnAddPlu(w,r)
				case "cmdSavePlu":
					objReturn.IntCode, objReturn.SJSON  = pkgAdmin.FnSavePlu(w,r,sJSON)
				case "cmdfProcessLogFile":
					objReturn.IntCode, objReturn.SJSON = pkgAdmin.FnProcessLogFile(w,r)
				case  "cmdMoveProduct":
					objReturn.IntCode, objReturn.SJSON = pkgAdmin.FnMoveProduct(w,r,sJSON);
				case  "cmdOrphans":
					objReturn.IntCode, objReturn.SJSON = pkgAdmin.FnOrphans(w,r,sJSON);
				case  "cmdModifyProduct":                    
					objReturn.IntCode, objReturn.SJSON = pkgAdmin.FnModifyProduct(w,r,sJSON);
				case  "cmdAddLocation":
					objReturn.IntCode, objReturn.SJSON = pkgAdmin.FnAddLocation(w,r,sJSON);
				case  "cmdModifyLocation":
					objReturn.IntCode, objReturn.SJSON = pkgAdmin.FnModifyLocation(w,r,sJSON);
				case "cmdGetProdLocation":
					objReturn.IntCode, objReturn.SJSON = pkgAdmin.FnGetProdLocation(w,r,sJSON);
				case "cmdFloorReport":
					objReturn.IntCode, objReturn.SJSON = pkgAdmin.FnReport(w,r,sJSON)
							// Admin end
				case "cmdGetTransactionId":
					objReturn.IntCode, objReturn.SJSON = cmdGetTransactionId(w,r,sJSON);
				case "andSaveInput":
					objReturn.IntCode, objReturn.SJSON = andSaveInput(w,r,sJSON);
				case "cmdSaveInput":
					objReturn.IntCode, objReturn.SJSON = cmdSaveInput(w,r,sJSON);
				case "cmdTest":
					objReturn.IntCode, objReturn.SJSON = cmdDoTest(w,r,sJSON);
				default: objReturn.IntCode = -9999	
			}
		}
	}
	objJasonReturn, localError := json.Marshal(objReturn)
	pkgError.BfHandleError(pkgError.KLOGERROR, sFunc+"->"+objReturn.StrCommand, localError)
	w.Write(objJasonReturn)
}



func fnRootL(w http.ResponseWriter, r *http.Request, sLang string) {
	sFunc := gPkgID
	sFunc += "fnRoot"

	// root just displays shrinkout page see /client/ for functionality	
	pkgSession.PutLang(w, r, sLang)
	fnGetMain0(w, r ) 
}

func fnRoot(w http.ResponseWriter, r *http.Request) {
	objError := pkgSession.GetLang(w, r)
	if objError.ICode < 0 {
		objError.SInfo = "english"
	}
	fnRootL(w, r, objError.SInfo)
}
func fnRootS(w http.ResponseWriter, r *http.Request) {
	fnRootL(w, r, "spanish")
}
func fnRootE(w http.ResponseWriter, r *http.Request) {
	fnRootL(w, r, "english")
}

/*----------------  TEMPLATES START -------------------------------*/
func fnTemplatesRead(w http.ResponseWriter, r *http.Request) {
	vfReadTemplates()
	fnRoot(w, r)
}

func vfReadTemplates() {
	fmt.Println(pkgTemplate.IfReadTemplate(ksLangPath))
}


func bfSessionTemplateIn(w http.ResponseWriter, sLang, sTemplateName string, objI interface{}) bool {
	return bfSessionTemplateSR(w, sLang, sTemplateName, objI)
}

func bfSessionTemplateSR(w http.ResponseWriter, sLang, sTemplateName string, objI interface{}) bool {
	sFunc := gPkgID
	sFunc += "bfSessionTemplateSR"
	
	var localError error
	sTemplate := pkgTemplate.SfGetTemplate(sLang, sTemplateName)
	if len(sTemplate) > 0 {
		if t, err := template.New(sTemplateName).Parse(sTemplate); err == nil {
			localError = t.Execute(w, objI)
		} else {
			localError = err
		}
	}
	return pkgError.BfHandleError(pkgError.KLOGERROR, sFunc, localError)
}

/*----------------  TEMPLATES END  -------------------------------*/

/*----------------  FUNCTIONS START -------------------------------*/
func fnErrorTemplate(w http.ResponseWriter, r *http.Request, sCallerKey, sErrorMessageKey, sMessage string) {
	type clsLocalE struct {
		SERRORMESSAGEKEY string
		SMESSAGE         string
		SCALLER          string
	}

	var objLocalE clsLocalE
	objLocalE.SERRORMESSAGEKEY = sErrorMessageKey
	objLocalE.SMESSAGE = sMessage
	objLocalE.SCALLER = sCallerKey
	bfSessionTemplateIn(w, getLanguage(w, r), "error.tmpl", objLocalE)
}

func vfDeviceNotSet(w http.ResponseWriter, r *http.Request, iCode int, sSender string) {
	fnErrorTemplate(w, r, sSender, "~!DeviceNotSetKey@#", strconv.Itoa(iCode))
}

func vfErrorCode(w http.ResponseWriter, r *http.Request, iCode int, sSender string) {
	fnErrorTemplate(w, r, sSender, "~!ErrorCodeKey@#", strconv.Itoa(iCode))
}

func vfNotFound(w http.ResponseWriter, r *http.Request, sSender string) {
	fnErrorTemplate(w, r, sSender, "~!NotFoundKey@#", "")
}

func vfNothingToShow(w http.ResponseWriter, r *http.Request, sSender string) {
	fnErrorTemplate(w, r, sSender, "~!NothingToShowKey@#", "")
}

func vfMustBeLoggedIn(w http.ResponseWriter, r *http.Request, sSender string) {
	fnErrorTemplate(w, r, sSender, "~!MustBeLoggedInKey@#", "")
}

func vfUrlQuery(w http.ResponseWriter, r *http.Request, iError int, sSender, sUrl, sError string) {
	fnErrorTemplate(w, r, sSender, "~!UrlProblemKey@#", "<p>"+strconv.Itoa(iError)+" "+
		sError+"</p><p>"+sUrl+"</p>")
}

func fnNotFound(w http.ResponseWriter, r *http.Request) {
	vfNotFound(w, r, r.URL.String())
}

func ifFromUrl(w http.ResponseWriter, r *http.Request, sSender, sUrl string) int {
	iCode := -1
	objURL, bOk := r.URL.Query()[sUrl]
	if bOk {
		ic, localError := strconv.Atoi(objURL[0])
		if localError == nil {
			iCode = ic
		} else {
			vfUrlQuery(w, r, ic, sSender, localError.Error(), sUrl)
		}
	}
	return iCode
}

/*----------------  FUNCTIONS END -------------------------------*/

/*----------------  ROUTER FUNCTIONS START -------------------------------*/


func fnTest(w http.ResponseWriter, r *http.Request) {
	
	IntCode, _ := cmdFloorLogout(w,r)	
	fmt.Println("Test. Loging out ", IntCode)
}


func fnColorsGetPut(w http.ResponseWriter, r *http.Request,sColors string)(int, string) {
	type tColors struct{
		IdBody,IdHeaderC,IdNavBrandB,IdMainNavB,IdBannerH1C,IdBannerH5C,IdItem1B,IdItem2B,IdItemAB,IdItem1H5C,
		IdItem2H5C,IdItemAH5C,IdItem1IC,IdItem2IC,IdItemAIC,IdFooterMainB,IdFooterMainLC,
		IdSectionTitleH2,IdSectionTitleP,IdTeamCBH5,IdTeamCBP,IdTeamHBH5,IdTeamHBP,IdTeamSDH5,IdTeamSDP,
		IdTeamCBB,IdTeamHBB,IdTeamSDB	 string		 
	}
	
	var( 		
		objColors tColors
	)
	sFile := "./setupcolors.txt"
	iReturn := -1
	if len(sColors) < 1{ // load colors from file		
		s,localError := pkgUtil.GetFile(sFile)
		if localError == nil{
			sColors = string(s)
		}else{
			sColors = `{"idPage":"#1a9b5a","idHeader":"#ffffff","idLogo":"#1a9b5a","idNavBar":"#ffffff","idMainNav":"#275e90","idBanner":"#ffffff","idInnovate":"#1a9b5a","idData":"#1a9b5a","idActive":"#275e90","idInnovateC":"#ffffff","idDataC":"#ffffff","idActiveC":"#ffffff","idInnovateI":"#ffffff","idDataI":"#fe7800","idActiveI":"#fe7800"}`
		}
		localError =  json.Unmarshal([]byte(sColors), &objColors)
		fmt.Println("------------------",localError)
		if localError == nil{
			fnGetSetupColors(w,r ,objColors)
		}
		if localError != nil{
			iReturn = 0
		}
	}else{
		if ioutil.WriteFile(sFile, []byte(sColors), 0777) == nil{
			sCssFile := "./static/stylesheets/local/main0.css"
			arLines,localError := pkgUtil.FileToArLines(sCssFile)
			if localError == nil{
				var(
					arCopy []string
				 )
				localError =  json.Unmarshal([]byte(sColors), &objColors)
			
				if localError == nil{
					bOn := true
					for i := 0; i < len(arLines); i++ {
						if bOn{
							bOn = strings.TrimSpace(arLines[i]) != "/* Colors Start */"
							if bOn{
								arCopy = append(arCopy,arLines[i])
							}
						}else{
							bOn = strings.TrimSpace(arLines[i]) == "/* Colors End */"
							if bOn{
								i++
							}
						}	
					}
				}
				sAppend := `
	body {
	background-color:` + objColors.IdBody + `;
	}
	
	h1, h2, h3, h4, h5, h6 {
	color:` + objColors.IdHeaderC + `;
	}

	.navbar-brand{
		background-color:` + objColors.IdNavBrandB + `;
	}

	.main-nav{
		background-color:` + objColors.IdMainNavB + `;
	}
	
	.banner .content-block h1{
		color:` + objColors.IdBannerH1C + `;
	}
	.banner .content-block h5 {
		color:` + objColors.IdBannerH5C + `;
	}

	.about .about-block .about-item.one{
		background:` + objColors.IdItem1B + `;
	.about .about-block .about-item.two{
		background:` + objColors.IdItem2B + `;
	}
	.about .about-block .about-item.active {
		background:` + objColors.IdItemAB + `;
	}
	.about .about-block .about-item.one .content h5,
	.about .about-block .about-item.one .content p {
		color:` + objColors.IdItem1H5C + `;
	}
	.about .about-block .about-item.two .content h5,
	.about .about-block .about-item.two .content p {
		color:` + objColors.IdItem2H5C + `;
	}
	.about .about-block .about-item.active .content h5,
	.about .about-block .about-item.active .content p {
		color:` + objColors.IdItemAH5C + `;
	}

	.about .about-block .about-item.active .icon i {
		color:` + objColors.IdItemAIC + `;
	}

	.about .about-block .about-item.one .icon i {
		color:` + objColors.IdItem1IC + `;
	}	
	.about .about-block .about-item.two .icon i {
		color:` + objColors.IdItem2IC + `;
	}	


	.footer-main {
		background:` + objColors.IdFooterMainB + `;
	}
	.footer-main a{
		color:` + objColors.IdFooterMainLC + `;
	}

	.section-title h2{
		color:` + objColors.IdSectionTitleH2 + `;
	}
	.section-title p{
		color:` + objColors.IdSectionTitleP + `;
	}
	.team-member.cb p{
		color:` + objColors.IdTeamCBH5 + `;
	}
	.team-member.cb h5{
		color:` + objColors.IdTeamCBP + `;
	}
	.team-member.hb p{
		color:` + objColors.IdTeamHBH5 + `;
	}
	.team-member.hb h5{
		color:` + objColors.IdTeamHBP + `;
	}
	.team-member.sd p{
		color:` + objColors.IdTeamSDH5 + `;
	}
	.team-member.sd h5{
		color:` + objColors.IdTeamSDP + `;
	}
	.team-member.cb{
		background:` + objColors.IdTeamCBB + `;
	}
	.team-member.hb{
		background:` + objColors.IdTeamHBB + `;
	}
	.team-member.sd{
		background:` + objColors.IdTeamSDB + `;
	}`
				arCopy = append(arCopy,"/* Colors Start */")
				arCopy = append(arCopy,sAppend)
				arCopy = append(arCopy,"/* Colors End */")
				if pkgUtil.ArLinesToFile(sCssFile,arCopy) == nil{
					iReturn = 0
				}								
			}			
		}
	}
	return iReturn,""
}

func fnSetupColors(w http.ResponseWriter, r *http.Request) {
	fnColorsGetPut(w,r,"")
}
		

func PstTest(w http.ResponseWriter, r *http.Request) {
}


/*----------------  ROUTER FUNCTIONS END -------------------------------*/


/*============================== ORIGINAL SHRINKOUT START ================================*/


/*----------------  GLOBAL VARIABLES START -------------------------------*/
/*----------------  GLOBAL VARIABLES END -------------------------------*/


func vfSetLanguage(pObjUser * pkgUser.ClsBaseUser){
    switch (pObjUser.ILanguage){
        case 2: pObjUser.SLanguage = "spanish"
       default: pObjUser.ILanguage = 1
                pObjUser.SLanguage = "english"    
    }
}


/*----------------  INIT END -------------------------------*/


/*----------------  UTILITY START -------------------------------*/

// Database returns inReturnCode~jsonString
func isf_firstTILde(strParam string) (int, string){
	var sReturn string
	iReturn := -1
	intIndex := strings.Index(strParam,"~")
	
	if(intIndex > 0){
		sReturn = strParam[(intIndex + 1):]
		iReturn,_ = strconv.Atoi(strParam[:intIndex])
	}else{
		sReturn = ""
	}
	
	return iReturn, sReturn
}

func strToInt(s string, intAlternate int) int{
    iReturn, localError := strconv.Atoi(s)
    if localError != nil{
        iReturn = intAlternate
    }
    return iReturn
}


/*----------------  UTILITY END -------------------------------*/

/*----------------  SCREENS START -------------------------------*/

func fnGetDownload(w http.ResponseWriter, r *http.Request,sFilePath string){
    sFunc := "fnGetDownload"
    data, err := ioutil.ReadFile(sFilePath )
    if pkgError.BfHandleError(pkgError.KLOGWARN | pkgError.KLOGPRINT,sFunc,err){
        w.Header().Set("Content-Type", "application/octet-stream")
        w.Header().Set("Content-Disposition", "attachment; filename=" + "carlos.txt")
        w.Header().Set("Content-Transfer-Encoding", "binary")
        w.Header().Set("Expires", "0")
        http.ServeContent(w, r, sFilePath, time.Now(), bytes.NewReader(data))
    }
}

func fnDownload(w http.ResponseWriter, r *http.Request){
    fnGetDownload(w,r,"./static/downloads/carlos.txt")
}

func fnClient(w http.ResponseWriter, r *http.Request){
	sFunc := gPkgID
    sFunc += "fnClient"

   	if r.Method=="GET"{         

        bLoggedIn := false
                    
        cbu,iErr := pkgUser.GetClsBaseUser(w,r) 
        
        if iErr == 0{
            bLoggedIn = len(cbu.SFirst) > 0
        }

	    if bLoggedIn{
            bLoggedIn = false
		    bOk,sJSON := pkgDatabase.BsfExecuteDbFunc("SP_STOREBYEMP(" + strconv.Itoa(cbu.IUid) + ")")
		    if bOk{
				fnGetMain1(w, r, sJSON, cbu)
			}
		    /*
            if bOk{
                sTemplate := pkgTemplate.PstrTemplate(ksDefaultLang,"main1.tmpl")
                if len(sJson) > 0{
                    type objStruct struct{ 
                            Stores []int
                            Inv int
                        }
                    var storeStruct objStruct
                    sStores := ""

                   	err := json.Unmarshal([]byte(sJson),&storeStruct)	

	                if pkgError.BfHandleError(pkgError.KLOGERROR,sFunc,err){
                        for _,iStore := range(storeStruct.Stores){
                            if len(sStores) > 0{
                                sStores += ","
                            }
                            sStores +=  strconv.Itoa(iStore)
                        }
                        sTemplate = strings.Replace(sTemplate,"<!--CarlosStoreNo-->",sStores ,-1)
                    }                    
                }

		        if t,err := template.New("foo").Parse(sTemplate); err == nil{
                    err = t.Execute(w,cbu)     
                }                             
            }
            */
		}else{
			// cd,iE := pkgUser.GetClsDevice(w,r) 
			_,iE := pkgUser.GetClsDevice(w,r) 
            if iE < 0{
                vfDeviceNotSet(w,r,-1,"ShrinkOut Root Module.fnClient")              
            }else{
                fnGetXMain0(w,r)
            }            
		}			            
		    
    }
}

func fnAdmin(w http.ResponseWriter, r *http.Request){
    sTemplate := pkgTemplate.PstrTemplate(ksDefaultLang,"admin0.tmpl")
    if t,err := template.New("foo").Parse(sTemplate); err == nil{
        err = t.Execute(w,nil)     
    }else{
        fmt.Println("fnAdmin",err)
    }
}

func fnSpc(w http.ResponseWriter, r *http.Request){
	sFunc := gPkgID
	sFunc += "fnSpc"
	
	m := r.URL.Query()
	s := m["p_storeproduceId"][0]
	if len(s) > 0{
		iHigh := -1
		iLow  := -1
		bOk,sRet := pkgDatabase.BsfExecuteDbFunc("SP_GETSPC("+ s + ",0)")		
		if !bOk{
			fmt.Println("OOO==>",s,"<--OPS ")
		}
		
		type obj struct{
			SpcId,Ucl,Lcl,IType int
			Points [][]int	
		}
		
		var myObj obj
		json.Unmarshal([]byte(sRet),&myObj)
		
		fmt.Println(myObj.SpcId)
		sArray := ""
		
		for i := len(myObj.Points)-1; i > -1; i--{
			iPt := (myObj.Points[i][2] * 100)/myObj.Points[i][1]
			if iPt > iHigh{
				iHigh = iPt
			}
			
			if (iLow < 0) || (iPt < iLow){
				iLow = iPt
			}
			if len(sArray) > 0{
				sArray += ","
			}
			sArray += strconv.Itoa(iPt)
		}
		if len(sArray) > 0{
			sArray = "var arTest = [" + sArray + "];var arParams = [" + strconv.Itoa(iHigh) + "," + 
			strconv.Itoa(iLow) + "," + strconv.Itoa(myObj.Ucl) + "," + strconv.Itoa(myObj.Lcl) + "]";
			
			fmt.Println(sArray)
		}

		sTemplate := pkgTemplate.PstrTemplate(ksDefaultLang,"spc.tmpl")		
		sTemplate = strings.Replace(sTemplate,"<!-- CarlosSpcArray-->",sArray ,-1)
		if t,err := template.New("foo").Parse(sTemplate); err == nil{
			pkgError.VfDebug(sFunc,"OK");
			err = t.Execute(w,nil)     
		}else{
			fmt.Println("fnSpc",err)
		  
		}
	};
}

func fnDeviceSetup(w http.ResponseWriter, r *http.Request){
    pkgUser.SetDevice(w,r)
}

func fnReport(w http.ResponseWriter, r *http.Request){
    sTemplate := pkgTemplate.PstrTemplate(ksDefaultLang,"report.tmpl")
    if t,err := template.New("foo").Parse(sTemplate); err == nil{
        err = t.Execute(w,nil)     
    }else{
        fmt.Println("fnReport",err)
    }
}

func fnReadTemplates(w http.ResponseWriter, r *http.Request){
	vfReadTemplates()
	fnRoot(w,r)
}

func fnDoShrink(w http.ResponseWriter, r *http.Request){    
    var bLoggedIn bool


    cbu,iErr := pkgUser.GetClsBaseUser(w,r); if iErr == 0{
            bLoggedIn = len(cbu.SFirst) > 0
    }

    if bLoggedIn{
        iStoreNoId := pkgUser.IfGetStoreNoId(w,r)
        r.ParseForm()
		bOk,sCodes := pkgDatabase.BsfExecuteDbFunc("SP_CODELIST(" + strconv.Itoa(iStoreNoId) + ")")
		if bOk{
            bOk,spId := pkgDatabase.BsfExecuteDbFunc("SP_GETSTOREPRODUCEIDS('" + r.Form["nmPluId"][0] + "')")
            if bOk{
                bOk,sTid := pkgDatabase.BsfExecuteDbFunc("SP_GETTRANSACTIONID('')")    

                // Plu does not exist and will be logged
                if len(spId) < 4{
                    spId = `[[0,0,""]]`
                }
                
                if bOk{
					fnGetDoShrink(w,r,sCodes,spId,sTid,iStoreNoId, cbu)
                }
            }
        } 
    }else{
        http.Redirect(w, r, "/", 301)
    }
    
}

func fnLogin(w http.ResponseWriter, r *http.Request,sFunc, s0,s1 string)int{
    var iUser int

    err := r.ParseForm() 
    if pkgError.BfHandleError(pkgError.KLOGERROR,sFunc,err){ 
        sPassword := r.PostFormValue("nmPassword")
        sUser := r.PostFormValue("nmUserId")
        sTemplate := s0
        var cbu pkgUser.ClsBaseUser
        
        iUser,_ = pkgUser.LoginUser(w,r,sUser,sPassword)
        fmt.Println("fnLogin ",iUser, sUser,sPassword,s0)
        if(iUser < 0){
            sTemplate = pkgTemplate.PstrTemplate(ksDefaultLang,s0)
        }else{
            fmt.Println("fnLogin")
            c,iErr := pkgUser.GetClsBaseUser(w,r); if iErr == 0{
                cbu = c
                if cbu.IAdminCode > 15{
                    sTemplate = pkgTemplate.PstrTemplate(ksDefaultLang,s1)
                }
            }
        }
        if t,err := template.New("foo").Parse(sTemplate); err == nil{
            err = t.Execute(w,cbu)     
        }else{
            fmt.Println("fnLogin Error",err)
        }
    }

    return iUser
}

func fnAdminLogin(w http.ResponseWriter,r *http.Request){    
    fnLogin(w,r,"fnAdminLogin","admin0.tmpl","admin1.tmpl") 
}

/*----------------  COMMAND ROUTER END ---------------------*/

/*-------------------- Commands start -------------------------------*/

func cmdDoTest(w http.ResponseWriter,r *http.Request,sJson string)(int, string){    
    sFunc := "cmdDoTest"
    fmt.Println(sFunc)
	sRetJson := ""
    iReturn := -1;

    iStoreNo := 1
    iUserId := 115
    spId := "3064"

	bOk,sCodes := pkgDatabase.BsfExecuteDbFunc("SP_CODELIST(" + strconv.Itoa(iStoreNo) + ")")
    if bOk{
        iReturn = -2
        bOk,iTid := pkgDatabase.BsfExecuteDbFunc("SP_GETTRANSACTIONID('')")
        if bOk{
            iReturn = -3
            bOk,spId := pkgDatabase.BsfExecuteDbFunc("SP_GETSTOREPRODUCEIDS('" + spId + "')")
            if bOk{
                iReturn = 0
                if len(spId) < 4{
                    spId = `[[0,0,""]]`
                }
                sRetJson = `{"iEmpId":` + strconv.Itoa(iUserId) + `,"iStoreNo":`  + strconv.Itoa(iStoreNo) + 
                    `,"Codes":` + sCodes + `,"TransId":`  + iTid + `,"arProduce":`  + spId + `}`;
            }
        }
    }

    fmt.Println(sFunc,sRetJson)

	return iReturn, sRetJson
}

// note html via fnDoSrhink
func andDoShrink(w http.ResponseWriter,r *http.Request,sJson string)(int, string){
    sFunc := "andDoShrink"
	sRetJson := ""
    iReturn := -1;
    dat,err := mapParam(sJson)
  
	if pkgError.BfHandleError(pkgError.KLOGERROR,sFunc,err){

        iStoreNo := int(dat["IstoreNo"].(float64))
        iUserId := int(dat["iUserId"].(float64))
        spId :=  string(dat["sPid"].(string))
        sName :=  string(dat["sName"].(string))
        sPname :=  string(dat["sPname"].(string))
        
		fmt.Println(sFunc + "/" , dat ,	iUserId,iStoreNo,spId,sName,sPname);

		bOk,sCodes := pkgDatabase.BsfExecuteDbFunc("SP_CODELIST(" + strconv.Itoa(iStoreNo) + ")")
        if bOk{
            iReturn = -2
            bOk,iTid := pkgDatabase.BsfExecuteDbFunc("SP_GETTRANSACTIONID('')")
            if bOk{
                iReturn = -3
                bOk,arPr := pkgDatabase.BsfExecuteDbFunc("SP_GETSTOREPRODUCEIDS('" + spId + "')")
                if bOk{
                    iReturn = 0
                    if len(spId) < 4{
                        spId = `[[0,0,""]]`
                    }
                    sRetJson = `{"IntCode":` + strconv.Itoa(iReturn) + `,"iEmpId":` + strconv.Itoa(iUserId) + `,"spId":"`  + spId + 
						`","iStoreNo":` + strconv.Itoa(iStoreNo) +	`,"sName":"`  + sName + `","spName":"`  + sPname + 
						`","iTransId":`  + iTid + `,"Codes":` + sCodes + `,"arProduce":`  + arPr + `}`;                
				}
            }
        }
    }
    fmt.Println(sRetJson)
	return iReturn, sRetJson
}


func andFloorLogin(w http.ResponseWriter,r *http.Request,sJson string)(int, string){
	sRetJson := ""    
	iReturn,dat := genFloorLogin(w,r,sJson)
   
    if(iReturn > -1){
        iReturn,sRetJson = pkgUser.AndLoginUser(w,r, dat["UserID"].(string),"")
	}
	return iReturn, sRetJson
}


func cmdFloorLogin(w http.ResponseWriter,r *http.Request,sJson string)(int, string){
	sFunc := gPkgID
	sFunc += "cmdFloorLogin"
	sRetJson := ""    
    
	iReturn,dat := genFloorLogin(w,r,sJson)
	pkgError.VfDebug(sFunc, iReturn, sJson)
    if(iReturn > -1){
        iReturn,sRetJson = pkgUser.LoginUser(w,r, dat["UserID"].(string),"")
	}

	return iReturn, sRetJson
}

func genFloorLogin(w http.ResponseWriter,r *http.Request,sJson string)(int, map[string]interface{}){
	sFunc := "genFloorLogin"
	iReturn := -1
   
	var dat map[string]interface{}
	
	localError := json.Unmarshal([]byte(sJson),&dat)
	
	if pkgError.BfHandleError(pkgError.KLOGERROR,sFunc,localError){
        iReturn = 0
	}
	return iReturn, dat

}

func cmdFloorLogout(w http.ResponseWriter,r *http.Request)(int, string){
    pkgSession.ClearSession(w,r)
    return 0, ""
}

func mapParam(sJson string)(map[string]interface{}, error){

    var dat map[string]interface{}	
	err := json.Unmarshal([]byte(sJson),&dat)

    return dat,err
}


func andSetStore(w http.ResponseWriter, r *http.Request, sJson string)(int, string){
    sFunc := "andSetStore"
    iRet,sRetJson := genSetStore(w, r, sJson,sFunc)

    // sRetJson = `{"setstore":{"user":[-1,-1,-1],"dept":[[1,0,"Produce"],[2,1,"Bakery"]],"produce":[[6,1,"Salad Case"]]}}`
    return iRet,sRetJson
}

func cmdSetStore(w http.ResponseWriter, r *http.Request, sJson string)(int, string){
    sFunc := "cmdSetStore"
    return genSetStore(w, r, sJson,sFunc)
}
  
/* SP_SETSTORE Output        
    F1 t_produce.f_commodity,F2 t_produce.f_uid, F3 t_storedept.f_uid,
    F4 t_storedept.f_order,F5 t_storeloc.f_uid,F7 t_storedept.f_name, F8 t_storeloc.f_name))) as jr
        
*/       
    func genSetStore(w http.ResponseWriter, r *http.Request, sJson,sFunc string)(int, string){
    iReturn := -1
    sReturn := ""
    iPlu := -1

    dat,err := mapParam(sJson)

	if pkgError.BfHandleError(pkgError.KLOGERROR,sFunc,err){
        iReturn = 0
        iStoreNo := int(dat["IstoreNo"].(float64))
        bOk,sJsonStore := pkgDatabase.BsfExecuteDbFunc("SP_SETSTORE(" + strconv.Itoa(iStoreNo) + ")")
        if bOk{
            type objStore struct{            
            	F1 string
            	F2 string
            	F3 int
            	F4 int
                F5 int
                F6 string   
                F7 string
            }

            var x map[string][][]objStore

            iDept := -1
            sDept    := ""
            sLoc     := ""
            sProduce := ""
            iLastLoc := -1

            err := json.Unmarshal([]byte(sJsonStore),&x)

            if err == nil{
                for _, objArrays := range(x["Store"]){
                    for _,objArray := range(objArrays){
                        if iDept != objArray.F3{          
                            if !strings.Contains(sDept,"[" + strconv.Itoa(objArray.F3) + ","){
                                iDept = objArray.F3
                                if len(sDept) > 0{
                                    sDept += ","
                                }
                                sDept = sDept + "[" + strconv.Itoa(iDept) + "," + strconv.Itoa(objArray.F4) +
                                    `,"` + objArray.F6 + `"]`
                            }
                        }

                        if (iLastLoc != objArray.F5){           
                            if !strings.Contains(sLoc,"[" + strconv.Itoa(objArray.F5) + ","){
                                if len(sLoc) > 0{
                                    sLoc += ","
                                }
                                iLastLoc = objArray.F5
                                sLoc = sLoc + "[" + strconv.Itoa(iLastLoc) + "," + strconv.Itoa(iDept) + 
                                    `,"` + objArray.F7 + `"]`
                            }
                        }
                        
                        if len(objArray.F1) > 0{
                            if len(sProduce) > 0{
                                    sProduce += `,`
                            }
             
                            sProduce += `["` + objArray.F1 + `","` + objArray.F2 +  `",` + strconv.Itoa(objArray.F3) + 
                                `,` + strconv.Itoa(objArray.F5) + `]`
                        
                        }                        
                    }
                }

                iToggleCode := -1
                iStoreNo := -1
                cbu,iErr := pkgUser.GetClsBaseUser(w,r); if iErr == 0{
                    iToggleCode = cbu.IToggleCode 
                    iPlu = cbu.ISpId
                }

                sReturn = `{"setstore":{"user":[` + strconv.Itoa(iToggleCode) + "," + 
                    strconv.Itoa(iStoreNo) + "," +
                    strconv.Itoa(iPlu) + `],` 
                sReturn = sReturn + `"dept":[` + sDept + `],"loc":[` + sLoc + `],"produce":[` + sProduce + `]}}`
                iReturn = 0
            }
        }
    }
    return iReturn,sReturn
}

func cmdSaveInput(w http.ResponseWriter, r *http.Request, sJson string)(int, string){
    iReturn := -1
    sRetJson := ""
	cbu,iErr := pkgUser.GetClsBaseUser(w,r)
	if iErr == 0{
		sUserId := strconv.Itoa(cbu.IUid) 
		sStoreId := strconv.Itoa(pkgUser.IfGetStoreNoId(w,r))
        iReturn,sRetJson = genSaveInput(w,r,sJson,sUserId,sStoreId)    
    }

    return iReturn, sRetJson
}

func andSaveInput(w http.ResponseWriter, r *http.Request, sJson string)(int, string){
    return genSaveInput(w,r,sJson,"15","1")    
}

func genSaveInput(w http.ResponseWriter, r *http.Request, sJson,sUserId,sStoreId string)(int, string){
	sFunc := gPkgID
    sFunc += "genSaveInput"
    iReturn := -1
    sRetJson := ""
    
	type objProduceId struct {
		Spid int
		Boxes int 	
		Units int	
		IMsecs  int64
		Shrink [][]int
	}
	
	type objProduceData struct{
		ItransId int
		BoxTransactionTimeInSecs int64
		ArSpId []objProduceId
	}
	
	var x objProduceData 
	sSpcCode := "0"			// Carlos forcing code 0 need to change this to assign right code.
	localError := json.Unmarshal([]byte(sJson),&x)
	if pkgError.BfHandleError(pkgError.KLOGERROR,sFunc,localError){
		
		iReturn = 0;
		for iz := 0; iz < len(x.ArSpId); iz++{
			if iReturn == 0{						
				timeNow := time.Now()
				iSecs := x.BoxTransactionTimeInSecs
				timeTemp := timeNow.Add(time.Duration(-iSecs) * time.Second)
				sTime := timeTemp.Format("2006-01-02 15:04:05")
				
				bOk,sRet := pkgDatabase.BsfExecuteDbFunc("SP_ADDBOXES('" + sTime + "'," + sUserId  + "," + sStoreId  + "," + 
					strconv.Itoa(x.ArSpId[iz].Spid) + "," + strconv.Itoa(x.ArSpId[iz].Boxes) +
					"," + strconv.Itoa(x.ArSpId[iz].Units) + "," + strconv.Itoa(x.ItransId) + ")")
				if(!bOk){
					iReturn = -2;
					sRetJson = sRet
				}else{
					for iz1 := 0; iz1 < len(x.ArSpId[iz].Shrink); iz1++{
						if iReturn == 0{
							iSecs := x.ArSpId[iz].Shrink[iz1][2]
							fmt.Println("-------------------->",iSecs)
							timeTemp := timeNow.Add(time.Duration(-iSecs) * time.Second)
							sTime := timeTemp.Format("2006-01-02 15:04:05")

							fmt.Println("SP_SHRINKDONE('"+ sTime + "'	," + sUserId  + "," + sStoreId  + "," + 
								strconv.Itoa(x.ArSpId[iz].Spid) +  "," + strconv.Itoa(x.ArSpId[iz].Shrink[iz1][0]) + "," + 
								strconv.Itoa(x.ArSpId[iz].Shrink[iz1][1]) + "," + strconv.Itoa(x.ItransId)  + "," + 
								sSpcCode + ")")		
								
							bOk,sRet := pkgDatabase.BsfExecuteDbFunc("SP_SHRINKDONE('"+ sTime + "'	," + sUserId  + "," + sStoreId  + "," + 
								strconv.Itoa(x.ArSpId[iz].Spid) +  "," + strconv.Itoa(x.ArSpId[iz].Shrink[iz1][0]) + "," + 
								strconv.Itoa(x.ArSpId[iz].Shrink[iz1][1]) + "," + strconv.Itoa(x.ItransId)  + "," + 
								sSpcCode + ")")		
							
							if(!bOk){
								iReturn = -3;
								sRetJson = sRet
							}
							sRetJson = sRet
						}
					}
				}
			}else{
				break;
			}						
		}		
	}else{
		fmt.Println(x,localError);	
	}
    return iReturn,sRetJson
}


func cmdShrinkLogToHtml(w http.ResponseWriter, r *http.Request)(int,string){    
    // sFunc := "cmdShrinkLogToHtml"
    iReturn := -1
    sReturn := ""


    iStoreNoId := pkgUser.IfGetStoreNoId(w,r);
    if iStoreNoId > 0{
        bOk,sRet := pkgDatabase.BsfExecuteDbFunc("SP_GETSHRINKLOG(" + 
            strconv.Itoa(iStoreNoId)+ ")")
        if bOk{
            iReturn = 0
            sReturn = sRet
        }else{
            iReturn  = -3
        }
    }else{
        iReturn = -2;
    }     
    
    return iReturn,sReturn
}


func cmdGetTransactionId(w http.ResponseWriter, r *http.Request, sJson string)(int,string){
    sFunc := "cmdGetTransactionId"
    iReturn := -1
    sReturn := ""
    sComments := ""

    var dat map[string]interface{}
    err := json.Unmarshal([]byte(sJson),&dat)

    if pkgError.BfHandleError(pkgError.KLOGERROR | pkgError.KLOGPRINT ,sFunc,err){
        sComments = dat["sComments"].(string)

        bOk, sReturn := pkgDatabase.BsfExecuteDbFunc("SP_GETTRANSACTIONID(" + sComments + ")")

        if bOk{
            if len(sReturn) > 0{
                iReturn = 0
            }
        }
    }
    return iReturn,sReturn
}
/*-------------------- Commands End -------------------------------*/
    
/*----------------------- URL QUERIES START ----------------------------*/
func fnGetMain0(w http.ResponseWriter, r *http.Request) {
	bfSessionTemplateIn(w, getLanguage(w, r), "main0.tmpl", nil)
}

func fnGetXMain0(w http.ResponseWriter, r *http.Request) {
	bfSessionTemplateIn(w, getLanguage(w, r), "xmain0.tmpl", nil)
}


func fnGetDevice(w http.ResponseWriter, r *http.Request) {
	bfSessionTemplateIn(w, getLanguage(w, r), "device.tmpl", nil)
}

func fnGetSetupColors(w http.ResponseWriter, r *http.Request,objColors interface{}) {
	bfSessionTemplateIn(w, getLanguage(w, r), "setupcolors.tmpl", objColors)
}

func fnGetMain1(w http.ResponseWriter, r *http.Request,sJSON string, cbu pkgUser.ClsBaseUser) {
	sFunc := gPkgID
	sFunc += "fnGetMain1"

	if len(sJSON) > 0{		
        type objStruct struct{ 
                Stores []int
                Inv int
            }
        var storeStruct objStruct
        sStores := ""
        err := json.Unmarshal([]byte(sJSON),&storeStruct)	
        
	    if pkgError.BfHandleError(pkgError.KLOGERROR,sFunc,err){
            for _,iStore := range(storeStruct.Stores){
                if len(sStores) > 0{
                    sStores += ","
                }
                sStores +=  strconv.Itoa(iStore)
            }
            sTemplate := pkgTemplate.PstrTemplate(ksDefaultLang,"main1.tmpl")
            sTemplate = strings.Replace(sTemplate,"<!--CarlosStoreNo-->",sStores ,-1)
			if t,err := template.New("foo").Parse(sTemplate); err == nil{
				err = t.Execute(w,cbu)     
            }                             
        }                    

	}else{
		bfSessionTemplateIn(w, getLanguage(w, r), "main1.tmpl", nil)
	}
}

func fnGetDoShrink(w http.ResponseWriter, r *http.Request,sCodes, spId,sTid string,iStoreNoId int, cbu pkgUser.ClsBaseUser) {
	sFunc := gPkgID
	sFunc += "fnGetDoShrink"
	
    sTemplate := pkgTemplate.PstrTemplate(ksDefaultLang,"doshrink.tmpl")
    sTemplate = strings.Replace(sTemplate,"<!--CarlosCode-->",sCodes,-1)
    sTemplate = strings.Replace(sTemplate,"<!--CarlosStoreProduceIds-->",spId,-1)
    sTemplate = strings.Replace(sTemplate,"<!--CarlosTransactionId-->",sTid,-1)
    type stUser struct{
        First       string
        Last        string
        PluText     string
        PluID       string
        SpId        int
        IToggleCode int
        IStore      int
        IDept       int
        ILoc        int
    }            

    iToggleCode,_ := strconv.Atoi(r.Form["nmToggleCode"][0])
    iDept,_ := strconv.Atoi(r.Form["nmDept"][0])
    iLoc,_ := strconv.Atoi(r.Form["nmLoc"][0])
    iSpId,_ := strconv.Atoi(r.Form["nmSpId"][0])

    cbu.IToggleCode = iToggleCode

    pkgUser.PutClsBaseUser(w,r,cbu);

    vUser := stUser {cbu.SFirst, cbu.SLast,r.Form["nmProduceNm"][0],
    r.Form["nmPluId"][0],iSpId, iToggleCode, iStoreNoId,iDept,iLoc}
    
	if t,err := template.New("foo").Parse(sTemplate); err == nil{
        err = t.Execute(w,vUser)     
    }else{
        fmt.Println("fnDoShrink Error",err)
    }
    
}

/*============================== ORIGINAL SHRINKOUT END ================================*/

  
