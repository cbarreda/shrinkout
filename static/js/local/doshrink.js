var kiTotal;
var kiTotal;
var kiTotal;
var kiDataIndex;
var gDataArray = [];
var gShrinkData;
var gITimeStampSecs;

function vf_dispatchpage()
{

    kiTotal = 0;
	kiDataIndex = 0;
	gITimeStampSecs = Date.now();
	
    vfAlphaNumInit("idQty",null,false,true);

	$("#idTestBtn").unbind('click').click(vfTestBtn);
    
	$("#idPlus").unbind('click').click(vfPlus);
	$("#idMinus").unbind('click').click(vfMinus);
	$("#idDone").unbind('click').click(vfDone);
    $("#idSoBtn").unbind('click').click(vfShrinkOut);   
    $("#idBxBtn").unbind('click').click(vfAddBoxes);    	
    $("#idUnBtn").unbind('click').click(vfAddUnits);    	
	$("#idRCodeBtn").unbind('click').click(vfReasonCode);
	$("#idVendorBtn").unbind('click').click(vfVendor);
    $("#idTableDiv").hide();
    $("#idData").attr("iMode",-1); // < 0 Boxes not entered,0 dont' ask for boxes,  > 0 save Boxes 
    gShrinkData = new produceData($("#idData").attr("iTransactionId"));
    vfCodeSelChanged(0);
    bfVendorSelChanged(0);
};

function vfTestBtn(){
      document.body.requestFullscreen();
};

function vfClearInput(){
	kiTotal = 0;
	$("#idQty").html("_");
	$("#idTotal").html("_");
};

function shrinkLine(pCode,pQty)
{
	this.ISecs = parseInt((Date.now() - gITimeStampSecs)/1000);
	this.iQty =  pQty;
	this.iCode = pCode;
};

function storeProduceId(tId)
{
	this.arShrink = [];
	this.iId = tId;
	this.iBoxes = 0;
	this.iUnits = 0;
	this.ISecs = 0;
	
	this.VfAddBoxes = function(pBoxes)
	{
		this.iBoxes = pBoxes;
		this.ISecs = parseInt((Date.now() - gITimeStampSecs)/1000);
		vfClearInput();
	};	
	
	this.VfAddUnits = function(pUnits)
	{
		this.iUnits = pUnits;
		this.ISecs = parseInt((Date.now() - gITimeStampSecs)/1000);
		vfClearInput();
	};	
	
	this.IfFindCode = function(pCode){
		var iIndex = this.arShrink.length;
		
		while(--iIndex > -1)
			if(this.arShrink[iIndex].iCode == pCode)
				break;
		return iIndex;
	};

	this.IfGetQty = function(pCode){
		var iQty = -1;
		var iIndex = this.IfFindCode(pCode);
		if(iIndex > -1)
			iQty = this.arShrink[iIndex].iQty;
			
		return iQty;	
	};
	
	this.BfCheckVendor = function(){		
		var bReturn = (this.iBoxes + this.iUnits) > 0;
		
		if (!bReturn){	
			var iIndex = this.arShrink.length;
			bReturn = true;			// in case there is no shrink
			while((--iIndex > -1) && bReturn)
				if (this.arShrink[iIndex].iQty > 0)
					bReturn = false;	// Shrink but no boxes or units
		}
		return bReturn;	
	};
	
	this.IfAddShrink = function(pCode,pQty)
	{
		 var iIndex = this.IfFindCode(pCode);
		
		if(iIndex < 0)
		{
			this.arShrink.push(new shrinkLine(pCode,pQty));
			iIndex = this.arShrink.length;
		}else		
			this.arShrink[iIndex].iQty = pQty;
			
		vfClearInput();
		return iIndex;
	};
	
	this.SfGetShrink = function(){
		var sReturn = "";
		var iIndex = this.arShrink.length;
		while(--iIndex > -1)
		{
			if(sReturn.length > 0)
				sReturn += ",";
			sReturn += "[" + this.arShrink[iIndex].iCode + "," + this.arShrink[iIndex].iQty + "," + 
				this.arShrink[iIndex].ISecs  + "]";
		};
		if(sReturn.length > 0)
			sReturn = "[" + sReturn + "]";
		
		return sReturn;
	};
	
	this.sfGetData = function(){
		var sReturn = "";
		 if((this.iBoxes + this.iUnits) > 0){
			sReturn = '{"spid":' + this.iId + ',"boxes":'  + this.iBoxes + ',"units":' + this.iUnits + ',"shrink":'	+ 
				this.SfGetShrink() + ',"BoxTransactionTimeInSecs":' + this.ISecs + '}';		
		 };
		 return sReturn;		 
	}
};
function produceData(tId)
{
	this.curSPIndx = -1;
	this.CodeId = -1;
	this.arSpId = [];
	this.transactionId = tId;
	
	this.IfAddSpId = function(pId)
	{
		var iIndex = this.arSpId.length;
		while(--iIndex > -1)
			if(this.arSpId[iIndex].iId == pId)
				break;
		this.curSPIndx = pId;
		if (iIndex < 0)
		{
			iIndex = this.arSpId.length;
			this.arSpId.push(new storeProduceId(pId));
		};
		this.curSPIndx = iIndex;
		this.VfDispIo();
		return iIndex;
	};
	
	this.VfDispIo = function(){
		var iSo = -1;
		var iBx = -1;
		var iUn = -1;
		
		if (this.curSPIndx > -1)
		{
			iSo = this.arSpId[this.curSPIndx].IfGetQty($("#idData").attr("iCode"));
			iBx = this.arSpId[this.curSPIndx].iBoxes;
			iUn = this.arSpId[this.curSPIndx].iUnits;
		};
		if(iSo < 0)
			iSo = "__";
		if(iBx < 0)
			iBx = "__";
		if(iUn < 0)
			iUn = "__";
			
		$("#idSoTxt").html(iSo);
		$("#idBxTxt").html(iBx);
		$("#idUnTxt").html(iUn);
	};
	this.IfAddShrink = function(pCode,pQty){
		iReturn = this.arSpId[this.curSPIndx].IfAddShrink(pCode,pQty);
		this.VfDispIo();
		return iReturn;
	};
	this.VfAddBoxes = function(pQty){
		iReturn = this.arSpId[this.curSPIndx].VfAddBoxes(pQty);
		this.VfDispIo();
		return iReturn;
	};
	this.VfAddUnits = function(pQty){
		iReturn = this.arSpId[this.curSPIndx].VfAddUnits(pQty);
		this.VfDispIo();
		return iReturn;		
	};
	this.IfCheckSave = function(){
		var iIndex = this.arSpId.length;
		while(--iIndex > -1)
			if (!this.arSpId[iIndex].BfCheckVendor())
				break;
				
		return iIndex ;	// > -1 Error		
	};
	this.fnSaveStr=function(){		
		var sReturn = "";		
		var sTemp = "";
		var iIndex = this.arSpId.length;
		
		while(--iIndex > -1){
			
			var sTemp1 = this.arSpId[iIndex].sfGetData();
			
			if(sTemp1.length > 0){
				if(sTemp.length > 0)
					sTemp += ",";
				sTemp += sTemp1;
			}
				
		}
			
		if(sTemp.length > 0)
			sReturn += '{"ItransId":' + this.transactionId+ ',"ArSpId":[' + sTemp + ']}';	
		
		return sReturn;
	}
};

