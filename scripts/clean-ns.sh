#!/bin/bash

# Helper to cleanup k8ssandra (or any other specified) namespace.
# Change  the filter as needed for a 'contains' match.
# A deleted-namespaces.txt file will contain a history of deletions.

NS_CONTAINS_FILTER="k8ssandra"
kubectl get namespace > namespaces.txt

while read p; do
  if [[ "$p" == *"$NS_CONTAINS_FILTER"* ]]; then
    echo "$p"
  fi
done <namespaces.txt

echo "Sure you want to delete these namespaces (type: yes)?"
read user_confirm

if [[ "$user_confirm" == "yes" ]]; then
   dt=`date ` 
   echo "$dt" > deleted-namespaces.txt
   while read p; do
     if [[ "$p" == *"$NS_CONTAINS_FILTER"* ]]; then
       target_ns=`echo "$p" | awk '{print $1}'`
       echo "Deleting ns: $target_ns"
       kubectl delete namespace/"$target_ns"
       echo "Deleted ns: $target_ns" >> deleted-namespaces.txt
     fi
   done <namespaces.txt
fi
