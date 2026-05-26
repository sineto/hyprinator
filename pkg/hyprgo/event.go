package hyprgo

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

func NewEventClient(socket string) *EventClient {
	if socket == "" {
		panic("sign is empty")
	}

	conn, err := net.Dial("unix", socket)
	if err != nil {
		panic(err)
	}

	return &EventClient{conn: conn}
}

func (e *EventClient) Receive(ctx context.Context)([]ReceiveData, error) {
	buf := make([]byte, BUFSIZE)

	n, err := e.readWithContext(ctx, buf)
	if err != nil {
		return nil, err
	}

	buf = buf[:n]

	var recv []ReceiveData

	rawEvents := strings.Split(string(buf), "\n")

	for _, event := range rawEvents {
		if event == "" {
			continue
		}

		split := strings.Split(event, ">>")
		if split[0] == "" || split[1] == "" || split[1] == "," {
			continue
		}

		recv = append(recv, ReceiveData{
			Type: EventType(split[0]),
			Data: split[1],
		})
	}

	return recv, nil
}

func (e *EventClient) Subscribe(ctx context.Context, ev EventHandler, events ...EventType) error {
	for {
		if err := e.receiveAndProcessEvent(ctx, ev, events...); err != nil {
			return fmt.Errorf("event processing: %w", err)
		}
	}
}

func (e *EventClient) readWithContext(ctx context.Context, buf [] byte) (n int, err error) {
	done := make(chan struct{})

	go func() {
		n, err = e.conn.Read(buf)
		close(done)
	}()

	select {
	case <-done:
		return n, err
	case <- ctx.Done():
		err = e.conn.SetReadDeadline(time.Now())
		if err != nil {
			return 0, err
		}

		defer func ()  {
			if e := e.conn.SetReadDeadline(time.Time{}); e != nil {
				err = errors.Join(err, e)
			}
		}()
		
		return 0, errors.Join(err, ctx.Err())
	}
}

func (e *EventClient) receiveAndProcessEvent(ctx context.Context, ev EventHandler, events ...EventType) error {
	msg, err := e.Receive(ctx)
	if err != nil {
		return err
	}

	for _, data := range msg {
		e.processEvent(ev, data, events)
	}

	return nil
}

func (e *EventClient) processEvent(ev EventHandler, msg ReceiveData, events []EventType) {
	for _, event := range events {
		rawEvent := strings.Split(string(msg.Data), ",")

		if msg.Type == event {
			switch event {
			case EventOpenWindow:
				ev.OpenWindow(OpenWindow{
					Address: rawEvent[0],
					Class: rawEvent[1],
					Title: rawEvent[3],
				})
			}
		}
	}
}

