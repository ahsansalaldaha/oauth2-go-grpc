## Authentication module/service features:

1. Set password complexity by admin (min. length, min. combination Aa$1-9, do not allow name in password)
2. Log all login events (success, failed) with key information like source IP, device, datetime,...
3. Enable/disable send notification of every login activity
4. Set login session timeout/expiry
5. Allow admin to lock specific account from login
6. User (and admin) can see active sessions and revoke any of them
7. Allow admin to enable/disable 2FA (WhatsApp, SMS, email, authenticator app)
8. Currently we have implemented Single Sign On for ADFS. we need to ensure including this mechanism
9. Login by temporary link
10. Multiple unsuccessfully requests logging. Lock the account until supervisor unlocks.
11. Session timeout/Expiry - event should be pushed for other microservices to listen and do whatever they need
12. When admin will lock an account all active session will immediately be expired
