
#!/bin/bash
# Bootnodes
RED='\033[1;31m'
NC='\033[0m' # No Color
bin="./bin"

options=("18.217.9.146" "18.191.42.39" "18.223.101.201" "Upload" "Quit")
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
     idx=$((i+1))
      echo "******** bootnode "
     rsync -ave 'ssh -i block.pem' --info=progress2 --timeout=5 --rsync-path="sudo rsync" $bin/bootnode  ubuntu@${options[$i]}:/home/release/bootnode
     echo "------------------"
     echo "***** key${idx}.bin"
     rsync -ave 'ssh -i block.pem' --info=progress2 --timeout=5 --rsync-path="sudo rsync" $bin/key$idx.bin  ubuntu@${options[$i]}:/home/release/key.bin
     echo "------------------"
     
    echo "- DONE -"
done
}
echo "**********************************************************"
PS3='Select bootnodes:' 
select opt in "${options[@]}"
do
    case $opt in
        "18.217.9.146")
	    clear
	echo "**** ${opt} ****"
            ssh -o ConnectTimeout=3 -i block.pem ubuntu@$opt
	    print_list
	    ;;
        "18.191.42.39")
	    clear
	echo "**** ${opt} ****"
            ssh -o ConnectTimeout=3 -i block.pem ubuntu@$opt
           print_list
            ;;
        "18.223.101.201")
	    clear
	echo "**** ${opt} ****"
            clear
            ssh -o ConnectTimeout=3 -i block.pem ubuntu@$opt
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
