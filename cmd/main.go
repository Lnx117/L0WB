package main

import (
	l0wb "L0WB"
	"L0WB/pkg/handler"
	"L0WB/pkg/repository"
	"L0WB/pkg/service"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	service := service.NewService(repos)
	handlers := handler.NewHandler(service)

	service.Cache = service.ReadOrdersData.ReadAllOrdersData()

	go start(service)

	srv := new(l0wb.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error while server running: %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("../configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func logCloser(c io.Closer) {
	if err := c.Close(); err != nil {
		logrus.Printf("close error: %s", err)
	}
}

func startSubscriber(s *service.Service, conn stan.Conn) {
	var err error

	_, err = conn.Subscribe("counter", func(msg *stan.Msg) {

		content, _ := s.ParseJSON.ParseJSON(msg.Data)
		s.SaveOrderData.SaveOrderData(content)
		s.Cache[content.Order_uid] = content
		// Print the value and whether it was redelivered.
		fmt.Printf("seq = %d [redelivered = %v] mes= %s \n", msg.Sequence, msg.Redelivered, msg.Data)

		// Add jitter..
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

		// Mark it is done.
		msg.Ack()

	}, stan.DurableName("i-will-remember"), stan.MaxInflight(100), stan.SetManualAckMode())

	if err != nil {
		logrus.Print(err)
	}
}

func start(s *service.Service) {
	if err := run(s); err != nil {
		logrus.Fatal(err)
	}
}

func run(s *service.Service) error {
	conn, err := stan.Connect(
		viper.GetString("natStreaming.clusterId"),
		viper.GetString("natStreaming.clientId"),
		stan.NatsURL(viper.GetString("natStreaming.natsURL")),
	)
	if err != nil {
		return err
	}
	defer logCloser(conn)

	done := make(chan struct{})
	time.Sleep(time.Duration(rand.Intn(4000)) * time.Millisecond)
	startSubscriber(s, conn)
	<-done

	return nil
}
