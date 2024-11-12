Done
-----------------------
1.Describe the MySQL CR and Try restarting all the pods

2.Increase mysql variable max_connections


TODO
------------
1.Also try tuning mysql for memory optimization

2.Check Slow Query log file here /var/log/mysql/mysql-slow.log

3.Reason for alert: Innodb_log_waits (MySQLInnoDBLogWaits)

4.Try reconfiguring innodb_log_buffer_size (MySQLInnoDBLogWaits)

5.Increase mysql variable open_files_limit (MySQLTooManyOpenFiles)

6.Scale MySQL using KubeDB Scaling OpsRequest
(MySQLHighQPS | MySQLHighIncomingBytes | MySQLHighOutgoingBytes
)