NAME=mqttListener
# gopher
MQTT_HOST?=192.168.8.180
TOPIC=zigbee2mqtt/aquara-door-01
PAYLOAD_OPEN={"contact":false}
PAYLOAD_CLOSE={"contact":true}
PUB=mosquitto_pub -u $(MQTT_USER) -P $(MQTT_PASS) -h $(MQTT_HOST) -t '$(TOPIC)' -m

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

pub/open:
	$(PUB) '$(PAYLOAD_OPEN)'

pub/close:
	$(PUB) '$(PAYLOAD_CLOSE)'
