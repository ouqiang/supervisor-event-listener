project_name:=supervisor-eventlistener
project_version:=1.2.3
root_dir := $(abspath $(CURDIR))
build_dir := $(root_dir)/build
GOPATH := ${HOME}/go


clean:
	rm -fr $(build_dir)


build-bydocker:
	sudo docker run -it --rm \
		-v $(GOPATH)/:/root/go/ \
		-v $(root_dir)/:/$(project_name) \
		-w /$(project_name)/ \
		golang:1.17.0-buster \
		make build


build:
	GO111MODULE=on go build -o $(project_name) ./main.go


release: clean build-bydocker
	mkdir -p                       $(build_dir)/$(project_name)/
	mv ./supervisor-eventlistener  $(build_dir)/$(project_name)/
	cp ./conf/config.toml          $(build_dir)/$(project_name)/
	cd $(build_dir) && tar -zcvf $(project_name)-$(project_version).tar.gz $(project_name)
	@echo ...created $(build_dir)/$(project_name)-$(project_version).tar.gz
	@echo ...done.


log:
	tmux new-session -d -s dev
	tmux split-window -t "dev:0"
	tmux split-window -t "dev:0.0" -h
	tmux split-window -t "dev:0.2" -h
	tmux send-keys -t "dev:0.0" "bash -c 'tail -f /tmp/supervisor-eventlistener.log'" Enter
	tmux send-keys -t "dev:0.1" "bash -c 'sudo supervisorctl tail -f $(project_name)'" Enter
	tmux set-option -g mouse on
	tmux attach -t dev
	tmux kill-session -t dev


test-integration:
	go build 
	sudo supervisorctl stop $(project_name)
	sudo cp ./$(project_name) /usr/local/bin/
	sudo cp ./tests/supervisor-app.ini /etc/supervisor.d/
	sudo supervisorctl remove $(project_name)
	sudo supervisorctl update $(project_name)
	sudo supervisorctl tail -f $(project_name) stderr


.PHONY: clean build build-bydocker release log test-integration
