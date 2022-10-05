package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"shrinkout/pkg/pkgDatabase"
	"shrinkout/pkg/pkgError"
	"shrinkout/pkg/pkgSession"
	"shrinkout/pkg/pkgTemplate"
	"shrinkout/pkg/pkgUser"
	"shrinkout/pkg/pkgUtil"

)

/*----------------  GLOBAL VARIABLES START -------------------------------*/

var (
	kConfigNode   pkgUtil.ClsConfigNode
	ksLangPath    = "./lang"
	ksDefaultLang = "english"

	gKey  = []byte("cb!bp!demo@te$ecret%Key")
	
	store = sessions.NewCookieStore(gKey)

	router = mux.NewRouter()

	gPkgID = "bpdemo."
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

		fs := http.FileServer(http.Dir("static"))
		http.Handle("/static/", http.StripPrefix("/static/", fs))

		fmt.Println("Template Read, ", pkgTemplate.IfReadTemplate(ksLangPath))

		router.NotFoundHandler = http.HandlerFunc(fnNotFound)

		router.HandleFunc("/fnRouter/", fnRouter)
		router.HandleFunc("/fnMultiPart/", fnMultiPart)

		router.HandleFunc("/test/", fnTest)
		router.HandleFunc("/readTemplate", fnTemplatesRead)
		router.HandleFunc("/", fnRoot)
		router.HandleFunc("/english/", fnRootE)
		router.HandleFunc("/spanish/", fnRootS)

		router.HandleFunc("/cmdacctcreate/", fnAcctCreate)
		router.HandleFunc("/cmdacctedit/", fnAcctEdit)
		router.HandleFunc("/cmdAboutUs/", fnAboutUs)

		http.Handle("/", router)
		localError = http.ListenAndServe(kConfigNode.ListenPort, nil)

		if localError != nil {
			pkgError.VfMark("func main", localError)
		}

	} else {
		pkgError.VfMark("func main", localError)
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
			case "PstLogin":
				objReturn.IntCode, objReturn.SJSON = PstLogin(w, r, sJSON)
				break
			case "PstPassForgot":
				objReturn.IntCode, objReturn.SJSON = PstPassForgot(w, r, sJSON)
				break
			case "PstPassChange":
				objReturn.IntCode, objReturn.SJSON = PstPassChange(w, r, sJSON)
				break
			case "PstAccountDel":
				objReturn.IntCode, objReturn.SJSON = PstAccountDel(w, r, sJSON)
				break
			}
		} else {
			switch objReturn.StrCommand {
			case "PstLogout":
				objReturn.IntCode, objReturn.SJSON = PstLogout(w, r)
				break
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

	pkgSession.PutLang(w, r, sLang)

	if r.Method == "GET" {

		bLoggedIn := false
		_, cbu := pkgUser.IfUserBaseClsGet(w, r)
		if cbu.IUid > 0 {
			bLoggedIn = len(cbu.SFirst) > 0
		}
		if bLoggedIn {			
			/* application may want to set page elements
			   ie: iCode, objMain1 := pkgBook.VfSetPage1(w, r, cbu)
			 */
			 iCode := 0
			 objMain1 := cbu
			
			if iCode > -1 {
				sTemplate := pkgTemplate.SfGetTemplate(sLang, "main1.tmpl")
				if len(sTemplate) > 0 {
					t := template.New("main1")
					t.Parse(sTemplate)
					t.Execute(w, objMain1)
				}
			}
		} else {
			fnGetMain0(w,r)
		}
	}
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

func fnAcctEditCreate(w http.ResponseWriter, r *http.Request, iFlag int) {
	var cbu pkgUser.ClsBaseUser

	_, cbu = pkgUser.IfUserBaseClsGet(w, r)
	if (iFlag > 0) || (cbu.IUid > 0) {
		arBirthday := strings.Split(cbu.SBirthday, " ")
		cbu.SBirthday = arBirthday[0]
		if len(cbu.SProfPic) < 1 {
			cbu.SProfPic = "./static/images/userSmall.png"
		}
	}

	bfSessionTemplateIn(w, getLanguage(w, r), "acctcreate.tmpl", cbu)
}

func fnAcctEdit(w http.ResponseWriter, r *http.Request) {
	fnAcctEditCreate(w, r, 1)
}

func fnAcctCreate(w http.ResponseWriter, r *http.Request) {
	fnAcctEditCreate(w, r, 0)
}


func fnAboutUs(w http.ResponseWriter, r *http.Request) {
	bfSessionTemplateIn(w, getLanguage(w, r), "aboutus.tmpl", nil)
}

/*----------------  FUNCTIONS END -------------------------------*/

/*----------------  ROUTER FUNCTIONS START -------------------------------*/

func fnMultiPart(w http.ResponseWriter, r *http.Request) {
	type clsReturn struct {
		IntCode    int
		StrCommand string
		SJSON      string
	}

	var (
		objReturn  clsReturn
		cbu        pkgUser.ClsBaseUser
		localError error
		sReturn    []byte
	)

	sFunc := gPkgID
	sFunc += "fnMultiPart"

	r.ParseMultipartForm(0)

	objReturn.StrCommand = r.FormValue("MPCMD")
	objReturn.SJSON = `{"command":"?"}`

	if objReturn.StrCommand != "mpAcctCreate" {
		objReturn.IntCode, cbu = pkgUser.IfUserBaseClsGet(w, r)
	}

	if objReturn.IntCode > -1 {
		switch objReturn.StrCommand {
		case "mpAcctCreate", "mpAcctEdit":
			objReturn.IntCode, objReturn.SJSON = pkgUser.AccountCreateEdit(w, r, cbu)
			break
		}

		sReturn, localError = json.Marshal(objReturn)

		pkgError.BfHandleError(pkgError.KLOGERROR, sFunc+"->"+objReturn.StrCommand, localError)

		w.Write(sReturn)
	}
}

func fnTest(w http.ResponseWriter, r *http.Request) {
	pkgSession.PrintSession(w, r)
}

func PstTest(w http.ResponseWriter, r *http.Request) {
}


func PstAccountDel(w http.ResponseWriter, r *http.Request, sJSON string) (int, string) {
	return pkgUser.AccountDelete(w, r, sJSON)
}

func PstLogin(w http.ResponseWriter, r *http.Request, sJSON string) (int, string) {
	return pkgUser.LoginUser(w, r, sJSON)
}

func PstPassForgot(w http.ResponseWriter, r *http.Request, sJSON string) (int, string) {
	return pkgUser.ForgotPass(w, r, sJSON)
}

func PstPassChange(w http.ResponseWriter, r *http.Request, sJSON string) (int, string) {
	return pkgUser.ChangePass(w, r, sJSON)
}

func PstLogout(w http.ResponseWriter, r *http.Request) (int, string) {
	return pkgUser.ClearSession(w, r)
}

/*----------------  ROUTER FUNCTIONS END -------------------------------*/


/*----------------------- URL QUERIES START ----------------------------*/
func fnGetMain0(w http.ResponseWriter, r *http.Request) {
	bfSessionTemplateIn(w, getLanguage(w, r), "main0.tmpl", nil)
}

/*----------------------- URL QUERIES END  ----------------------------*/
