## Golang oAuth2 using gRPC and Docker

### Authentication module/service features:
- [x] Set password complexity by admin (min. length, min. combination Aa$1-9, do not allow name in password)
- [x] Change Postgres to MySQL
- [x] Log all login events (success, failed) with key information like source IP, device, datetime,...
- [x] Enable/disable send notification of every login activity
- [x] Set login session timeout/expiry
- [ ] Allow admin to lock specific account from login
- [ ] User (and admin) can see active sessions and revoke any of them
- [ ] Allow admin to enable/disable 2FA (WhatsApp, SMS, email, authenticator app)
- [ ] Login by temporary link
- [x] Multiple unsuccessfully requests logging. Lock the account until supervisor unlocks.
- [ ] Session timeout/Expiry - event should be pushed for other microservices to listen and do whatever they need
- [ ] When admin will lock an account all active session will immediately be expired
