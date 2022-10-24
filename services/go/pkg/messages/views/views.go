package views

// CustomerNotificationPreferences is the recharge representation of these preferences
// TODO the dates in here need to be formatted correctly when we determine how recharge does their datetimes
type CustomerNotificationPreferenceOptIn struct {
	Status           string `json:"status"`
	LastOptInAt      string `json:"last_opt_in_at"`
	LastOptOutAt     string `json:"last_opt_out_at"`
	LastOptInSource  string `json:"last_opt_in_source"`
	LastOptOutSource string `json:"last_opt_out_source"`
}

type CustomerNotificationPreferenceChannel struct {
	Transactional *CustomerNotificationPreferenceOptIn `json:"transactional"`
	Promotional   *CustomerNotificationPreferenceOptIn `json:"promotional"`
	// TODO this needs to change to "repurchase" tag value
	Repurchase *CustomerNotificationPreferenceOptIn `json:"replenishment"`
}

type CustomerNotificationPreferences struct {
	NotificationPreferences struct {
		SMS   *CustomerNotificationPreferenceChannel `json:"sms"`
		Email *CustomerNotificationPreferenceChannel `json:"email"`
	} `json:"notification_preferences"`
}

type Repurchase struct {
	ProductVariantIDs       []string             `json:"product_variant_ids"`
	WorkflowID              string               `json:"workflow_id"`
	DiscountPercentage      int                  `json:"discount_percentage"`
	ReminderIntervalSeconds int64                `json:"reminder_interval_seconds"`
	Templates               *RepurchaseTemplates `json:"templates"`
	Links                   *RepurchaseLinks     `json:"links"`
}
type RepurchaseTemplates struct {
	Welcome  string `json:"welcome"`
	Reminder string `json:"reminder"`
}
type RepurchaseLinks struct {
	Purchase string `json:"purchase"`
}
type StoreRepurchases struct {
	StoreID     string        `json:"store_id"`
	Repurchases []*Repurchase `json:"repurchase_specs"`
}
type RenderedMessage struct {
	Message string `json:"message"`
}
type SelectedRepurchaseReminder struct {
	//SelectedReminder *values.RepurchaseReminder `json:"selected_reminder"`
}

type Customer struct {
	CustomerID string   `json:"customer_id"`
	StoreID    string   `json:"store_id"`
	SafeHours  []string `json:"safe_hours"`
}
