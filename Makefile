NAME=mqttListener

.PHONY: build
build:
	cd cmd/$(NAME) ;\
	env GOOS=linux go build -o ../../$(NAME) 

.PHONY: deploy
deploy: build
	ansible-playbook -i ./inventory.ini main.yml
