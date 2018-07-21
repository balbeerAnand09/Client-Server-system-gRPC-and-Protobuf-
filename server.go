package main

import (
	"log"
	"net"
	pb "protocols"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type server struct{}

var user = make(map[string]bool)
var ch = make([]pb.Text, 0)

func (s *server) Enter(ctx context.Context, u *pb.Str) (*pb.Str, error) {
	if user[u.Noti] == true {
		return &pb.Str{Noti: "USER ALREADY EXIST"}, nil
	}
	user[u.Noti] = true
	log.Println(u.Noti, "Entered the room.")
	return &pb.Str{Noti: "Congo, you have entered the room"}, nil
}

func (s *server) Send(ctx context.Context, t *pb.Text) (*pb.Ack, error) {
	ch = append(ch, *t)
	log.Println("In Send: From:", t.Msg.From, "To:", t.Msg.To, "Mess:", t.Msg.Mess)
	return &pb.Ack{Done: true}, nil
}

func (s *server) Recieve(ctx context.Context, t *pb.Text) (*pb.Text, error) {
	for i, txt := range ch {
		if txt.Msg.To == t.Msg.To {
			t.Msg.From = txt.Msg.From
			t.Msg.Mess = txt.Msg.Mess
			log.Println("From:", t.Msg.From, "To:", t.Msg.To, "Mess:", t.Msg.Mess)
			//ch = append(ch[:i], ch[i+1:]...)
			ch[i] = ch[len(ch)-1]
			ch[len(ch)-1] = pb.Text{}
			ch = ch[:len(ch)-1]
			return t, nil
		}
	}
	return t, nil
}

func main() {
	lis, err := net.Listen("tcp", ":500")
	if err != nil {
		log.Fatalf("Port is not listening: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChatterServer(s, &server{})
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
