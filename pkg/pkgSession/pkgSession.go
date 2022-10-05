// go get github.com/gorilla/sessions
package pkgSession

import (
	"fmt"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"

	"shrinkout/pkg/pkgError"
)

var (
	gKey         = []byte("cb!ShrinkOut$ecret%Key")
	gSessionName = "cb$e$$ionShrinkOut"
	gStore       = sessions.NewCookieStore(gKey)
	gPkgId       = "pkgSession/"
)

func PutLang(w http.ResponseWriter, r *http.Request, sLang string) error {
	return PutString(w, r, "language", sLang)
}

func GetLang(w http.ResponseWriter, r *http.Request) pkgError.ClsError {
	return GetString(w, r, "language")
}

func PutString(w http.ResponseWriter, r *http.Request, sJsonKey, sJson string) error {
	sFunc := gPkgId + "PutString"
	objSession, localError := gStore.Get(r, gSessionName)

	if pkgError.BfHandleError(pkgError.KLOGERROR|pkgError.KLOGPRINT, sFunc, localError) {
		objSession.Values[sJsonKey] = sJson

		localError = objSession.Save(r, w)
	}
	return localError

}

func PrintSession(w http.ResponseWriter, r *http.Request) {

	objSession, localError := gStore.Get(r, gSessionName)
	fmt.Println(gPkgId+"PrintSession", objSession, localError)
}

func GetString(w http.ResponseWriter, r *http.Request, sJsonKey string) pkgError.ClsError {
	sFunc := gPkgId + "GetString"
	var objError pkgError.ClsError
	objError.ICode = -1
	objSession, le := gStore.Get(r, gSessionName)

	objError.LocalError = le
	if pkgError.BfHandleError(pkgError.KLOGERROR|pkgError.KLOGPRINT, sFunc, objError.LocalError) {
		s, bOk := objSession.Values[sJsonKey]

		if bOk {
			objError.SInfo = s.(string)
			objError.ICode = 0
		} else {
			objError.Set(-2, "Get JSON "+sJsonKey+" FAILED")
		}
	}

	return objError
}

func PutJson(w http.ResponseWriter, r *http.Request, sJsonKey, sJson string) error {
	return PutString(w, r, sJsonKey, sJson)
}

func GetJson(w http.ResponseWriter, r *http.Request, sJsonKey string) pkgError.ClsError {
	return GetString(w, r, sJsonKey)
}

func ClearSession(w http.ResponseWriter, r *http.Request) {
	objSession, _ := gStore.Get(r, gSessionName)
	objSession.Options.MaxAge = -1
	objSession.Save(r, w)
}

/* Added from original ShrinkOut start */
func PutDevice(w http.ResponseWriter, r *http.Request, sJsonKey, sJson string)error{
    sFunc := gPkgId + "SaveDevice"
    objDevice := sessions.NewCookieStore(gKey)
    
    objSession, localError := objDevice.Get(r,gSessionName)

    if pkgError.BfHandleError(pkgError.KLOGERROR | pkgError.KLOGPRINT, sFunc,localError){
        objSession.Values[sJsonKey] = sJson

        localError = objSession.Save(r,w)
    }
    return localError
}

func GetDevice(w http.ResponseWriter, r *http.Request,sJsonKey string)(string,error){
    sFunc := gPkgId + "GetDevice"
    var sJson string

    objDevice := sessions.NewCookieStore(gKey)
    objSession, localError := objDevice.Get(r,gSessionName)
    if pkgError.BfHandleError(pkgError.KLOGERROR | pkgError.KLOGPRINT, sFunc,localError){
        s,bOk := objSession.Values[sJsonKey]
    
        if bOk{
            return s.(string),nil
        }

        return "",errors.New("Get JSON " + sJsonKey + " FAILED")
    }

    return sJson, localError
}
/* Added from original ShrinkOut End*/
