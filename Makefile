project_name:=supervisor-event-listener
project_version:=1.2.2
root_dir := $(abspath $(CURDIR))
build_dir := $(root_dir)/build
GOPATH := ${HOME}/go


.PHONY: clean
clean:
	rm -fr $(build_dir)

.PHONY: build-bydocker
build-bydocker:
	sudo docker run -it --rm \
		-v $(GOPATH)/:/root/go/ \
		-v $(root_dir)/:/$(project_name) \
		-w /$(project_name)/ \
		golang:1.16.2-buster \
		make build


.PHONY: build
build:
	GO111MODULE=on go build -o $(project_name) ./$(project_name).go


.PHONY: release
release: clean build-bydocker
	mkdir -p                            $(build_dir)/$(project_name)/
	mv ./supervisor-event-listener      $(build_dir)/$(project_name)/
	cp ./supervisor-event-listener.toml $(build_dir)/$(project_name)/
	cd $(build_dir) && tar -zcvf $(project_name)-$(project_version).tar.gz $(project_name)
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


test-integration:
	go build 
	sudo supervisorctl stop supervisor-event-listener
	sudo cp ./supervisor-event-listener /usr/local/bin/
	sudo cp ./tests/supervisor-app.ini /etc/supervisor.d/
	sudo supervisorctl remove supervisor-event-listener
	sudo supervisorctl update supervisor-event-listener
	sudo supervisorctl tail -f supervisor-event-listener stderr

