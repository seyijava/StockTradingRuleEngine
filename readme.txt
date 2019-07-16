

HDFS
=====================================================================================================
cd $HADOOP_HOME
bin/start-all.sh
hdfs dfs -mkdir /bigdata/fintech/lambda/datalake/stockquoteData




Kafka 
=======================================================================================================
cd $Flume_HOME
bin/kafka-topics --create  --zookeeper localhost:2181  --partitions 1 --replication-factor 1 --topic MarketDataTopic

bin/kafka-topics --create  --zookeeper localhost:2181  --partitions 1 --replication-factor 1 --topic TradeOrderTopic


Flume
=====================================================================================================
cd $Flume_HOME
bin/flume-ng agent -n lambda -c conf -f conf/lambda_hdfs_flume_conf.properties


  
  
  
Spark
==========================================================================================================  
spark-submit --class com.bigdataconcept.bigdata.spark.trading.TradeRuleEngine  --master local[2]  trading-rule-engine-jar-with-dependencies.jar


 