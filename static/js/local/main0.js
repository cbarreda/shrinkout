"use strict";

function vfSetupScroll(){
    $(document).on('ready', function () {
        // -----------------------------
        //  On Scroll Resize Nav
        // -----------------------------
        $(window).scroll(function () {
            if ($('.main-nav').offset().top > 100) {
                $('.main-nav').removeClass('large');
            } else {
                $('.main-nav').addClass('large');
            }
        });
        // -----------------------------
        //  On Click Smooth scrool
        // -----------------------------
        // animation scroll js
        var html_body = $('html, body');
        $('.scrollTo').on('click', function () { //use scrollTo class in any HTML tag for scrolling
            if (location.pathname.replace(/^\//, '') === this.pathname.replace(/^\//, '') && location.hostname === this.hostname) {
                var target = $(this.hash);
                target = target.length ? target : $('[name=' + this.hash.slice(1) + ']');
                if (target.length) {
                    html_body.animate({
                        scrollTop: target.offset().top - 50
                    }, 1500, 'easeInOutExpo');
                    return false;
                }
            }
        });

        // easeInOutExpo Declaration
        jQuery.extend(jQuery.easing, {
            easeInOutExpo: function (x, t, b, c, d) {
                if (t === 0) {
                    return b;
                }
                if (t === d) {
                    return b + c;
                }
                if ((t /= d / 2) < 1) {
                    return c / 2 * Math.pow(2, 10 * (t - 1)) + b;
                }
                return c / 2 * (-Math.pow(2, -10 * --t) + 2) + b;
            }
        });
        // -----------------------------
        //  Testimonial Slider
        // -----------------------------
        $('.testimonial-slider').slick({
            autoplay: false
        });
        // -----------------------------
        //  Screenshot Slider
        // -----------------------------
        $('.screenshot-slider').slick({
            dots: true,
            slidesToShow: 3,
            centerMode: true,
            infinite: false,
            responsive: [{
                breakpoint: 768,
                settings: {
                    arrows: false
                }
            }]
        });
        // -----------------------------
        //  Video Replace
        // -----------------------------
        $('.video-box span.icon').click(function () {
            var video = '<iframe allowfullscreen src="' + $(this).attr('data-video') + '"></iframe>';
            $(this).replaceWith(video);
        });
        // -----------------------------
        //  Team Progress Bar
        // -----------------------------
        /* $('.team').waypoint(function () {
            $('.progress').each(function () {
                $(this).find('.progress-bar').animate({
                    width: $(this).attr('data-percent')
                });
            });
            this.destroy();
        }, {
            offset: 100
        }); */


        // -----------------------------
        //  Count Up
        // -----------------------------
        function counter() {
            var oTop;
            if ($('.count').length !== 0) {
                oTop = $('.count').offset().top - window.innerHeight;
            }
            if ($(window).scrollTop() > oTop) {
                $('.count').each(function () {
                    var $this = $(this),
                        countTo = $this.attr('data-count');
                    $({
                        countNum: $this.text()
                    }).animate({
                        countNum: countTo
                    }, {
                        duration: 1000,
                        easing: 'swing',
                        step: function () {
                            $this.text(Math.floor(this.countNum));
                        },
                        complete: function () {
                            $this.text(this.countNum);
                        }
                    });
                });
            }
        }
        $(window).on('scroll', function () {
            counter();
        });

    });
}



function vf_dispatchpage()
{
	$("#idMain0Login").unbind('click').click(bfLogin);     
	$("#idMain0ForgotP").unbind('click').click(bfForgotP);     
	$('#idPasswordShow').prop('checked', false);
	$("#idPasswordShow").on("click",vfPasswordShow)
	$("#idTest").unbind('click').click(bfTest);       
	$("#idSelectLanguage").on("change",vfLangChange);
	$("#idMain1Logout").unbind('click').click(bfLogout);
	
	$("#idMain1ProfEdit").unbind('click').click(bfProfEdit);
	

	$(".clsRowAccept").each(function(){
		var formatter = new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD',});
				
		let objA = $(this).find(".clsAccept");
		let iCents = objA.attr("iCents");
		$(this).find("#idBPay").text(gsfFormatDollars(iCents/100,'USD'));			
	});	

	vfSetupScroll();
	$("#idMain0LoginEmail").focus();
};

/********* Account Start */
function bfForgotP(){
	let sEmail = 	$("#idMain0LoginEmail").val().trim().toLowerCase();
	let sError = "";
	let sPassword = "";

	if (gifIsValidEmail(sEmail) < 0)
        alert(sfLocalMessage(0));
    else
    {
		let sJson = '{"email":"' + sEmail + '"}';
			
        vf_postJson("PstPassForgot",sJson);
	};
	
	return(false);
}

function bfPassChange(){
	alert("Password Change")
}

function bfProfEdit(){
	alert("Profile Edit");
}
function bfLogout(){
     vf_postNoJson("PstLogout");
}
/********* Account End */

function bfLogin(){
	let sEmail = 	$("#idMain0LoginEmail").val().trim().toLowerCase();
	let sError = "";
	let sPassword = "";

	if (gifIsValidEmail(sEmail) < 0)
		sError = sfLocalMessage(0);
	else
	{
		sPassword = gsf_removeQuotes($("#idMain0PassLogin").val()).trim();
		if (sPassword.length < 1)
			sError = sfLocalMessage(1);
	};

    if (sError.length > 0)
        alert(sError);
    else
    {
		let sJson = '{"email":"' + sEmail + '","password":"' + sPassword + '","remember":';
		
		if ($("#idMain0RememberLogin").is(":checked"))
			sJson += "1";
		else
			sJson += "0";
			
        vf_postJson("PstLogin",sJson +'}');
	};
	
	return(false);
}
function vfPasswordShow(){

	if ($(this).prop('checked'))
		$("#idMain0PassLogin").attr('type', 'text'); 
	else
		$("#idMain0PassLogin").attr('type', 'password'); 
}

function bfTest(){
	alert("bfTest")
}
function vfLangChange(){
	let arLang = ["english","spanish"];
	let iValue = this.value-1
	if (iValue > -1)
		vf_redirect(arLang[iValue] + "/");
};
