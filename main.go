package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/fatih/color"
	"github.com/wassimbj/lnk/utils"
)

var (
	invalidNumOfParams = errors.New("Expetced at least 2 params")
	invalidSubCmds     = errors.New("expected 'new' or 'list' subcommands")
	invalidLink        = errors.New("invalid link")
)

const (
	dataFileName = "_lnk_data.txt"
	usage        = `

	lnk <new|list> [url]

		- new: creates a new link

			url: required params, which is the actual url of the page you want to save

			e.g: lnk new https://cool-blog/go-generics-is-out-(joke)
			

		- list: Lists the saved links

			click <ENTER> to show more links
	`
)

func main() {
	args := os.Args[1:]
	var dataFilePath string
	var homeDir, err = os.UserHomeDir()
	if err != nil {
		utils.PrintMsg("Can't get the User Home Dir")
	}

	// create lnk dir in the user home dir if it does not exist
	if _, err := os.Stat(path.Join(homeDir, "lnk")); os.IsNotExist(err) {
		os.Mkdir(path.Join(homeDir, "lnk"), fs.FileMode(os.O_APPEND|os.O_RDONLY|os.O_CREATE))
	}

	dataFilePath = path.Join(homeDir, "lnk", dataFileName)

	switch args[0] {
	case "new":
		if len(args) <= 1 {
			utils.PrintMsg("error", "\t %s", invalidNumOfParams.Error())
			os.Exit(1)
		}
		link := args[1]
		if link == "" {
			color.Red("\n\n %s \n\n", invalidLink.Error())
			os.Exit(1)
		}

		var title, status = utils.GetTitleOfLink(link)

		if status == 404 {
			utils.PrintMsg("error", "\t 404 ERROR, page not found")
			os.Exit(1)
		}

		// if title is not found or there is no internet connection to make the request
		if status > 202 || title == "" {
			utils.PrintMsg("error", "\t We couldn't find the title of this link")
			color.Cyan("\n Write something and press double ENTER to move on: \n")
			fmt.Scan(&title)
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				title = title + " " + scanner.Text()
				if scanner.Text() == "" {
					break
				}
			}
		}

		f, ferr := utils.OpenFile(dataFilePath)
		if ferr != nil {
			utils.PrintMsg("error", "\t Error when opening the data file, %s", ferr.Error())
			os.Exit(1)
		}

		// the "~~" separate the link from the title
		data := link + " ~~ " + title
		saveErr := utils.AppendToFile(f, data)
		if saveErr != nil {
			utils.PrintMsg("error", "We couldn't save the link, ERROR: %s", saveErr.Error())
			os.Exit(1)
		}
		utils.PrintMsg("success", "\t Success ! Link is saved")

	case "list":
		f, _ := os.Open(dataFilePath)
		defer f.Close()
		scanner := bufio.NewScanner(f)
		var i = 0
		var max = 2

		// var moreMsg string
		color.New(color.BlinkSlow, color.Bold, color.BgHiGreen).Print("\n Saved Links \n")
		// moreMsg = "> more..."
		for i = 0; i < max && scanner.Scan(); i++ {
			// formatting the display
			if i == 0 {
				fmt.Print("\n\n")
			}
			idxOfSep := strings.Index(scanner.Text(), "~~")
			link := scanner.Text()[0 : idxOfSep-1]
			title := scanner.Text()[idxOfSep+3:]

			// like a card
			color.New(color.CrossedOut, color.FgHiBlack).Println("-----------------------------------------------------------------")
			color.New(color.Bold, color.FgHiBlue).Printf("\r [%s]", strings.TrimSpace(link))
			color.New(color.FgHiMagenta, color.Italic).Printf("\n\n * %s \n", title)
			color.New(color.CrossedOut, color.FgHiBlack).Println("-----------------------------------------------------------------")

			if i == max-1 {
				consoleReader := bufio.NewReaderSize(os.Stdin, 1)
				fmt.Print("> more...")
				input, _ := consoleReader.ReadByte()
				if input == 13 {
					max += 2
				}
			}
		}

	case "usage":
		fmt.Printf("\n %s \n", usage)

	default:
		utils.PrintMsg("error", "\t %s", invalidSubCmds.Error())
		os.Exit(1)
	}

}
