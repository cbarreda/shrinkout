/* Using pkgUser from original ShrinkOut Not from Boilerplate */
package pkgUser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"shrinkout/pkg/pkgDatabase"
	"shrinkout/pkg/pkgError"

	"shrinkout/pkg/pkgSession"
)

type ClsDevice struct {
	IUid       int
	SCust      string
	ICustId    int
	IStoreNoId int
	IDeptId    int
	IStore     int
	SDeviceMac string
	SCurrentP  string
}

type ClsBaseUser struct {
	IUid        int
	IAdminCode  int
	SFirst      string
	SLast       string
	SUser       string
	SLanguage   string // Used for Path to languages
	sCity       string
	IRemember   int
	IState      int
	ILanguage   int
	IToggleCode int
	ISpId       int // T_STOREPRODUCE.f_uid
}

var gPkgId = "pkgUser."

func IfSetPlu(w http.ResponseWriter, r *http.Request, iPlu int) int {
	iReturn := -1
	if iPlu < 0 {
		iReturn = 0
	} else {
		cbu, iErr := GetClsBaseUser(w, r)
		if iErr > -1 {
			cbu.ISpId = iPlu
			PutClsBaseUser(w, r, cbu)
			iReturn = 0
		}
	}
	return iReturn
}

func IfGetStoreNoId(w http.ResponseWriter, r *http.Request) int {
	iStoreNoId := -1

	cd, iErr := GetClsDevice(w, r)

	if iErr == 0 {
		iStoreNoId = cd.IStoreNoId
	}
	return iStoreNoId
}

func BfIsAdmin(w http.ResponseWriter, r *http.Request) bool {

	cbu, iErr := GetClsBaseUser(w, r)

	bReturn := (iErr == 0)

	if bReturn {
		bReturn = cbu.IAdminCode > 15
	}

	return bReturn
}

func PutClsDevice(w http.ResponseWriter, r *http.Request, cd ClsDevice) error {
	sFunc := gPkgId + "PutClsDevice"
	arBytes, localError := json.Marshal(cd)
	if pkgError.BfHandleError(pkgError.KLOGERROR|pkgError.KLOGPRINT, sFunc, localError) {
		return pkgSession.PutDevice(w, r, "ClsDevice", string(arBytes))
	}
	return localError
}

func GetClsDevice(w http.ResponseWriter, r *http.Request) (ClsDevice, int) {
	sFunc := gPkgId + "GetClsDevice"
	iReturn := -1
	var cd ClsDevice
	s, localError := pkgSession.GetDevice(w, r, "ClsDevice")

	if pkgError.BfHandleError(pkgError.KLOGERROR|pkgError.KLOGPRINT, sFunc, localError) {
		localError := json.Unmarshal([]byte(s), &cd)
		if pkgError.BfHandleError(pkgError.KLOGWARN, sFunc, localError) {
			iReturn = 0
		}
	}
	return cd, iReturn
}

func DispClsDevice(cd ClsDevice, sFunc string) {

	fmt.Println("Func->", sFunc, "\nIUid->", cd.IUid, "\nSCust->", cd.SCust, "\nICustId->", cd.ICustId, "\nIStoreNoId->", cd.IStoreNoId, "\nIDeptId->", cd.IDeptId, "\nIStore->", cd.IStore, "\nSDeviceMac->", cd.SDeviceMac, "\nSCurrentP->", cd.SCurrentP)
}

func PutClsBaseUser(w http.ResponseWriter, r *http.Request, cbu ClsBaseUser) error {
	sFunc := gPkgId + "PutClsBaseUser"

	arBytes, localError := json.Marshal(cbu)
	if pkgError.BfHandleError(pkgError.KLOGERROR|pkgError.KLOGPRINT, sFunc, localError) {
		return pkgSession.PutJson(w, r, "ClsBaseUser", string(arBytes))
	}

	return localError
}

func GetClsBaseUser(w http.ResponseWriter, r *http.Request) (ClsBaseUser, int) {
	var cbu ClsBaseUser
	sFunc := gPkgId + "GetClsBaseUser"
	iReturn := -1
	objError := pkgSession.GetJson(w, r, "ClsBaseUser")

	if pkgError.BfHandleError(pkgError.KLOGERROR|pkgError.KLOGPRINT, sFunc, objError.LocalError) {
		objError.LocalError = json.Unmarshal([]byte(objError.SInfo), &cbu)

		if pkgError.BfHandleError(pkgError.KLOGWARN, sFunc, objError.LocalError) {
			iReturn = 0
		}
	}

	return cbu, iReturn
}

/* SP_SETDEVICE returns {"code":iCode,"pass":"sNewPass","store":iStore};
   i_code
   <ul>
   <li>-1->Form Parse Failed,<br>
   <li>-2-> parse of nmDept or nmStoreNoId failed,
   <li>-3->SP_SETDEVICE( return failure
   <li>-4->Unmarshal failed
   <li>-5->parse of nmCust or nmCustId failed
   <li>0 New device inserted
   <li>-6->save cookie failed
   <li>< -99 Database Error (Here database errors are multiplied by 100).
   </ul>
*/

