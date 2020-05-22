LDFLAGS:=-w
BUILD_DIR:=./build/
PROJECT_NAME:=supervisor-event-listener
VERSION:=1.1.1

install:
	cp ./supervisor-event-listener /usr/local/bin/


test-integration:
	go build 
	sudo supervisorctl stop supervisor-event-listener
	sudo cp ./supervisor-event-listener /usr/local/bin/
	sudo cp ./tests/supervisor-app.ini /etc/supervisor.d/
	sudo supervisorctl remove supervisor-event-listener
	sudo supervisorctl update supervisor-event-listener
	sudo supervisorctl tail -f supervisor-event-listener stderr


clean:
	rm -fr $(BUILD_DIR)


release: 
	GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS)
	rm -fr                             $(BUILD_DIR)/$(PROJECT_NAME)/
	mkdir -p                           $(BUILD_DIR)/$(PROJECT_NAME)/
	mv ./supervisor-event-listener     $(BUILD_DIR)/$(PROJECT_NAME)/
	cp ./supervisor-event-listener.ini $(BUILD_DIR)/$(PROJECT_NAME)/
	cd $(BUILD_DIR) && tar -zcvf $(PROJECT_NAME)-$(VERSION).tar.gz $(PROJECT_NAME)
	@echo ...done.


log:
	tmux new-session -d -s dev
	tmux split-window -t "dev:0"
	tmux split-window -t "dev:0.0" -h
	tmux split-window -t "dev:0.2" -h
	tmux send-keys -t "dev:0.0" "bash -c 'tail -f /tmp/supervisor-event-listener.log'" Enter
	tmux send-keys -t "dev:0.1" "bash -c 'sudo supervisorctl tail -f supervisor-event-listener'" Enter
	tmux set-option -g mouse on
	tmux attach -t dev
	tmux kill-session -t dev

