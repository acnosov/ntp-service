package ntp

import (
	"context"
	"fmt"
	"github.com/aibotsoft/ntp-service/pkg/config"
	"github.com/beevik/ntp"
	"go.uber.org/zap"
	"os/exec"
	"time"
)

type Ntp struct {
	cfg *config.Config
	log *zap.Logger
	ctx context.Context
}

func NewNtp(cfg *config.Config, log *zap.Logger, ctx context.Context) *Ntp {
	return &Ntp{cfg: cfg, log: log, ctx: ctx}

}
func SetSystemDate(newTime time.Time) error {
	dateString := newTime.Format(time.RFC3339Nano)
	//fmt.Printf("Setting system date to: %s, %s\n", dateString, newTime)
	err := exec.Command("date", "--set", dateString).Run()
	return err
}

func (n *Ntp) SyncTime() error {
	start := time.Now()
	ntpTime, err := n.GetNtpTime()
	if err != nil {
		return err
	}
	//n.log.Info("ntp_time", zap.Any("resp", ntpTime))

	msg := "check_only"
	if !n.cfg.Service.CheckOnly {
		if ntpTime.ClockOffset > time.Millisecond || ntpTime.ClockOffset < -time.Millisecond {
			err = SetSystemDate(ntpTime.Time.Add(ntpTime.RTT))
			if err != nil {
				return fmt.Errorf("SetSystemDate_error: %w", err)
			}
			msg = "time_synced"
		} else {
			msg = "time_ok"
		}
	}
	n.log.Info(msg,
		zap.Time("ntp", ntpTime.Time),
		zap.Duration("offset", ntpTime.ClockOffset),
		zap.Duration("rtt", ntpTime.RTT),
		zap.Duration("elapsed", time.Since(start)),
	)
	return nil
}
func (n *Ntp) GetNtpTime() (ntpTime *ntp.Response, err error) {
	for _, h := range n.cfg.Service.Hosts {
		ntpTime, err = ntp.QueryWithOptions(h, ntp.QueryOptions{})
		if err == nil {
			return
		}
		n.log.Warn("query_ntp_server_error", zap.Error(err), zap.String("host", h))
	}
	return
}
func (n *Ntp) Run() error {
	_, err := exec.LookPath("date")
	if err != nil {
		return fmt.Errorf("date binary not found, cannot set system date: %w", err)
	}
	err = n.SyncTime()
	if err != nil {
		return fmt.Errorf("sync_time_error: %w", err)
	}
	syncTick := time.Tick(n.cfg.Service.SyncPeriod)

	for {
		select {
		case <-n.ctx.Done():
			return n.ctx.Err()
		case <-syncTick:
			err := n.SyncTime()
			if err != nil {
				n.log.Error("sync_time_error", zap.Error(err))
			}
		}
	}
}
func (n *Ntp) Close() error {
	return nil
}
