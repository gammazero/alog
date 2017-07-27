package alog

// ****************************************************************************
// Commented out to prevent this package from bringing in unneeded
// dependencies.  Uncomment to run benchmark.
//

/*
import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"testing"
	"time"
)

// smallFields is a small size data set for benchmarking
var smallFields = map[string]interface{}{
	"foo":   "bar",
	"baz":   "qux",
	"one":   "two",
	"three": "four",
}

func BenchmarkLogger1(b *testing.B) {
	logger := NewText(ioutil.Discard, InfoLevel, time.RFC3339, "")
	entry := logger.WithFields(Fields(smallFields))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			entry.Info("aaa")
		}
	})
}

func BenchmarkLogrus1(b *testing.B) {
	logger := logrus.Logger{
		Out:       ioutil.Discard,
		Level:     logrus.InfoLevel,
		Formatter: &logrus.TextFormatter{DisableColors: true},
	}
	logger.SetNoLock()
	entry := logger.WithFields(logrus.Fields(smallFields))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			entry.Info("aaa")
		}
	})
}

func BenchmarkLogger2(b *testing.B) {
	lg := NewText(ioutil.Discard, InfoLevel, time.RFC3339, "")
	lg2 := NewText(ioutil.Discard, InfoLevel, time.RFC3339, "")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ent := lg.WithFields(map[string]interface{}{
			"cat":   "mammal",
			"whale": "mammal",
			"moose": "mammal",
		})
		ent.WithFields(map[string]interface{}{
			"shark":  "fish",
			"guppy":  "fish",
			"batray": "fish",
		}).Info("Animal Types")
		lg.Info("for", "score", "and", "seven", "years", "ago")

		lg2.Warn("Hello world")
		lg2.Errorf("Numbers: %d %d %d", 1, 2, 3)
		lg2.WithFields(map[string]interface{}{
			"foo":   "bar",
			"count": 42,
		}).Error("Some Fields")
		ent = lg2.WithFields(map[string]interface{}{
			"cat":   "mammal",
			"whale": "mammal",
			"moose": "mammal",
		})
		ent.WithFields(map[string]interface{}{
			"shark":  "fish",
			"guppy":  "fish",
			"batray": "fish",
		}).Info("Animal Types")
		lg2.Info("for", "score", "and", "seven", "years", "ago")
		lg.WithField("color", "red").Warn("favorite color")
	}
	lg.Close()
}

func BenchmarkLogrus2(b *testing.B) {
	lg := logrus.New()
	lg.SetNoLock()
	lg.Formatter = &logrus.JSONFormatter{TimestampFormat: time.RFC3339}
	lg.Out = ioutil.Discard

	lg2 := logrus.New()
	lg.SetNoLock()
	lg2.Formatter = &logrus.JSONFormatter{TimestampFormat: time.RFC3339}
	lg2.Out = ioutil.Discard

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lg.Warn("Hello world")
		lg.Errorf("Numbers: %d %d %d", 1, 2, 3)
		lg.WithFields(map[string]interface{}{
			"foo":   "bar",
			"count": 42,
		}).Error("Some Fields")
		ent := lg.WithFields(map[string]interface{}{
			"cat":   "mammal",
			"whale": "mammal",
			"moose": "mammal",
		})
		ent.WithFields(map[string]interface{}{
			"shark":  "fish",
			"guppy":  "fish",
			"batray": "fish",
		}).Info("Animal Types")
		lg.Info("for", "score", "and", "seven", "years", "ago")

		lg2.Warn("Hello world")
		lg2.Errorf("Numbers: %d %d %d", 1, 2, 3)
		lg2.WithFields(map[string]interface{}{
			"foo":   "bar",
			"count": 42,
		}).Error("Some Fields")
		ent = lg2.WithFields(map[string]interface{}{
			"cat":   "mammal",
			"whale": "mammal",
			"moose": "mammal",
		})
		ent.WithFields(map[string]interface{}{
			"shark":  "fish",
			"guppy":  "fish",
			"batray": "fish",
		}).Info("Animal Types")
		lg2.Info("for", "score", "and", "seven", "years", "ago")
		lg.WithField("color", "red").Warn("favorite color")
	}
}
*/
