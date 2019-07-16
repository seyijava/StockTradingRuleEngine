package com.bigdataconcept.bigdata.spark.trading

import java.sql.Timestamp
import java.io.Serializable

  case class StockMovingAvg(symbol: String, avgBidPrice: Double, avgAskPrice: Double) extends Serializable
  case class StockQuote(symbol: String,shortName: String, qouteType: String, regularMarketChangePercent: Double, regularMarketPreviousClose: Double,
                        regularMarketPrice: Double, regularMarketTime: Long, regularMarketChange: Double,regularMarketOpen: Double,regularMarketDayHigh: Double,
                        regularMarketDayLow: Double, regularMarketVolume: Long, bid: Double,ask: Double,bidSize: Long,askSize: Long, preMarketPrice: Double,
                         preMarketChangePercent: Double, preMarketTime: Long, postMarketChange: Double,postMarketChangePercent: Double,postMarketTime: Long, fiftyTwoWeekLowChange: Double,
                         fiftyTwoWeekLowChangePercent: Double, fiftyTwoWeekHighChange: Double, fiftyTwoWeekHighChangePercent: Double, fiftyTwoWeekLow: 
                         Double,fiftyDayAverage: Double,fiftyDayAverageChange: Double, fiftyDayAverageChangePercent: Double, twoHundredDayAverage: Double,
                         twoHundredDayAverageChange: Double,twoHundredDayAverageChangePercent: Double,averageDailyVolume3Month: Double,currency: String,tradeable: Boolean,
                         exchangeDataDelayedBy: Long,sourceInterval: Long,exchangeTimezoneName: String,exchangeTimezoneShortName: String,market: String,exchange: String, quoteTimeStamp: Timestamp)
                         extends Serializable                          
  case class TradeRule(accountNumber:String,symbol:String,sellPrice:Double,buyPrice: Double,sellQuantity: Long,buyQuantity: Long, tradeType: Int) extends Serializable
  case class TradeOrder(accountNumber: String, symbol: String, quantity: Long, price: Double, orderType: String) extends Serializable