#!/bin/bash
# Miners
RED='\033[1;31m'
GREEN='\033[01;32m'
NC='\033[0m' # No Color
bin="./bin"

options=(  "18.224.121.61" "18.224.159.84" "18.224.168.178" "18.224.198.158" "Upload" "Start All" "Stop All" "Remove chain data" "Quit")
echo "Availible miners"
print_list(){
for ((i = 0; i < ${#options[@]}; ++i)); do
    # bash arrays are 0-indexed
    position=$(( $i + 1 ))
    echo "$position) ${options[$i]}"
done
}       
upload(){

for ((i = 0; i < ${#options[@]}-5; ++i)); do
    echo -e "* ${RED} ${options[$i]} ${NC} *"
     ssh -i block.pem -o ConnectTimeout=5  ubuntu@${options[$i]} sudo mkdir -p /home/release
     ssh -i block.pem -o ConnectTimeout=5  ubuntu@${options[$i]} sudo rm -rv /home/release/*  
     echo "* UPLOADING"
     rsync -ave 'ssh -i block.pem' --info=progress2  --timeout=5 --rsync-path="sudo rsync" $bin/gess  ubuntu@${options[$i]}:/home/release/gess
     rsync -ave 'ssh -i block.pem' --info=progress2  --timeout=5 --rsync-path="sudo rsync" $bin/pass.txt  ubuntu@${options[$i]}:/home/release/pass.txt
    echo -e "* ${GREEN} DONE ${NC} *"
done
}
start_process(){

for ((i = 0; i < ${#options[@]}-5; ++i)); do
   echo -e "* Starting node  ${RED} ${options[$i]} ${NC} *"
    ssh -i block.pem -o ConnectTimeout=5  ubuntu@${options[$i]} sudo pkill gess
    ipstr=$( ssh -i block.pem -o ConnectTimeout=5  ubuntu@${options[$i]} hostname -I)
    ip4="$(echo "${ipstr}" | sed -e 's/[[:space:]]*$//')"

 #  ssh -i block.pem -o ConnectTimeout=5 -Y ubuntu@${options[$i]} sudo rm -rv 	/home/ubuntu/.essentia
  # ssh -i block.pem -o ConnectTimeout=5 -Y ubuntu@${options[$i]} sudo rm -rv 	/home/ubuntu/.esshash
   echo  --password /home/release/pass.txt account new
   ssh -i block.pem -o ConnectTimeout=5 ubuntu@${options[$i]} sudo /home/release/gess  --testnet --password /home/release/pass.txt account new
   echo sudo  /home/release/gess   --rpc  --mine --minerthreads=1 -cache=2048 --gcmode=archive --etherbase '0xc129cfc9844fc110b5639bb7357b788c516251c9' --testnet  --rpc --rpcaddr $ip4 --nat extip:${options[$i]}
   ssh -i block.pem -o ConnectTimeout=5 ubuntu@${options[$i]} nohup sudo  /home/release/gess   --rpc  --mine --minerthreads=1 -cache=2048 --gcmode=archive --etherbase '0xc129cfc9844fc110b5639bb7357b788c516251c9' --testnet  --rpc --rpcaddr $ip4 --nat extip:${options[$i]}  > gess.out 2>&1 &
    echo -e "* ${GREEN} DONE ${NC} *"
 done  
}
remove_chain(){

for ((i = 0; i < ${#options[@]}-5; ++i)); do
   echo -e "*  NODE  ${RED} ${options[$i]} ${NC} *"
    ssh -i block.pem -o ConnectTimeout=5  ubuntu@${options[$i]} sudo pkill gess
    ipstr=$( ssh -i block.pem -o ConnectTimeout=5  ubuntu@${options[$i]} hostname -I)
    ip4="$(echo "${ipstr}" | sed -e 's/[[:space:]]*$//')"
   echo "Press Ctrl+C to continue loading ..."
   echo  ssh -i block.pem -o ConnectTimeout=5 -Y ubuntu@${options[$i]} sudo rm -rv 	/home/ubuntu/.essentia
   ssh -i block.pem -o ConnectTimeout=5 -Y ubuntu@${options[$i]} sudo rm -rv 	/home/ubuntu/.essentia
   echo  ssh -i block.pem -o ConnectTimeout=5 -Y ubuntu@${options[$i]} sudo rm -rv 	/home/ubuntu/.esshash
   ssh -i block.pem -o ConnectTimeout=5 -Y ubuntu@${options[$i]} sudo rm -rv 	/home/ubuntu/.esshash
    echo -e "* ${GREEN} DONE ${NC} *"
 done  
}
stop_process(){

for ((i = 0; i < ${#options[@]}-5; ++i)); do
   echo -e "* Starting node  ${RED} ${options[$i]} ${NC} *"
    ssh -i block.pem -o ConnectTimeout=5  ubuntu@${options[$i]} sudo pkill gess
    ipstr=$( ssh -i block.pem -o ConnectTimeout=5  ubuntu@${options[$i]} hostname -I)
    ip4="$(echo "${ipstr}" | sed -e 's/[[:space:]]*$//')"
   echo "Press Ctrl+C to continue loading ..."
   echo  ssh -i block.pem -o ConnectTimeout=5  ubuntu@${options[$i]} sudo pkill gess 
   ssh -i block.pem -o ConnectTimeout=5  ubuntu@${options[$i]} sudo pkill gess 
    echo -e "* ${GREEN} DONE ${NC} *"
 done  
}
echo "**********************************************************"
PS3='Select:' 
select opt in "${options[@]}"
do
    case $opt in
        "18.224.121.61")
	    clear
	echo "**** ${opt} ****"
            ssh -i block.pem ubuntu@$opt
           print_list
            ;;
        "18.224.159.84")
	    clear
	echo "**** ${opt} ****"
            clear
            ssh -i block.pem ubuntu@$opt
            print_list
            ;; 
            "18.224.168.178")
	    clear
	echo "**** ${opt} ****"
            clear
            ssh -i block.pem ubuntu@$opt
            print_list
            ;; 
          "18.224.198.158")
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
	"Stop All")
		stop_process
		print_list
	    ;;	
    "Start All")
        start_process
		print_list
	    ;;	
	"Remove chain data")
        remove_chain
		print_list
	;;
        "Quit")
            break
            ;;
        *) echo "invalid option $REPLY";;
    esac
done
