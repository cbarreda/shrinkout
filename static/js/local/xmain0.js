function vf_dispatchpage()
{
    vfAlphaNumInit("idMain0UserId",null,false,true);
    $("#idLogin").unbind('click').click(bfLogin);
};

function bfLogin(){
    var sVal = $("#idMain0UserId").text().trim()
    if (sVal.length > 1)
        sVal = sVal.substring(0,sVal.length -1);

    var sError = ""
    var sJson = ""

    if (/^\d+$/.test(sVal))
        sJson = '{"UserID":"' + sVal + '","idPassword":""}';
    else
        sError = sVal + " is not a valid Number"
    
    if (sError.length > 0)
        alert(sError);
    else
        vf_postJson("cmdFloorLogin",sJson);

	return(false);
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
			case "cmdFloorLogin":
                gvfGoHome();
			    break;
        }
    }
};
