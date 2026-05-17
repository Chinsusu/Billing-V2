package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Chinsusu/Billing-V2/internal/app"
	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/checkout"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/payment"
	"github.com/Chinsusu/Billing-V2/internal/modules/rbac"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
	"github.com/Chinsusu/Billing-V2/internal/platform/config"
	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
	"github.com/Chinsusu/Billing-V2/internal/platform/logger"
	"github.com/Chinsusu/Billing-V2/internal/platform/secrets"
)

type databaseOpener func(ctx context.Context, cfg platformdb.Config) (*sql.DB, error)

type apiRuntime struct {
	api     *app.API
	cleanup func() error
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "api exited: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return err
	}

	log := logger.New(os.Stdout, cfg.LogLevel)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	runtime, err := newRuntime(ctx, cfg, log, platformdb.Open)
	if err != nil {
		return err
	}
	defer func() {
		if err := runtime.close(); err != nil {
			log.Error("api cleanup failed", logger.String("module", "api"), logger.String("error", err.Error()))
		}
	}()

	return runtime.api.Run(ctx)
}

func newRuntime(ctx context.Context, cfg config.Config, log *logger.Logger, openDatabase databaseOpener) (*apiRuntime, error) {
	options := app.APIOptions{}
	cleanup := func() error { return nil }

	if cfg.DatabaseDSN != "" {
		conn, err := openDatabase(ctx, platformdb.Config{
			DriverName: platformdb.DefaultDriverName,
			DSN:        cfg.DatabaseDSN,
		})
		if err != nil {
			return nil, fmt.Errorf("open api database: %w", err)
		}
		cleanup = conn.Close
		authService, err := newAuthService(conn, cfg)
		if err != nil {
			_ = cleanup()
			return nil, err
		}
		options.AuthRoutes = newAuthRoutes(authService, cfg)
		options.SessionMiddleware = identity.SessionMiddleware(identity.SessionMiddlewareOptions{
			CookieName:            cfg.SessionCookieName,
			Resolver:              authService,
			RequireAdminTwoFactor: true,
		})
		options.AccountRoutes = newAccountRoutes(conn)
		options.AuditRoutes = newAuditRoutes(conn)
		options.CatalogRoutes = newCatalogRoutes(conn)
		options.CheckoutRoutes = newCheckoutRoutes(conn)
		options.InvoiceRoutes = newInvoiceRoutes(conn)
		options.JobsRoutes = newJobsRoutes(conn)
		credentialCipher, err := newServiceCredentialCipher(cfg)
		if err != nil {
			_ = cleanup()
			return nil, err
		}
		options.OrderRoutes = newOrderRoutesWithCredentialCipher(conn, credentialCipher)
		options.PaymentRoutes = newPaymentRoutes(conn)
		options.WalletRoutes = newWalletRoutes(conn)
	}

	api, err := app.NewAPIWithOptions(cfg, log, options)
	if err != nil {
		_ = cleanup()
		return nil, err
	}
	return &apiRuntime{
		api:     api,
		cleanup: cleanup,
	}, nil
}

func newAccountRoutes(executor platformdb.Executor) app.RouteRegistrar {
	service := identity.NewAdminReadService(tenant.NewPostgresStore(executor), identity.NewPostgresUserStore(executor))
	authorizer := rbac.NewStoreAuthorizer(rbac.NewPostgresStore(executor))
	return identity.NewAdminReadHTTPHandlerWithOptions(service, identity.AdminReadHTTPHandlerOptions{
		AdminMiddleware:    accountAuthMiddleware(authorizer, rbac.PermissionTenantView, rbac.RiskLow, adminRouteActorTypes),
		ResellerMiddleware: accountAuthMiddleware(authorizer, rbac.PermissionTenantView, rbac.RiskLow, resellerRouteActorTypes),
	})
}

