ipstr=$(hostname -I)
ip4="$(echo "${ipstr}" | sed -e 's/[[:space:]]*$//')"
boot='essnode://e1d615ef7c55f0a43049ad1bacd3646148b13939d61042e580eb5891c179ccb0b291e287b73e250c59c0817c984e49cb09cfd62bee7b20c50251e546a7aa8a5c@'"$ip4"':51901'
echo $boot
echo "Your host is :" $ip4
echo "*** Node A ***"
./gess --testnet --networkid 666  --port 51906  --bootnodes $boot --datadir="/home/essdeveloper/.essentia2" --verbosity 5 console
