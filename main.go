package main

import (
	"log"
	"time"
	"os"
	"os/exec"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/glamour"
	readability "github.com/go-shiori/go-readability"
	"github.com/spf13/cobra"
)

var NoPager, NoLinks bool

func main() {
	rootCmd := &cobra.Command{
		Use:  "rdr [url]",
		Run:  rootCmdHandler,
		Args: cobra.ExactArgs(1),
	}

	rootCmd.Flags().BoolVar(&NoPager, "no-pager", false, "Don't pipe output to a pager")
	rootCmd.Flags().BoolVar(&NoLinks, "no-links", false, "Don't display any links")

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

	out, err := glamour.Render(markdown, "light")

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
