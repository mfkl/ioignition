# Simple helper to create secrets
# I always forget the command
.PHONY: new_secret
new_secret:
	openssl rand -base64 64
