package pkg

import (
	"fmt"
	"os"
)

var Clicki = `                                                                         
        CCCCCCCCCCCCClllllll   iiii                      kkkkkkkk             iiii  
     CCC::::::::::::Cl:::::l  i::::i                     k::::::k            i::::i 
   CC:::::::::::::::Cl:::::l   iiii                      k::::::k             iiii  
  C:::::CCCCCCCC::::Cl:::::l                             k::::::k                   
 C:::::C       CCCCCC l::::l iiiiiii     cccccccccccccccc k:::::k    kkkkkkkiiiiiii 
C:::::C               l::::l i:::::i   cc:::::::::::::::c k:::::k   k:::::k i:::::i 
C:::::C               l::::l  i::::i  c:::::::::::::::::c k:::::k  k:::::k   i::::i 
C:::::C               l::::l  i::::i c:::::::cccccc:::::c k:::::k k:::::k    i::::i 
C:::::C               l::::l  i::::i c::::::c     ccccccc k::::::k:::::k     i::::i 
C:::::C               l::::l  i::::i c:::::c              k:::::::::::k      i::::i 
C:::::C               l::::l  i::::i c:::::c              k:::::::::::k      i::::i 
 C:::::C       CCCCCC l::::l  i::::i c::::::c     ccccccc k::::::k:::::k     i::::i 
  C:::::CCCCCCCC::::Cl::::::li::::::ic:::::::cccccc:::::ck::::::k k:::::k   i::::::i
   CC:::::::::::::::Cl::::::li::::::i c:::::::::::::::::ck::::::k  k:::::k  i::::::i
     CCC::::::::::::Cl::::::li::::::i  cc:::::::::::::::ck::::::k   k:::::k i::::::i
        CCCCCCCCCCCCClllllllliiiiiiii    cccccccccccccccckkkkkkkk    kkkkkkkiiiiiiii
`

func PrintInfo(version, date, commit string) {
	if version == "" {
		version = "N/A"
	}
	if date == "" {
		date = "N/A"
	}
	if commit == "" {
		commit = "N/A"
	}
	fmt.Println(Clicki)
	fmt.Fprintf(os.Stdout, "Build version: %s\n", version)
	fmt.Fprintf(os.Stdout, "Build date:    %s\n", date)
	fmt.Fprintf(os.Stdout, "Build commit:  %s\n", commit)
}
