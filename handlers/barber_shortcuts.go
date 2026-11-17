package handlers

import (
	"context"

	"barber-booking/storage"
)

// shortcutCounts is the JSON shape the live poll uses to refresh the numbers
// under the provider home shortcuts without a page reload.
type shortcutCounts struct {
	Bookings      int `json:"bookings"`
	ProductOrders int `json:"orders"`
	Products      int `json:"products"`
	Hairstyles    int `json:"hairstyles"`
	Messages      int `json:"messages"`
}

func loadProviderCounts(ctx context.Context, barberID int) (storage.ProviderCounts, error) {
	return storage.ProviderShortcutCounts(ctx, barberID)
}

func providerDashboardItems(ctx context.Context, barberID int) ([]dashboardItem, error) {
	counts, err := loadProviderCounts(ctx, barberID)
	if err != nil {
		return nil, err
	}
	return []dashboardItem{
		{Key: "bookings", Label: "Bookings", Path: "/barber/bookings", Count: counts.Bookings},
		{Key: "orders", Label: "Product orders", Path: "/barber/orders", Count: counts.ProductOrders},
		{Key: "products", Label: "Products", Path: "/barber/products", Count: counts.Products},
		{Key: "hairstyles", Label: "Hairstyles", Path: "/barber/hairstyles", Count: counts.Hairstyles},
		{Key: "messages", Label: "Messages", Path: "/barber/chats", Count: counts.Messages},
	}, nil
}

func shortcutCountsForBarber(ctx context.Context, barberID int) (*shortcutCounts, error) {
	counts, err := loadProviderCounts(ctx, barberID)
	if err != nil {
		return nil, err
	}
	return &shortcutCounts{
		Bookings:      counts.Bookings,
		ProductOrders: counts.ProductOrders,
		Products:      counts.Products,
		Hairstyles:    counts.Hairstyles,
		Messages:      counts.Messages,
	}, nil
}
