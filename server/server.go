package server

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/hebo/chompy/config"
	"github.com/hebo/chompy/downloader"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/afero"
)

// Server handles HTTP routes
type Server struct {
	fs              afero.Fs
	router          *echo.Echo
	downloadsDir    string
	downloader      downloader.Downloader
	playlistSyncURL string
	maxSize         int
	cleanup         chan (struct{})
}

const videosIndexPath = "/videos"

// New creates a new Server
func New(cfg config.Config, fs afero.Fs) Server {
	srv := Server{
		fs:              fs,
		downloadsDir:    cfg.DownloadsDir,
		playlistSyncURL: cfg.PlaylistSyncURL,
		maxSize:         cfg.MaxSize,
		cleanup:         make(chan struct{}),
	}

	srv.downloader = downloader.New(cfg.DownloadsDir, cfg.Format, srv.triggerCleanup)

	funcMap := template.FuncMap{
		"escape":    url.PathEscape,
		"humanTime": humanize.Time,
	}

	t := &tmpl{
		templates: template.Must(template.New("main").Funcs(funcMap).ParseGlob("public/views/*.html")),
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
	e.GET("/download", srv.downloadVideo)
	e.POST("/download", srv.downloadVideo)
	e.POST("/touch", srv.touch)

	fSrv := http.FileServer(http.Dir(srv.downloadsDir))
	e.GET(videosIndexPath+"/*", echo.WrapHandler(http.StripPrefix(videosIndexPath, fSrv)))
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
	return c.Redirect(http.StatusTemporaryRedirect, "/videos")
}

type videoFile struct {
	Filename string
	Created  time.Time
	Size     int64
}

func (s *Server) videosList(c echo.Context) error {
	vids, err := getVideoFiles(s.downloadsDir, createdDesc)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

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

func (s *Server) touch(c echo.Context) error {
	filename := c.FormValue("filename")
	log.Printf("filename: %v", filename)

	err := touch(s.downloadsDir, filename)
	if os.IsNotExist(err) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid filename")
	} else if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}
