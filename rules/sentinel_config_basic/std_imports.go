package sentinel_config_basic

var importToVersion = map[string]string{
	// These are the "base" standard imports. Although they've always existed, we set the version to 0.19.0 as that's when
	// Sentinel would error if you try to override a standard import.
	"base64":   "v0.19.0",
	"decimal":  "v0.19.0",
	"http":     "v0.19.0",
	"json":     "v0.19.0",
	"runtime":  "v0.19.0",
	"sockaddr": "v0.19.0",
	"strings":  "v0.19.0",
	"time":     "v0.19.0",
	"types":    "v0.19.0",
	"units":    "v0.19.0",
	"version":  "v0.19.0",
	// The collection based imports were introduced in 0.26.0
	"collection":       "v0.26.0",
	"collection/maps":  "v0.26.0",
	"collection/lists": "v0.26.0",
}
