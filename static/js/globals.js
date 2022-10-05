jQuery(document).ready(function(){vf_dispatchpage();});

var gs_post = "http://localhost:8080/";//UrlStringFlag
var gb_debug = false;//DebugFlag


function vf_redirect(strUrl)
{
	window.location = gs_post + strUrl;
}

function gvf_windowOpen(strName,strUrl){
	if (strName.length < 1)
		strName = "_blank";
	window.open(gs_post + strUrl, strName)
}

function gvfGoHome(){
	window.location = gs_post + "/client";
}

function gbf_isInt(value) {
  if (isNaN(value)) {
    return false;
  }
  var x = parseFloat(value);
  return (x | 0) === x;
}
