#
# /Makefile
#

test:
	go test ...

pcap:
	sudo tcpdump -ni lo -s 0 -w gprc_test_2024_11_08.pcap -v port 50052