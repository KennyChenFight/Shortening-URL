package job

import (
	"fmt"
	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/golib/loglib"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type Job interface {
	Name() string
	TimerFormat() string
	Work() (map[string]interface{}, *business.Error)
}

func NewManager(jobs []Job, logger *loglib.Logger) *Manager {
	return &Manager{
		cron:   cron.New(),
		jobs:   jobs,
		logger: logger,
	}
}

type Manager struct {
	cron   *cron.Cron
	jobs   []Job
	logger *loglib.Logger
}

func (m *Manager) Start() {
	for _, job := range m.jobs {
		_, err := m.cron.AddFunc(job.TimerFormat(), func() {
			m.logger.Info(fmt.Sprintf("job:%s start to run this round", job.Name()))
			result, err := job.Work()
			if err != nil {
				m.logger.Error(fmt.Sprintf("job:%s fail this round", job.Name()), zap.Error(err))
				return
			}
			var fields []zap.Field
			for key, value := range result {
				fields = append(fields, zap.Any(key, value))
			}
			m.logger.Info(fmt.Sprintf("job:%s success this round", job.Name()), fields...)
		})
		if err != nil {
			panic(err)
		}
	}
	m.cron.Start()
}

func (m *Manager) Stop() {
	ctx := m.cron.Stop()
	select {
	case <-ctx.Done():
	}
}
