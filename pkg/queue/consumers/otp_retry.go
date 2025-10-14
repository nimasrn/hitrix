package consumers

import (
	"fmt"
	"log"
	"time"

	"github.com/latolukasz/fluxaorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/pkg/queue/streams"
	"github.com/coretrix/hitrix/service/component/otp"
)

type OTPRetryConsumer struct {
	ormService      fluxaorm.Context
	maxRetries      int
	gatewayRegistry map[string]otp.IOTPSMSGateway
}

func NewOTPRetryConsumer(ormService fluxaorm.Context, maxRetries int, gatewayRegistry map[string]otp.IOTPSMSGateway) *OTPRetryConsumer {
	return &OTPRetryConsumer{ormService: ormService, maxRetries: maxRetries, gatewayRegistry: gatewayRegistry}
}

func (c *OTPRetryConsumer) GetQueueName() string {
	return streams.StreamMsgRetryOTP
}

func (c *OTPRetryConsumer) GetGroupName(suffix *string) string {
	return streams.GetGroupName(c.GetQueueName(), suffix)
}

func (c *OTPRetryConsumer) Consume(ormService fluxaorm.Context, event fluxaorm.Event) error {
	log.Println(".")

	retryDTO := &otp.RetryDTO{}
	event.Unserialize(retryDTO)

	if retryDTO == nil || retryDTO.Gateway == "" {
		return nil
	}

	otpTrackerEntity, found := fluxaorm.GetByID[entity.OTPTrackerEntity](ormService, retryDTO.OTPTrackerEntityID)
	if !found {
		return nil
	}

	if otpTrackerEntity.MaxRetriesReached {
	}

	RetryOTP(ormService, c.gatewayRegistry, retryDTO, otpTrackerEntity, c.maxRetries)

	return nil
}

func RetryOTP(
	ormService fluxaorm.Context,
	gatewayRegistry map[string]otp.IOTPSMSGateway,
	retryDTO *otp.RetryDTO,
	otpTrackerEntity *entity.OTPTrackerEntity,
	maxRetries int,
) {
	retryAfter := time.Second / 2

	retryCount := 1
	for retryCount <= maxRetries {
		retryCount++

		gateway, ok := gatewayRegistry[retryDTO.Gateway]
		if !ok {
			panic(fmt.Sprintf("gateway %s not found in registry", retryDTO.Gateway))
		}

		var err error

		otpTrackerEntity.GatewaySendRequest, otpTrackerEntity.GatewaySendResponse, err = gateway.SendOTP(retryDTO.Phone, retryDTO.Code)
		if err == nil {
			otpTrackerEntity.GatewaySendStatus = entity.OTPTrackerGatewaySendStatusSent
		}

		otpTrackerEntity.RetryCount = retryCount - 1
		if retryCount == maxRetries {
			otpTrackerEntity.MaxRetriesReached = true
		}

		fluxaorm.EditEntity(ormService, otpTrackerEntity)
		ormService.Flush()

		if otpTrackerEntity.GatewaySendStatus == entity.OTPTrackerGatewaySendStatusSent {
			break
		}

		time.Sleep(retryAfter)
		retryAfter = retryAfter * 2
	}
}
