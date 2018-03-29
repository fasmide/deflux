# deflux
deflux connects to deCONZ rest api and listens for sensor updates and write these to InfluxDB.

deCONZ supports a variaty of Zigbee sensors but have no historical data about their values - with deflux you'll be able to store all these measurements in influxdb where they can be queried from the command line or graphical tools such as grafana. 

## Usage

Start off by `go get`'ting deflux:

```
go get github.com/fasmide/deflux
```

deflux tries to read `$(pwd)/deflux.yml` or `/etc/deflux.yml` in that order, if both fails it will try to discover deCONZ with their webservice and output a configuration sample to stdout. 

Hint: if you've temporarily unlocked the deconz gateway, it should be able to fill in the api key by it self, this needs some testing though...

First run generates a sample configuration:

```
$ deflux
2018/03/29 13:51:02 no configuration could be found: could not read configuration: 
open /home/fas/go/src/github.com/fasmide/deflux/deflux.yml: no such file or directory
open /etc/deflux.yml: no such file or directory
2018/03/29 13:51:03 unable to pair with deconz: unable to pair with deconz: link button not pressed, please fill out APIKey manually
2018/03/29 13:51:03 Outputting default configuration, save this to /etc/deflux.yml
deconz:
  addr: http://192.168.1.90:8080/api
  apikey: ""
influxdb:
  addr: http://127.0.0.1:8086/
  username: change me
  password: change me
  useragent: Deflux
influxdbdatabase: deconz
```

Save the sample configuration and edit it to your needs, then run again

```
$ deflux 
2018/03/29 13:52:06 Using configuration /home/fas/go/src/github.com/fasmide/deflux/deflux.yml
2018/03/29 13:52:06 Connected to deCONZ at http://192.168.1.90:8080/api
2018/03/29 13:57:06 recv: {"e":"changed","id":"7","r":"sensors","state":{"buttonevent":1004,"lastupdated":"2018-03-29T11:57:06"},"t":"event"}
2018/03/29 13:57:06 SensorStore updated, found 17 sensors
2018/03/29 13:57:07 Saved 1 records to influxdb
2018/03/29 13:57:12 recv: {"e":"changed","id":"7","r":"sensors","state":{"buttonevent":1005,"lastupdated":"2018-03-29T11:57:12"},"t":"event"}
2018/03/29 13:57:13 Saved 1 records to influxdb
2018/03/29 13:58:23 recv: {"config":{"battery":100,"on":true,"reachable":true,"temperature":2000},"e":"changed","id":"6","r":"sensors","t":"event"}
2018/03/29 13:58:23 not adding event to influx batch: this event (*event.EmptyState:lumi.sensor_wleak.aq1) has no time series data
2018/03/29 14:00:39 recv: {"e":"changed","id":"16","r":"sensors","state":{"lastupdated":"2018-03-29T12:00:39","temperature":2238},"t":"event"}
2018/03/29 14:00:39 recv: {"e":"changed","id":"17","r":"sensors","state":{"humidity":2598,"lastupdated":"2018-03-29T12:00:39"},"t":"event"}
2018/03/29 14:00:40 Saved 2 records to influxdb
2018/03/29 14:03:46 recv: {"e":"changed","id":"1","r":"sensors","state":{"lastupdated":"2018-03-29T12:03:46","temperature":2232},"t":"event"}
2018/03/29 14:03:46 recv: {"e":"changed","id":"2","r":"sensors","state":{"humidity":2615,"lastupdated":"2018-03-29T12:03:46"},"t":"event"}
2018/03/29 14:03:46 recv: {"e":"changed","id":"3","r":"sensors","state":{"lastupdated":"2018-03-29T12:03:46","pressure":1004},"t":"event"}
2018/03/29 14:03:47 Saved 3 records to influxdb
```

It does have some rough edges that i'll hopefully be working on - now you should be able to find these sensor measurements in influxdb

## Influxdb

Sensor values are added as influxdb values and tagged with sensor type, id and name

```
> select * from deflux order by time desc limit 25;
name: deflux
time                buttonevent daylight humidity name               pressure status temperature type
----                ----------- -------- -------- ----               -------- ------ ----------- ----
1522325026000000000                               Terrasse                           22.32       ZHATemperature
1522325026000000000                      26.15    Terrasse                                       ZHAHumidity
1522325026000000000                               Terrasse           1004                        ZHAPressure
1522324839000000000                               Kælder lab                         22.38       ZHATemperature
1522324839000000000                      25.98    Kælder lab                                     ZHAHumidity
1522324632000000000 1005                          lumi.sensor_switch                             ZHASwitch
1522324626000000000 1004                          lumi.sensor_switch                             ZHASwitch
1522320725000000000                               Kælder lab                         21.88       ZHATemperature
1522320725000000000                      26.36    Kælder lab                                     ZHAHumidity
1522320092000000000                      26.28    Terrasse                                       ZHAHumidity
1522320092000000000                               Terrasse                           22.1        ZHATemperature
1522320092000000000                               Terrasse           1004                        ZHAPressure
1522315309000000000                      31.43    Kælder gang                                    ZHAHumidity
1522315309000000000                               Kælder gang                        20.79       ZHATemperature
1522315299000000000                      39.91    Kælder gang                                    ZHAHumidity
1522315299000000000                               Kælder gang                        20.76       ZHATemperature
``` 

## Grafana

TODO: As soon as i have a few weeks of sensor data i'll put some graph examples and a getting started dashboard

## Notes
I'm in possession of Temperature, Humidity, Pressure, Water flood, Fire alarm and a few buttons - all Xiaomi branded and as such dont know if all other sensors will just work, there is properly a lot deconz event i don't account for, these should be easily added though.