/* Default Page */
var GobjData;

function vf_dispatchpage()
{
    vfAlphaNumInit("idFilter",vfFilterChanged,true,true);
	$("#idPluBtn").unbind('click').click(vfPlu);    
	$("#idDeptSelBtn").unbind('click').click(vfDeptClick);
	$("#idGoBtn").unbind('click').click(vfGo);    
	$("#idMatchBtn").unbind('click').click(vfMatch);    
	$("#idToggleBtn").unbind('click').click(vfToggle);    
	$("#idLogout").unbind('click').click(bfLogout);    
	$("#idStoreNoSel").unbind('click').click(vfStoreNoClick);
	$("#idDeptLocBtn").unbind('click').click(vfLocClick);       
    $("#idStoreNoSel").text(ifGetStoreButtonIndex(-1));
    $("#idItemSel").on("change",vfItemSelChange)
    $("#idMain1Letters").on("change",vfMain1LettersChange)
    vfSetStore(0);
};

function vfPlu(){
    $("#idSelTable").hide();
    $("#idByAlphaNum").show();
    $("#idGo").show();
    return false;
};

function vfMain1LettersChange(){
    vfProdClick($("#idMain1Letters").val())
};

function vfItemSelChange(){
    $("#idGo").attr("iIndex", $("#idItemSel").val()) ;
};

function vfGo()
{
    var sError = "";
    var iIndex = $("#idGo").attr("iIndex") * 1;
    var iFlag = iIndex;

    if (iFlag < 0)
    {
        var iTemp = sfGetTargetText() * 1;
        if (iTemp > 0)
        {
            if(confirm("Log " + iTemp + "?")== true)
            {
                iIndex = 0 - iTemp;
                $("#idPluId").val(iIndex);
                iFlag = 0;
            }
            else
                iIndex = 0;        
        }else{
            iIndex = 0;
            alert(sfTranslate("~!nothingtosavekey@#"));
        };
    }
    
    if(iFlag > -1)
        vfSubmitMain1Form(iIndex,0);
     return false;
};

function vfByPlu()
{
    $("#idByTable").show();
    $("#idByAlphaNum").hide();
    vfHideAlpha(true,true);
    vfHideAlpha(false,false);
    vfPlu();
}

function vfByAlpha()
{
    $("#idByTable").hide();       
    $("#idByAlphaNum").show();
    vfHideAlpha(false,true);
    vfHideAlpha(true,false);
}

function vfToggle()
{
    $("#idGo").show();
    if($("#idByTable").is(":visible"))
        vfByAlpha();    
    else
        vfByPlu();
	return false;
}

function sfGetItem(arAr,iId,iCol){
    var iIndex = arAr.length;
    var sReturn = "";
    while(--iIndex > -1)
        if(iId == arAr[iIndex][0])
        {
            sReturn = arAr[iIndex][2]   
            break
        }
    
    return sReturn;
}

function vfMatch()
{

  	var sKey = sfGetTargetText().toLowerCase();

    $("#idGo").attr("iIndex", -1)            
    
    if (GobjData)
    {
  	    var sKey = sfGetTargetText().toLowerCase();
        var objData = GobjData.setstore.produce;
        var sSelect = "";

        for (var iIndex = 0;iIndex < objData.length;iIndex++)
            if(objData[iIndex][0].toLowerCase().includes(sKey))
            {
                if (sSelect.length == 0)
                    $("#idGo").attr("iIndex", iIndex);
                    
                sSelect += "<Option Value = '" + iIndex + "'>" + objData[iIndex][0] + "/" +
                    sfGetItem(GobjData.setstore.loc,objData[iIndex][3],2) ;
                } 
                
        $("#idItemSel").empty().append(sSelect);
    };
	return false;
};


function vfFilterChanged(){    
  	var sKey = sfGetTargetText();

    $("#idGo").attr("iIndex", -1)            
    
    if (GobjData)
    {
        var sSelect = "";
        var iItemsSelected = 0;

        var objData = GobjData.setstore.produce;

        var arKey;
        var iIndex = 0;    
        
        if (gbf_isInt(sKey))
            arKey = 1;
        else
            arKey = 0;

        while(iIndex < objData.length)        
        {
            var sCompare = objData[iIndex][arKey].substring(0,sKey.length);

            if(sCompare.toLowerCase() == sKey.toLowerCase())
            {
                if (iItemsSelected > 10)          
                    iIndex = objData.length;
                else
                { 
                    if (sSelect.length == 0)
                        $("#idGo").attr("iIndex", iIndex);
                        
                    sSelect += "<Option Value = '" + iIndex + "'>" + objData[iIndex][0] + "/" +
                        sfGetItem(GobjData.setstore.loc,objData[iIndex][3],2) ;                              

                    iItemsSelected++;
                };
            };
            iIndex++;
        }; 

        $("#idItemSel").empty().append(sSelect);
    };
    
};

