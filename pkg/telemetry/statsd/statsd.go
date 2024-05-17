// This file only wraps all the statsd client call to be used directly in the form of
// telemetry.Gauge() instead of telemetry.statsd.Gauge()
// It does NOT implement all the function from the statsd.ClientInterface interface
// because some of these are never going to be used in this application
package statsd

import (
	"time"

	"github.com/DataDog/KubeHound/pkg/config"
	"github.com/DataDog/KubeHound/pkg/telemetry/log"
	"github.com/DataDog/KubeHound/pkg/telemetry/tag"
	"github.com/DataDog/datadog-go/v5/statsd"
)

var (
	// statsdClient is the global statsd statsdClient.
	statsdClient statsd.ClientInterface
)

// just to make sure we have a client that does nothing by default
func init() {
	statsdClient = &NoopClient{}
}

func Setup(cfg *config.KubehoundConfig) error {
	statsdURL := cfg.Telemetry.Statsd.URL
	log.I.Infof("Using %s for statsd URL", statsdURL)

	var err error
	tags := tag.GetBaseTags()
	for tk, tv := range cfg.Telemetry.Tags {
		tags = append(tags, tag.MakeTag(tk, tv))
	}

	statsdClient, err = statsd.New(statsdURL,
		statsd.WithTags(tags))

	// In case we don't have a statsd url set or DD_DOGSTATSD_URL env var, we just want to continue, but log that we aren't going to submit metrics.
	if err != nil || statsdClient == nil {
		log.I.Warn("No metrics collector has been setup. All metrics submission are going to be NOOPmmm.")
		statsdClient = &NoopClient{}

		return err
	}

	return nil
}

// Count tracks how many times something happened per second.
func Count(name string, value int64, tags []string, rate float64) error {
	if statsdClient == nil {
		return nil
	}

	return statsdClient.Count(name, value, tags, rate)
}

// Gauge measures the value of a metric at a particular time.
func Gauge(name string, value float64, tags []string, rate float64) error {
	if statsdClient == nil {
		return nil
	}

	return statsdClient.Gauge(name, value, tags, rate)
}

// Incr is just Count of 1
func Incr(name string, tags []string, rate float64) error {
	if statsdClient == nil {
		return nil
	}

	return statsdClient.Incr(name, tags, rate)
}

// Decr is just Count of -1
func Decr(name string, tags []string, rate float64) error {
	if statsdClient == nil {
		return nil
	}

	return statsdClient.Decr(name, tags, rate)
}

// Histogram tracks the statistical distribution of a set of values.
func Histogram(name string, value float64, tags []string, rate float64) error {
	if statsdClient == nil {
		return nil
	}

	return statsdClient.Histogram(name, value, tags, rate)
}

// Event sends the provided Event.
func Event(event *statsd.Event) error {
	if statsdClient == nil {
		return nil
	}

	return statsdClient.Event(event)
}

// SimpleEvent sends an event with the provided title and text.
func SimpleEvent(title string, text string) error {
	if statsdClient == nil {
		return nil
	}

	return statsdClient.SimpleEvent(title, text)
}

// Set counts the number of unique elements in a group.
func Set(name string, value string, tags []string, rate float64) error {
	if statsdClient == nil {
		return nil
	}

	return statsdClient.Set(name, value, tags, rate)
}

// Timing sends timing information, it is an alias for TimeInMilliseconds
func Timing(name string, value time.Duration, tags []string, rate float64) error {
	if statsdClient == nil {
		return nil
	}

	return statsdClient.Timing(name, value, tags, rate)
}

// TimingDist sends dt in milliseconds as a distribution (p50-p99)
func TimingDist(name string, dt time.Duration, tags []string, rate float64) error {
	if statsdClient == nil {
		return nil
	}

	const secToMillis = 1000

	return statsdClient.Distribution(name, dt.Seconds()*secToMillis, tags, rate)
}

// TimeInMilliseconds sends timing information in milliseconds.
func TimeInMilliseconds(name string, value float64, tags []string, rate float64) error {
	if statsdClient == nil {
		return nil
	}

	return statsdClient.TimeInMilliseconds(name, value, tags, rate)
}

// Distribution tracks accurate global percentiles of a set of values.
func Distribution(name string, value float64, tags []string, rate float64) error {
	if statsdClient == nil {
		return nil
	}

	return statsdClient.Distribution(name, value, tags, rate)
}

// Flush flushes any pending stats in the statsd client.
func Flush() error {
	if statsdClient == nil {
		return nil
	}

	return statsdClient.Flush()
}

func IsClosed() bool {
	if statsdClient == nil {
		return false
	}

	return statsdClient.IsClosed()
}

func Close() error {
	if statsdClient == nil {
		return nil
	}

	return statsdClient.Close()
}
