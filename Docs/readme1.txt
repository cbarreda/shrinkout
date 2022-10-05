PLU 3064 for spc chart.
t_produce.f_uid 3064
select  * from t_storeproduce where f_produceid = '3064';
 t_storeproduce.f_each 12 for vendor 1 (default)
 t_storeproduce.f_each 22 for vendor 2 (Sunny Farms)
 t_storeproduce.f_each 12 for vendor 3 (Green Pastures)

SPC in t_spc and t_produce.f_spcid
SP_ADDBOXES stores qty entered
	T_STOCKEDBOXES.f_tunits = (t_storeproduce.f_each * p_boxes) + p_units;
 
 
Id 1: password DREW for set device.
for  3064 ucl = 8, lcl = 1 (See t_storeproduce.f_each above to insure scrap is around 1-8%)

	

Entered: Default Vendor (5  boxes 10 units) 060 stocked   3 scrapped  5%
	     Sunny Farms    (10 boxes 20 Units) 240 Stocked  24 scrapped  10%
	     Green Pastures (5  boxes 10 units) 060 stocked   1 scrapped 

Set phone so app knows what store it is in.
1.- Go to https://shrinkout.com/device
	a.- Admin ID: 1
	b.- password: DREW
	c.- Store ID: 1
	d.- Dept ID: 1
	e.- MAC: (Eventually it will be a real mac but for now anything works)
	f.- Desc: (Just type anything)
	g.- Click on Setup.
	h.- Should return:OK 0
2.- Go https://shrinkout.com/client
	a) Using the cyanide buttons type 115 and then sign in
	a1) On Android: click on the 3 dots on the upper right and then click on Add to home screen to create a shorcut on your android phone.)
	a2) On Iphone: clickon box with uparrow at the bottom of the Iphone and then click on Add to home screen to create a shorcut on your android phone.)
3.- Using the app:
	After signing in: Use Toggle button on the far right to Toggle between number key and alpha key.
	- Use number key to enter the Price Lookup Code (PLU) or the alpha key to search product by Name
	- Once a product is selected click Go.
	- Default allows you to select a farmer.
	- Code allows you to choose a reason code to shrink.
	- Use + or - to set a quantity and then choose shrinkOut, boxes or units.
	- Click done when done.
Bug: For some reason loging out is unregistering the device so have to do step 1.- (I'll fix the bug)	


	
