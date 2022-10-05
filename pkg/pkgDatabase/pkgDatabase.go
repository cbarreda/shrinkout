package pkgDatabase

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"shrinkout/pkg/pkgError"
	"shrinkout/pkg/pkgUtil"
	"database/sql"

	_ "github.com/lib/pq"
)

var gPkgId = "pkgDatabase."

var m_strDb string

func SetDatabase(s string) {
	m_strDb = s
}

func setupDb() (*sql.DB, bool) {

	sFunc := gPkgId
	sFunc += "setupDb"

	objDb, localError := sql.Open("postgres", m_strDb)

	if len(m_strDb) < 1 {
		localError = errors.New("Empty Database Path")
	}
	return objDb, pkgError.BfHandleError(pkgError.KLOGERROR, sFunc, localError)
}

func openDatabase(strQuery string) pkgError.ClsError {
	var objError pkgError.ClsError

	objDb, bOk := setupDb()

	if bOk {
		defer objDb.Close()
		objError.LocalError = objDb.QueryRow("select " + strQuery).Scan(&objError.SInfo)
		if objError.LocalError == nil {
			objError.ICode = 0
		} else {
			objError.ICode = -2
			objError.SInfo = "Database:" + m_strDb + ";Query: " + strQuery + ";" + objError.LocalError.Error()
		}
	} else {
		objError.Set(-1, "Database setup failed")
	}

	return objError
}

func codeDataGet(sJSON, sLeftDelim, sRightDelim string) pkgError.ClsError {
	var objError pkgError.ClsError
	objError.Set(-1, "")
	arLines := strings.Split(sJSON, ",")
	if len(arLines) > 0 {
		if len(arLines) > 1 {
			objError.SInfo = pkgUtil.DataAsStringGet(sJSON, sLeftDelim, sRightDelim, `"`)
		} else {
			arLines = strings.Split(arLines[0], "}")
		}
		arLines = strings.Split(arLines[0], ":")
		if len(arLines) > 1 {
			objError.ICode, objError.LocalError = strconv.Atoi(arLines[1])
			if objError.LocalError != nil {
				objError.ICode = -2
			}
		}
	}

	return objError
}

func CodeDataExecuteDb(sCallFunc, strQuery string) pkgError.ClsError {

	objError := openDatabase(strQuery)

	if objError.ICode > -1 {
		objError = codeDataGet(objError.SInfo, "", "")
	}

	if objError.ICode < 0 {
		pkgError.BfHandleError(pkgError.KLOGERROR, "CodeDataExecuteDb->"+sCallFunc+", Query:"+strQuery, errors.New("Code:"+strconv.Itoa(objError.ICode)+", Error: "+objError.SInfo))
	}
	return objError
}

func CodeDataUnmarshal(sCallFunc, strQuery string, objInterface interface{}) pkgError.ClsError {

	objError := CodeDataExecuteDb(sCallFunc, strQuery)
	if objError.ICode > -1 {
		objError.Unmarshal(objInterface, "")
	}

	if objError.ICode < 0 {
		pkgError.BfHandleError(pkgError.KLOGERROR, "CodeDataExecuteDb->"+sCallFunc, errors.New("Code:"+strconv.Itoa(objError.ICode)+", Error: "+objError.SInfo))
	}
	return objError
}

func BsfExecuteDbFunc(strFunc string) (bool, string) {
	return bsfExecuteDbQuery("Select " + strFunc)
}


func bsfExecuteDbQuery(strQuery string) (bool, string) {
	sFunc := gPkgId
	sFunc += "bsfExecuteDbQuery"

	strFunc := "bsfExecuteDbQuery: " + strQuery
	var strResult string

	objDb, bOk := setupDb()

	if bOk {
		defer objDb.Close()
		localError := objDb.QueryRow(strQuery).Scan(&strResult)
		if localError != nil {
			bOk = pkgError.BfHandleError(pkgError.KLOGERROR, strFunc, localError)
		}
	}

	return bOk, strResult
}


func DbCodeStringDataGet(sJSON string) (int, string, error) {
	sFunc := gPkgId
	sFunc += "DbCodeStringDataGet"

	type x struct {
		Code int
		Data string
	}
	var x1 x

	localError := json.Unmarshal([]byte(sJSON), &x1)

	if localError != nil {
		x1.Code = -1
	}

	return x1.Code, x1.Data, localError
}

func DbCodeDataGet(sJSON string) (int, interface{}, error) {
	sFunc := gPkgId
	sFunc += "DbCodeDataGet"

	type x struct {
		Code int
		Data interface{}
	}
	var x1 x

	localError := json.Unmarshal([]byte(sJSON), &x1)

	if localError != nil {
		x1.Code = -1
	}
	return x1.Code, x1.Data, localError
}

// use when Data is a json object as in {} not a data type as in int 1 or string 'a..'
func DbCodeDataMapGet(sJSON string) (int, map[string]interface{}, error) {
	sFunc := gPkgId
	sFunc += "DbCodeDataMapGet"

	var r map[string]interface{}

	iCode, iInterface, localError := DbCodeDataGet(sJSON)
	if localError == nil && iCode > -1 {
		r = iInterface.(map[string]interface{})
	}
	return iCode, r, localError
}

// use when Data is a json object as in {} not a data type as in int 1 or string 'a..'
func DbCodeDataArMapGet(sJSON string) (int, []map[string]interface{}, error) {
	sFunc := gPkgId
	sFunc += "DbCodeDataArMapGet"

	type arX struct {
		Code int
		Data []map[string]interface{}
	}

	var x1 arX

	localError := json.Unmarshal([]byte(sJSON), &x1)
	if localError != nil {
		x1.Code = -1
	}

	return x1.Code, x1.Data, localError
}

func IfLimitOffset(sJSON string) (int, int) {
	type clsLO struct {
		ILimit, IOffset int
	}
	var objLO clsLO
	if json.Unmarshal([]byte(sJSON), &objLO) != nil {
		objLO.ILimit = -1
		objLO.IOffset = -1
	}
	return objLO.ILimit, objLO.IOffset
}



/* Added from the original shrinkout */

func SfCsvTable(sQuery string)(error,string){
    sFunc := "SfCsvTable: " 

    var sResult string
    var sTemp string
    var localError error
  
    objDb,bOk := setupDb()
    if bOk{
        defer objDb.Close()
        rows,err := objDb.Query("select csv_table('" + sQuery + "')")
        defer rows.Close() 
        if pkgError.BfHandleError(pkgError.KLOGERROR,sFunc,localError){
            for rows.Next(){
                if err := rows.Scan(&sTemp); err == nil{
                    sResult = sResult + sTemp                    
                }else{
                    localError = err                 
                }
            }
        }else{
            localError = err
        }
   }
   // pkgError.VfPrintln("SRESULT",sResult)
   return localError,sResult    
}

func SfHtmlTable(sQuery string)(error,string){
    sFunc := "sfHtmlTable: " 

    var sResult string
    var sTemp string
    var localError error
  
    objDb,bOk := setupDb()
    if bOk{
        defer objDb.Close()
        rows,err := objDb.Query("select html_table('" + sQuery + "')")
        defer rows.Close() 
        if pkgError.BfHandleError(pkgError.KLOGERROR,sFunc,localError){
            for rows.Next(){
                if err := rows.Scan(&sTemp); err == nil{
                    sResult = sResult + sTemp                    
                }else{
                    localError = err                 
                }
            }
        }else{
            localError = err
        }
   }
	
   return localError,sResult    
}
