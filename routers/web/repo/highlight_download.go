package repo

import (
<<<<<<< HEAD
=======
	"html/template"
>>>>>>> 75358a09f8 (main 최신화 (#113))
	"io"
	"net/http"
	"unicode/utf8"

<<<<<<< HEAD
	"html/template"

=======
>>>>>>> 75358a09f8 (main 최신화 (#113))
	"code.gitea.io/gitea/modules/highlight"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/services/context"
)

type HighlightResponse struct {
	Html  []template.HTML `json:"html"`
	IsBin bool            `json:"isBin"`
}

func SingleHighlightDownload(ctx *context.Context) {
	blob, _ := getBlobForEntry(ctx)
	if blob == nil {
		ctx.JSONErrorf("blob is null")
		return
	}
	dataRc, err := blob.DataAsync()
	if err != nil {
		ctx.JSONErrorf("cannot retreiving datarc: %v", err)
		return
	}
	defer func() {
		if err = dataRc.Close(); err != nil {
			log.Error("blob close err: %v", err)
		}
	}()
	rawCode, err := io.ReadAll(dataRc)
	if err != nil {
		ctx.JSONErrorf("error on io.readall : %v", err)
		return
	}

	isText := utf8.ValidString(string(rawCode))
	if !isText {
		ctx.JSON(http.StatusOK, HighlightResponse{IsBin: true})
		return
	}

	htmlCode, _, err := highlight.File(blob.Name(), "", rawCode)
	if err != nil {
		ctx.JSONErrorf("error on highlight %v", err)
		return
	}

	ctx.JSON(http.StatusOK, HighlightResponse{Html: htmlCode, IsBin: false})
}
