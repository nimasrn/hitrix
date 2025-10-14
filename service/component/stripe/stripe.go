package stripe

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/account"
	"github.com/stripe/stripe-go/v72/accountlink"
	portalsession "github.com/stripe/stripe-go/v72/billingportal/session"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/refund"
	"github.com/stripe/stripe-go/v72/setupintent"
	"github.com/stripe/stripe-go/v72/sub"
	"github.com/stripe/stripe-go/v72/webhook"

	"github.com/coretrix/hitrix/service/component/app"
)

const Env = "env"

type Stripe struct {
	webhookSecrets map[string]string
	appService     *app.App
}

func NewStripe(token string, webhookSecrets map[string]string, appService *app.App) *Stripe {
	stripe.Key = token

	return &Stripe{
		webhookSecrets: webhookSecrets,
		appService:     appService,
	}
}

func (s *Stripe) CreateAccount(accountParams *stripe.AccountParams) (*stripe.Account, error) {
	if accountParams.Metadata == nil {
		accountParams.Metadata = map[string]string{Env: s.appService.Mode}
	} else {
		accountParams.Metadata[Env] = s.appService.Mode
	}

	return account.New(accountParams)
}

func (s *Stripe) UpdateAccount(accountID string, accountParams *stripe.AccountParams) (*stripe.Account, error) {
	if accountParams.Metadata == nil {
		accountParams.Metadata = map[string]string{Env: s.appService.Mode}
	} else {
		accountParams.Metadata[Env] = s.appService.Mode
	}

	return account.Update(accountID, accountParams)
}

func (s *Stripe) GetAccount(accountID string, params *stripe.AccountParams) (*stripe.Account, error) {
	return account.GetByID(accountID, params)
}

func (s *Stripe) CreateCustomer(customerParams *stripe.CustomerParams) (*stripe.Customer, error) {
	if customerParams.Metadata == nil {
		customerParams.Metadata = map[string]string{Env: s.appService.Mode}
	} else {
		customerParams.Metadata[Env] = s.appService.Mode
	}

	return customer.New(customerParams)
}

func (s *Stripe) UpdateCustomer(customerID string, customerParams *stripe.CustomerParams) (*stripe.Customer, error) {
	if customerParams.Metadata == nil {
		customerParams.Metadata = map[string]string{Env: s.appService.Mode}
	} else {
		customerParams.Metadata[Env] = s.appService.Mode
	}

	return customer.Update(customerID, customerParams)
}

func (s *Stripe) CreateCheckoutSession(checkoutSessionParams *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error) {
	if checkoutSessionParams.Metadata == nil {
		checkoutSessionParams.Metadata = map[string]string{Env: s.appService.Mode}
	} else {
		checkoutSessionParams.Metadata[Env] = s.appService.Mode
	}

	return session.New(checkoutSessionParams)
}

func (s *Stripe) CreateBillingPortalSession(billingPortalSessionParams *stripe.BillingPortalSessionParams) (*stripe.BillingPortalSession, error) {
	return portalsession.New(billingPortalSessionParams)
}

func (s *Stripe) GetSubscription(subscriptionID string, params *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	return sub.Get(subscriptionID, params)
}

func (s *Stripe) CreateSubscription(subscriptionParams *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	if subscriptionParams.Metadata == nil {
		subscriptionParams.Metadata = map[string]string{Env: s.appService.Mode}
	} else {
		subscriptionParams.Metadata[Env] = s.appService.Mode
	}

	return sub.New(subscriptionParams)
}

func (s *Stripe) UpdateSubscription(subscriptionID string, subscriptionParams *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	if subscriptionParams.Metadata == nil {
		subscriptionParams.Metadata = map[string]string{Env: s.appService.Mode}
	} else {
		subscriptionParams.Metadata[Env] = s.appService.Mode
	}

	return sub.Update(subscriptionID, subscriptionParams)
}

func (s *Stripe) CancelSubscription(subscriptionID string, subscriptionCancelParams *stripe.SubscriptionCancelParams) (*stripe.Subscription, error) {
	if subscriptionCancelParams.Metadata == nil {
		subscriptionCancelParams.Metadata = map[string]string{Env: s.appService.Mode}
	} else {
		subscriptionCancelParams.Metadata[Env] = s.appService.Mode
	}

	return sub.Cancel(subscriptionID, subscriptionCancelParams)
}

func (s *Stripe) CreateSetupIntent(setupIntentParams *stripe.SetupIntentParams) (*stripe.SetupIntent, error) {
	if setupIntentParams.Metadata == nil {
		setupIntentParams.Metadata = map[string]string{Env: s.appService.Mode}
	} else {
		setupIntentParams.Metadata[Env] = s.appService.Mode
	}

	return setupintent.New(setupIntentParams)
}

func (s *Stripe) CreateAccountLink(accountLinkParams *stripe.AccountLinkParams) (*stripe.AccountLink, error) {
	return accountlink.New(accountLinkParams)
}

