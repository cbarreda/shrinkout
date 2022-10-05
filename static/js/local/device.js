
function vf_dispatchpage(){    
	$("#idPasswordShow").on("click",vfPasswordShow)
};

function CheckfnDeviceSetup(){
    var bOk;
    bOk = (gs_cust != undefined) && (gs_custId != undefined);
    
    if (bOk)
        bOk = (gs_cust.length >0 ) && (parseInt(gs_custId) > 0);
    if (!bOk)
		alert(sfTranslate("~!configerror@#"));

    if (bOk){
        $("#idCustId").val(gs_custId);
        $("#idCust").val(gs_cust);
    }
    
    return bOk;
}

function sf_errorCode(ar_temp){
    switch(ar_temp.IntCode){
        case -1: sError = sfTranslate("~!devicesetuperror@#");break;
        case -12: sError = sfTranslate("~!devicesetuperror@#");break;
        case -13: sError = sfTranslate("~!devicesetuperror@#");break;
		default: sError = sfTranslate("~!unknownerror@# ") + ar_temp.IntCode; break;
    }

    return sError
}
function vfPasswordShow(){

	if ($(this).prop('checked'))
		$("#idMain0PassLogin").attr('type', 'text'); 
	else
		$("#idMain0PassLogin").attr('type', 'password'); 
}
