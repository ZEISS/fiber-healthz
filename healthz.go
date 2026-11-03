// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber

package healthz

import (
	"github.com/gofiber/fiber/v3"
)

// response represents the response to a health check.
type response struct {
	Status string `json:"status" xml:"status" msgpack:"status" cbor:"status"`
}

const (
	defaultSuccessMessage = "OK"
	defaultFailureMessage = "Service Unavailable"
	defaultSuccessStatus  = fiber.StatusOK
	defaultFailureStatus  = fiber.StatusServiceUnavailable
)

const (
	// FormatText returns a plain text response (default behavior).
	FormatText int = iota
	// FormatJSON returns a JSON response.
	FormatJSON
	// FormatXML returns an XML response.
	FormatXML
	// FormatMsgPack returns a MessagePack response.
	FormatMsgPack
	// FormatCBOR returns a CBOR response.
	FormatCBOR
)

const (
	// Liveness is the conventional path for a liveness check.
	// Register the middleware on this path to expose it.
	Liveness = "/livez"
	// Readiness is the conventional path for a readiness check.
	// Register the middleware on this path to expose it.
	Readiness = "/readyz"
	// Startup is the conventional path for a startup check.
	// Register the middleware on this path to expose it.
	Startup = "/startupz"
	// Healthz is the conventional path for a health check.
	// Register the middleware on this path to expose it.
	Healthz = "/healthz"
)

// ProbeFunc defines a function to probe the health of the application.
type ProbeFunc func(c fiber.Ctx) bool

// Success is a probe function that always returns true.
func Success(fiber.Ctx) bool {
	return true
}

// Failure is a probe function that always returns false.
func Failure(fiber.Ctx) bool {
	return false
}

// Config caputes the configuration for running the goth middleware.
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	Next func(c fiber.Ctx) bool
	// Probe defines a function to probe the health of the application.
	Probe ProbeFunc
	// ResponseFormat defines the response format for the health check.
	ResponseFormat int
	// SuccessMessage defines the message to return when the probe succeeds.
	SuccessMessage string
	// SuccessStatus defines the status code to return when the probe succeeds.
	SuccessStatus int
	// FailureMessage defines the message to return when the probe fails.
	FailureMessage string
	// FailureStatus defines the status code to return when the probe fails.
	FailureStatus int
}

// DefaultConfig returns the default configuration for the healthz middleware.
var DefaultConfig = Config{
	Probe:          Success,
	ResponseFormat: FormatText,
	SuccessStatus:  defaultSuccessStatus,
	SuccessMessage: defaultSuccessMessage,
	FailureMessage: defaultFailureMessage,
	FailureStatus:  defaultFailureStatus,
}

// New returns a new instance of the goth middleware with the given configuration.
func New(config ...Config) fiber.Handler {
	cfg := configDefault(config...)

	return func(c fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		if c.Method() != fiber.MethodGet && c.Method() != fiber.MethodHead {
			return c.Next()
		}

		ok := cfg.Probe(c)
		statusCode := cfg.SuccessStatus
		statusMsg := cfg.SuccessMessage

		if !ok {
			statusCode = cfg.FailureStatus
			statusMsg = cfg.FailureMessage
		}

		c.Status(statusCode)

		switch cfg.ResponseFormat {
		case FormatJSON:
			return c.JSON(response{Status: statusMsg})
		case FormatXML:
			return c.XML(response{Status: statusMsg})
		case FormatMsgPack:
			return c.MsgPack(response{Status: statusMsg})
		case FormatCBOR:
			return c.CBOR(response{Status: statusMsg})
		default: // FormatText
			return c.SendString(statusMsg)
		}
	}
}

// NewProbe returns a new healthz middleware with the given probe function.
func NewProbe(probe ProbeFunc) fiber.Handler {
	return New(Config{Probe: probe})
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return DefaultConfig
	}

	// Override default config
	cfg := config[0]

	if cfg.Next == nil {
		cfg.Next = DefaultConfig.Next
	}

	if cfg.Probe == nil {
		cfg.Probe = DefaultConfig.Probe
	}

	if cfg.ResponseFormat == 0 {
		cfg.ResponseFormat = DefaultConfig.ResponseFormat
	}

	if cfg.SuccessMessage == "" {
		cfg.SuccessMessage = DefaultConfig.SuccessMessage
	}

	if cfg.SuccessStatus == 0 {
		cfg.SuccessStatus = DefaultConfig.SuccessStatus
	}

	if cfg.FailureMessage == "" {
		cfg.FailureMessage = DefaultConfig.FailureMessage
	}

	if cfg.FailureStatus == 0 {
		cfg.FailureStatus = DefaultConfig.FailureStatus
	}

	return cfg
}
