function vf_post(s_post,psCommand, strJson){if(strJson == undefined) strJson = '';$.post(gs_post + s_post + "/"  ,{sCommand:psCommand,sData:strJson},vf_ajax)};	
function vf_postJson(psCommand,strJson){vf_post('fnRouter',psCommand,strJson);}
function vf_postNoJson(psCommand){vf_post('fnRouter',psCommand,'');}
function vf_postCodeData(psCommand,piCode,psData){
	if (psData.length < 1)
		psData = 'null';
	vf_postJson(psCommand,'{"Code":' + piCode + ',"Data":' + psData + '}' );
}

function vf_multiPost(sData,vfMultiResponse){
	   $.ajax({
        type: "POST",
        enctype: 'multipart/form-data',
        url: "/fnMultiPart/",
        data: sData,
        processData: false,
        contentType: false,
        cache: false,
        timeout: 600000,
        success: function (data) {vfMultiResponse(0,data);},
        error: function (e) {vfMultiResponse(-1,e.responseText);}
    });

};




