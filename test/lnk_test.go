package test

import (
	"os"
	"testing"

	"github.com/wassimbj/lnk/cli"
	"github.com/wassimbj/lnk/utils"
)

const linkToCreate = "https://quickref.me"
const dataPath = "./test_lnk_data.txt"

func TestNewCmd(t *testing.T) {
	t.Run("lnk new", func(t *testing.T) {
		nerr := cli.NewLnk(linkToCreate, dataPath)
		if nerr != nil {
			t.Errorf("NewLnk Error: %s", nerr.Error())
		}
	})
}

func TestOpenFile(t *testing.T) {
	t.Run("open file", func(t *testing.T) {
		_, err := utils.OpenFile("test_lnk_data.txt", os.O_RDONLY)

		if err != nil {
			t.Errorf(err.Error())
		}
	})
}

func TestGetDataFilePath(t *testing.T) {
	t.Run("data file path", func(t *testing.T) {
		_, err := utils.GetDataFilePath("_lnk_data.txt")
		if err != nil {
			t.Errorf(err.Error())
		}
	})
}

func TestGetLinkTitle(t *testing.T) {
	t.Run("get link title", func(t *testing.T) {
		_, status := utils.GetLinkTitle(linkToCreate)

		if status == 500 {
			t.Errorf("Error can't get the title")
		}
	})
}
