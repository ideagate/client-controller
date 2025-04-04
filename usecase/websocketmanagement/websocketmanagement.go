package websocketmanagement

import (
	"context"
	"encoding/json"
	"log"
	"time"

	entitywebsocket "github.com/bayu-aditya/ideagate/backend/core/model/entity/websocket"
	"github.com/gorilla/websocket"
)

type IWebsocketManagementUsecase interface {
	WorkerSubscriber()
	WorkerPublisher()
	Close() error
}

func NewWebsocketManagement(wsConn *websocket.Conn, router IRouter) IWebsocketManagementUsecase {
	return &websocketManagement{
		conn:          wsConn,
		usecaseRouter: router, // TODO
	}
}

type websocketManagement struct {
	conn          *websocket.Conn
	usecaseRouter IRouter
}

func (w *websocketManagement) WorkerSubscriber() {
	ctx := context.Background()

	for {
		_, message, err := w.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			continue
		}

		// unmarshal event request in JSON format
		var eventRequest entitywebsocket.Event
		if err = json.Unmarshal(message, &eventRequest); err != nil {
			log.Println("unmarshal event request:", err)
			continue
		}
		log.Printf("recv: %+v", eventRequest)

		response, err := w.usecaseRouter.Switch(ctx, eventRequest.Type, eventRequest.Data)
		if err != nil {
			// TODO handle error
			log.Println("error:", err)
			// TODO construct error into event response and send error to server
		}

		// construct event response in JSON format
		eventResponse := eventRequest
		eventResponse.Data = response

		eventResponseJson, err := json.Marshal(eventResponse)
		if err != nil {
			log.Println("Marshal event response:", err)
			continue
		}

		if err = w.conn.WriteMessage(websocket.TextMessage, eventResponseJson); err != nil {
			log.Println("write:", err)
			continue
		}
	}
}

func (w *websocketManagement) WorkerPublisher() {
	for {
		time.Sleep(2 * time.Second)
		if err := w.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			log.Println("write:", err)
		}
	}
}

func (w *websocketManagement) Close() error {
	// Cleanly close the connection by sending a close message and then waiting for the server to close the connection.
	err := w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}

	if err = w.conn.Close(); err != nil {
		return err
	}

	return nil
}
