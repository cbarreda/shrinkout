"use strict";

function vf_dispatchpage(){
};

function vfResetColors(){
	vfSetColors(["#1a9b5a","#ffffff","#1a9b5a","#1a9b5a","#ffffff","#ffffff","#1a9b5a","#1a9b5a","#275e90","#ffffff",
	"#ffffff","#ffffff","#fe7800","#fe7800","#ffffff","#1a9b5a","#ffffff",
	"#ffffff","#ffffff","#ffffff","#ffffff","#ffffff","#ffffff","#ffffff","#ffffff",
	"#1a9b5a","#1a9b5a","#1a9b5a"]);
	vfSave();
}

function arOBJ(){
	return ["IdBody","IdHeaderC","IdNavBrandB","IdMainNavB","IdBannerH1C","IdBannerH5C","IdItem1B","IdItem2B",
	"IdItemAB","IdItem1H5C","IdItem2H5C","IdItemAH5C","IdItem1IC","IdItem2IC","IdItemAIC","IdFooterMainB","IdFooterMainLC",
	"IdSectionTitleH2","IdSectionTitleP","IdTeamCBH5","IdTeamCBP","IdTeamHBH5","IdTeamHBP","IdTeamSDH5","IdTeamSDP",
	"IdTeamCBB","IdTeamHBB","IdTeamSDB"];
}

function vfSetColors(arColor){
	var arValue = arOBJ();
	
	if(arValue.length != arColor.length)
		alert("Colors passed does not match array.")
	else{
		for (var i=0; i < arColor.length; i++){
			$("#" + arValue[i]).val(arColor[i])
		}
	}
		
}

function vfSave(){
	var arValue = arOBJ();
	var sCss = "";
	for(var i = 0; i < arValue.length; i++){
		var sValue = $("#" +arValue[i]).val();
		if (sValue.length > 0){
			if (sCss.length > 0)
				sCss = sCss + ',';
			sCss = sCss + '"' +arValue[i] + '":"' + sValue + '"';
		}
	}
	vf_postJson("PstSaveColors","{" + sCss + "}")
}

function vf_ajax(s_data)
{
    var ar_temp = $.parseJSON(s_data);
    if(ar_temp.IntCode < 0)
        alert(sfLocalMessage(ar_temp.IntCode));
    else
    {
		// sfLocalMessage(iCode) where iCode > 99 from main1 template, < 100 from main0 template
        switch(ar_temp.StrCommand)
        {
			case "PstSaveColors":
				alert("Save Colors");
				window.location="/setupcolors/";
			break;			
        }
    }
};