func newAuditRoutes(executor platformdb.Executor) app.RouteRegistrar {
	store := audit.NewPostgresStore(executor)
	service := audit.NewService(store)
	authorizer := rbac.NewStoreAuthorizer(rbac.NewPostgresStore(executor))
	return audit.NewHTTPHandlerWithOptions(service, audit.HTTPHandlerOptions{
		AdminMiddleware: auditAuthMiddleware(authorizer, rbac.PermissionAuditView, rbac.RiskHigh, adminRouteActorTypes),
	})
}

func newCatalogRoutes(executor platformdb.Executor) app.RouteRegistrar {
	store := catalog.NewPostgresStore(executor)
	service := catalog.NewService(store)
	authorizer := rbac.NewStoreAuthorizer(rbac.NewPostgresStore(executor))
	return catalog.NewHTTPHandlerWithOptions(service, catalog.HTTPHandlerOptions{
		AdminReadMiddleware:      catalogAuthMiddleware(authorizer, rbac.PermissionCatalogView, rbac.RiskLow, adminRouteActorTypes),
		AdminManageMiddleware:    catalogAuthMiddleware(authorizer, rbac.PermissionCatalogManage, rbac.RiskHigh, adminRouteActorTypes),
		ResellerViewMiddleware:   catalogAuthMiddleware(authorizer, rbac.PermissionCatalogView, rbac.RiskLow, resellerRouteActorTypes),
		ResellerManageMiddleware: catalogAuthMiddleware(authorizer, rbac.PermissionCatalogManage, rbac.RiskMedium, resellerRouteActorTypes),
		ClientMiddleware:         catalogAuthMiddleware(authorizer, rbac.PermissionCatalogView, rbac.RiskLow, clientRouteActorTypes),
	})
}

func newCheckoutRoutes(executor platformdb.Executor) app.RouteRegistrar {
	orderStore := order.NewPostgresStore(executor)
	invoiceStore := invoice.NewPostgresStore(executor)
	service := checkout.NewService(invoice.NewServiceWithOrderReader(invoiceStore, orderStore))
	authorizer := rbac.NewStoreAuthorizer(rbac.NewPostgresStore(executor))
	return checkout.NewHTTPHandlerWithOptions(service, checkout.HTTPHandlerOptions{
		ClientMiddleware: checkoutAuthMiddleware(authorizer, rbac.PermissionOrderCreate, rbac.RiskMedium, clientRouteActorTypes),
	})
}

func newOrderRoutes(executor platformdb.Executor) app.RouteRegistrar {
	return newOrderRoutesWithCredentialCipher(executor, nil)
}

