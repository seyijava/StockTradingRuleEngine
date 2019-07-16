package com.bigdataconcept.bigdata.spark.trading
import com.redis._
import scala.Serializable
import com.google.gson.Gson
import com.google.gson.GsonBuilder
import org.apache.spark.broadcast.Broadcast
import com.google.gson.Gson
import org.apache.log4j.Logger


/*
 * Oluwaseyi Otun 
 * 
 */
class TradeOrderRulesExecutor() extends Serializable {
  
    private final val BUY = "BUY"
    private final val SELL = "SELL"

  @transient lazy val logger = Logger.getLogger(getClass.getName) 
   
   
   def executeTradeRule(stockMovingAVGList: Iterator[StockMovingAvg], redisClient: RedisClient,kafkaTopic: String, kafkaSink: Broadcast[com.bigdataconcept.bigdata.spark.trading.KafkaSink]) 
   {
         logger.info("Executing Trade Rule With Moving Average\n\n")
         stockMovingAVGList.foreach(stockMovingAVG => {
           val tradeRules = redisClient.lrange("TradeRule",0,-1)
            tradeRules.map(rule => rule.foreach(
               rule => evaluateTradeRuleToMovingAVGAndSendTradeOrder(convertJsonDataToTradeRule(rule.get),stockMovingAVG,kafkaTopic,kafkaSink)))
         } )
          
   }
   
   
   private def evaluateTradeRuleToMovingAVGAndSendTradeOrder(rule: TradeRule,stockMovingAVG: StockMovingAvg,kakfaTopic: String, kafkaSink: Broadcast[com.bigdataconcept.bigdata.spark.trading.KafkaSink]) : Unit = 
   {
     // logger.info("Evaluating Trade Rule against Simple Moving Averag" + stockMovingAVG.symbol + "\n");
       if(rule.symbol.endsWith(stockMovingAVG.symbol))
       {
             rule.tradeType match
             {
               case 0 =>
               if(stockMovingAVG.avgBidPrice  < rule.buyPrice)
               {
                   val tradeOrder = new TradeOrder(rule.accountNumber   , rule.symbol , rule.buyQuantity , rule.buyPrice ,BUY)
                   
                   val tradeOrderRequest = new Gson().toJson(tradeOrder)
                   logger.info(String.format("Sending Buy Trade Order To Kakfa Topic for onward Processsing via the exchange [%s] \n\n", tradeOrderRequest))
                   kafkaSink.value.send(kakfaTopic,tradeOrderRequest)
               }
               
               case 1 =>
               if(stockMovingAVG.avgAskPrice  > rule.sellPrice)
               {
            	   val tradeOrder = new TradeOrder(rule.accountNumber, rule.symbol , rule.sellQuantity  , rule.sellPrice  ,SELL)
            	   
            	   val tradeOrderRequest = new Gson().toJson(tradeOrder)
            	   
            	   logger.info(String.format("Sending Sell Trade Order To Kakfa Topic for onward Processsing via the exchange [%s] \n\n", tradeOrderRequest))
            	   
            	   kafkaSink.value.send(kakfaTopic,tradeOrderRequest)  
               }
             }
       }
      
      
   }
   
   
  def convertJsonDataToTradeRule(message: String) : TradeRule =
  {
    
     logger.info("Trade Rule \n\n" + message)
     
      var gson = new GsonBuilder().setDateFormat("yyyy-MM-dd'T'HH:mm:ss").create();
      var tradeRule = gson.fromJson(message, classOf[TradeRule]);
      
      return tradeRule
  }
}