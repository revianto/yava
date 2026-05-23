# WhatsApp Integration Package

Unified interface untuk mengirim WhatsApp messages melalui berbagai provider.

## 📁 Structure

```
whatsapp/
├── meta/          (Meta WhatsApp Business API)
│   ├── client.go    - API client
│   ├── message.go   - Message builders
│   └── service.go   - High-level functions
├── twilio/        (Twilio WhatsApp API)
│   ├── client.go    - API client
│   ├── message.go   - Message builders
│   └── service.go   - High-level functions
└── README.md      (This file)
```

## 🔧 Configuration

### Using Twilio (Current)

Add to `.env`:
```bash
TWILIO_ACCOUNT_SID=your_account_sid
TWILIO_AUTH_TOKEN=your_auth_token
TWILIO_WHATSAPP_NUMBER=whatsapp:+14155238886
```

### Using Meta WhatsApp Business API

Add to `.env`:
```bash
WHATSAPP_API_TOKEN=your_api_token
WHATSAPP_PHONE_ID=your_phone_id
```

## 📤 Usage Examples

### Send Password Reset (Twilio)

```go
import "github.com/revianto/app/whatsapp/twilio"

// Send async (non-blocking)
twilio.SendAsync(func() error {
    return twilio.SendPasswordReset(phone, resetCode)
})

// Or send sync (blocking)
if err := twilio.SendPasswordReset(phone, resetCode); err != nil {
    log.Println("Failed to send reset code:", err)
}
```

### Send Password Reset (Meta)

```go
import "github.com/revianto/app/whatsapp/meta"

// Same interface, different provider
meta.SendAsync(func() error {
    return meta.SendPasswordReset(phone, resetCode)
})
```

### Send OTP

```go
twilio.SendOTP(phone, "123456")
```

### Send Welcome Message

```go
meta.SendWelcome(phone, "John Doe")
```

### Send Custom Text

```go
twilio.SendText(phone, "Your custom message here")
```

## 🔄 Switching Providers

Both packages implement the same interface. To switch providers:

**Current (Twilio):**
```go
import "github.com/revianto/app/whatsapp/twilio"

twilio.SendPasswordReset(phone, code)
```

**Switch to Meta:**
```go
import "github.com/revianto/app/whatsapp/meta"

meta.SendPasswordReset(phone, code)
```

## 📝 Features

### Both Providers Support

- ✅ Password Reset Messages
- ✅ OTP Verification
- ✅ Welcome Messages
- ✅ Custom Text Messages
- ✅ Async (non-blocking) sending
- ✅ Error handling
- ✅ Automatic phone normalization

### Twilio Specific

- Template-based messages (pre-defined in Twilio)
- ContentSID for template mapping
- Content variables for dynamic content

### Meta Specific

- Pre-approved templates
- Template components with parameters
- Language support (en_US, id, etc.)

## 🧪 Testing

```bash
# Test with Twilio
curl -X POST http://localhost:3001/api/admin/id/auth/user/forgot-password \
  -H 'Content-Type: application/json' \
  -d '{"phone": "+6282134497226"}'

# Check WhatsApp for message
```

## 📊 Provider Comparison

| Feature | Twilio | Meta |
|---------|--------|------|
| Setup Complexity | Easy | Medium |
| Free Tier | Yes (sandbox) | Yes (1000 msgs/month) |
| Per Message Cost | $0.01 | $0.03-0.05 |
| Rate Limiting | Built-in | Manual |
| Support | Excellent | Good |
| Sandbox Testing | Yes | Limited |
| API Version | v1 | v22.0 |

## 🚀 Production Checklist

- [ ] Choose provider (Twilio or Meta)
- [ ] Add credentials to `.env`
- [ ] Configure templates in provider dashboard
- [ ] Update template SIDs in code
- [ ] Implement code storage & verification
- [ ] Add rate limiting
- [ ] Enable audit logging
- [ ] Test end-to-end flow
- [ ] Setup monitoring/alerts

## 📚 Related Documentation

- [FORGOT_PASSWORD_SETUP.md](../../FORGOT_PASSWORD_SETUP.md)
- [RESET_PASSWORD_FLOW.md](../../RESET_PASSWORD_FLOW.md)

## 🔗 Provider Links

- [Twilio WhatsApp API](https://www.twilio.com/whatsapp)
- [Meta WhatsApp Business API](https://developers.facebook.com/docs/whatsapp/cloud-api)
