package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	apiKey := os.Getenv("GPT3_API_KEY")
	if apiKey == "" {
		panic("Missing GPT3 API KEY")
	}
	config := openai.DefaultConfig(apiKey)
	httpClient := &http.Client{
		Timeout: time.Second * 600,
	}
	sock5Proxy := os.Getenv("SOCK5_PROXY")
	if sock5Proxy != "" {
		proxyUrl, _ := url.Parse(sock5Proxy)
		httpClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}
	config.HTTPClient = httpClient
	client := openai.NewClientWithConfig(config)
	ctx := context.Background()

	rootCmd := &cobra.Command{
		Use:   "chatgpt",
		Short: "Chat with ChatGPT in console.",
		Run: func(cmd *cobra.Command, args []string) {
			scanner := bufio.NewScanner(os.Stdin)
			quit := false

			for !quit {
				fmt.Print("輸入你的問題(quit 離開): ")
				if !scanner.Scan() {
					break
				}
				question := scanner.Text()
				questionParam := validateQuestion(question)
				switch questionParam {
				case "quit":
					quit = true
				case "":
					continue

				default:
					GetResponse(client, ctx, questionParam)
				}
			}
		},
	}

	log.Fatal(rootCmd.Execute())
}

func GetResponse(client *openai.Client, ctx context.Context, quesiton string) {
	req := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		MaxTokens:   3000,
		Temperature: 0,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: quesiton,
			},
		},
		Stream: true,
	}
	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			//fmt.Println("\nStream finished")
			fmt.Println("\n")
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}

		fmt.Printf(response.Choices[0].Delta.Content)
	}
}

func validateQuestion(question string) string {
	quest := strings.Trim(question, " ")
	keywords := []string{"", "loop", "break", "continue", "cls", "exit", "block"}
	for _, x := range keywords {
		if quest == x {
			return ""
		}
	}
	return quest
}
