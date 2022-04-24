package main

import (
        "log"
        "strconv"

        "github.com/streadway/amqp"
				"github.com/gin-gonic/gin"
)

func failOnError(err error, msg string) {
        if err != nil {
                log.Fatalf("%s: %s", msg, err)
        }
}

func main() {
        conn, err := amqp.Dial("amqp://user:password@rabbitmq:5672/")
        failOnError(err, "Failed to connect to RabbitMQ")
        defer conn.Close()

        ch, err := conn.Channel()
        failOnError(err, "Failed to open a channel")
        defer ch.Close()

        err = ch.ExchangeDeclare(
                "msgs",   // name
                "fanout", // type
                true,     // durable
                false,    // auto-deleted
                false,    // internal
                false,    // no-wait
                nil,      // arguments
        )
        failOnError(err, "Failed to declare an exchange")

				d := 0

				r := gin.Default()
				
				r.GET("/send", func(c *gin.Context) {
					d += 1
					err = ch.Publish(
						"msgs", // exchange
						"",     // routing key
						false,  // mandatory
						false,  // immediate
						amqp.Publishing{
										ContentType: "text/plain",
										Body:        []byte("msg " + strconv.Itoa(d) ),
					})
					failOnError(err, "Failed to publish a message")
				
					log.Printf(" [x] Sent %s", "msg " + strconv.Itoa(d))

					c.JSON(200, gin.H{
						"message": ("msg send " + strconv.Itoa(d)),
					})
				})
				r.Run() 
}
