## MQTT sound player

This repo contains a golang tool that listens to MQTT topics and plays a sound. I use it 
at home to detect when my door outside opens.

Besides the code for the tool (`/cmd/mqttListener`), it also contains some ansible bits to 
deploy the tool as a service in raspberry pi (tested on version 3).

## Usage

1. Modify the config.yml to your liking.
2. Modify the inventory.ini to point to your server (or servers).
3. Run `make deploy`

## References

This is based on the work of [kkentzo](https://github.com/kkentzo/deployment-ansible-systemd-demo). Kudos to him.
