package pkgAdmin

import(
	"encoding/json"
    "net/http"	
	"strconv"
	
    "shrinkout/pkg/pkgError" 	
    "shrinkout/pkg/pkgDatabase"

    "shrinkout/pkg/pkgUser"

    "fmt"
)

var gPkg = "pkgAdmin."

func FnReport(w http.ResponseWriter, r *http.Request, sJson string)(int,string) {
    sFunc  := gPkg + "FnReport"
    iReturn := -1
    iStoreNoId := pkgUser.IfGetStoreNoId (w,r)
    sRetJson := ""
    var dat map[string]interface{}    
    sSql := ""
    localError := json.Unmarshal([]byte(sJson),&dat)
	if pkgError.BfHandleError(pkgError.KLOGERROR | pkgError.KLOGPRINT ,sFunc,localError){
        sFrom := dat["from"].(string)
        sTo := dat["to"].(string)
        iType := int(dat["type"].(float64))
        iReturn = 0
        sDate := ` and (s.f_timestamp >=''`+sFrom + `'') and (s.f_timestamp <=''` + sTo + `'') `;
        if(iType == 1) {
            sSql = ``
        } else {
            sSql = ``
        }
        switch iType{
            case 1: sSql = 
                `select  
                    u.f_first  ||  u.f_last Employee,o.f_produceid PLU, o.f_each Ea, o.f_lb LB,
                    o.f_casecost CaseCost, o.f_retail RetailCost, o.f_unitcost UnitCost,
                    s.f_qty Qty,s.f_code,c.f_reason, 
                    p.f_commodity Description
                 from 
                    t_pshrink s, t_user u,t_pcode c, t_storeproduce o,t_produce p, t_vendor v
                 where 	
                    (s.f_userid = u.f_uid) and (s.f_code = c.f_uid) and (s.f_storeNoId = ` + strconv.Itoa(iStoreNoId) + `)
                    and s.f_storeproduceid = o.f_uid 	
                    and o.f_produceid = p.f_uid 
                    and o.f_vendorid = v.f_uid` + sDate + `order by s.f_timestamp`                
                break;
            case 2: sSql = 
                `select count(*),sum(f_qty) "Total Qty",
                     date_trunc(''day'',s.f_timestamp) "Date", u.f_first,u.f_last 
                    from t_user u, t_pshrink s where (u.f_uid = s.f_userid) 
                    ` + sDate +`  
                      group by "Date", f_first,f_last order by "Date" desc`
                break;               
            case 3: sSql =
                `select  
                    o.f_produceid PLU, o.f_each Ea, o.f_lb LB,
                    o.f_casecost CaseCost, o.f_retail RetailCost, o.f_unitcost UnitCost,
                    sum(s.f_qty) Qty, 
                    p.f_commodity Description 
                 from 
                    t_pshrink s, t_storeproduce o, t_produce p  
                 where 
                    (s.f_storeNoId = ` + strconv.Itoa(iStoreNoId) + `)
                     and (s.f_storeproduceid = o.f_uid) and (o.f_produceid = p.f_uid) 
                    ` + sDate +`                    
                  group by o.f_uid,p.f_uid
                  order by p.f_uid`                    
                break;
            case  4: sSql = 
                `select  
                    o.f_produceid PLU, o.f_each Ea, o.f_lb LB,
                    o.f_casecost CaseCost, o.f_retail RetailCost, o.f_unitcost UnitCost,
                    sum(s.f_qty) Qty, 
                    c.f_uid,c.f_reason, 
                    p.f_commodity Description
                 from
                    t_pshrink s,  t_storeproduce o, t_produce p,t_pcode c
                 where (s.f_storeNoId = ` + strconv.Itoa(iStoreNoId) + `)
                    and (s.f_storeproduceid = o.f_uid) and (o.f_produceid = p.f_uid)
                  and (s.f_code = c.f_uid) 
                  ` + sDate +`                  
                  group by o.f_uid,p.f_uid, c.f_uid 
                  order by o.f_uid`
                break;
            case 5: sSql = sSql + 
                `select  
                    v.f_uid VendorId,v.f_name VendorName,
                    o.f_produceid PLU, o.f_each Ea, o.f_lb LB,
                    o.f_casecost CaseCost, o.f_retail RetailCost, o.f_unitcost UnitCost,
                    sum(s.f_qty) Qty, 
                    c.f_uid,c.f_reason, 
                    p.f_commodity Description        
                 from 
                    t_pshrink s,  t_storeproduce o, t_produce p,t_pcode c, t_vendor v
                    where 
                    (s.f_storeNoId = ` + strconv.Itoa(iStoreNoId) + `)
                     and (s.f_storeproduceid = o.f_uid) and (o.f_produceid = p.f_uid)
                    and (s.f_code =device c.f_uid) and (o.f_vendorid = v.f_uid)
                    ` + sDate +`                  
                    group by v.f_uid, o.f_uid,p.f_uid ,c.f_uid
                    order by v.f_name`
                break;
            case 6: sSql = `select u.f_uid UserId, u.f_first FirstName,u.f_last LastName,`+
                `l.f_itemid PLU,l.f_qty Qty,c.f_reason ShrinkOutput,l.f_timestamp DateTime`+
                ` from t_pshrinklog l, t_user u, t_pcode c`+
                ` where (l.f_storeNoId = ` + strconv.Itoa(iStoreNoId) + `)
                    and (u.f_uid = l.f_userid)  and (c.f_uid = l.f_code)`+
                ` order by u.f_first`
                break;
            case 7: sSql = `select f_deptid Dept, f_devicemac MAC, f_devicedesc Descr 
                    from t_device where f_storenoid=` + strconv.Itoa(iStoreNoId) + 
                    ` order by f_devicedesc`
                break;
        }
        
        if len(sSql) > 0{
            localError,sRet := pkgDatabase.SfHtmlTable(sSql);
	        if pkgError.BfHandleError(pkgError.KLOGERROR,sFunc,localError){
                sRetJson = sRet
            }else{
                iReturn = -2;
            }

        }
    }
    return iReturn,sRetJson
}

