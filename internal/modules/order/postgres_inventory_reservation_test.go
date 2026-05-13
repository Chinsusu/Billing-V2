package order

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestReserveInventoryArgsNormalizeAndValidate(t *testing.T) {
	expiresAt := time.Date(2026, 5, 13, 10, 5, 0, 0, time.FixedZone("ICT", 7*60*60))
	args, err := reserveInventoryArgs(ReserveInventoryInput{
		OrderID:          "order-1",
		TenantID:         tenant.ID("tenant-1"),
		ProviderSourceID: catalog.ProviderSourceID("source-1"),
		ExpiresAt:        expiresAt,
	})
	if err != nil {
		t.Fatalf("expected reserve args: %v", err)
	}
	if len(args) != 5 || args[3] != 1 {
		t.Fatalf("unexpected reserve args: %#v", args)
	}
	normalizedExpiry, ok := args[4].(time.Time)
	if !ok || !normalizedExpiry.Equal(expiresAt.UTC()) {
		t.Fatalf("expected UTC expiry, got %#v", args[4])
	}
}

func TestReserveInventoryArgsRejectsBadQuantity(t *testing.T) {
	_, err := reserveInventoryArgs(ReserveInventoryInput{
		OrderID:          "order-1",
		TenantID:         tenant.ID("tenant-1"),
		ProviderSourceID: catalog.ProviderSourceID("source-1"),
		Quantity:         -1,
		ExpiresAt:        time.Now(),
	})
	if !errors.Is(err, ErrReservationQuantityInvalid) {
		t.Fatalf("expected quantity error, got %v", err)
	}
}

func TestReserveInventorySQLUsesConditionalCounterUpdate(t *testing.T) {
	for _, clause := range []string{
		"UPDATE provider_inventory inventory",
		"inventory.capacity_total - inventory.reserved_count - inventory.allocated_count >= $4",
		"inventory.status = 'active'",
		"NOT EXISTS (SELECT 1 FROM existing)",
		"INSERT INTO order_reservations",
		"'reserved'",
		"UNION ALL",
	} {
		if !strings.Contains(reserveInventorySQL, clause) {
			t.Fatalf("expected %q in reserve inventory SQL: %s", clause, reserveInventorySQL)
		}
	}
}

func TestExpireReservationsSQLReleasesReservedCountIdempotently(t *testing.T) {
	for _, clause := range []string{
		"reservation.status = 'reserved'",
		"reservation.expires_at < $2",
		"reservation_expired",
		"GREATEST(0, inventory.reserved_count - expired_totals.quantity)",
		"UPDATE provider_inventory inventory",
		"SELECT COUNT(*)::int FROM expired",
	} {
		if !strings.Contains(expireReservationsSQL, clause) {
			t.Fatalf("expected %q in expire reservations SQL: %s", clause, expireReservationsSQL)
		}
	}
}

func TestExpireReservationsArgsRequireTenantAndNow(t *testing.T) {
	_, err := expireReservationsArgs(ExpireReservationsInput{TenantID: tenant.ID("tenant-1")})
	if !errors.Is(err, ErrReservationExpiryMissing) {
		t.Fatalf("expected now error, got %v", err)
	}
	_, err = expireReservationsArgs(ExpireReservationsInput{Now: time.Now()})
	if !errors.Is(err, tenant.ErrTenantIDMissing) {
		t.Fatalf("expected tenant error, got %v", err)
	}
}
