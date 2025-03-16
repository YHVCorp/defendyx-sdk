package helpers

import "fmt"

func PrintBanner() {
	banner := `

oooooooooo.              .o88o.                             .o8              ooooooo  ooooo 
 888     Y8b             888  "                            "888                8888    d8'  
 888      888  .ooooo.  o888oo   .ooooo.  ooo. .oo.    .oooo888  oooo    ooo    Y888..8P    
 888      888 d88'  88b  888    d88'  88b  888P"Y88b  d88'  888    88.  .8'       8888'     
 888      888 888ooo888  888    888ooo888  888   888  888   888     88..8'      .8PY888.    
 888     d88' 888    .o  888    888    .o  888   888  888   888      888'      d8'   888b   
o888bood8P'    Y8bod8P' o888o    Y8bod8P' o888o o888o  Y8bod88P"     .8'     o888o  o88888o 
                                                                 .o..P'                     
                                                                  Y8P'       
																                 
`
	fmt.Println(banner)
}
