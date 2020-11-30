package server

import (
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/hebo/chompy/downloader"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Server handles RSS and other HTTP routes
type Server struct {
	router       *echo.Echo
	downloadsDir string
	downloader   downloader.Downloader
}

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
	e.GET("/videos", srv.getVideo)
	e.POST("/download", srv.downloadVideo)

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
	}

	return c.JSON(http.StatusOK, res)
}

func (s *Server) getVideo(c echo.Context) error {
	filename := c.QueryParam("filename")
	if filename == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing filename")
	}

	path := path.Join(s.downloadsDir, filename)
	if !fileExists(path) {
		return echo.NewHTTPError(http.StatusNotFound, "no file found")
	}

	http.ServeFile(c.Response().Writer, c.Request(), path)
	return nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
