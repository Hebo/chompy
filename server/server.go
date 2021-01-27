package server

import (
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hebo/chompy/config"
	"github.com/hebo/chompy/downloader"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server handles HTTP routes
type Server struct {
	router          *echo.Echo
	downloadsDir    string
	downloader      downloader.Downloader
	playlistSyncURL string
}

const videosIndexPath = "/videos"

// New creates a new Server
func New(cfg config.Config) Server {
	srv := Server{
		downloadsDir:    cfg.DownloadsDir,
		playlistSyncURL: cfg.PlaylistSyncURL,
		downloader:      downloader.New(cfg.DownloadsDir, cfg.Format),
	}

	t := &tmpl{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	// Echo instance
	e := echo.New()
	e.Renderer = t

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())

	// Routes
	e.Static("/assets", "public/assets")

	e.GET("/", srv.index)
	e.GET(videosIndexPath+"/", srv.videosList)
	e.POST("/download", srv.downloadVideo)

	fs := http.FileServer(http.Dir(srv.downloadsDir))
	e.GET(videosIndexPath+"/*", echo.WrapHandler(http.StripPrefix(videosIndexPath, fs)))
	e.GET(videosIndexPath, func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, videosIndexPath+"/")
	})

	srv.router = e
	return srv
}

type tmpl struct {
	templates *template.Template
}

func (t *tmpl) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// Serve starts the HTTP server and background tasks
func (s Server) Serve(port int) {
	if err := s.startWorkers(); err != nil {
		s.router.Logger.Fatal("Failed to start tasks:", err)
	}

	portString := ":" + strconv.Itoa(port)
	s.router.Logger.Fatal(s.router.Start(portString))
}

func (s *Server) index(c echo.Context) error {
	return c.String(http.StatusOK, "Chompy is ready to eat!")
}

type videoFile struct {
	Filename string
	Created  time.Time
}

func (s *Server) videosList(c echo.Context) error {
	files, err := ioutil.ReadDir(s.downloadsDir)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var vids []videoFile
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		vids = append(vids, videoFile{Filename: file.Name(), Created: file.ModTime()})
	}

	sort.Slice(vids, func(i, j int) bool { return vids[i].Created.After(vids[j].Created) })
	return c.Render(http.StatusOK, "videos_list.html", vids)
}

type downloadRequest struct {
	URL    string `json:"url" form:"url"`
	Format string `json:"format" form:"format"`
}

type downloadResponse struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
}

func (s *Server) downloadVideo(c echo.Context) error {
	req := new(downloadRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if req.URL == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing url")
	}

	filePath, err := s.downloader.Download(req.URL, req.Format)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	_, filename := path.Split(filePath)
	res := downloadResponse{
		Filename: filename,
		Path:     path.Join(videosIndexPath, filename),
	}

	return c.JSON(http.StatusOK, res)
}
