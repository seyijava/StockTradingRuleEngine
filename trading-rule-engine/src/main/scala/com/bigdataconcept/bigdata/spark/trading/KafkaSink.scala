package com.bigdataconcept.bigdata.spark.trading

import java.util.concurrent.Future

import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.ProducerRecord;
import org.apache.kafka.clients.producer.ProducerConfig;
import org.apache.kafka.clients.producer.RecordMetadata





class KafkaSink(createProducer: () => KafkaProducer[String, String]) extends Serializable {

  lazy val producer = createProducer()

  def send(topic: String, key: String, value: String): Future[RecordMetadata] =
    producer.send(new ProducerRecord[String, String](topic, key, value))

  def send(topic: String, value: String): Future[RecordMetadata] =
    producer.send(new ProducerRecord[String, String](topic, value))
}


object KafkaSink {
  def apply(config: java.util.Properties): KafkaSink = {
    val f = () => {
     val producer =  new KafkaProducer[String, String](config)
      sys.addShutdownHook {
        producer.close()
      }

      producer
    }
    new KafkaSink(f)
  }
}