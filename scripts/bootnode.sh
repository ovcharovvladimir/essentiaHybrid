ipstr=$(hostname -I)
ip4="$(echo "${ipstr}" | sed -e 's/[[:space:]]*$//')"
echo "Your host is :" $ip4
echo "***Bootnode***"
chmod +x bootnode
./bootnode -nodekey=key1.bin  -addr=$ip4:51901

