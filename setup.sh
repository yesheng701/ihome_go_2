#启动redisserver
redis-server ./conf/redis.conf

#启动fastdfs - tracker
fdfs_trackerd ./conf/tracker.conf restart
#启动fastdfs - storage
fdfs_storaged /root/workspace/go/src/ihome_idlefish/conf/storage.conf restart