func newOrderRoutesWithCredentialCipher(executor platformdb.Executor, credentialCipher secrets.Cipher) app.RouteRegistrar {
	store := order.NewPostgresStore(executor)
	service := order.NewServiceWithOptions(order.ServiceOptions{
		Store:                  store,
		Credentials:            store,
		Audit:                  audit.NewService(audit.NewPostgresStore(executor)),
		CredentialCipher:       credentialCipher,
		CredentialRevealLimits: order.NewPostgresCredentialRevealRateLimiter(executor),
	})
	authorizer := rbac.NewStoreAuthorizer(rbac.NewPostgresStore(executor))
	return order.NewHTTPHandlerWithOptions(service, order.HTTPHandlerOptions{
		AdminMiddleware:                    orderAuthMiddleware(authorizer, rbac.PermissionOrderView, rbac.RiskLow, adminRouteActorTypes),
		AdminManageMiddleware:              orderAuthMiddleware(authorizer, rbac.PermissionOrderManage, rbac.RiskHigh, adminRouteActorTypes),
		AdminServiceMiddleware:             orderAuthMiddleware(authorizer, rbac.PermissionServiceView, rbac.RiskLow, adminRouteActorTypes),
		AdminServiceSuspendMiddleware:      orderAuthMiddleware(authorizer, rbac.PermissionServiceSuspend, rbac.RiskHigh, adminRouteActorTypes),
		AdminServiceUnsuspendMiddleware:    orderAuthMiddleware(authorizer, rbac.PermissionServiceUnsuspend, rbac.RiskHigh, adminRouteActorTypes),
		AdminServiceTerminateMiddleware:    orderAuthMiddleware(authorizer, rbac.PermissionServiceTerminate, rbac.RiskCritical, adminRouteActorTypes),
		AdminCredentialMiddleware:          orderAuthMiddleware(authorizer, rbac.PermissionServiceReveal, rbac.RiskHigh, adminRouteActorTypes),
		ResellerMiddleware:                 orderAuthMiddleware(authorizer, rbac.PermissionOrderView, rbac.RiskLow, resellerRouteActorTypes),
		ResellerServiceMiddleware:          orderAuthMiddleware(authorizer, rbac.PermissionServiceView, rbac.RiskLow, resellerRouteActorTypes),
		ResellerServiceSuspendMiddleware:   orderAuthMiddleware(authorizer, rbac.PermissionServiceSuspend, rbac.RiskHigh, resellerRouteActorTypes),
		ResellerServiceUnsuspendMiddleware: orderAuthMiddleware(authorizer, rbac.PermissionServiceUnsuspend, rbac.RiskHigh, resellerRouteActorTypes),
		ResellerServiceTerminateMiddleware: orderAuthMiddleware(authorizer, rbac.PermissionServiceTerminate, rbac.RiskCritical, resellerRouteActorTypes),
		ResellerCredentialMiddleware:       orderAuthMiddleware(authorizer, rbac.PermissionServiceReveal, rbac.RiskHigh, resellerRouteActorTypes),
		ClientMiddleware:                   orderAuthMiddleware(authorizer, rbac.PermissionOrderCreate, rbac.RiskMedium, clientRouteActorTypes),
		ClientServiceMiddleware:            orderAuthMiddleware(authorizer, rbac.PermissionServiceView, rbac.RiskLow, clientRouteActorTypes),
		ClientServiceRenewMiddleware:       orderAuthMiddleware(authorizer, rbac.PermissionServiceRenew, rbac.RiskMedium, clientRouteActorTypes),
		ClientCredentialMiddleware:         orderAuthMiddleware(authorizer, rbac.PermissionServiceView, rbac.RiskHigh, clientRouteActorTypes),
	})
}

func newServiceCredentialCipher(cfg config.Config) (secrets.Cipher, error) {
	if cfg.EncryptionKey == "" {
		return nil, nil
	}
	cipher, err := secrets.NewAESGCMCipher(cfg.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("configure service credential cipher: %w", err)
	}
	return cipher, nil
}

func newInvoiceRoutes(executor platformdb.Executor) app.RouteRegistrar {
	store := invoice.NewPostgresStore(executor)
	service := invoice.NewService(store)
	authorizer := rbac.NewStoreAuthorizer(rbac.NewPostgresStore(executor))
	return invoice.NewHTTPHandlerWithOptions(service, invoice.HTTPHandlerOptions{
		AdminMiddleware:    invoiceAuthMiddleware(authorizer, rbac.PermissionWalletView, rbac.RiskLow, adminRouteActorTypes),
		ResellerMiddleware: invoiceAuthMiddleware(authorizer, rbac.PermissionWalletView, rbac.RiskLow, resellerRouteActorTypes),
		ClientMiddleware:   invoiceAuthMiddleware(authorizer, rbac.PermissionWalletView, rbac.RiskLow, clientRouteActorTypes),
	})
}

