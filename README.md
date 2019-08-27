# DDNS Tool for Gehirn DNS [![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/hyperium/hyper/master/LICENSE)
Created with reference to [Gehirn DNS API Documentation](https://support.gehirn.jp/apidocs/gis/dns/index.html).   

## Preparation
* Create target A record on WEBUI.  
```
(Example)
www.example.com    360    A    192.168.254.1
```
* Check the zone_id information for the zone to which you just added a record.  
* Create API Token, and check Authorization header's value (Basic ********************).
* Create a yaml file anywhere you like.
```
(Example)
Authorization: "Basic abcdefghijklmnopqrstuvwxyz"
ZoneID: "12345678-9999-abcd-efgh-ijklmnopqrst"
DomainName: "www.example.com"
```
## Run
* Get & build
```
$ git https://github.com/nozomi1773/gehirn-ddns-tool.git
$ cd gehirn-ddns-tool/cmd && go build . && mv cmd /root/tools/ddns
```
* Run Test this command
```
$ ./ddns -f ~/tools/ddns-config.yaml
```
* Set Cron
```
$ crontab -e
0-59 * * * * /root/tools/ddns -f "/root/tools/config.yaml"
```
