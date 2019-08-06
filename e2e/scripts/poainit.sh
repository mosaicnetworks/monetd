#!/bin/bash

monetcli accounts import --file ~/.giverny/networks/crowdfundnet/keystore/Amelia.json --default

monetcli poa init -h 172.77.5.10 --pwd ~/.giverny/networks/crowdfundnet/keystore/Amelia.txt
