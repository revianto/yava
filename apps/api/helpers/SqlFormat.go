package helpers

// OrderNoSQL returns a SQL expression that formats an order id as "ORD-00001".
// Usage: helpers.OrderNoSQL("tr_orders.id") or helpers.OrderNoSQL("o.id")
func OrderNoSQL(idExpr string) string {
	return `'ORD-' || LPAD(` + idExpr + `::text, 5, '0')`
}
