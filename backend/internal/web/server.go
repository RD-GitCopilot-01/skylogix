package web

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	config     Config
	httpServer *http.Server
	router     Router
	middleware []Middleware
}

type Config struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	TLS          TLSConfig     `json:"tls"`
	CORS         CORSConfig    `json:"cors"`
}

type TLSConfig struct {
	Enabled  bool   `json:"enabled"`
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}

type CORSConfig struct {
	Enabled        bool     `json:"enabled"`
	AllowedOrigins []string `json:"allowed_origins"`
	AllowedMethods []string `json:"allowed_methods"`
	AllowedHeaders []string `json:"allowed_headers"`
}

type Router interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	GET(path string, handler HandlerFunc)
	POST(path string, handler HandlerFunc)
	PUT(path string, handler HandlerFunc)
	DELETE(path string, handler HandlerFunc)
	Use(middleware ...Middleware)
	Group(prefix string) RouterGroup
}

type RouterGroup interface {
	GET(path string, handler HandlerFunc)
	POST(path string, handler HandlerFunc)
	PUT(path string, handler HandlerFunc)
	DELETE(path string, handler HandlerFunc)
	Use(middleware ...Middleware)
	Group(prefix string) RouterGroup
}

type HandlerFunc func(c Context) error

type Context interface {
	Request() *http.Request
	Response() ResponseWriter
	Param(key string) string
	Query(key string) string
	Bind(obj interface{}) error
	JSON(code int, obj interface{}) error
	String(code int, format string, values ...interface{}) error
	File(filepath string) error
	Set(key string, value interface{})
	Get(key string) interface{}
}

type ResponseWriter interface {
	http.ResponseWriter
	Status() int
	Size() int
	Written() bool
}

type Middleware func(HandlerFunc) HandlerFunc

func NewServer(config Config) *Server {
	return &Server{
		config:     config,
		router:     NewGinRouter(),
		middleware: make([]Middleware, 0),
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
	}
	
	fmt.Printf("Starting web server on %s\n", addr)
	
	if s.config.TLS.Enabled {
		return s.httpServer.ListenAndServeTLS(s.config.TLS.CertFile, s.config.TLS.KeyFile)
	}
	
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	
	fmt.Println("Stopping web server...")
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Use(middleware ...Middleware) {
	s.middleware = append(s.middleware, middleware...)
	for _, m := range middleware {
		s.router.Use(m)
	}
}

func (s *Server) GET(path string, handler HandlerFunc) {
	s.router.GET(path, handler)
}

func (s *Server) POST(path string, handler HandlerFunc) {
	s.router.POST(path, handler)
}

func (s *Server) PUT(path string, handler HandlerFunc) {
	s.router.PUT(path, handler)
}

func (s *Server) DELETE(path string, handler HandlerFunc) {
	s.router.DELETE(path, handler)
}

func (s *Server) Group(prefix string) RouterGroup {
	return s.router.Group(prefix)
}

func (s *Server) SetupRoutes(handlers *Handlers) {
	if s.config.CORS.Enabled {
		s.Use(CORSMiddleware(s.config.CORS))
	}
	
	s.Use(LoggingMiddleware())
	
	s.Use(RecoveryMiddleware())
	
	s.GET("/health", handlers.HealthCheck)
	
	api := s.Group("/api/v1")
	
	storage := api.Group("/storage")
	storage.GET("/backends", handlers.GetStorageBackends)
	storage.POST("/backends", handlers.CreateStorageBackend)
	storage.GET("/backends/:name", handlers.GetStorageBackend)
	storage.PUT("/backends/:name", handlers.UpdateStorageBackend)
	storage.DELETE("/backends/:name", handlers.DeleteStorageBackend)
	storage.GET("/backends/:name/health", handlers.GetStorageBackendHealth)
	storage.GET("/backends/:name/metrics", handlers.GetStorageBackendMetrics)
	storage.POST("/backends/:name/test", handlers.TestStorageBackend)
	
	storage.POST("/migrate", handlers.MigrateData)
	storage.GET("/metrics", handlers.GetAggregatedMetrics)
	storage.GET("/health", handlers.GetAggregatedHealth)
	
	config := api.Group("/config")
	config.GET("/", handlers.GetConfiguration)
	config.PUT("/", handlers.UpdateConfiguration)
	config.POST("/backup", handlers.BackupConfiguration)
	config.POST("/restore", handlers.RestoreConfiguration)
	config.GET("/history", handlers.GetConfigurationHistory)
	
	snmp := api.Group("/snmp")
	snmp.GET("/status", handlers.GetSNMPStatus)
	snmp.POST("/start", handlers.StartSNMP)
	snmp.POST("/stop", handlers.StopSNMP)
	snmp.GET("/mib", handlers.GetMIBObjects)
	snmp.POST("/trap", handlers.SendSNMPTrap)
	
	database := api.Group("/database")
	database.GET("/status", handlers.GetDatabaseStatus)
	database.GET("/stats", handlers.GetDatabaseStats)
	database.POST("/backup", handlers.BackupDatabase)
	database.POST("/restore", handlers.RestoreDatabase)
	database.POST("/cleanup", handlers.CleanupDatabase)
	
	monitoring := api.Group("/monitoring")
	monitoring.GET("/alerts", handlers.GetAlerts)
	monitoring.POST("/alerts/:id/acknowledge", handlers.AcknowledgeAlert)
	monitoring.GET("/metrics/storage", handlers.GetStorageMetricsHistory)
	monitoring.GET("/metrics/health", handlers.GetHealthHistory)
	
	audit := api.Group("/audit")
	audit.GET("/logs", handlers.GetAuditLogs)
	
	s.GET("/ws", handlers.WebSocketHandler)
	
	s.router.Use(StaticFileMiddleware("/", "./web/dist"))
}

func CORSMiddleware(config CORSConfig) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			origin := c.Request().Header.Get("Origin")
			
			allowed := false
			for _, allowedOrigin := range config.AllowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					allowed = true
					break
				}
			}
			
			if allowed {
				c.Response().Header().Set("Access-Control-Allow-Origin", origin)
			}
			
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
			
			if c.Request().Method == "OPTIONS" {
				return c.String(http.StatusOK, "")
			}
			
			return next(c)
		}
	}
}

func LoggingMiddleware() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			start := time.Now()
			
			err := next(c)
			
			duration := time.Since(start)
			status := c.Response().Status()
			
			fmt.Printf("[%s] %s %s %d %v\n",
				start.Format("2006-01-02 15:04:05"),
				c.Request().Method,
				c.Request().URL.Path,
				status,
				duration,
			)
			
			return err
		}
	}
}

func RecoveryMiddleware() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Panic recovered: %v\n", r)
					c.JSON(http.StatusInternalServerError, map[string]interface{}{
						"error": "Internal server error",
					})
				}
			}()
			
			return next(c)
		}
	}
}

func StaticFileMiddleware(prefix, root string) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			return next(c)
		}
	}
}
