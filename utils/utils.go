package utils

import (
	"bufio"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/fatih/color"
)

// create it of it doesn't exist
func OpenFile(name string, perms int) (*os.File, error) {
	// perms = os.O_APPEND|os.O_CREATE|os.O_WRONLY
	// os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_RDONLY
	f, err := os.OpenFile(name, perms, 0644)

	return f, err
}

func AppendToFile(f *os.File, data string) error {
	defer f.Close()
	if _, err := f.WriteString(data + "\n"); err != nil {
		return err
	}

	return nil

}

func GetLinkTitle(url string) (string, int) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 500
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", resp.StatusCode
	}

	scanner := bufio.NewScanner(resp.Body)

	// searching the first 100 lines for the title
	var title string
	for i := 0; scanner.Scan() && i < 100; i++ {
		// fmt.Println(scanner.Text())
		if strings.Index(scanner.Text(), "<title") != -1 {
			title = scanner.Text()
			break
		}
	}
	if title == "" {
		return "", 202
	}

	idxOfCloseTitle := strings.Index(title, "</title>")
	// <title> can have attrs so we are only searching for the first letters of the tag
	idxOfTitle := strings.Index(title, "<title")
	idxOfOnlyTitle := strings.Index(title[idxOfTitle:idxOfCloseTitle], ">") + 1

	pureTitle := title[idxOfTitle+idxOfOnlyTitle : idxOfCloseTitle]

	return strings.TrimSpace(pureTitle), 200

}

func GetDataFilePath(dataFileName string) (string, error) {
	var dataFilePath string
	var homeDir, err = os.UserHomeDir()
	if err != nil {
		return "", errors.New("Can't get the User Home Dir")
	}

	// create lnk dir in the user home dir if it does not exist
	if _, err := os.Stat(path.Join(homeDir, "lnk")); os.IsNotExist(err) {
		mkdirErr := os.Mkdir(path.Join(homeDir, "lnk"), fs.FileMode(os.O_APPEND|os.O_RDONLY|os.O_CREATE))

		if mkdirErr != nil {
			return "", mkdirErr
		}
	}

	dataFilePath = path.Join(homeDir, "lnk", dataFileName)

	return dataFilePath, nil

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
