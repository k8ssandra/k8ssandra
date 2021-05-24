#!/bin/sh

# Collect various outputs to ease up investigation in case of integration test failure
mkdir -p $ARTIFACTS_DIR
kubectl get namespaces > $ARTIFACTS_DIR/namespaces.txt
k8ssandra_ns=$(kubectl get namespaces|grep k8ssandra|cut -d' ' -f1)

kubectl cluster-info dump --namespaces $k8ssandra_ns -o yaml --output-directory $ARTIFACTS_DIR

# List all objects from the k8ssandra namespace
kubectl get all -n $k8ssandra_ns > $ARTIFACTS_DIR/k8ssandra_all.txt

# Describe the cassandradatacenter resource
kubectl describe cassandradatacenter/dc1 -n $k8ssandra_ns > $ARTIFACTS_DIR/cassandradc_dc1_describe.txt
kubectl get cassandradatacenter/dc1 -o yaml -n $k8ssandra_ns > $ARTIFACTS_DIR/cassandradc_dc1.txt

# Extract backup information
for backup in $(kubectl get CassandraBackup -n $k8ssandra_ns|grep -v "NAME"|cut -d' ' -f1); do
    echo "Storing artifacts for backup $backup..."
    kubectl describe CassandraBackup/$backup -n $k8ssandra_ns > $ARTIFACTS_DIR/backup_${backup}_describe.txt
    kubectl get CassandraBackup/$backup -o yaml -n $k8ssandra_ns > $ARTIFACTS_DIR/backup_${backup}.txt
done

# Extract restore information
for restore in $(kubectl get CassandraRestore -n $k8ssandra_ns|grep -v "NAME"|cut -d' ' -f1); do
    echo "Storing artifacts for restore $restore..."
    kubectl describe CassandraRestore/$restore -n $k8ssandra_ns > $ARTIFACTS_DIR/restore_${restore}_describe.txt
    kubectl get CassandraRestore/$restore -o yaml -n $k8ssandra_ns > $ARTIFACTS_DIR/restore_${restore}.txt
done

