package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	gpt3 "github.com/PullRequestInc/go-gpt3"
	"github.com/spf13/cobra"
)

func GetResponse(client gpt3.Client, ctx context.Context, quesiton string) {
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			quesiton,
		},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0),
	}, func(resp *gpt3.CompletionResponse) {
		fmt.Print(resp.Choices[0].Text)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(13)
	}
	fmt.Printf("\n")
}

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func main() {
	log.SetOutput(new(NullWriter))
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		panic("Missing API KEY")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)
	rootCmd := &cobra.Command{
		Use:   "chatgpt",
		Short: "Chat with ChatGPT in console.",
		Run: func(cmd *cobra.Command, args []string) {
			var question string
			if len(args) == 0 {
				b, err := io.ReadAll(os.Stdin)
				if err != nil {
					log.Fatal(err)
				}
				question = string(b)
			} else {
				question = args[0]
			}
			GetResponse(client, ctx, question)
		},
	}

	rootCmd.Execute()
}
