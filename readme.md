#  OpenRC on the Go  
  
###  OpenRC service status checkers  
  
- cli-checker - use to check services statuses from CLI (simple, JSON, monitor mode) - use without args to see usage  
- telegram-checker - use with Telegram bot to request check of specified services - add bot token and password to config.ini and write to bot /openrc-check {password} to trigger checking  
- ntfysh-notifier - periodically check statuses and use ntfy.sh to send notification if service changes status  
  
 Generate template config file using ./cli-checker --cfg=./config.ini  
  