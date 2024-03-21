package main

import (
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/glamour"
	readability "github.com/go-shiori/go-readability"
	"github.com/spf13/cobra"
)

var NoPager, NoLinks, Raw bool

func main() {
	rootCmd := &cobra.Command{
		Use:  "rdr [url]",
		Run:  rootCmdHandler,
		Args: cobra.MaximumNArgs(1),
	}

	rootCmd.Flags().BoolVar(&NoPager, "no-pager", false, "Don't pipe output to a pager")
	rootCmd.Flags().BoolVar(&NoLinks, "no-links", false, "Don't display any links")
	rootCmd.Flags().BoolVar(&Raw, "raw", false, "Just raw")

	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

func rootCmdHandler(cmd *cobra.Command, args []string) {
	var article readability.Article
	var err error

	if len(args) > 0 && args[0] != "-" {
		article, err = readability.FromURL(args[0], 30*time.Second)
	} else {
		var inputReader io.Reader = cmd.InOrStdin()
		url, err1 := url.Parse("https://example.com")
		if err1 != nil {
			log.Fatalf("failed to parse url: %v", err1)
		}
		article, err = readability.FromReader(inputReader, url)
	}
	if err != nil {
		log.Fatalln(err)
	}

	converter := md.NewConverter("", true, nil)
	converter.Remove("figure")

	if NoLinks {
		converter.AddRules(
			md.Rule{
				Filter: []string{"a"},
				Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
					return md.String(content)
				},
			},
		)
	}

	markdown, err := converter.ConvertString(article.Content)
	if err != nil {
		log.Fatal(err)
	}

	if Raw {
		os.Stdout.WriteString(markdown)
		return
	}

	out, err := glamour.RenderWithEnvironmentConfig(markdown)

	if err != nil {
		log.Fatal(err)
	}

	if NoPager {
		os.Stdout.WriteString(out)
		return
	}

	f, err := os.CreateTemp("", `rdr-page-*`)
	if err != nil {
		log.Fatal(err)
	}
	name := f.Name()
	_, err = f.WriteString(out)
	defer f.Close()
	defer os.Remove(name)
	if err != nil {
		log.Fatal(err)
	}
	execCmd := exec.Command("less", "-R", name)
	execCmd.Stdout = os.Stdout
	execCmd.Stdin = os.Stdin
	execCmd.Stderr = os.Stderr
	execCmd.Run()
}
