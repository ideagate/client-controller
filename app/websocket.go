package main

import (
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bayu-aditya/ideagate/backend/client/controller/adapter/kubernetes"
	usecasewebsocketmanagement "github.com/bayu-aditya/ideagate/backend/client/controller/usecase/websocketmanagement"
	"github.com/bayu-aditya/ideagate/backend/client/controller/usecase/workermanagement"
	"github.com/bayu-aditya/ideagate/backend/core/utils/log"
	"github.com/gorilla/websocket"
)

func NewWebsocketClient() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/event/ws", RawQuery: "project_id=abc"}
	log.Info("connecting to %s", u.String())

	reqHeaders := http.Header{"connection_id": []string{"123456789"}}
	websocketConn, resp, err := websocket.DefaultDialer.Dial(u.String(), reqHeaders)
	if err != nil {
		log.Error("dial error: %v", err)
		if resp != nil {
			log.Error("response status: %s", resp.Status)
			for key, values := range resp.Header {
				for _, value := range values {
					log.Error("response header: %s: %s", key, value)
				}
			}
		}
		log.Fatal("dial failed")
	}
	defer websocketConn.Close()

	// Initialize Adapters
	adapterKubernetes, err := kubernetes.New()
	if err != nil {
		log.Fatal("failed to create kubernetes adapter: %v", err)
	}

	// Initialize Usecases
	usecaseWorkerManagement := workermanagement.New(adapterKubernetes)
	usecaseRouter := usecasewebsocketmanagement.NewRouter(usecaseWorkerManagement)
	usecaseWebsocketManagement := usecasewebsocketmanagement.NewWebsocketManagement(websocketConn, usecaseRouter)

	done := make(chan struct{})

	go func() {
		for {
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Error("Recovered in WorkerSubscriber: %+v", r)
					}
				}()
				usecaseWebsocketManagement.WorkerSubscriber()
			}()
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Error("Recovered in WorkerPublisher: %v", r)
					}
				}()
				usecaseWebsocketManagement.WorkerPublisher()
			}()
			time.Sleep(time.Second)
		}
	}()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	usecaseWebsocketManagement.Close()
	log.Info("Shutting down controller...")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := websocketConn.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Info("write: %v", err)
				return
			}
		case <-interrupt:
			log.Info("interrupt")
			// Cleanly close the connection by sending a close message and then waiting for the server to close the connection.
			err := websocketConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Info("write close: %v", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
