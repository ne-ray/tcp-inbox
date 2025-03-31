package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/ne-ray/tcp-inbox/pkg/nticlient"
	snn "github.com/ne-ray/tcp-inbox/pkg/scanersplitter"
)

func main() {
	cmd := &cli.Command{
		Name:  "NTI Client",
		Usage: "Client for NTI (Word of Wisdom network)",
		Commands: []*cli.Command{
			{
				Name:    "post",
				Aliases: []string{"p"},
				Usage:   "post paragraph Word of Wisdom",
				Action:  postWoWfunc,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "server",
						Aliases: []string{"s"},
						Value:   "127.0.0.1",
						Usage:   "server NTI for send paragraph",
					},
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Value:   19777,
						Usage:   "port NTI for send paragraph",
					},
					&cli.IntFlag{
						Name:     "line",
						Aliases:  []string{"l"},
						Required: true,
						Usage:    "line of Word of Wisdom book",
					},
					&cli.IntFlag{
						Name:     "chapter",
						Aliases:  []string{"c"},
						Required: true,
						Usage:    "chapter of Word of Wisdom book",
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func postWoWfunc(ctx context.Context, cmd *cli.Command) error {
	s := cmd.String("server")
	p := int(cmd.Int("port"))
	l := int(cmd.Int("line"))
	c := int(cmd.Int("chapter"))

	fmt.Println("Enter Part of Word of Wisdom book:")

	scanner := snn.New(os.Stdin, []byte{'\n', '\n'})

	var t string
	if scanner.Scan() {
		t = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	n, err := nticlient.New(s, p)
	if err != nil {
		return err
	}

	fmt.Print("Handshake for select support protocol - ")
	if err := n.SelectSupportProtocols(); err != nil {
		fmt.Print("Error\n")

		return err
	}
	fmt.Print("OK\n")

	fmt.Print("Request PoW data for calculate - ")
	if err := n.RequestPoW(); err != nil {
		fmt.Print("Error\n")

		return err
	}
	fmt.Print("OK\n")

	fmt.Print("Calculate PoW - ")
	if err := n.CalculatePoW(); err != nil {
		fmt.Print("Error\n")

		return err
	}
	fmt.Print("OK\n")

	fmt.Print("Send data - ")
	if err := n.Post(l, c, t); err != nil {
		fmt.Print("Error\n")

		return err
	}
	fmt.Print("OK\n")

	fmt.Println("Send part of Word of Wisdom - success")

	return nil
}
