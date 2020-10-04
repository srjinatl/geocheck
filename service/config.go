package service

import (
    "github.com/gin-gonic/gin"
    "github.com/srjinatl/geocheck/log"
)

// GeoCheckConfig info for GeoCheck service
type GeoCheckConfig struct {
    log *log.Logger
    dbDir string
    dbFileName string
    router *gin.Engine
}

// NewGeoCheckConfig creates new config for the service
func NewGeoCheckConfig() *GeoCheckConfig {
    return &GeoCheckConfig{}
}

func (c *GeoCheckConfig) WithLogger(logger *log.Logger) *GeoCheckConfig {
    c.log = logger
    return c
}

func (c *GeoCheckConfig) WithDbDir(dbDirName string) *GeoCheckConfig {
    c.dbDir = dbDirName
    return c
}

func (c *GeoCheckConfig) WithDbFileName(dbFileName string) *GeoCheckConfig {
    c.dbFileName = dbFileName
    return c
}

func (c *GeoCheckConfig) WithRouter(router *gin.Engine) *GeoCheckConfig {
    c.router = router
    return c
}

