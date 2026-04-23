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
	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/payment"
	"github.com/Chinsusu/Billing-V2/internal/modules/rbac"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
	"github.com/Chinsusu/Billing-V2/internal/platform/config"
	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
	"github.com/Chinsusu/Billing-V2/internal/platform/logger"
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
		options.CatalogRoutes = newCatalogRoutes(conn)
		options.InvoiceRoutes = newInvoiceRoutes(conn)
		options.OrderRoutes = newOrderRoutes(conn)
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

func newCatalogRoutes(executor platformdb.Executor) app.RouteRegistrar {
	store := catalog.NewPostgresStore(executor)
	service := catalog.NewService(store)
	authorizer := rbac.NewStoreAuthorizer(rbac.NewPostgresStore(executor))
	return catalog.NewHTTPHandlerWithOptions(service, catalog.HTTPHandlerOptions{
		AdminMiddleware:          catalogAuthMiddleware(authorizer, rbac.PermissionCatalogManage, rbac.RiskHigh),
		ResellerViewMiddleware:   catalogAuthMiddleware(authorizer, rbac.PermissionCatalogView, rbac.RiskLow),
		ResellerManageMiddleware: catalogAuthMiddleware(authorizer, rbac.PermissionCatalogManage, rbac.RiskMedium),
		ClientMiddleware:         catalogAuthMiddleware(authorizer, rbac.PermissionCatalogView, rbac.RiskLow),
	})
}

func newOrderRoutes(executor platformdb.Executor) app.RouteRegistrar {
	store := order.NewPostgresStore(executor)
	service := order.NewService(store)
	authorizer := rbac.NewStoreAuthorizer(rbac.NewPostgresStore(executor))
	return order.NewHTTPHandlerWithOptions(service, order.HTTPHandlerOptions{
		AdminMiddleware:         orderAuthMiddleware(authorizer, rbac.PermissionOrderView, rbac.RiskLow),
		AdminManageMiddleware:   orderAuthMiddleware(authorizer, rbac.PermissionOrderManage, rbac.RiskHigh),
		AdminServiceMiddleware:  orderAuthMiddleware(authorizer, rbac.PermissionServiceView, rbac.RiskLow),
		ClientMiddleware:        orderAuthMiddleware(authorizer, rbac.PermissionOrderCreate, rbac.RiskMedium),
		ClientServiceMiddleware: orderAuthMiddleware(authorizer, rbac.PermissionServiceView, rbac.RiskLow),
	})
}

func newInvoiceRoutes(executor platformdb.Executor) app.RouteRegistrar {
	store := invoice.NewPostgresStore(executor)
	service := invoice.NewService(store)
	authorizer := rbac.NewStoreAuthorizer(rbac.NewPostgresStore(executor))
	return invoice.NewHTTPHandlerWithOptions(service, invoice.HTTPHandlerOptions{
		AdminMiddleware:  invoiceAuthMiddleware(authorizer, rbac.PermissionWalletView, rbac.RiskLow),
		ClientMiddleware: invoiceAuthMiddleware(authorizer, rbac.PermissionWalletView, rbac.RiskLow),
	})
}

func newPaymentRoutes(executor platformdb.Executor) app.RouteRegistrar {
	store := payment.NewPostgresStore(executor)
	service := payment.NewService(store)
	authorizer := rbac.NewStoreAuthorizer(rbac.NewPostgresStore(executor))
	return payment.NewHTTPHandlerWithOptions(service, payment.HTTPHandlerOptions{
		AdminMiddleware:  paymentAuthMiddleware(authorizer, rbac.PermissionWalletView, rbac.RiskLow),
		ClientMiddleware: paymentAuthMiddleware(authorizer, rbac.PermissionWalletView, rbac.RiskLow),
	})
}

func newWalletRoutes(executor platformdb.Executor) app.RouteRegistrar {
	store := wallet.NewPostgresStore(executor)
	service := wallet.NewService(store)
	authorizer := rbac.NewStoreAuthorizer(rbac.NewPostgresStore(executor))
	return wallet.NewHTTPHandlerWithOptions(service, wallet.HTTPHandlerOptions{
		AdminMiddleware:  walletAuthMiddleware(authorizer, rbac.PermissionWalletView, rbac.RiskLow),
		ClientMiddleware: walletAuthMiddleware(authorizer, rbac.PermissionWalletView, rbac.RiskLow),
	})
}

func orderAuthMiddleware(authorizer rbac.Authorizer, permission rbac.Permission, risk rbac.RiskLevel) order.RouteMiddleware {
	return chainOrderMiddleware(
		wrapOrderMiddleware(identity.HeaderActorMiddleware),
		order.RouteMiddleware(rbac.RequirePermission(authorizer, permission, risk)),
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

func invoiceAuthMiddleware(authorizer rbac.Authorizer, permission rbac.Permission, risk rbac.RiskLevel) invoice.RouteMiddleware {
	return chainInvoiceMiddleware(
		wrapInvoiceMiddleware(identity.HeaderActorMiddleware),
		invoice.RouteMiddleware(rbac.RequirePermission(authorizer, permission, risk)),
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

func paymentAuthMiddleware(authorizer rbac.Authorizer, permission rbac.Permission, risk rbac.RiskLevel) payment.RouteMiddleware {
	return chainPaymentMiddleware(
		wrapPaymentMiddleware(identity.HeaderActorMiddleware),
		payment.RouteMiddleware(rbac.RequirePermission(authorizer, permission, risk)),
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

func walletAuthMiddleware(authorizer rbac.Authorizer, permission rbac.Permission, risk rbac.RiskLevel) wallet.RouteMiddleware {
	return chainWalletMiddleware(
		wrapWalletMiddleware(identity.HeaderActorMiddleware),
		wallet.RouteMiddleware(rbac.RequirePermission(authorizer, permission, risk)),
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

func catalogAuthMiddleware(authorizer rbac.Authorizer, permission rbac.Permission, risk rbac.RiskLevel) catalog.RouteMiddleware {
	return chainCatalogMiddleware(
		wrapCatalogMiddleware(identity.HeaderActorMiddleware),
		catalog.RouteMiddleware(rbac.RequirePermission(authorizer, permission, risk)),
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