function vfShrinkOut(){
	var iQty = sfGetTargetText()  * 1;
    
    if (iQty > 0)
		kiTotal += iQty;
    
	gShrinkData.IfAddShrink($("#idData").attr("iCode") ,kiTotal);	
	return false;	
};	

function vfAddBoxes(){
	var iQty = sfGetTargetText()  * 1;
    
    if (iQty > 0)
		kiTotal += iQty;
	gShrinkData.VfAddBoxes(kiTotal);		
	return false;	
};
function vfAddUnits(){
	var iQty = sfGetTargetText()  * 1;
    
    if (iQty > 0)
		kiTotal += iQty;
    
	gShrinkData.VfAddUnits(kiTotal);		
	return false;	
};

function vfReasonCode(){
    vfCodeSelChanged(1);
    $("#idTableDiv").show();
    $("#idNumDiv").hide();
    return false;
};

function vfVendor(){
    if(bfVendorSelChanged(1))
    {
        $("#idTableDiv").show();
        $("#idNumDiv").hide();
    }
    return false;
};


function vfDone(){
	var iError = gShrinkData.IfCheckSave();
	if (iError < 0){
		var sJson = gShrinkData.fnSaveStr();
		if (sJson.length > 0)
			vf_postJson("cmdSaveInput",sJson);
		else
			iError = 0;
	}
	else
	{
		alert(sfTranslate("~!nothingenteredkey@#"));
		iError = 0;
	}
	return false;
	gvfGoHome();

};

function vfBoxes(){
    $("#idData").attr("iMode",1);
    $("#idPlus").hide();
    $("#idMinus").hide();   
    $("#idTotal").hide();   
    $("#idQty").text("_");
    $("#idDone").text("BOXES");
    return false;
}


//------------- carlos replace this functions with a function that sends gShrinkData to the middle layer START

function vfSaveLine(sCommand){
    var iSpId =  $("#idData").attr("storeProduceId") * 1;

    if (iSpId < 0)
        alert(sfTranslate("~!selectvendorkey@#"));
    else{

        var sJson = '{"qty":' + kiTotal + ',"code":' +  $("#idData").attr("iCode") +
             ',"storeProduceId":' + iSpId +
            ',"pluid":' +  $("#idData").attr("pluid") + 
            ',"transid":' +  $("#idData").attr("iTransactionId") + '}';
    
        vf_postJson(sCommand,sJson);
    }    
};

