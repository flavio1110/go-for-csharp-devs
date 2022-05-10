package MessageExchange

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	messageExchange "github.com/flavio1110/server_client/internal/server_client_pb"
	"google.golang.org/grpc/metadata"
)

type MessageExchangeServer struct {
	messageExchange.UnimplementedMessageExchagerServiceServer
	mutex       sync.Mutex
	connections map[string]chan *messageExchange.Message
}

func NewServer() *MessageExchangeServer {
	return &MessageExchangeServer{
		connections: make(map[string]chan *messageExchange.Message),
	}
}

func (s *MessageExchangeServer) Connect(stream messageExchange.MessageExchagerService_ConnectServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok || len(md.Get("client_id")) == 0 {
		ctx.Done()
		return fmt.Errorf("Client_Id not defined in metadata")
	}

	client_id := md.Get("client_id")[0]
	ch := make(chan *messageExchange.Message)

	fmt.Printf("[%s] connected\n", client_id)
	defer fmt.Printf("[%s] disconnected\n", client_id)

	s.mutex.Lock()
	s.connections[client_id] = ch
	defer func() {
		delete(s.connections, client_id)
	}()
	s.mutex.Unlock()

	go s.setupSender(client_id, ctx, stream)
	go s.sendPing(client_id, ctx)
	for {
		message, err := stream.Recv()

		if err := s.handleMessage(client_id, message, err); err != nil {
			fmt.Printf("Error to process message %v - err %v\n", message, err)
			break
		}
	}

	return nil
}

func (s *MessageExchangeServer) handleMessage(client_id string, message *messageExchange.Message, err error) error {
	if err == io.EOF {
		return nil
	}

	if err != nil {
		return err
	}

	destination, ok := s.connections[message.Destination]

	if !ok || client_id == message.Destination {
		return nil
	}

	destination <- &messageExchange.Message{
		Type:    messageExchange.MessageType_Terminate,
		Payload: fmt.Sprintf("%s SENT: %s", client_id, message.Payload),
	}

	// s.connections[client_id] <- &messageExchange.Message{
	// 	Type:    messageExchange.MessageType_Ack,
	// 	Payload: "Message received by server",
	// }

	return nil
}

func (s *MessageExchangeServer) setupSender(client_id string, ctx context.Context, stream messageExchange.MessageExchagerService_ConnectServer) {
	for {
		select {
		case message := <-s.connections[client_id]:
			fmt.Printf("[%s]Sending: %v\n", client_id, message)
			err := stream.Send(message)

			if err != nil {
				fmt.Println("Failed to send message", err)
			}
		case <-ctx.Done():
			fmt.Printf("[%s] Sender shutdown\n", client_id)
			return
		}
	}
}

func (s *MessageExchangeServer) sendPing(client_id string, context context.Context) {
	for {
		select {
		case <-time.Tick(120 * time.Second):
			s.connections[client_id] <- &messageExchange.Message{
				Type:    messageExchange.MessageType_Payload,
				Payload: "Ping 120",
			}
		case <-context.Done():
			fmt.Printf("[%s]Won't send more ping. Connection cancelled\n", client_id)
			return
		}
	}
}
