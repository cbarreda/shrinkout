-- setup start
change config file. key, database, etc.
change boilerplate.go to newname.go ie boilerplate.go to editmywords.go
change gPkgID = "boilerplate." to gPkgID = "newname."
change pkgSession Vars: gKey and gSessionName
set dirs to be made in pkgUser.fnMakeDirs

./bpsetup boilerplate newname  changes:
	- all imports from /boilerplate to /newname
	- ./masterTemplates/boilerplate folder to ./masterTemplates/newName
	- ./masterTemplates/boilerplate/boilerplate.config to ./masterTemplates/newName/newname.config
	- ./templates/boilerplate to ./templates/newName
rm go.mod and then go mod init


-- makeTemplate Start
	-- masterTemplate folders define global tags in the config file in masterTemplates
		defined tags start: Change ~!dcbootstrapjs to change version of Bootstrap/Jquery
			~!dcbootstrapjs@# 
			:~!host@# sep:
			~!dcbootstrapCss@#
			~!dcglobalCss@#
			~!dcglobalJs@#
			~!dcgoogleAnalytics@# defined as <!-- Google analytics --> change if using Google Analytics.
			~!dctopload@#  defined as <!-- top load files here --> (Use this one for things that must load first.)
		undefined tags start. Must be defined in local templates. Use <!-- --> if they don't apply to local template
			~!dclocalJs@#		
			~!dclocalCSS@#
			~!dctitle@#
			~!dcLocalPageHeader@# (Page headers)
			~!dcglobalbodystart@# <body but could be <body class="body-wrapper"> or something like that.
-- local template folders define tags for local pages.

makeTemplate cleanconfig removes unused local tags, adds to the local config file undefined tags.

makeTemplate makeall creates the template using the config files to replace tags with content.
	it also modifies global.js based on the go.config file

use sf_errorCode and sfTranslate(sKey) in localjs.tmpl (localjs.tmp will be incorporated into the page by makeTemplate makeTemplate


Finally: Update/change
----------------
./postgres/create
./postgres/functions
./postgres/deploy

build newname.go