func FnAddPlu(w http.ResponseWriter, r *http.Request)(int,string){
    iReturn := -1
    sReturn := ""
    iStoreNoId := pkgUser.IfGetStoreNoId(w,r);
    if iStoreNoId > 0{
        bOk, sRet := pkgDatabase.BsfExecuteDbFunc("SP_GETSTORELOC(" + 
            strconv.Itoa(iStoreNoId)+ ")")
        if bOk{
            iReturn = 0
            sReturn = sRet
        }
    }

    return iReturn,sReturn
}

func FnSavePlu(w http.ResponseWriter, r *http.Request, sJson string)(int,string){
    iReturn := -1
    sReturn := ""    
    bOk, sRet := pkgDatabase.BsfExecuteDbFunc("SP_ADDPLU" + sJson  )
    if bOk{
        iReturn = 0
        sReturn = sRet
    }
    return iReturn,sReturn
}

func FnProcessLogFile(w http.ResponseWriter, r *http.Request)(int,string){
    sReturn := ""
    iReturn := -1
    iStoreNoId := pkgUser.IfGetStoreNoId(w,r);
    if iStoreNoId > 0{
        bOk, sRet := pkgDatabase.BsfExecuteDbFunc("SP_PROCESSLOGFILE(" + 
            strconv.Itoa(iStoreNoId) + ")")
        if bOk{
            iReturn = 0
            sReturn = sRet
        }
        
    }
    return iReturn,sReturn
}

// sJson [selectedIndex,qty,direction,t_inventoryline.f_uid]
func FnGetProdLocation(w http.ResponseWriter, r *http.Request, sJson string)(int,string){
    iReturn := -1
    sRetJson := ""
    bOk,sRet := pkgDatabase.BsfExecuteDbFunc("SP_GETPRODLOCATION(" +
             strconv.Itoa(pkgUser.IfGetStoreNoId(w,r)) + ")")
	if bOk{
	 	iReturn = 0
        sRetJson = sRet
    }else{
	    iReturn = -2  
    }
    return iReturn,sRetJson
}

func FnModifyProduct(w http.ResponseWriter,r *http.Request, sJson string)(int, string){
	iReturn := -1
    sRetJson := ""

    var m map[string]interface{}
    err := json.Unmarshal([]byte(sJson),&m)
    fmt.Println("FnModifyProduct",err)
    if err == nil{
    
        bOk,sRet:= pkgDatabase.BsfExecuteDbFunc("SP_MODIFYPRODUCT(" + 
            strconv.Itoa(int(m["id"].(float64))) + ",'" + 
            strconv.Itoa(int(m["prodId"].(float64))) + "','" + 
            m["name"].(string) + "'," + 
            strconv.Itoa(int(m["loc"].(float64))) +")")
            

        if bOk{
            iReturn = 0
            sRetJson = sRet
        }else{
            iReturn = -2
        }
    }

	return iReturn, sRetJson
}

