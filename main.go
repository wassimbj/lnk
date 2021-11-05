package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/wassimbj/lnk/cli"
	"github.com/wassimbj/lnk/utils"
)

var (
	invalidNumOfParams = errors.New("Expetced at least 2 params")
	invalidSubCmds     = errors.New("invalid subcommands")
	invalidLink        = errors.New("invalid link")
)

const (
	dataFileName = "_lnk_data.txt"
	limit        = 4
	separator    = "~~"
	usage        = `
lnk <new | list | rm> [url|ID]

	* new <url>: creates a new link

	* list: Lists the saved links

	* rm <ID>: remove a saved link

`
)

func main() {
	args := os.Args[1:]

	dataFilePath, dpError := utils.GetDataFilePath(dataFileName)
	if dpError != nil {
		utils.PrintMsg(dpError.Error())
		os.Exit(1)
	}

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
		err := cli.NewLnk(link, dataFilePath)

		if err != nil {
			utils.PrintMsg("error", err.Error())
			os.Exit(1)
		}

		utils.PrintMsg("success", "\t Success ! Link is saved")

	case "list":
		f, _ := utils.OpenFile(dataFilePath, os.O_RDONLY)
		defer f.Close()
		scanner := bufio.NewScanner(f)
		var i = 0

		color.New(color.BlinkSlow, color.Bold, color.BgHiGreen).Print("\n Saved Links \n")
		for i = 0; scanner.Scan(); i++ {
			// formatting the display
			if i == 0 {
				fmt.Print("\n\n")
			}

			// separate the link and title (link description)
			idxOfSep := strings.Index(scanner.Text(), separator)
			link := scanner.Text()[0 : idxOfSep-1]
			title := scanner.Text()[idxOfSep+3:]

			// like a card
			color.New(color.CrossedOut, color.FgHiBlack).Println("-----------------------------------------------------------------")
			color.New(color.Bold, color.FgHiBlue).Printf("\r %d - %s", i+1, strings.TrimSpace(link))
			color.New(color.FgHiMagenta, color.Italic).Printf("\n\n * %s \n", title)
			color.New(color.CrossedOut, color.FgHiBlack).Println("-----------------------------------------------------------------")

		}

	case "rm":
		if len(args) <= 1 {
			utils.PrintMsg("error", "\t %s", invalidNumOfParams.Error())
			os.Exit(1)
		}

		linkId := args[1]

		fileContent, _ := ioutil.ReadFile(dataFilePath)
		deletedLink := ""
		// find the link to delete and remove it from the array
		dataArr := strings.Split(string(fileContent), "\n")
		for i := 0; i < len(dataArr); i++ {
			linkIdInt, _ := strconv.Atoi(linkId)
			if i+1 == linkIdInt {
				deletedLink = strings.Split(dataArr[i], separator)[0]
				dataArr = append(dataArr[:i], dataArr[i+1:]...)
			} else if i == len(dataArr) {
				utils.PrintMsg("error", "the id doesn't exist, try lnk list to see links with its id")
			}
		}

		// re-write the file with the new data
		newData := strings.Join(dataArr, "\n")
		f, _ := utils.OpenFile(dataFilePath, os.O_WRONLY|os.O_TRUNC)

		if err := utils.AppendToFile(f, newData); err != nil {
			utils.PrintMsg("error", "Can't delete the link, Error: "+err.Error())
		}

		utils.PrintMsg("success", "\t Removed ! ["+strings.TrimSpace(deletedLink)+"]")

	case "usage":
		fmt.Printf("\n %s \n", usage)

	default:
		utils.PrintMsg("error", "\t", invalidSubCmds.Error())
		os.Exit(1)
	}

}
