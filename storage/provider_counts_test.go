package storage

import (
	"context"
	"testing"

	"barber-booking/models"
)

func TestProviderShortcutCountsFollowActivity(t *testing.T) {
	Reset()
	t.Cleanup(Reset)

	barber := newBarber(t, "shortcut_owner@example.com")
	customer := newCustomer(t, "shortcut_client@example.com")
	style := newHairstyle(t, barber.ID)

	counts, err := ProviderShortcutCounts(context.Background(), barber.ID)
	if err != nil {
		t.Fatal(err)
	}
	if counts.Hairstyles != 1 || counts.Products != 0 || counts.Bookings != 0 || counts.ProductOrders != 0 || counts.Messages != 0 {
		t.Fatalf("initial counts: %+v", counts)
	}

	product, err := CreateProduct(models.Product{BarberID: barber.ID, Name: "Pomade", Category: "Care", Price: 1250, Quantity: 2, IsAvailable: true})
	if err != nil {
		t.Fatal(err)
	}

	first, err := CreateBooking(newBooking(customer, barber.ID, style.ID))
	if err != nil {
		t.Fatal(err)
	}
	counts, err = ProviderShortcutCounts(context.Background(), barber.ID)
	if err != nil {
		t.Fatal(err)
	}
	if counts.Bookings != 1 || counts.Products != 1 {
		t.Fatalf("counts after booking/product: %+v", counts)
	}

	if err := MarkProviderBookingsSeen(context.Background(), barber.ID); err != nil {
		t.Fatal(err)
	}
	second := newBooking(customer, barber.ID, style.ID)
	second.Time = "11:00"
	secondBooking, err := CreateBooking(second)
	if err != nil {
		t.Fatal(err)
	}
	counts, err = ProviderShortcutCounts(context.Background(), barber.ID)
	if err != nil {
		t.Fatal(err)
	}
	if counts.Bookings != 1 {
		t.Fatalf("new booking not counted after opening schedule: %+v", counts)
	}
	if err := SetBookingStatus(secondBooking.ID, models.StatusAccepted, barber.ID, 0); err != nil {
		t.Fatal(err)
	}
	counts, err = ProviderShortcutCounts(context.Background(), barber.ID)
	if err != nil {
		t.Fatal(err)
	}
	if counts.Bookings != 0 {
		t.Fatalf("accepted booking still counted: %+v", counts)
	}
	_ = first

	order, err := CreateProductOrder(models.ProductOrder{CustomerID: customer.ID, CustomerName: customer.FullName, Status: models.OrderPending})
	if err != nil {
		t.Fatal(err)
	}
	item, err := CreateProductOrderItem(models.ProductOrderItem{OrderID: order.ID, ProductID: product.ID, BarberID: barber.ID, ProductName: product.Name, Price: product.Price, Quantity: 1, Subtotal: product.Price})
	if err != nil {
		t.Fatal(err)
	}
	counts, err = ProviderShortcutCounts(context.Background(), barber.ID)
	if err != nil {
		t.Fatal(err)
	}
	if counts.ProductOrders != 1 {
		t.Fatalf("pending order not counted: %+v", counts)
	}
	if err := FulfillOrderItem(context.Background(), item.ID, barber.ID, models.FulfillmentReady); err != nil {
		t.Fatal(err)
	}
	if err := FulfillOrderItem(context.Background(), item.ID, barber.ID, models.FulfillmentCollected); err != nil {
		t.Fatal(err)
	}
	counts, err = ProviderShortcutCounts(context.Background(), barber.ID)
	if err != nil {
		t.Fatal(err)
	}
	if counts.ProductOrders != 0 {
		t.Fatalf("collected order still counted: %+v", counts)
	}

	if _, err := AddMessage(context.Background(), models.Message{BarberID: barber.ID, CustomerID: customer.ID, SenderRole: models.SenderCustomer, SenderName: customer.FullName, Text: "Hello"}); err != nil {
		t.Fatal(err)
	}
	counts, err = ProviderShortcutCounts(context.Background(), barber.ID)
	if err != nil {
		t.Fatal(err)
	}
	if counts.Messages != 1 {
		t.Fatalf("unread message not counted: %+v", counts)
	}
	if err := MarkChatSeen(context.Background(), barber.ID, customer.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := AddMessage(context.Background(), models.Message{BarberID: barber.ID, CustomerID: customer.ID, SenderRole: models.SenderBarber, SenderName: barber.FullName, Text: "Hi"}); err != nil {
		t.Fatal(err)
	}
	counts, err = ProviderShortcutCounts(context.Background(), barber.ID)
	if err != nil {
		t.Fatal(err)
	}
	if counts.Messages != 0 {
		t.Fatalf("provider reply counted as unread: %+v", counts)
	}
	if _, err := AddMessage(context.Background(), models.Message{BarberID: barber.ID, CustomerID: customer.ID, SenderRole: models.SenderCustomer, SenderName: customer.FullName, Text: "One more"}); err != nil {
		t.Fatal(err)
	}
	counts, err = ProviderShortcutCounts(context.Background(), barber.ID)
	if err != nil {
		t.Fatal(err)
	}
	if counts.Messages != 1 {
		t.Fatalf("new client message not counted: %+v", counts)
	}
}
