/*
Package gdw implements a wrapper for Google Driwe API for Go (Golang).
*/
package gdw

import (
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
)

//SendByUrl download file from specified url and send it to Google Drive.
func SendByUrl(u string, f *drive.File, config *oauth2.Config, token *oauth2.Token) (file *drive.File, err error) {
	response, err := http.Get(u)
	if err != nil {
		return
	}
	defer response.Body.Close()

	ctx := context.Background()
	client := config.Client(ctx, token)
	srv, err := drive.New(client)
	if err != nil {
		return
	}

	file, err = srv.Files.Create(f).Media(response.Body).Do()
	return
}
