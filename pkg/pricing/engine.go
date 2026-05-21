package pricing

import (
	"math"
	"time"
)

const (
	BookingFeeIDR  int64 = 5000
	HourlyRateIDR  int64 = 5000
	OvernightFeeIDR int64 = 20000
)

// Result holds the output of the pricing calculation
type Result struct {
	ParkingFeeIDR   int64
	OvernightFeeIDR int64
	TotalIDR        int64
	DurationHours   int
	IsOvernight     bool
}

// Calculate computes the parking fee based on session duration.
//
// Rules:
//  1. duration = checked_out_at - checked_in_at
//  2. hours = ceil(duration.Minutes() / 60)  — each started hour counts
//  3. parking_fee = hours * 5000
//  4. overnight_fee = 20000 * number of midnights crossed
//  5. total = parking_fee + overnight_fee
//
// Note: booking_fee is charged separately at reservation time and is NOT
// included in the invoice total calculated here.
func Calculate(checkedInAt, checkedOutAt time.Time) Result {
	duration := checkedOutAt.Sub(checkedInAt)
	minutes := duration.Minutes()
	if minutes < 0 {
		minutes = 0
	}

	hours := int(math.Ceil(minutes / 60))
	if hours == 0 {
		hours = 1 // minimum 1 hour
	}

	parkingFee := int64(hours) * HourlyRateIDR

	nights := countMidnightsCrossed(checkedInAt, checkedOutAt)
	overnightFee := int64(nights) * OvernightFeeIDR
	isOvernight := nights > 0

	total := parkingFee + overnightFee

	return Result{
		ParkingFeeIDR:   parkingFee,
		OvernightFeeIDR: overnightFee,
		TotalIDR:        total,
		DurationHours:   hours,
		IsOvernight:     isOvernight,
	}
}

// countMidnightsCrossed returns the number of calendar day boundaries (00:00)
// between checkedInAt and checkedOutAt.
func countMidnightsCrossed(checkedInAt, checkedOutAt time.Time) int {
	// Truncate both to the start of their calendar day
	inDay := time.Date(checkedInAt.Year(), checkedInAt.Month(), checkedInAt.Day(), 0, 0, 0, 0, checkedInAt.Location())
	outDay := time.Date(checkedOutAt.Year(), checkedOutAt.Month(), checkedOutAt.Day(), 0, 0, 0, 0, checkedOutAt.Location())

	days := int(outDay.Sub(inDay).Hours() / 24)
	return days
}
