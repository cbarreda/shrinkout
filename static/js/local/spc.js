/* Default Page */
var objCanvas ;
var objCtx ;
var objCanvas = document.getElementById("myCanvas");
var pointCL  ; // = parseInt(objCanvas.height/2);
var pointUCL ; // = parseInt(pointCL - (pointCL/2));
var pointLCL ; // = parseInt(pointCL + (pointCL/2));
var giOffset = 0;
var giOffsetX;
var giOffsetY;
var giUCL;
var giCL;
var giLCL;
var giHigh;
var giLow;
var giOffset = 100;

function vfHome(){
	gvfGoHome();
};

function vf_dispatchpage()
{
	$("#idHomeBtn").unbind('click').click(vfHome);
	objCanvas = document.getElementById("myCanvas");
	objCtx = objCanvas.getContext("2d");
	objCanvas.width = $("#canvasDiv").width();
	objCanvas.height = $("#canvasDiv").height();
	objCtx.fillStyle = "#FFFFFF";		
	objCtx.fillRect(0,0,objCanvas.width,objCanvas.height);		
	vfDrawDiagonal(0,0,objCanvas.width,objCanvas.height);
	vfDrawChart();
};

function vfPointRun(arAr,iCounter,iz){
	var iReturn = 0;
	var sColor = "";
	if(arAr[iz][0] > giUCL){
		if (iCounter < 0)
			iCounter--;
		else{
			if (iCounter > 4)
			{
				iReturn = iCounter;
				sColor = "green";
				iCounter = -1;
			};
		}		
	}
	else{
		if(arAr[iz][0] < giLCL){
			if(iCounter > 0)
				iCounter++;
			else{
				if (iCounter < -4)
				{
					iReturn = abs(iCounter);
					sColor = "red";
					iCounter = 1;
				};
			}
		}else
			iCounter = 0;
	};
	
	
	while(--iReturn > -1){
		arAr[iz - iReturn][2] = sColor;
	}
	
	return iCounter;
};

function vfDrawArray(arPoints,arParams){
	var arAr = new Array();
	var sColor = "";
	var iCounter = 0;
	vfSetup(arParams, arPoints.length);
	for(var iz = 0; iz < arPoints.length;iz++){
		if (arPoints[iz] > giUCL)
			sColor = "red";
		else if (arPoints[iz] < giLCL)
			sColor = "green";
		else
			 sColor = "blue";	
			 
		arAr.push([arPoints[iz],sColor]);
	};
	
	var iz = 0;
	var iz1 = 0;
	while(iz < arAr.length){
		iz1 = iz;
		while(
			(++iz1 < arAr.length) && 
			( 
			
				(
					(arAr[iz][0] == giCL) && (arAr[iz1][0] == giCL)
				) || 					 	
				(
					(arAr[iz][0] >giCL) && (arAr[iz1][0] > giCL)
				) || 
				(
					(arAr[iz][0] < giCL) && (arAr[iz1][0] < giCL)
				) 
			))
		
		if (iz1 == arAr.length)
			iz1--;
			
		if((iz1 - iz) > 4)
		{
			if((arAr[iz][0] == giCL) || (arAr[iz][0] < giCL))
				sColor = "green";
			else
				sColor = "red";
				
			while(iz < iz1)
				arAr[iz++][1] = sColor;
		
		};
		
		iz = iz1;
	}
	
	
	for(var iz = 0; iz < arAr.length;iz++)
		vfDrawPoint(iz + 1,arAr[iz][0], (iz == 0),arAr[iz][1]);
};

function vfDrawDiagonal(x,y,x1,y1){
	objCtx.beginPath();
	objCtx.fillStyle = "#000000";		
	objCtx.moveTo(x,y);
	objCtx.lineTo(x1,y1);
	objCtx.stroke();		
};


function vfDrawLine(x,y,iValue){
	
	objCtx.font = "10px Georgia";
	
	var yv = ifPointY(y);
	objCtx.beginPath();
	objCtx.fillStyle = "#000000";		
	objCtx.moveTo(ifPointX(x),yv);
	objCtx.lineTo(objCanvas.width,yv);
	objCtx.stroke();		
	objCtx.moveTo(ifPointX(x),ifPointY(y));
	objCtx.fillText(iValue,ifPointX(x),ifPointY(y));
	
};



function vfSetup(arParams,iPoints){
		// objCanvas.width = window.innerWidth;
		// objCanvas.height = window.innerHeight;
		giHigh = arParams[0],
		giLow = arParams[1];
			
		giUCL = arParams[2];
		giLCL = arParams[3];
		giCL = (giUCL + giLCL)/2;
		if (giUCL > giHigh)
			giHigh = giUCL;
		if (giLCL < giLow)
			giLow = giLCL;
			
		giOffsetY = parseInt((objCanvas.height - giOffset)/ (giHigh - giLow)) ;
		giOffsetX = parseInt(objCanvas.width /iPoints);
		
		objCtx.fillStyle = "#FFFFFF";		
		objCtx.fillRect(0,0,objCanvas.width,objCanvas.height);		
		vfDrawLine(0,giUCL,giUCL);
		vfDrawLine(0,giCL ,giCL);
		vfDrawLine(0,giLCL,giLCL);		
};

function ifPointX(iX){
	return (iX  * giOffsetX);
}

function ifPointY(iY){
	return  ((giHigh - iY)  * giOffsetY) + (giOffset/2);
}

function vfDrawPoint(iX,iY,iMove,sColor){
			
		if(iMove == 1)
			objCtx.moveTo(ifPointX(iX),ifPointY(iY));
		else
			objCtx.lineTo(ifPointX(iX),ifPointY(iY));
			
		objCtx.stroke();	
		objCtx.beginPath();
		objCtx.fillStyle = sColor;
		objCtx.arc(ifPointX(iX),  ifPointY(iY),10,0,2*Math.PI);
		objCtx.fill();	

};


