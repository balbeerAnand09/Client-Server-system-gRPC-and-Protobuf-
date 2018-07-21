package main

import (
	"fmt"
	"log"
	pb "protocols"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func input(ctx context.Context, usr string, c *pb.ChatterClient) {
	for {
		var (
			ms string
			to string
		)
		fmt.Println("Send Message To: ")
		fmt.Scanln(&to)
		fmt.Println("Enter Message: ")
		fmt.Scanln(&ms)
		msg := &pb.Text{Msg: &pb.TextMail{From: usr, To: to, Mess: ms}}
		(*c).Send(ctx, msg)
	}
}

func main() {
	conn, err := grpc.Dial("localhost:500", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Connection Not established: %v", err)
	}
	defer conn.Close()
	// creating the server instance
	c := pb.NewChatterClient(conn)

	ctx := context.Background()
	//defer cancel()

	var usr string
	for {
		fmt.Printf("Enter Username: ")
		fmt.Scanln(&usr)

		str := &pb.Str{Noti: usr}
		r, err := c.Enter(ctx, str)

		if err != nil {
			log.Fatalf("Could not connect: %v", err)
			continue
		}
		er := &pb.Str{Noti: "USER ALREADY EXIST"}
		if r.Noti == er.Noti {
			log.Println("Could not connect:", er.Noti)
			continue
		}
		fmt.Println(r.Noti)
		break
	}
	txt := &pb.Text{Msg: &pb.TextMail{From: "", To: usr, Mess: ""}}

	go input(ctx, usr, &c)
	for {
		txt, err := c.Recieve(ctx, txt)
		//fmt.Println("Sd")
		if err == nil && txt.Msg.From != "" {
			fmt.Println("\tFrom:", txt.Msg.From, "To:", txt.Msg.To, "Mess:", txt.Msg.Mess)
		}
	}
}
