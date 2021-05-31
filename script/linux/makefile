install:
	cp ddns-go /usr/bin/ddns-go
	mkdir -p /etc/ddns-go
	if [ ! -f /etc/ddns-go/config.yaml ]; then cp ./config.yaml.default /etc/ddns-go/config.yaml; fi;
	cp ./ddns-go.service /etc/systemd/system/ddns-go.service
.PHONY: install