package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Terisback/TorrServer/server/torr"
	"github.com/Terisback/TorrServer/server/torr/state"
	"github.com/Terisback/TorrServer/server/web/api/utils"
)

// play godoc
//
//	@Summary		Play given torrent referenced by hash
//	@Description	Play given torrent referenced by hash.
//
//	@Tags			API
//
//	@Param			hash		path	string	true	"Torrent hash"
//	@Param			id			path	string	true	"File index in torrent"
//
//	@Produce		application/octet-stream
//	@Success		200	"Torrent data"
//	@Router			/play/{hash}/{id} [get]
func play(c *gin.Context) {
	hash := c.Param("hash")
	indexStr := c.Param("id")
	notAuth := c.GetBool("auth_required") && c.GetString(gin.AuthUserKey) == ""

	if hash == "" || indexStr == "" {
		c.AbortWithError(http.StatusNotFound, errors.New("link should not be empty"))
		return
	}

	spec, err := utils.ParseLink(hash)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tor := torr.GetTorrent(spec.InfoHash.HexString())
	if tor == nil && notAuth {
		c.Header("WWW-Authenticate", "Basic realm=Authorization Required")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if tor == nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New("error get torrent"))
		return
	}

	if tor.Stat == state.TorrentInDB {
		tor, err = torr.AddTorrent(spec, tor.Title, tor.Poster, tor.Data)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	if !tor.GotInfo() {
		c.AbortWithError(http.StatusInternalServerError, errors.New("timeout connection torrent"))
		return
	}

	// find file
	index := -1
	if len(tor.Files()) == 1 {
		index = 1
	} else {
		ind, err := strconv.Atoi(indexStr)
		if err == nil {
			index = ind
		}
	}
	if index == -1 { // if file index not set and play file exec
		c.AbortWithError(http.StatusBadRequest, errors.New("\"index\" is wrong"))
		return
	}

	tor.Stream(index, c.Request, c.Writer)
}
