package service

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"nhooyr.io/websocket"
)

const connSlow = "connection too slow to keep up with messages"
const subscriberMessageBuffer = 16

func (s *Service) BinWatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	binKey, err := uuid.Parse(vars["binKey"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	bin, ok := s.dataStore.GetBin(binKey)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "")
	err = s.subscribe(r.Context(), c, bin.BinKey)
	if errors.Is(err, context.Canceled) {
		return
	}
	if isCloseStatus(err, websocket.StatusNormalClosure, websocket.StatusGoingAway) {
		return
	}
	if err != nil {
		log.Println(err)
		return
	}
}

type subscriber struct {
	msgs      chan []byte
	closeSlow func()
}

func (s *Service) subscribe(ctx context.Context, c *websocket.Conn, binKey uuid.UUID) error {
	go ping(ctx, c)
	ctx = c.CloseRead(ctx)
	sub := &subscriber{
		msgs: make(chan []byte, subscriberMessageBuffer),
		closeSlow: func() {
			c.Close(websocket.StatusPolicyViolation, connSlow)
		},
	}
	s.addSubscriber(sub, binKey)
	defer s.deleteSubscriber(sub, binKey)
	for {
		select {
		case msg := <-sub.msgs:
			err := writeTimeout(ctx, time.Second*5, c, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *Service) addSubscriber(sub *subscriber, binKey uuid.UUID) {
	s.subscribersMutex.Lock()
	defer s.subscribersMutex.Unlock()
	subs, ok := s.subscribers[binKey]
	if !ok {
		return
	}
	subs[sub] = struct{}{}
}

func (s *Service) deleteSubscriber(sub *subscriber, binKey uuid.UUID) {
	s.subscribersMutex.Lock()
	defer s.subscribersMutex.Unlock()
	subs, ok := s.subscribers[binKey]
	if !ok {
		return
	}
	delete(subs, sub)
}

func (s *Service) publish(msg []byte, binKey uuid.UUID) {
	s.subscribersMutex.Lock()
	defer s.subscribersMutex.Unlock()
	s.publishLimiter.Wait(context.Background())
	subs, ok := s.subscribers[binKey]
	if !ok {
		return
	}
	for sub := range subs {
		select {
		case sub.msgs <- msg:
		default:
			go sub.closeSlow()
		}
	}
}

func isCloseStatus(err error, expectedCodes ...websocket.StatusCode) bool {
	closeStatus := websocket.CloseStatus(err)
	for _, code := range expectedCodes {
		if closeStatus == code {
			return true
		}
	}
	return false
}

func ping(ctx context.Context, c *websocket.Conn) {
	for {
		for range time.Tick(time.Second * 60) {
			if err := c.Ping(ctx); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return c.Write(ctx, websocket.MessageText, msg)
}
