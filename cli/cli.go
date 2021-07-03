package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/wassimbj/lnk/utils"
)

func NewLnk(link, dataFilePath string) error {
	var title, status = utils.GetLinkTitle(link)

	if status == 404 {
		return errors.New("404 ERROR, page not found")
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

	f, ferr := utils.OpenFile(dataFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE)
	if ferr != nil {
		return errors.New(ferr.Error())
	}

	// the "~~" separate the link from the title
	data := link + " ~~ " + title
	saveErr := utils.AppendToFile(f, data)
	if saveErr != nil {
		return errors.New(saveErr.Error())
	}

	return nil

}
