ipstr=$(hostname -I)
ip4="$(echo "${ipstr}" | sed -e 's/[[:space:]]*$//')"
echo "Ypur host is :" $ip4
echo "***Bootnode***"
./bootnode -nodekey=key1.bin  -addr=$ip4:51901