func newJobsRoutes(executor platformdb.Executor) app.RouteRegistrar {
	store := jobs.NewPostgresStore(executor)
	service := jobs.NewServiceWithAudit(store, audit.NewService(audit.NewPostgresStore(executor)))
	authorizer := rbac.NewStoreAuthorizer(rbac.NewPostgresStore(executor))
	return jobs.NewHTTPHandlerWithOptions(service, jobs.HTTPHandlerOptions{
		AdminMiddleware:             jobsAuthMiddleware(authorizer, rbac.PermissionOrderView, rbac.RiskLow, adminRouteActorTypes),
		AdminSummaryMiddleware:      jobsAuthMiddleware(authorizer, rbac.PermissionProvisioningJobView, rbac.RiskLow, adminRouteActorTypes),
		AdminRetryMiddleware:        jobsAuthMiddleware(authorizer, rbac.PermissionProvisioningJobRetry, rbac.RiskHigh, adminRouteActorTypes),
		AdminManualReviewMiddleware: jobsAuthMiddleware(authorizer, rbac.PermissionManualReviewResolve, rbac.RiskHigh, adminRouteActorTypes),
		AdminCancelMiddleware:       jobsAuthMiddleware(authorizer, rbac.PermissionManualReviewResolve, rbac.RiskHigh, adminRouteActorTypes),
		ResellerMiddleware:          jobsAuthMiddleware(authorizer, rbac.PermissionOrderView, rbac.RiskLow, resellerRouteActorTypes),
	})
}

func newPaymentRoutes(executor platformdb.Executor) app.RouteRegistrar {
	store := payment.NewPostgresStore(executor)
	service := payment.NewServiceWithAudit(store, audit.NewService(audit.NewPostgresStore(executor)))
	authorizer := rbac.NewStoreAuthorizer(rbac.NewPostgresStore(executor))
	return payment.NewHTTPHandlerWithOptions(service, payment.HTTPHandlerOptions{
		AdminMiddleware:    paymentAuthMiddleware(authorizer, rbac.PermissionWalletView, rbac.RiskLow, adminRouteActorTypes),
		ResellerMiddleware: paymentAuthMiddleware(authorizer, rbac.PermissionWalletView, rbac.RiskLow, resellerRouteActorTypes),
		ClientMiddleware:   paymentAuthMiddleware(authorizer, rbac.PermissionWalletView, rbac.RiskLow, clientRouteActorTypes),
	})
}

func newWalletRoutes(executor platformdb.Executor) app.RouteRegistrar {
	store := wallet.NewPostgresStore(executor)
	service := wallet.NewServiceWithAudit(store, audit.NewService(audit.NewPostgresStore(executor)))
	authorizer := rbac.NewStoreAuthorizer(rbac.NewPostgresStore(executor))
	return wallet.NewHTTPHandlerWithOptions(service, wallet.HTTPHandlerOptions{
		AdminMiddleware:           walletAuthMiddleware(authorizer, rbac.PermissionWalletView, rbac.RiskLow, adminRouteActorTypes),
		AdminReviewMiddleware:     walletAuthMiddleware(authorizer, rbac.PermissionWalletTopupApprove, rbac.RiskHigh, adminRouteActorTypes),
		AdminAdjustmentMiddleware: walletAuthMiddleware(authorizer, rbac.PermissionWalletAdjustment, rbac.RiskCritical, adminRouteActorTypes),
		ResellerMiddleware:        walletAuthMiddleware(authorizer, rbac.PermissionWalletView, rbac.RiskLow, resellerRouteActorTypes),
		ResellerReviewMiddleware:  walletAuthMiddleware(authorizer, rbac.PermissionWalletTopupApprove, rbac.RiskHigh, resellerRouteActorTypes),
		ClientMiddleware:          walletAuthMiddleware(authorizer, rbac.PermissionWalletView, rbac.RiskLow, clientRouteActorTypes),
	})
}

func auditAuthMiddleware(authorizer rbac.Authorizer, permission rbac.Permission, risk rbac.RiskLevel, allowedActorTypes []identity.ActorType) audit.RouteMiddleware {
	return chainAuditMiddleware(
		wrapAuditMiddleware(identity.HeaderActorMiddleware),
		audit.RouteMiddleware(permissionMiddleware(authorizer, permission, risk, allowedActorTypes)),
	)
}

