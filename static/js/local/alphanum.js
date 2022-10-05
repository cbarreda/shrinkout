// call vfAlphaNumInit(targetTextId) where targetTextId is the id of the text to display keystrokes

var gTargetTextId;
var gfTargetTextChanged

function vfAlphaNumInit(targetTextId,lftargetTextChanged, bHasAlpha, bHasNum){
	gTargetTextId = '#' + targetTextId;
	gfTargetTextChanged  = lftargetTextChanged;

	$(".clsLetter").unbind('click').click(vfAlphaNumClick);
	$("#idSpace").unbind('click').click(vfAlphaNumSpace);
	$("#idClr").unbind('click').click(vfAlphaNumClear);
	$("#idDel").unbind('click').click(vfAlphaNumDel);
	$("#idDot").unbind('click').click(vfDot);
	$("#idSlash").unbind('click').click(vfSlash);

	if(bHasAlpha)
		vf_doletters(true);		
}

function bfAddSingle(sLetter){
	var sText = sfGetTargetText();
	if(!sText.includes(sLetter))
		vfAddLetter(sLetter);
	return false;
}

function vfAddLetter(sLetter){	
	 $(gTargetTextId).text(sfGetTargetText() + sLetter + '_');
	
	if (gfTargetTextChanged != null)
		gfTargetTextChanged();	
}

function vfAlphaNumClick(){	
	var sLetter = $(this).text();
	if(sLetter == '>')
		vf_doletters($("#idA1c1").text()=='Q');	
	else
		vfAddLetter(sLetter);	
	return false;
};

function vfAlphaNumClear(){
	$(gTargetTextId).text("_");
	if (gfTargetTextChanged != null)
		gfTargetTextChanged();	
	return false;
};

function sfGetTargetText(){
	var sReturn = $(gTargetTextId).text();
	if(sReturn.length < 1)
		sReturn = "_";
	return sReturn.substring(0,sReturn.length -1);	
};

function vfAlphaNumDel(){
	var sText = sfGetTargetText();
	if(sText.length > 0)
		sText=sText.substring(0,sText.length -1);

	$(gTargetTextId).text(sText + "_");		
	if (gfTargetTextChanged != null)
		gfTargetTextChanged();	
	return false;
};

function vfDot(){
	return bfAddSingle('.');
};

function vfSlash(){
	return bfAddSingle('/');
};

function vfAlphaNumSpace(){
	return bfAddSingle(' ');
};


function vfHideAlpha(bHide,bAlpha){
	if(bAlpha)
	{
		if(bHide)
			$("#idAlphaSection").hide();
		else
			$("#idAlphaSection").show();
	}
	else
	{
		if(bHide)
			$("#idNumSection").hide();
		else
			$("#idNumSection").show();
	};
};

function vf_doletters(bAlpha){
	var keys = "QWERTYUIOPASDFGHJKLZXCVBNM";
	var iCol =1;
	var iRow = 1;
	
	if (bAlpha)
		keys = keys.split('').sort().join('');
	
	for ( var iIndex = 0; iIndex < keys.length; iIndex++){						
		$("#idA" + iRow + "c" + iCol).text(keys[iIndex]);
		
		iCol++;
		if(iCol > 6)
		{
			iRow ++;
			iCol = 1;
		};
	}
};
