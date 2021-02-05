#!/bin/sh

# Collect various outputs to ease up investigation in case of integration test failure
mkdir -p $ARTIFACTS_DIR
kubectl get namespaces > $ARTIFACTS_DIR/namespaces.txt
k8ssandra_ns=$(kubectl get namespaces|grep k8ssandra|cut -d' ' -f1)

# List all objects from the k8ssandra namespace
kubectl get all -n $k8ssandra_ns > $ARTIFACTS_DIR/k8ssandra_all.txt

# Describe the cassandradatacenter resource
kubectl describe cassandradatacenter/dc1 -n $k8ssandra_ns > $ARTIFACTS_DIR/cassandradc_describe.txt

# Extract logs from the Cassandra pods
for pod in $(kubectl get pods -n $k8ssandra_ns|grep dc1-default-sts|cut -d' ' -f1); do
    echo "Storing artifacts for pod $pod..."
    kubectl logs pod/$pod cassandra -n $k8ssandra_ns > $ARTIFACTS_DIR/${pod}_cassandra.log || echo "can't extract cassandra log"
    kubectl logs pod/$pod server-system-logger -n $k8ssandra_ns > $ARTIFACTS_DIR/${pod}_system_log.log || echo "can't extract server-system-logger log"
    kubectl logs pod/$pod medusa -n $k8ssandra_ns > $ARTIFACTS_DIR/${pod}_medusa.log || echo "can't extract medusa logs"
done

# Extract backup information
for backup in $(kubectl get CassandraBackup -n $k8ssandra_ns|grep -v "NAME"|cut -d' ' -f1); do
    echo "Storing artifacts for backup $backup..."
    kubectl describe CassandraBackup/$backup -n $k8ssandra_ns > $ARTIFACTS_DIR/backup_${backup}_describe.txt
done