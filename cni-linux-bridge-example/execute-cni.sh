#Creating the netns
NETNS_NAME=ns1
ip netns add $NETNS_NAME
echo "Create network namespace:" $NETNS_NAME

# Build bin
go build -o example

#Add CNI args
sudo CNI_COMMAND=ADD CNI_CONTAINERID=ns1 \
CNI_NETNS=/var/run/netns/ns1 CNI_IFNAME=eth10 \
CNI_PATH=`pwd` \
./example < config


ip netns exec $NETNS_NAME ifconfig -a
# Print links and grep for its name
ip addr | grep test