func SetDevice(w http.ResponseWriter, r *http.Request) {
	sFunc := "SetDevice"
	iReturn := -1
	bOk := true
	err := r.ParseForm()
	if pkgError.BfHandleError(pkgError.KLOGERROR, sFunc, err) {
		iReturn = -2
		sDept := r.PostFormValue("nmDept")
		sStoreNoId := r.PostFormValue("nmStoreNoId")
		if len(sDept) < 1 {
			sDept = "0"
		}
		sSp := r.PostFormValue("nmUserId") + ",'" +
			r.PostFormValue("nmPassword") + "'," +
			sStoreNoId + "," +
			sDept + ",'" +
			r.PostFormValue("nmMac") + "','" +
			r.PostFormValue("nmDesc") + "'"

		var sJson string
		bOk, sJson = pkgDatabase.BsfExecuteDbFunc("SP_SETDEVICE(" + sSp + ")")
		fmt.Println(bOk, "SP_SETDEVICE("+sSp+")"+","+sJson+"-")
		if bOk {
			iReturn = -3
			var m map[string]interface{}

			err = json.Unmarshal([]byte(sJson), &m)

			if pkgError.BfHandleError(pkgError.KLOGERROR, sFunc, err) {
				iReturn = -4
				iUid := int(m["code"].(float64))
				if iUid < 0 {
					iReturn = iUid * 100 // Need standard strategy to handle db codes. Here using DB > 99
				} else {
					var cd ClsDevice
					cd.SCust = r.PostFormValue("nmCust")
					cd.ICustId, err = strconv.Atoi(r.PostFormValue("nmCustId"))
					if err == nil {
						iReturn = -5
						cd.IUid = iUid
						cd.IStoreNoId, _ = strconv.Atoi(sStoreNoId)
						cd.IDeptId, _ = strconv.Atoi(sDept)
						cd.IStore = int(m["store"].(float64))
						// cd.SDeviceMac  string
						cd.SCurrentP = m["pass"].(string)
						err = PutClsDevice(w, r, cd)
						if err == nil {
							iReturn = 0
						}
					}
				}
			}
		}
	}
	if err != nil {
		fmt.Println("Error:", err)
		sCode := `<ul>
    <li>-1->Form Parse Failed,<br>
    <li>-2-> parse of nmDept or nmStoreNoId failed, 
    <li>-3->SP_SETDEVICE( return failure
    <li>-4->Unmarshal failed
    <li>-5->parse of nmCust or nmCustId failed
    <li>-6->save cookie failed
    <li>0 New device inserted
    <li>< -99 Database Error (Here database errors are multiplied by 100). 
    </ul>`
		fmt.Fprintf(w, "%d<br>%s", iReturn, sCode)
	} else {
		fmt.Fprintf(w, "OK %d", iReturn)
	}
}

func AndLoginUser(w http.ResponseWriter, r *http.Request, sEmp, sPassword string) (int, string) {
	iStoreId := 1 // in Session on html. Need to figure out for android

	iReturn, sRetJson := genUser(w, r, sEmp, sPassword, iStoreId)
	return iReturn, sRetJson
}

func LoginUser(w http.ResponseWriter, r *http.Request, sEmp, sPassword string) (int, string) {
	sRetJson := ""
	iReturn := IfGetStoreNoId(w, r) // gets store id

	if iReturn > -1 {
		iReturn, sRetJson = genUser(w, r, sEmp, sPassword, iReturn)
	}
	return iReturn, sRetJson
}

/* Returns: -1 No Device Id, -2 Db problem validating user, -3 Invalid user or password
 */
func genUser(w http.ResponseWriter, r *http.Request, sEmp, sPassword string, iStoreId int) (int, string) {
	sRet := ""
	sRetJson := ""
	iReturn := -55
	iLoginCode := 1 // store produce. Should get it from databaser

	var bOk bool

	sFunc := gPkgId + "genUser"

	sRet = "SP_LOGIN(" + strconv.Itoa(iLoginCode) + "," + strconv.Itoa(iStoreId) + "," + sEmp + ",'" + sPassword + "')"
	pkgError.VfDebug(sFunc, sRet)
	bOk, sRet = pkgDatabase.BsfExecuteDbFunc(sRet)

	if bOk {
		if len(sRet) > 0 {
			sRetJson = sRet
			iReturn = -4

			var det map[string]interface{}
			pkgError.VfDebug(sFunc, 2, sRet)
			localError := json.Unmarshal([]byte(sRet), &det)

			if pkgError.BfHandleError(pkgError.KLOGERROR, sFunc, localError) {
				var cbu ClsBaseUser
				cbu.SFirst = det["f_first"].(string)
				cbu.SLast = det["f_last"].(string)
				cbu.IUid = int(det["f_uid"].(float64))
				cbu.IAdminCode = int(det["f_admincode"].(float64))

				PutClsBaseUser(w, r, cbu)
				iReturn = 0
			}

		}
	}

	return iReturn, sRetJson
}
