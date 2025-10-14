package registry

import (
	"github.com/latolukasz/fluxaorm"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	"github.com/coretrix/hitrix/service/component/sentry"
	"github.com/coretrix/hitrix/service/component/slack"
)

func ServiceProviderErrorLogger() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.ErrorLoggerService,
		Build: func(ctn di.Container) (interface{}, error) {
			var sentryService sentry.ISentry
			var slackAPIService slack.Slack

			sentryServiceInterface, err := ctn.SafeGet(service.SentryService)
			if err == nil {
				sentryService = sentryServiceInterface.(sentry.ISentry)
			}

			slackAPIServiceInterface, err := ctn.SafeGet(service.SlackService)
			if err == nil {
				slackAPIService = slackAPIServiceInterface.(slack.Slack)
			}

			appService := ctn.Get(service.AppService).(*app.App)
			return errorlogger.NewRedisErrorLogger(
				appService,
				ctn.Get(service.ORMGlobalService).(fluxaorm.Context),
				ctn.Get(service.ORMGlobalService).(fluxaorm.Context).Engine().Redis(appService.RedisPools.Persistent),
				slackAPIService,
				sentryService,
				service.RequestBodyKey,
			), nil
		},
	}
}
