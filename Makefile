NAME=mqttListener
MQTT_HOST=192.168.8.180
TOPIC=shellies/shellydw2-1/sensor/state

.PHONY: build
build:
	cd cmd/$(NAME) ;\
	env GOARCH=arm GOOS=linux go build -o ../../$(NAME) 

.PHONY: deploy
deploy: build
	ansible-playbook -i ./inventory.ini main.yml --tags=sounds
	ansible-playbook -i ./inventory.ini main.yml
	rm -f $(NAME)

.PHONY: deploy
deploy-no-sounds: build
	ansible-playbook -i ./inventory.ini main.yml
	rm -f $(NAME)

pub:
	mosquitto_pub -u $(MQTT_USER) -P $(MQTT_PASS) -t '$(TOPIC)' -m 'open' -h $(MQTT_HOST)