func FnModifyLocation(w http.ResponseWriter,r *http.Request, sJson string)(int, string){
	iReturn := -1
    sRetJson := ""

    var m map[string]interface{}
    err := json.Unmarshal([]byte(sJson),&m)
   
    if err == nil{
    
        bOk,sRet:= pkgDatabase.BsfExecuteDbFunc("SP_MODIFYLOCATION(" + 
            strconv.Itoa(int(m["id"].(float64))) + ",'" + 
            m["name"].(string) + "')")
            

        if bOk{
            iReturn = 0
            sRetJson = sRet
        }else{
            iReturn = -2
        }
    }
	return iReturn, sRetJson
}


func FnAddLocation(w http.ResponseWriter,r *http.Request, sJson string)(int, string){
	iReturn := -1
    sRetJson := ""

    var m map[string]interface{}
    err := json.Unmarshal([]byte(sJson),&m)
   
    if err == nil{   
        bOk,sRet:= pkgDatabase.BsfExecuteDbFunc("SP_ADDLOCATION(" + 
            strconv.Itoa(int(m["idDept"].(float64))) + ",'" + 
            m["sName"].(string) + "')")        

        if bOk{
            iReturn = 0
            sRetJson = sRet
        }else{
            iReturn = -2
        }
    }
	return iReturn, sRetJson
}

func FnOrphans(w http.ResponseWriter,r *http.Request, sJson string)(int, string){
    iReturn := -1
    sReturn := ""
    var iLocId int
    err := json.Unmarshal([]byte(sJson),&iLocId)
    if err==nil{
        bOk,sRet:= pkgDatabase.BsfExecuteDbFunc("SP_ORPHANS(" + strconv.Itoa(iLocId) + ")")        
        fmt.Println("FnOrphans","SP_ORPHANS(" + strconv.Itoa(iLocId) + ")")
        if bOk{
            iReturn = 0      
            sReturn = sRet     
        }
    }
    return iReturn, sReturn 
}

func FnMoveProduct(w http.ResponseWriter,r *http.Request, sJson string)(int, string){
	iReturn := -1
    sRetJson := ""
    sSql := ""
    var arParam [][]int // [[iFrom,iToo],[p1,p2..]]

    err := json.Unmarshal([]byte(sJson),&arParam)
   
    if err == nil{

        sInP := "("
        
        for i,v := range(arParam[1]){
            if i > 0{ 
                sInP += ","
            }
            sInP += "'" + strconv.Itoa(v) + "'"
        }
        sInP += ")"
        sIn :=  "and f_uid in " + sInP

        if ((arParam[0][1] < 0) && (arParam[0][0] < 0)) {
            switch (arParam[0][1]){
                case -100: 
                    sSql = "delete from t_produce where f_uid in " + sInP;
                break;
                case -101: 
                    sSql = "delete from t_storeloc where f_uid in " + sInP;
                break;
            }
        }else{
            if arParam[0][1] < 0{     // Remove Product
                sSql = "delete from t_producelocation where f_storelocationid=" +
                strconv.Itoa(arParam[0][0]) + sIn;

            }else if(arParam[0][0] < 0 ){ // add product
                sSql = "delete from t_producelocation where f_storelocationid= " +
                    strconv.Itoa(arParam[0][1]) + sIn + "; insert into t_producelocation(f_storelocationid, f_storeproduceid) values"
                for i,v := range(arParam[1]){
                    if i > 0{
                        sSql += ","
                    }
                    sSql += "(" + strconv.Itoa(arParam[0][1]) + "," + strconv.Itoa(v) + ")"
                }            

            }else{ // move product
                sSql = " update t_producelocation set f_storelocationid = "+ strconv.Itoa(arParam[0][1]) +
                    " where f_uid > 0 "  + sIn;
            }
        }
        
        if len(sSql) > 0{
            bOk,sRet:= pkgDatabase.BsfExecuteDbFunc(sSql)
            if bOk{
                iReturn = 0
                sRetJson = sRet
            }   
        }
    }

    fmt.Println(arParam[0][0],arParam[0][1],"FnMoveProduct",sSql)      
	return iReturn, sRetJson

}