// [storeNo,StoreNo...]
function ifGetStoreButtonIndex(iIndex){  
    if(iIndex < 0)  
        iIndex = $("#idStoreNoSel").attr("iIndex");
    else
        $("#idStoreNoSel").attr("iIndex",iIndex);

    var arStores = $("#idStoreNoSel").attr("arStore").split(',');
    return arStores[iIndex];
}

// [[DeptId,deptOrder,deptDesc],[DeptId,deptOrder,deptDesc]...]
function ifGetDeptButtonIndex(iIndex,iCol){  
    var sReturn = "?";
    var iStore = ifGetStoreButtonIndex(-1);

    if (GobjData)
    {
        
        var arAr = GobjData.setstore.dept;
        
        if(iIndex < 0)  
        {
            if(iIndex == -1)
                iIndex = $("#idDeptSel").attr("iIndex");
            else
            {
                var iMin = -1;
                for (var i = 0; i < arAr.length; i++){
                    if((iMin == -1) || (arAr[1] < iMin))
                        iIndex = i;
                }
            }
        }
        
        if((iIndex > -1) && (iIndex < arAr.length))
        {
            sReturn = arAr[iIndex][iCol];
            $("#idDeptSel").attr("iIndex",iIndex);
        };
       
    }
    return sReturn;
}

//[[locId,deptId,locName],[locId,deptId,locName]...]
function ifGetLocButtonIndex(iIndex,iCol){  
    var sReturn = "?";
    var iDept = ifGetDeptButtonIndex(-1,0);

    if (GobjData)
    {
         var arAr = GobjData.setstore.loc;
        
        if(iIndex < 0)  
        {
            if(iIndex == -1)
                iIndex = $("#idDeptLoc").attr("iIndex");
            else
            {
                for (iIndex = 0; iIndex < objData.length; iIndex++){
                    if (objData[iIndex][1] == iDept)
                        break;                    
                    }
            }
        }        

        if((iIndex > -1) && (arAr.length > 0))
        {        
            $("#idDeptLoc").attr("iIndex",iIndex);
            sReturn = arAr[iIndex][iCol];      
        }
        
    };
    return sReturn;
}

function vfStoreNoClick(){    
    var sHtml = "";
    var arAr = $("#idStoreNoSel").attr("arStore").split(',');    
 
    for(var iIndex = 0; iIndex < arAr.length;iIndex++)
            sHtml += "<tr><td><aclass = 'btn btn-primary clsAnchor' href='javascript:cmdSetStore(" + iIndex + ")'>" + arAr[iIndex] + "</td></tr>";

    if(sHtml.length > 0)
        $("#idSelTable").html("<table>" + sHtml + "</table>")    
}

function vfDeptClick(){    
    $("#idSelTable").show();
    $("#idByAlphaNum").hide();
    $("#idGo").hide();

    var sHtml = "";

    if (GobjData)
    {
        var arAr = GobjData.setstore.dept;

        for(var iIndex = 0; iIndex < arAr.length;iIndex++)
            sHtml += "<tr><td><a class = 'btn btn-primary clsAnchor' href='javascript:vfSetDept(" + iIndex + ")'>" + arAr[iIndex][2] + "</td></tr>";
    };
    if(sHtml.length > 0){
        $("#idSelTable").html(sHtml)    
    }
    
    return false;
    
};

function vfLocClick(){    
    $("#idSelTable").show();
    $("#idByAlphaNum").hide();
    $("#idGo").hide();
    
    var sHtml = "";
    if (GobjData)
    {
        var iDept = ifGetDeptButtonIndex(-1,0);
        var arAr = GobjData.setstore.loc;
        
        for(var iIndex = 0; iIndex < arAr.length;iIndex++)
            if (arAr[iIndex][1] == iDept)
                sHtml += "<tr><td><a class = 'btn btn-primary clsAnchor' href='javascript:vfSetLoc(" + iIndex  +  ")'>" + arAr[iIndex][2] + "</td></tr>";
    };
        
    if(sHtml.length > 0)
        $("#idSelTable").html("<table>" + sHtml + "</table>")            
    return false;
};