//------------- carlos replace this functions with a function that sends gShrinkData to the middle layer END


function adder(iPlus){
    var sQty = sfGetTargetText() ;
    var iQty ;

    if (sQty.length < 1)
        sQty = '1';
    else
    {
        var sDoit = '';
        var bDecimal = false;
        if(sQty[0] == '.')
            sQty = '0' + sQty;

        for(var i = 0; i < sQty.length;i++)
        {
            if(((sQty[i] == '.') && !bDecimal) || ( (sQty[i] >= '0') && (sQty[i] <= '9')))
                sDoit += sQty[i]; 
        };

        if(sQty.length < 1)
        {
            alert(sfsQty + sfTranslate(" ~!notvalidnumber@#"));
        }
        else
            sQty = sDoit;
    }

    if(sQty.length > 0)
    {
        var iQty = sQty  * 1;
                
        if(iPlus == 1)
            kiTotal += iQty;
        else if(iPlus == 0)        
            kiTotal -= iQty;

        if(kiTotal < 1)
            kiTotal = 0;

        $("#idQty").text("_");

        $("#idTotal").text(kiTotal);
    };
}

function vfPlus(){
	adder(1);    
	return false;
};
function vfMinus(){
	adder(0);
	return false;
};
function vfSaveVendor(iTable,iSpId,sVendor)
{
	gShrinkData.IfAddSpId(iSpId);
    $("#idData").attr("storeProduceId",iSpId);// this line carlos marked for deletion
    $("#idVendorBtn").text(sVendor);
    if(iTable > 0){
        $("#idNumDiv").show();
        $("#idTableDiv").hide();
    }
};

function bfVendorSelChanged(iTable)
{
    var sHtml = "";
    var sJson = $("#idData").attr("arStoreProduceId");
    var arAr = $.parseJSON(sJson);
    var iVendorIndex = arAr.length;

    if(iVendorIndex > 0)
	{
        if(iVendorIndex == 1)
            vfSaveVendor(0,arAr[0][0],"(" + arAr[0][1] + ")" + arAr[0][2]);
        else
        {
            vfSaveVendor(0,arAr[1][0],"(" + arAr[1][1] + ")" + arAr[1][2]);
			if( iTable != 0)	// iVendorIndex must be >= 2
			{
				while(--iVendorIndex > 0)
				{
					var sVendor = "(" + arAr[iVendorIndex][1] + ")" + arAr[iVendorIndex][2];
				
					sHtml += "<tr><td><a class = 'btn btn-primary clsAnchor' href='javascript:vfSaveVendor(1," + 
						arAr[iVendorIndex][0] + ",\"" + sVendor + "\")'>" + sVendor +  "</td></tr>";       
				}
				
				if(sHtml.length > 0)        
					$("#idShrinkTable").html("<table>" + sHtml + "</table>");
			}
		}
    
    };

    return (sHtml.length > 0);
};

function vfSaveCode(iTable,iCode,sReason)
{
    $("#idData").attr("iCode",iCode);
    $("#idRCodeBtn").text("Code: " + sReason);
    if(iTable > 0){
        $("#idNumDiv").show();
        $("#idTableDiv").hide();
		gShrinkData.VfDispIo();
    };
};

function vfCodeSelChanged(iTable)
{
    var sHtml = "";
    var sJson = $("#idData").attr("sReasonCodeJson");
    var arAr = $.parseJSON(sJson);
 
    for(var iIndex = 0; iIndex < arAr.length;iIndex++)
    {
        var sReason = "(" + arAr[iIndex].f2 + ")" + arAr[iIndex].f3;

            sHtml += "<tr><td><a class = 'btn btn-primary clsAnchor' href='javascript:vfSaveCode(1," + 
                arAr[iIndex].f1 + ",\"" + sReason + "\")'>" + sReason +  "</td></tr>";       
    }

    if(sHtml.length > 0)
    {
        $("#idShrinkTable").html("<table>" + sHtml + "</table>");
        vfSaveCode(iTable,arAr[0].f1,"(" + arAr[0].f2 + ")" + arAr[0].f3 );
    }
};

function vfSetBoxInput()
{
    $("#idData").attr("iMode",1);
    vfBoxes();              
}

function sf_errorCode(ar_temp){
    sError = "Unknown error " + ar_temp.IntCode;

    switch(ar_temp.IntCode){
    }

    return sError
}

function vf_ajax(s_data)
{
    var ar_temp = $.parseJSON(s_data);
    if(ar_temp.IntCode < 0)
        alert(sf_errorCode(ar_temp));
    else
    {
        switch(ar_temp.StrCommand)
        {
			case "cmdSaveInput":
				var objSpc = $.parseJSON(ar_temp.SJSON);
				if (objSpc.SPCID > 0){
					window.location = gs_post + "/fnSpc?p_storeproduceId=" + objSpc.p_storeproduceId;
				}else{
					gvfGoHome();
				}
			break;
        }
    }
};			


