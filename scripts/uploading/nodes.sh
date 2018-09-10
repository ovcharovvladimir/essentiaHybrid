#!/bin/bash
# Nodes
RED='\033[1;31m'
NC='\033[0m' # No Color
bin="./bin"

options=("18.188.202.224" "18.222.125.29" "18.216.229.30" "18.221.195.24"  "Upload" "Quit")
echo "Availible bootnodes:"
print_list(){
for ((i = 0; i < ${#options[@]}; ++i)); do
    # bash arrays are 0-indexed
    position=$(( $i + 1 ))
    echo "$position) ${options[$i]}"
done
}    
upload(){

for ((i = 0; i < ${#options[@]}-2; ++i)); do
    echo -e "* ${RED} ${options[$i]} ${NC} *"
     ssh -i block.pem -o ConnectTimeout=5  ubuntu@${options[$i]} sudo mkdir -p /home/release
     echo ". Files Uploading"
     rsync -ave 'ssh -i block.pem' --info=progress2  --timeout=5 --rsync-path="sudo rsync" $bin/gess  ubuntu@${options[$i]}:/home/release/gess
    echo "- DONE -"
done
}
echo "**********************************************************"
PS3='Select miners:' 
select opt in "${options[@]}"
do
    case $opt in
        "18.191.153.144")
	    clear
	echo "**** ${opt} ****"
            ssh -i block.pem ubuntu@$opt
	    print_list
	    ;;
        "18.223.214.204")
	    clear
	echo "**** ${opt} ****"
            ssh -i block.pem ubuntu@$opt
           print_list
            ;;
        "18.188.133.172")
	    clear
	echo "**** ${opt} ****"
            clear
            ssh -i block.pem ubuntu@$opt
            print_list
            ;; 
	"Upload")
		upload
		print_list
	    ;;	            
        "Quit")
            break
            ;;
        *) echo "invalid option $REPLY";;
    esac
done
