package repo

import (
	"io"
	"net/http"

	"html/template"

	"code.gitea.io/gitea/modules/highlight"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/services/context"
)

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
	htmlCode, _, err := highlight.File(blob.Name(), "", rawCode)
	if err != nil {
		ctx.JSONErrorf("error on highlight %v", err)
		return
	}

	response := struct {
		Html []template.HTML `json:"html"`
	}{
		Html: htmlCode,
	}
	ctx.JSON(http.StatusOK, response)
}
