package main

import (
	"fmt"
	"log"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/charmbracelet/glamour"
	readability "github.com/go-shiori/go-readability"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:  "read",
		Run:  rootCmdHandler,
		Args: cobra.ExactArgs(1),
	}

	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}

func rootCmdHandler(cmd *cobra.Command, args []string) {
	article, err := readability.FromURL(args[0], 30*time.Second)
	if err != nil {
		log.Fatalln(err)
	}

	converter := md.NewConverter("", true, nil)

	markdown, err := converter.ConvertString(article.Content)
	if err != nil {
		log.Fatal(err)
	}

	out, err := glamour.Render(markdown, "light")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(out)
}
