
BINRARY=apns_svr apnstest getp12file feedback_svr

all:$(BINRARY)

apns_svr: apnssvr.go anpsproc.go
	go build -o $@ $^


getp12file: getp12file.go
	go build -o $@ $^

apnstest: apnstest.go
	go build -o $@ $^

feedback_svr: feedback.go
	go build -o $@ $^

clean:
	rm $(BINRARY)

in:
	scp  $(BINRARY)  im@192.168.251.55:/home/im/bin/apns/
