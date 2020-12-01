package server

import (
	"net/http"
	"path"
	"strconv"

	"github.com/hebo/chompy/downloader"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server handles RSS and other HTTP routes
type Server struct {
	router       *echo.Echo
	downloadsDir string
	downloader   downloader.Downloader
}

const videosIndexPath = "/videos"

// New creates a new Server
func New(downloadsDir string) Server {
	srv := Server{
		downloadsDir: downloadsDir,
		downloader:   downloader.New(downloadsDir),
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", index)
	e.POST("/download", srv.downloadVideo)

	fs := http.FileServer(http.Dir(downloadsDir))
	e.GET(videosIndexPath+"/*", echo.WrapHandler(http.StripPrefix(videosIndexPath, fs)))

	srv.router = e
	return srv
}

// Serve starts the HTTP server
func (s Server) Serve(port int) {
	portString := ":" + strconv.Itoa(port)
	s.router.Logger.Fatal(s.router.Start(portString))
}

func index(c echo.Context) error {
	return c.String(http.StatusOK, "Chompy is ready to eat!")
}

type downloadRequest struct {
	URL string `json:"url" form:"url"`
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

	filePath, err := s.downloader.Download(req.URL)
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