// see produce.go.cmdSetStore
function vfProdClick(cChar){
    var sHtml = "";
    
    if (GobjData)
    {
        var iLoc =  ifGetLocButtonIndex(-1,0)

        var arAr = GobjData.setstore.produce;

        for(var iIndex = 0; iIndex < arAr.length;iIndex++)
            if ((arAr[iIndex][3] == iLoc) &&((cChar=='-') || (arAr[iIndex][0][0].toUpperCase() == cChar)))
            {
                var sItem = arAr[iIndex][0];
                sHtml += "<tr><td><a class = 'btn btn-primary clsAnchor' href='javascript:vfSubmitMain1Form(" + iIndex + ",1)'>" + sItem.substring(0,45) + "</td></tr>";
            }
    };

    if(sHtml.length > 0)
        $("#idSelTable").html("<table>" + sHtml + "</table>")            
};


function vfSetStore(iIndex){
    vf_postJson("cmdSetStore",'{"IstoreNo":' + ifGetStoreButtonIndex(iIndex) + '}');
};


function vfSetDept(iIndex){
    $("#idDeptLbl").text(ifGetDeptButtonIndex(iIndex,2));
    $("#idSelTable").html("");
    vfSetLoc(-1);
    vfProdClick('-');
};

function vfSetLoc(iIndex){
    $("#idSelTable").html("");
    $("#idDeptLoc").attr("iIndex",iIndex);
    $("#idLocLbl").text(ifGetLocButtonIndex(iIndex,2));
    vfProdClick('-');
};


function vfRememberMe(iPlu){
    $("#idByTable").hide();
    vfToggle();
    $("#idSelTable").show();
    
    var i = GobjData.setstore.produce.length;

    while(--i > -1)
        if(GobjData.setstore.produce[i][1] == iPlu)
            break;
    if (i > -1){
        var iLoc =   GobjData.setstore.produce[i][3];
        i = GobjData.setstore.loc.length;
        while (--i > -1)
            if (GobjData.setstore.loc[i][0] == iLoc)
                break;
        if(i>-1)
            vfSetLoc(i);
    }
    /*
    vfSetDept(iDept);
    */
}
// iCode 0 called from next button, iCode 1 called from location
function vfSubmitMain1Form(iIndex,iCode)
{
    var iDept = $("#idDeptSel").attr("iIndex");
    var iLoc = $("#idDeptLoc").attr("iIndex");
    if (iCode == 0)
    {
        iCode = 1;
        iDept = -1;
        iLoc = -1;
    };


    $("#idToggleCode").val(iCode);
    $("#idDept").val(iDept);
    $("#idLoc").val(iLoc);

    $("#idArray").val(iIndex);
    $("#idMain1Form").submit();
};

/* f1 t_storedept.f_name, f2 t_storeloc.f_name, f3 t_producelocation.f_storeproduceid, 
    f4 t_produce.f_commodity, f5 vendor id, f6 vendor f_name f7 t_storeloc.f_uid 
    f8 t_storeproduce.f_produceid, f9 t_storedept.f_order*/
function bfValidateMain1Form()
{
    var iId =  $("#idArray").val();
    var bOk = GobjData.setstore.produce.length > 0;

    if(bOk)
    {
        var arAr = GobjData.setstore.produce;
        var bOk = (iId < arAr.length) && (arAr.length > 0);
    };

    if (bOk)
    {   
        if (iId > -1)
        {
            $("#idProduceNm").val(arAr[iId][0]);           
            $("#idPluId").val(arAr[iId][1]);
            $("#idSpId").val(arAr[iId][4]);
        }else{
            $("#idProduceNm").val(0);
            $("#idPluId").val(0);
            $("#idSpId").val(0);
        }
    };   
    return bOk;
};


function bfLogout(){
    vf_postJson("cmdFloorLogout",'{"Command":"cmdFloorLogout"}');
	return(false);    
}

function vf_ajax(s_data)
{
    var ar_temp = $.parseJSON(s_data);
    if(ar_temp.IntCode < 0)
        alert( sfLocalMessage(ar_temp.IntCode));
    else
    {
        switch(ar_temp.StrCommand)
        {
			case "cmdFloorLogout":
                gvfGoHome();
			    break;
            case "cmdSetStore":                
                $("#idStoreNoSel").text(ifGetStoreButtonIndex(-1));
                $("#idSelTable").html("");
                $("#idStoreNoSel").attr("sJson",ar_temp.StrJson);                    
                if (ar_temp.SJSON.length > 0){
                    $("#idDeptLbl").text(ifGetDeptButtonIndex(-2,0));
                };
                GobjData =  $.parseJSON(ar_temp.SJSON);

                if( GobjData.setstore.user[2] > 0 )
                    vfRememberMe(GobjData.setstore.user[2]);
                else
                    vfRememberMe(-1);
                    
					$("#idSelTable").hide();
                break;
        }
    }
};				


