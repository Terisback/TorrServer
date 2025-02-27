package pages

import (
	"github.com/anacrolix/torrent/metainfo"
	"github.com/gin-gonic/gin"

	"github.com/Terisback/TorrServer/server/settings"
	"github.com/Terisback/TorrServer/server/torr"
	"github.com/Terisback/TorrServer/server/web/auth"
	"github.com/Terisback/TorrServer/server/web/pages/template"

	"golang.org/x/exp/slices"
)

func SetupRoute(route gin.IRouter) {
	authorized := route.Group("/", auth.CheckAuth())

	webPagesAuth := route.Group("/", func() gin.HandlerFunc {
		return func(c *gin.Context) {
			if slices.Contains([]string{"/site.webmanifest"}, c.FullPath()) {
				return
			}
			auth.CheckAuth()(c)
		}
	}())

	template.RouteWebPages(webPagesAuth)
	authorized.GET("/stat", statPage)
	authorized.GET("/magnets", getTorrents)
}

// stat godoc
//
//	@Summary		Stat server
//	@Description	Stat server.
//
//	@Tags			Pages
//
//	@Produce		text/plain
//	@Success		200	"Stats"
//	@Router			/stat [get]
func statPage(c *gin.Context) {
	torr.WriteStatus(c.Writer)
	c.Status(200)
}

// getTorrents godoc
//
//	@Summary		Get HTML of magnet links
//	@Description	Get HTML of magnet links.
//
//	@Tags			Pages
//
//	@Produce		text/html
//	@Success		200	"Magnet links"
//	@Router			/magnets [get]
func getTorrents(c *gin.Context) {
	list := settings.ListTorrent()
	http := "<div>"
	for _, db := range list {
		ts := db.TorrentSpec
		mi := metainfo.MetaInfo{
			AnnounceList: ts.Trackers,
		}
		// mag := mi.Magnet(ts.DisplayName, ts.InfoHash)
		mag := mi.Magnet(&ts.InfoHash, &metainfo.Info{Name: ts.DisplayName})
		http += "<p><a href='" + mag.String() + "'>magnet:?xt=urn:btih:" + mag.InfoHash.HexString() + "</a></p>"
	}
	http += "</div>"
	c.Data(200, "text/html; charset=utf-8", []byte(http))
}
