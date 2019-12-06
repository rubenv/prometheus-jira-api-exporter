package jiraapiexporter

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
)

func Run() error {
	var opts struct {
		Host   string `long:"host" description:"Jira host" env:"JIRA_HOST"`
		User   string `long:"user" description:"Jira user" env:"JIRA_USER"`
		Pass   string `long:"pass" description:"Jira pass" env:"JIRA_PASS"`
		Listen string `long:"listen" description:"Listen Address" default:":9090" env:"LISTEN"`

		UpdateInterval time.Duration `long:"update-interval" description:"Update interval" default:"15m" env:"UPDATE_INTERVAL"`
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			return nil
		}
		return err
	}
	if opts.Host == "" {
		return errors.New("Missing host!")
	}
	if opts.User == "" {
		return errors.New("Missing user!")
	}
	if opts.Pass == "" {
		return errors.New("Missing pass!")
	}

	every := time.Tick(opts.UpdateInterval)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)

	srv := &http.Server{
		Addr:    opts.Listen,
		Handler: promhttp.Handler(),
	}

	ctx := context.Background()
	ctx, cf := context.WithCancel(ctx)
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		err := srv.ListenAndServe()
		if err == http.ErrServerClosed {
			err = nil
		}
		return err
	})
	g.Go(func() error {
		<-stop
		cf()
		return nil
	})
	g.Go(func() error {
		<-ctx.Done()
		return srv.Shutdown(context.Background())
	})
	g.Go(func() error {
		for {
			err = Update(ctx, opts.Host, opts.User, opts.Pass)
			if err != nil {
				log.Printf("Update error: %s", err)
			}

			select {
			case <-ctx.Done():
				return nil
			case <-every:
			}
		}

	})
	return g.Wait()
}
