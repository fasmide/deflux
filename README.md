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

Sensor values are added as influxdb values and tagged with sensor type, id and name.
Different event types are stored in different measurements, meaning you will end up with multiple influxdb measurements:
```
> show measurements;
name: measurements
name
----
deflux_Daylight
deflux_ZHAHumidity
deflux_ZHAPressure
deflux_ZHATemperature

```
Example from deflux_ZHAHumidity
```
> select * from deflux_ZHAHumidity;
name: deflux_ZHAHumidity
time                humidity id name       type
----                -------- -- ----       ----
1523555448000000000 38.92    13 Kælder bad ZHAHumidity
1523556151000000000 40.88    13 Kælder bad ZHAHumidity
1523556658000000000 39.44    13 Kælder bad ZHAHumidity
1523557231000000000 55.48    2  Terrasse   ZHAHumidity
1523557476000000000 38.86    13 Kælder bad ZHAHumidity
1523557846000000000 56       2  Terrasse   ZHAHumidity
1523558273000000000 37.74    13 Kælder bad ZHAHumidity
``` 

## Grafana

TODO: As soon as i have a few weeks of sensor data i'll put some graph examples and a getting started dashboard

## Notes
I'm in possession of Temperature, Humidity, Pressure, Water flood, Fire alarm and a few buttons - all Xiaomi branded and as such dont know if all other sensors will just work, there is properly a lot of deconz events i don't account for, these should be easily added though.