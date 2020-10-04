package service

import (
    "os"

    "github.com/gin-contrib/static"
    "github.com/srjinatl/geocheck/maxmind"
    "github.com/srjinatl/geocheck/service/routers"
)

type GeoCheckService struct {
    config *GeoCheckConfig
    geoLocProvider maxmind.GeoLocationProvider
}

func NewGeoCheckService(cfg *GeoCheckConfig) *GeoCheckService {
    return &GeoCheckService{config: cfg}
}

func (s *GeoCheckService) Init() (err error) {
    // create geo location provider
    currDir, _ := os.Getwd()
    s.config.log.Zap.Debug(currDir)
    prov, err := maxmind.NewGeoLocationProvider(s.config.dbDir + s.config.dbFileName)
    if err != nil {
        return
    }
    s.geoLocProvider = prov

    // initialize routers for service
    s.initEndpoints()

    return
}

func (s *GeoCheckService) Shutdown() {
    s.config.log.Zap.Info("Shutting down service...")
    // close out geo provider
    if s.geoLocProvider != nil {
        s.geoLocProvider.Close()
    }
}

func (s *GeoCheckService) initEndpoints() {
    // Swagger endpoints
    s.config.router.GET("/explorer", routers.SwaggerExplorerRedirect)
    s.config.router.GET("/swagger/api", routers.SwaggerAPI)
    s.config.router.Use(static.Serve("/explorer/", static.LocalFile("./public/swagger-ui/", true)))

    // Static webpage content endpoints
    s.config.router.LoadHTMLGlob("public/*.html")
    s.config.router.Use(static.Serve("/", static.LocalFile("./public", false)))
    s.config.router.GET("/", routers.Index)
    s.config.router.NoRoute(routers.NotFoundError)
    s.config.router.GET("/500", routers.InternalServerError)

    // main service endpoint
    handler := routers.NewGeoCheckHandler(s.config.log, s.geoLocProvider)
    s.config.router.GET( "/geocheck", handler.GeoCheck)

    // Health endpoint
    s.config.router.GET("/health", routers.HealthGET)
}