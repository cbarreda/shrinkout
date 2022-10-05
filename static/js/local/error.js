"use strict";
function vf_dispatchpage(){	
	$("#idErrorDisplay").html(
		sfTranslate($("#idData").attr("ErrorCallerKey")) + 
		"<p>" +	sfTranslate($("#idData").attr("ErrorMessageKey")) + "</p>" +
		"<p>" + $("#idData").attr("Message") + "</p>"
	);
	
};

