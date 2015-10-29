# portmapper

A portmapper backed by etcd v2.

# Testing

* Set the environmental variable PORTMAPPER_ETCD_HOST="http://etcd-docker-ip"
* Run ``` docker-compose up ```
* go test