func (s *Stripe) GetPaymentIntent(paymentIntentID string, paymentIntentParams *stripe.PaymentIntentParams) (*stripe.PaymentIntent, error) {
	return paymentintent.Get(paymentIntentID, paymentIntentParams)
}

func (s *Stripe) CreatePaymentIntentMultiparty(
	paymentIntentParams *stripe.PaymentIntentParams,
	linkedAccountID string,
) (*stripe.PaymentIntent, error) {
	if paymentIntentParams.Metadata == nil {
		paymentIntentParams.Metadata = map[string]string{Env: s.appService.Mode}
	} else {
		paymentIntentParams.Metadata[Env] = s.appService.Mode
	}

	paymentIntentParams.SetStripeAccount(linkedAccountID)

	return paymentintent.New(paymentIntentParams)
}

func (s *Stripe) CreateRefundMultiparty(refundParams *stripe.RefundParams, linkedAccountID string) (*stripe.Refund, error) {
	if refundParams.Metadata == nil {
		refundParams.Metadata = map[string]string{Env: s.appService.Mode}
	} else {
		refundParams.Metadata[Env] = s.appService.Mode
	}

	refundParams.SetStripeAccount(linkedAccountID)

	return refund.New(refundParams)
}

func (s *Stripe) NewCheckoutSession(
	paymentMethods []string,
	mode string,
	successURL string,
	CancelURL string,
	lineItems []*stripe.CheckoutSessionLineItemParams,
	discounts []*stripe.CheckoutSessionDiscountParams,
) *stripe.CheckoutSession {
	params := &stripe.CheckoutSessionParams{
		Params:             stripe.Params{Metadata: map[string]string{Env: s.appService.Mode}},
		PaymentMethodTypes: stripe.StringSlice(paymentMethods),
		LineItems:          lineItems,
		Mode:               stripe.String(mode),
		SuccessURL:         stripe.String(successURL),
		CancelURL:          stripe.String(CancelURL),
		Discounts:          discounts,
	}

	checkoutSession, err := session.New(params)
	if err != nil {
		panic("failed creating new session for stripe checkout" + err.Error())
	}

	return checkoutSession
}

func (s *Stripe) ConstructWebhookEvent(reqBody []byte, signature string, webhookKey string) (stripe.Event, error) {
	secret, ok := s.webhookSecrets[webhookKey]
	if !ok {
		panic("stripe webhook secret [" + webhookKey + "] not found")
	}

	return webhook.ConstructEvent(reqBody, signature, secret)
}

type IStripe interface {
	CreateAccount(accountParams *stripe.AccountParams) (*stripe.Account, error)
	UpdateAccount(accountID string, accountParams *stripe.AccountParams) (*stripe.Account, error)
	GetAccount(accountID string, accountParams *stripe.AccountParams) (*stripe.Account, error)
	CreateCustomer(customerParams *stripe.CustomerParams) (*stripe.Customer, error)
	UpdateCustomer(customerID string, customerParams *stripe.CustomerParams) (*stripe.Customer, error)
	CreateCheckoutSession(checkoutSessionParams *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error)
	GetSubscription(subscriptionID string, subscriptionParams *stripe.SubscriptionParams) (*stripe.Subscription, error)
	CreateSubscription(subscriptionParams *stripe.SubscriptionParams) (*stripe.Subscription, error)
	UpdateSubscription(subscriptionID string, subscriptionParams *stripe.SubscriptionParams) (*stripe.Subscription, error)
	CancelSubscription(subscriptionID string, subscriptionCancelParams *stripe.SubscriptionCancelParams) (*stripe.Subscription, error)
	CreateSetupIntent(setupIntentParams *stripe.SetupIntentParams) (*stripe.SetupIntent, error)
	CreateBillingPortalSession(billingPortalSessionParams *stripe.BillingPortalSessionParams) (*stripe.BillingPortalSession, error)
	CreateAccountLink(accountLinkParams *stripe.AccountLinkParams) (*stripe.AccountLink, error)
	GetPaymentIntent(paymentIntentID string, paymentIntentParams *stripe.PaymentIntentParams) (*stripe.PaymentIntent, error)
	CreatePaymentIntentMultiparty(paymentIntentParams *stripe.PaymentIntentParams, linkedAccountID string) (*stripe.PaymentIntent, error)
	CreateRefundMultiparty(refundParams *stripe.RefundParams, linkedAccountID string) (*stripe.Refund, error)
	ConstructWebhookEvent(reqBody []byte, signature string, webhookKey string) (stripe.Event, error)
	NewCheckoutSession(
		paymentMethods []string,
		mode string,
		successURL string,
		CancelURL string,
		lineItems []*stripe.CheckoutSessionLineItemParams,
		discounts []*stripe.CheckoutSessionDiscountParams,
	) *stripe.CheckoutSession
}
