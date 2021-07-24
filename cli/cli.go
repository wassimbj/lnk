package cli

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/wassimbj/lnk/utils"
)

func NewLnk(link, dataFilePath string) error {
	var title, status = utils.GetLinkTitle(link)

	if status == 404 {
		return errors.New("404 ERROR, page not found")
	}

	// if title is not found or something went wrong
	if status > 202 || title == "" {
		utils.PrintMsg("error", fmt.Sprintf("Can't get the title, (HTTP ERROR) Status Code: %s", fmt.Sprint(status)+" "+http.StatusText(status)))
		color.Cyan("\n What this link is talking about ? (press Double Enter to move on): \n")
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
	saveErr := utils.AppendToFile(f, data+"\n")
	if saveErr != nil {
		return errors.New(saveErr.Error())
	}

	return nil

}
