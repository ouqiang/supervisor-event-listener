
install:
	sudo supervisorctl stop alert
	go build
	sudo cp ./supervisor-event-listener /usr/local/bin/
	sudo supervisorctl start alert

test:
	sudo supervisorctl start test

