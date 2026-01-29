#!/bin/env fish

# Assumes the script is running from /tmp/deployment
echo Installing
sudo rm -r /opt/curious-ape 2>/dev/null
sudo mkdir /opt/curious-ape
sudo mv /tmp/deployment/* /opt/curious-ape
sudo rm -r /tmp/deployment 2>/dev/null
sudo chown -R daniel:daniel /opt/curious-ape
echo litestream version: (/opt/curious-ape/litestream version)
sudo mkdir -p /usr/local/lib/systemd/system
sudo cp /opt/curious-ape/litestream.service /usr/local/lib/systemd/system/ape-litestream.service
sudo systemctl enable --now ape-litestream.service
sudo cp /opt/curious-ape/curious-ape.service /usr/local/lib/systemd/system/curious-ape.service
sudo systemctl enable --now curious-ape.service
echo Done
