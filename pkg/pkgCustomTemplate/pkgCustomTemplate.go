// change this package to provide customization ie: stripe

package pkgCustomTemplate

// anything that needs to be set before initialization
func FnPreInit()error{	
	return nil
}

// anything that needs to be set after initialization. Mode is local, remote, live
func FnAfterInit(iMode int)error{		
	return nil
}

// anything that needs to be set at end of initialization. Mode is local, remote, live
func FnPostInit(iMode int)error{	
	return nil
}

// If needed expand arFlags, arReplace
func ConfigArs(sUrl string)([]string, []string){	
	arFlags := []string{"//UrlStringFlag"}
	arReplace := []string{`var gs_post = "` + sUrl + `";`}

	return arFlags,arReplace
}
