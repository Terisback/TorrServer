package api

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/Terisback/TorrServer/server/rutor"
	"github.com/Terisback/TorrServer/server/rutor/models"
	sets "github.com/Terisback/TorrServer/server/settings"
)

// rutorSearch godoc
//
//	@Summary		Makes a rutor search
//	@Description	Makes a rutor search.
//
//	@Tags			API
//
//	@Param			query	query	string	true	"Rutor query"
//
//	@Produce		json
//	@Success		200	{array}	models.TorrentDetails	"Rutor torrent search result(s)"
//	@Router			/search [get]
func rutorSearch(c *gin.Context) {
	if !sets.BTsets.EnableRutorSearch {
		c.JSON(http.StatusBadRequest, []string{})
		return
	}
	query := c.Query("query")
	query, _ = url.QueryUnescape(query)
	list := rutor.Search(query)
	if list == nil {
		list = []*models.TorrentDetails{}
	}
	c.JSON(200, list)
}
