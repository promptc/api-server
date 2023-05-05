package cast

import (
	"github.com/promptc/api-server/interfaces"
	scheduler "github.com/promptc/openai-scheduler"
)

func SchedulerToOpenAIProvider(scheduler *scheduler.Scheduler) interfaces.OpenAIClientProvider {
	return &provider{client: scheduler}
}

type provider struct {
	client *scheduler.Scheduler
}

func (p *provider) GetClient() interfaces.OpenAIClient {
	return p.client.GetClient()
}
