package com.bigdataconcept.bigdata.spark.trading

import java.sql.Timestamp
import  org.apache.spark.rdd.RDD
import org.apache.spark.sql.SparkSession
import com.typesafe.config.{Config,ConfigFactory}
import org.apache.spark.SparkConf
import org.apache.spark.streaming.dstream.DStream
import org.apache.spark.streaming.StreamingContext
import org.apache.spark.streaming.Seconds
import org.apache.kafka.clients.consumer.ConsumerRecord
import org.apache.kafka.common.serialization.StringDeserializer
import org.apache.spark.streaming.kafka010._
import org.apache.spark.streaming.kafka010.LocationStrategies.PreferConsistent
import org.apache.spark.streaming.kafka010.ConsumerStrategies.Subscribe
import java.io.File
import org.apache.kafka.clients.producer._
import org.apache.kafka.common.serialization.LongSerializer
import org.apache.kafka.common.serialization.StringSerializer
import org.apache.spark.broadcast.Broadcast
import com.google.gson.Gson
import com.google.gson.GsonBuilder
import java.util.concurrent.Future
import org.apache.spark.sql.{DataFrame, Dataset, SparkSession}
import org.apache.spark.sql.functions._
import com.redis._
import com.redis.RedisClientPool
import org.apache.log4j.Logger



/**
 * @author Oluwaseyi Otun.
 */

object TradeRuleEngine 
{

  
   @transient lazy val logger = Logger.getLogger(getClass.getName)  
   
  case class StockEvent(symbol:String,bid: Double, ask: Double, quoteTimeStamp: Timestamp) extends Serializable
  
  case class AppConfig(kafkaBrokerAddress:  String, kafkaStockQouteTopic: String, kafkaConsumerGroupid: String,windowLenght: Long, slideInterval: Long, checkPointDir: String,
     kafkaTradeOrderTopic: String,influxDBHost: String, influxDBPort: String, redisHost: String, redisPort: Int)
  extends Serializable
  