func chainAuditMiddleware(middlewares ...audit.RouteMiddleware) audit.RouteMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for index := len(middlewares) - 1; index >= 0; index-- {
			if middlewares[index] == nil {
				continue
			}
			next = middlewares[index](next)
		}
		return next
	}
}

func wrapAuditMiddleware(middleware func(http.Handler) http.Handler) audit.RouteMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return middleware(http.HandlerFunc(next)).ServeHTTP
	}
}

func accountAuthMiddleware(authorizer rbac.Authorizer, permission rbac.Permission, risk rbac.RiskLevel, allowedActorTypes []identity.ActorType) identity.AdminReadRouteMiddleware {
	return chainAccountMiddleware(
		wrapAccountMiddleware(identity.HeaderActorMiddleware),
		identity.AdminReadRouteMiddleware(permissionMiddleware(authorizer, permission, risk, allowedActorTypes)),
	)
}

func chainAccountMiddleware(middlewares ...identity.AdminReadRouteMiddleware) identity.AdminReadRouteMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for index := len(middlewares) - 1; index >= 0; index-- {
			if middlewares[index] == nil {
				continue
			}
			next = middlewares[index](next)
		}
		return next
	}
}

func wrapAccountMiddleware(middleware func(http.Handler) http.Handler) identity.AdminReadRouteMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return middleware(http.HandlerFunc(next)).ServeHTTP
	}
}

func orderAuthMiddleware(authorizer rbac.Authorizer, permission rbac.Permission, risk rbac.RiskLevel, allowedActorTypes []identity.ActorType) order.RouteMiddleware {
	return chainOrderMiddleware(
		wrapOrderMiddleware(identity.HeaderActorMiddleware),
		order.RouteMiddleware(permissionMiddleware(authorizer, permission, risk, allowedActorTypes)),
	)
}

func chainOrderMiddleware(middlewares ...order.RouteMiddleware) order.RouteMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for index := len(middlewares) - 1; index >= 0; index-- {
			if middlewares[index] == nil {
				continue
			}
			next = middlewares[index](next)
		}
		return next
	}
}

func wrapOrderMiddleware(middleware func(http.Handler) http.Handler) order.RouteMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return middleware(http.HandlerFunc(next)).ServeHTTP
	}
}

func checkoutAuthMiddleware(authorizer rbac.Authorizer, permission rbac.Permission, risk rbac.RiskLevel, allowedActorTypes []identity.ActorType) checkout.RouteMiddleware {
	return chainCheckoutMiddleware(
		wrapCheckoutMiddleware(identity.HeaderActorMiddleware),
		checkout.RouteMiddleware(permissionMiddleware(authorizer, permission, risk, allowedActorTypes)),
	)
}

func chainCheckoutMiddleware(middlewares ...checkout.RouteMiddleware) checkout.RouteMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for index := len(middlewares) - 1; index >= 0; index-- {
			if middlewares[index] == nil {
				continue
			}
			next = middlewares[index](next)
		}
		return next
	}
}

func wrapCheckoutMiddleware(middleware func(http.Handler) http.Handler) checkout.RouteMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return middleware(http.HandlerFunc(next)).ServeHTTP
	}
}

func invoiceAuthMiddleware(authorizer rbac.Authorizer, permission rbac.Permission, risk rbac.RiskLevel, allowedActorTypes []identity.ActorType) invoice.RouteMiddleware {
	return chainInvoiceMiddleware(
		wrapInvoiceMiddleware(identity.HeaderActorMiddleware),
		invoice.RouteMiddleware(permissionMiddleware(authorizer, permission, risk, allowedActorTypes)),
	)
}

func chainInvoiceMiddleware(middlewares ...invoice.RouteMiddleware) invoice.RouteMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for index := len(middlewares) - 1; index >= 0; index-- {
			if middlewares[index] == nil {
				continue
			}
			next = middlewares[index](next)
		}
		return next
	}
}

