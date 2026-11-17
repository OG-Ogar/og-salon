package storage

import (
	"context"

	"barber-booking/models"
)

type chatSeenKey struct {
	barberID   int
	customerID int
}

// ProviderCounts is the activity behind the five home shortcuts.
//
// Bookings and Messages only count rows the provider has not opened yet.
// ProductOrders counts items that are still waiting to be handed over.
// Products and Hairstyles are the live catalogue totals.
type ProviderCounts struct {
	Bookings      int
	ProductOrders int
	Products      int
	Hairstyles    int
	Messages      int
}

// ProviderShortcutCounts returns the numbers shown under the provider home
// shortcuts. The reads are scoped to one barber and never leak another
// provider's activity.
func ProviderShortcutCounts(ctx context.Context, barberID int) (ProviderCounts, error) {
	ctx, cancel := databaseContext(ctx)
	defer cancel()

	if err := ctx.Err(); err != nil {
		return ProviderCounts{}, err
	}
	if barberID <= 0 {
		return ProviderCounts{}, ErrForbidden
	}

	var counts ProviderCounts
	if databaseEnabled() {
		err := db.QueryRowContext(ctx, `
select
  (select count(*) from public.bookings b
     where b.barber_id=$1 and b.status='Pending'
       and b.id > coalesce((select r.booking_id from degosalon_internal.provider_read_markers r where r.barber_id=$1), 0)),
  (select count(*) from public.product_order_items i
     where i.barber_id=$1 and i.fulfillment_status in ('Pending','Ready')),
  (select count(*) from public.products p
     where p.barber_id=$1 and not p.is_archived),
  (select count(*) from public.hairstyles h
     where h.barber_id=$1 and not h.is_archived),
  (select count(*) from public.messages m
     where m.barber_id=$1 and m.sender_role='customer'
       and m.id > coalesce((select r.last_message_id from degosalon_internal.chat_reads r
                            where r.barber_id=$1 and r.customer_id=m.customer_id), 0))`,
			barberID,
		).Scan(&counts.Bookings, &counts.ProductOrders, &counts.Products, &counts.Hairstyles, &counts.Messages)
		return counts, err
	}

	mu.RLock()
	defer mu.RUnlock()
	bookingMarker := providerBookingSeen[barberID]
	for _, booking := range bookings {
		if booking.BarberID == barberID && booking.Status == models.StatusPending && booking.ID > bookingMarker {
			counts.Bookings++
		}
	}
	for _, item := range productOrderItems {
		if item.BarberID == barberID && (item.FulfillmentStatus == models.FulfillmentPending || item.FulfillmentStatus == models.FulfillmentReady) {
			counts.ProductOrders++
		}
	}
	for _, product := range products {
		if product.BarberID == barberID && !product.Archived {
			counts.Products++
		}
	}
	for _, hairstyle := range hairstyles {
		if hairstyle.BarberID == barberID && !hairstyle.Archived {
			counts.Hairstyles++
		}
	}
	for _, message := range messages {
		if message.BarberID != barberID || message.SenderRole != models.SenderCustomer {
			continue
		}
		key := chatSeenKey{barberID: barberID, customerID: message.CustomerID}
		if message.ID > chatLastSeen[key] {
			counts.Messages++
		}
	}
	return counts, nil
}

// MarkProviderBookingsSeen remembers that the provider opened the schedule.
// The marker only moves forward, so a slow request cannot rewind it.
func MarkProviderBookingsSeen(ctx context.Context, barberID int) error {
	ctx, cancel := databaseContext(ctx)
	defer cancel()

	if err := ctx.Err(); err != nil {
		return err
	}
	if barberID <= 0 {
		return ErrForbidden
	}
	if databaseEnabled() {
		_, err := db.ExecContext(ctx, `
insert into degosalon_internal.provider_read_markers(barber_id, booking_id)
select $1, coalesce(max(id), 0) from public.bookings where barber_id=$1
on conflict (barber_id) do update
set booking_id=greatest(provider_read_markers.booking_id, excluded.booking_id),
    updated_at=now()`, barberID)
		return err
	}

	mu.Lock()
	defer mu.Unlock()
	marker := providerBookingSeen[barberID]
	for _, booking := range bookings {
		if booking.BarberID == barberID && booking.ID > marker {
			marker = booking.ID
		}
	}
	providerBookingSeen[barberID] = marker
	return nil
}

// MarkChatSeen remembers how far the provider has read one conversation.
func MarkChatSeen(ctx context.Context, barberID, customerID int) error {
	ctx, cancel := databaseContext(ctx)
	defer cancel()

	if err := ctx.Err(); err != nil {
		return err
	}
	if barberID <= 0 || customerID <= 0 {
		return ErrForbidden
	}
	if databaseEnabled() {
		_, err := db.ExecContext(ctx, `
insert into degosalon_internal.chat_reads(barber_id, customer_id, last_message_id)
select $1, $2, coalesce(max(id), 0) from public.messages where barber_id=$1 and customer_id=$2
on conflict (barber_id, customer_id) do update
set last_message_id=greatest(chat_reads.last_message_id, excluded.last_message_id),
    updated_at=now()`, barberID, customerID)
		return err
	}

	mu.Lock()
	defer mu.Unlock()
	key := chatSeenKey{barberID: barberID, customerID: customerID}
	marker := chatLastSeen[key]
	for _, message := range messages {
		if message.BarberID == barberID && message.CustomerID == customerID && message.ID > marker {
			marker = message.ID
		}
	}
	chatLastSeen[key] = marker
	return nil
}
