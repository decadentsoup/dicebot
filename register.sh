#!/bin/sh -e

# Load .env file if it is available.
[ -f .env ] && source .env

# Register our commands with Discord.
cat commands.json | http put \
        "https://discord.com/api/v10/applications/$DISCORD_APPLICATION_ID/commands" \
        "Authorization: Bot $DISCORD_APPLICATION_BOT_TOKEN"
