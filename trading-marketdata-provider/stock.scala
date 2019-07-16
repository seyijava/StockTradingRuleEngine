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
     
     val stockQuteStreaming = KafkaUtils.createDirectStream[String, String](streamingContext,PreferConsistent, Subscribe[String, String](topics, kafkaParams))
     
     val tradeOrderRulesExecutor = new  TradeOrderRulesExecutor()
    
     val kafkaSink = streamingContext.sparkContext.broadcast(KafkaSink(kafkaProducerConfig))
     
    
     val stockQuoteDstream =  stockQuteStreaming.map(stockQuote => parseIncomingStockQuoteMessage(stockQuote.value()))
     
     processStockQuote(stockQuoteDstream,kafkaSink,appConfig, tradeOrderRulesExecutor);
     
     streamingContext.start();
     
     
     streamingContext.awaitTermination();
    
     streamingContext.stop(stopSparkContext = true, stopGracefully = true)