    def main(args : Array[String]) {
    
     val bigdataPath = sys.env("BIGDATA_APP_HOME")
     val configFile = new File(bigdataPath + "/tradeEngine.conf")
     val configFactory = ConfigFactory.parseFile(configFile).getConfig("appConfig")
     val appConfig = AppConfig(configFactory.getString("kafka.host"),configFactory.getString("kafka.stockQouteTopic"),
          configFactory.getString("kafka.consumerGroup"),configFactory.getLong("spark.windowLenght"),configFactory.getLong("spark.slideInterval"),
          configFactory.getString("spark.checkPointDir"),configFactory.getString("kafka.tradeOrderTopic"),configFactory.getString("influxDB.host"),configFactory.getString("influxDB.port"),configFactory.getString("redis.host"),configFactory.getInt("redis.port"))
     val sparkconf = new SparkConf().setMaster("local[2]").setAppName("TradeRuleEngine")
     val streamingContext = new StreamingContext(sparkconf, Seconds(5))
     val kafkaParams = Map[String, Object]("bootstrap.servers" -> appConfig.kafkaBrokerAddress,
                    	"key.deserializer" -> classOf[StringDeserializer],
                    	"value.deserializer" -> classOf[StringDeserializer],
                    	"group.id" -> appConfig.kafkaConsumerGroupid,
                    	"auto.offset.reset" -> "latest",
                    	"enable.auto.commit" -> (false: java.lang.Boolean)
                       )
    val kafkaProducerConfig = {
    val p = new java.util.Properties()
    p.setProperty("bootstrap.servers", appConfig.kafkaBrokerAddress)
    p.setProperty("key.serializer", classOf[StringSerializer].getName)
    p.setProperty("value.serializer", classOf[StringSerializer].getName)
    p
  }
     val topics: Set[String] = appConfig.kafkaStockQouteTopic.split(",").map(_.trim).toSet
     val stockQuteDStream = KafkaUtils.createDirectStream[String, String](streamingContext,PreferConsistent, Subscribe[String, String](topics, kafkaParams))
   
     val tradeOrderRulesExecutor = new  TradeOrderRulesExecutor()
     val kafkaSink = streamingContext.sparkContext.broadcast(KafkaSink(kafkaProducerConfig))
     val stockQuoteEvents =  stockQuteDStream.map(stockQuote => parseIncomingStockQuoteMessage(stockQuote.value()))
     analyzeMovingAverage(stockQuoteEvents, kafkaSink, appConfig,tradeOrderRulesExecutor)
     streamingContext.start();
     streamingContext.awaitTermination();
     streamingContext.stop(stopSparkContext = true, stopGracefully = true)
  } 
  
        
       def executeTradeOrderRuleOnSMVG(stockMovingAvgRDD: RDD[StockMovingAvg], kafkaSink: Broadcast[KafkaSink], appConfig: AppConfig, tradeOrderRulesExecutor: TradeOrderRulesExecutor) {
         if(!stockMovingAvgRDD.isEmpty)
         {
             stockMovingAvgRDD.foreachPartition( stockMovingAVG=>{
              val redisClient = new RedisClient(appConfig.redisHost , 6379)
               tradeOrderRulesExecutor.executeTradeRule(stockMovingAVG, redisClient,appConfig.kafkaTradeOrderTopic,kafkaSink)
            } )
         }
         else
         {
           logger.info("RDD is Empty Waiting for Data\n\n")
         }
         
    }
  
  
  
  
    def analyzeMovingAverage(stockQuoteDstream: DStream[StockQuote], kafkaSink: Broadcast[KafkaSink], appConfig: AppConfig, tradeOrderRulesExecutor: TradeOrderRulesExecutor) {
      logger.info("Caclulating Simple Moving Average\n")
      stockQuoteDstream.foreachRDD(  
        rdd=> 
         if(!rdd.isEmpty)
         {
           val spark = SparkSession.builder.appName("StockQuoteMVG").getOrCreate()
           import spark.implicits._
           val stockEvent = rdd.map(stock => StockEvent(stock.symbol ,stock.bid, stock.ask, stock.quoteTimeStamp)).toDF.as[StockEvent];
           val stockQuoteSMVGDataFrame  = analyzeMovingAverageInSlidingWindow(stockEvent).groupBy("symbol").agg(avg("avgBid"),avg("avgAsk"))
           stockQuoteSMVGDataFrame.show
           val stockQouteSMVGRDD =  stockQuoteSMVGDataFrame.map(_.toSeq.toList match {
            case List(tick: String, avgBid: Double, avgAsk: Double) => StockMovingAvg(symbol = tick, avgBidPrice=avgAsk, avgAskPrice= avgAsk)
            })
            executeTradeOrderRuleOnSMVG( stockQouteSMVGRDD.toJavaRDD, kafkaSink: Broadcast[KafkaSink], appConfig: AppConfig, tradeOrderRulesExecutor: TradeOrderRulesExecutor)
         }
         else
         {
            logger.info("RDD is Empty Waiting for Data\n\n")
         }
         )
    }
  
    
     def  analyzeMovingAverageInSlidingWindow(df: Dataset[StockEvent]) : DataFrame  = {
      df
      .withWatermark("quoteTimeStamp", "5 seconds") // Ignore data if they are late by more than 5 seconds
      .groupBy(window(col("quoteTimeStamp"),"3 seconds","1 seconds"), col("symbol"))  //sliding window of size 4 seconds, that slides every 1 second
      .agg(avg("bid").alias("avgBid"), min("bid").alias("minBid"), max("bid").alias("maxBid"), count("bid").alias("bidCount"),
           avg("ask").alias("avgAsk"), min("ask").alias("minAsk"), max("ask").alias("maxAsk"), count("ask").alias("askCount"))
      .select("window.start", "window.end", "avgBid", "minBid",  "maxBid", "bidCount",
               "avgAsk", "minAsk", "maxAsk",  "bidCount","symbol")
  }
  
  
       
   
 
  def parseIncomingStockQuoteMessage(message: String) : StockQuote =
  {
      var gson = new GsonBuilder().setDateFormat("yyyy-MM-dd'T'HH:mm:ss").create();
      var stockQuote = gson.fromJson(message, classOf[StockQuote]);
      return stockQuote
  }
  
}