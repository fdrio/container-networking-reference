## Network namespaces

#### Network namespaces live in /var/run/netns/

Simple example that creates two netns that communicates with each other
1. Create two netns: 
	ip netns add nyc
	ip netns add sfo

2. Create a link of type veth pair (Imagine its a virtual cable)

	ip link add veth-sfo type veth peer name veth-nyc
	ip link list | grep veth # this lists the veth pairs

3. Set the two ends of the veth pair to each namespace
	
	ip link set veth-sfo netns sfo
	ip link set veth-nyc netns nyc
	ip link list | grep veth # this lists the veth pairs

4. Check that the end exists in either namespace
	
	ip netns exec sfo ip link
	ip netns exec nyo ip link

5. Add ip addresses to the interfaces and bring them up

	ip netns exec sfo ip address add 10.0.0.11/24 dev veth-sfo
	ip netns exec sfo ip link set veth-sfo up
	ip netns exec nyc ip address add 10.0.0.12/24 dev veth-nyc
	ip netns exec nyc ip link set veth-nyc up

6. Ping the interfaces to check either end is reachable
	ip netns exec sfo ping 10.0.0.12
	ip netns exec nyc ping 10.0.0.11


