package utils

import (
	"bufio"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
)

// create it of it doesn't exist
func OpenFile(name string) (*os.File, error) {
	// perms = os.O_APPEND|os.O_CREATE|os.O_WRONLY
	f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	return f, err
}

func AppendToFile(f *os.File, data string) error {
	defer f.Close()
	if _, err := f.WriteString(data + "\n"); err != nil {
		return err
	}

	return nil

}

func GetTitleOfLink(url string) (string, int) {
	resp, _ := http.Get(url)

	if resp.StatusCode >= 400 {
		return "", resp.StatusCode
	}

	// if resp.StatusCode == 404 {
	// 	return "", 404
	// }

	scanner := bufio.NewScanner(resp.Body)

	// title tag is always found on the first 10 lines
	var title string
	for i := 0; scanner.Scan() && i < 100; i++ {
		if strings.Index(scanner.Text(), "<title>") != -1 {
			title = scanner.Text()
			break
		}
	}

	if title == "" {
		return "", 202
	}

	idxOfTitle := strings.Index(title, "<title>")
	idxOfCloseTitle := strings.Index(title, "</title>")
	// 7 = len(<title>)
	pureTitle := title[idxOfTitle+7 : idxOfCloseTitle]

	return strings.TrimSpace(pureTitle), 200

}

// var msgtype = success|warninghttp
func PrintMsg(typ string, msg ...interface{}) {
	success := color.New(color.Bold, color.FgHiGreen).PrintlnFunc()
	error := color.New(color.Bold, color.FgHiRed).PrintlnFunc()

	if typ == "success" {
		success("\n --------------------------------------------------------------------- \n")
		success(msg...)
		success("\n --------------------------------------------------------------------- \n")
		return
	} else if typ == "error" {
		error("\n --------------------------------------------------------------------- \n")
		error(msg...)
		error("\n --------------------------------------------------------------------- \n")
		return
	}

}
