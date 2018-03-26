#!/bin/bash


DISK_DEV_NAME_SHORT=`lsblk |grep disk|awk '{print $1}'`

for i in $DISK_DEV_NAME_SHORT
do
    DISK_SIZE=`smartctl -i /dev/$i|grep -i "User Capacity:"|awk -F "[" ' {print $2}'i|awk -F "]" '{print $1}'`
    DISK_HEALTH=`smartctl -H /dev/$i|grep -i "smart health status"|awk '{print $4}'`

    echo "Disk: /dev/$i Size: $DISK_SIZE Health_Status: $DISK_HEALTH  "
done
