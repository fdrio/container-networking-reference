echo "Turning down bridge"
ip link set test down
echo "Deleting bridge"
brctl delbr test