func wrapInvoiceMiddleware(middleware func(http.Handler) http.Handler) invoice.RouteMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return middleware(http.HandlerFunc(next)).ServeHTTP
	}
}

func jobsAuthMiddleware(authorizer rbac.Authorizer, permission rbac.Permission, risk rbac.RiskLevel, allowedActorTypes []identity.ActorType) jobs.RouteMiddleware {
	return chainJobsMiddleware(
		wrapJobsMiddleware(identity.HeaderActorMiddleware),
		jobs.RouteMiddleware(permissionMiddleware(authorizer, permission, risk, allowedActorTypes)),
	)
}

func chainJobsMiddleware(middlewares ...jobs.RouteMiddleware) jobs.RouteMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for index := len(middlewares) - 1; index >= 0; index-- {
			if middlewares[index] == nil {
				continue
			}
			next = middlewares[index](next)
		}
		return next
	}
}

func wrapJobsMiddleware(middleware func(http.Handler) http.Handler) jobs.RouteMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return middleware(http.HandlerFunc(next)).ServeHTTP
	}
}

func paymentAuthMiddleware(authorizer rbac.Authorizer, permission rbac.Permission, risk rbac.RiskLevel, allowedActorTypes []identity.ActorType) payment.RouteMiddleware {
	return chainPaymentMiddleware(
		wrapPaymentMiddleware(identity.HeaderActorMiddleware),
		payment.RouteMiddleware(permissionMiddleware(authorizer, permission, risk, allowedActorTypes)),
	)
}

func chainPaymentMiddleware(middlewares ...payment.RouteMiddleware) payment.RouteMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for index := len(middlewares) - 1; index >= 0; index-- {
			if middlewares[index] == nil {
				continue
			}
			next = middlewares[index](next)
		}
		return next
	}
}

func wrapPaymentMiddleware(middleware func(http.Handler) http.Handler) payment.RouteMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return middleware(http.HandlerFunc(next)).ServeHTTP
	}
}

func walletAuthMiddleware(authorizer rbac.Authorizer, permission rbac.Permission, risk rbac.RiskLevel, allowedActorTypes []identity.ActorType) wallet.RouteMiddleware {
	return chainWalletMiddleware(
		wrapWalletMiddleware(identity.HeaderActorMiddleware),
		wallet.RouteMiddleware(permissionMiddleware(authorizer, permission, risk, allowedActorTypes)),
	)
}

func chainWalletMiddleware(middlewares ...wallet.RouteMiddleware) wallet.RouteMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for index := len(middlewares) - 1; index >= 0; index-- {
			if middlewares[index] == nil {
				continue
			}
			next = middlewares[index](next)
		}
		return next
	}
}

func wrapWalletMiddleware(middleware func(http.Handler) http.Handler) wallet.RouteMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return middleware(http.HandlerFunc(next)).ServeHTTP
	}
}

func catalogAuthMiddleware(authorizer rbac.Authorizer, permission rbac.Permission, risk rbac.RiskLevel, allowedActorTypes []identity.ActorType) catalog.RouteMiddleware {
	return chainCatalogMiddleware(
		wrapCatalogMiddleware(tenant.HeaderContextMiddleware),
		wrapCatalogMiddleware(identity.HeaderActorMiddleware),
		catalog.RouteMiddleware(permissionMiddleware(authorizer, permission, risk, allowedActorTypes)),
	)
}

func chainCatalogMiddleware(middlewares ...catalog.RouteMiddleware) catalog.RouteMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for index := len(middlewares) - 1; index >= 0; index-- {
			if middlewares[index] == nil {
				continue
			}
			next = middlewares[index](next)
		}
		return next
	}
}

func wrapCatalogMiddleware(middleware func(http.Handler) http.Handler) catalog.RouteMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return middleware(http.HandlerFunc(next)).ServeHTTP
	}
}

func (runtime *apiRuntime) close() error {
	if runtime == nil || runtime.cleanup == nil {
		return nil
	}
	return runtime.cleanup()
}
