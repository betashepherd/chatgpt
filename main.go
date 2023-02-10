package main

import (
	"bufio"
	"context"
	"net/http"
	"net/url"

	//"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PullRequestInc/go-gpt3"
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
		//js, _ := json.Marshal(resp)
		//fmt.Println(string(js))
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
	apiKey := os.Getenv("GPT3_API_KEY")
	if apiKey == "" {
		panic("Missing GPT3 API KEY")
	}
	httpClient := &http.Client{
		Timeout: time.Second * 600,
	}
	sock5Proxy := os.Getenv("SOCK5_PROXY") // socks5://127.0.0.1:2080
	if sock5Proxy != "" {
		proxyUrl, _ := url.Parse(sock5Proxy)
		httpClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}
	ctx := context.Background()
	options := []gpt3.ClientOption{
		gpt3.WithDefaultEngine(gpt3.TextDavinci003Engine),
		gpt3.WithTimeout(time.Second * 600),
		gpt3.WithHTTPClient(httpClient),
	}
	client := gpt3.NewClient(apiKey, options...)
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
