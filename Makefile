LDFLAGS:=-w
BUILD_DIR:=./build/
PROJECT_NAME:=supervisor-event-listener
VERSION:=1.1.1

install:
	cp ./supervisor-event-listener /usr/local/bin/


test-integration:
	sudo supervisorctl stop supervisor-event-listener
	go build 
	sudo cp ./supervisor-event-listener /usr/local/bin/
	sudo cp ./tests/supervisor-app.ini /etc/supervisor.d/
	sudo supervisorctl start  supervisor-event-listener
	sudo supervisorctl update
	sudo supervisorctl start sleep-then-exit


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

