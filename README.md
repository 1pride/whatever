# Discord Voice Channel Bot

A Discord bot that provides dynamic voice channel creation and management features.

## Features

- **Dynamic Voice Channels**: Create temporary voice channels with `/sala` command
- **Voice Limit Management**: Change user limits in voice channels with `/voice-limite` command
- **Auto-cleanup**: Temporary channels are automatically deleted when empty

## Setup Instructions

### 1. Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
DISCORD_BOT_TOKEN=your_bot_token_here
DISCORD_APP_ID=your_application_id_here
DISCORD_GUILD_ID=your_guild_id_here
```

#### How to get these values:

1. **Bot Token**: 
   - Go to [Discord Developer Portal](https://discord.com/developers/applications)
   - Create a new application or select existing one
   - Go to "Bot" section
   - Copy the token

2. **Application ID**:
   - In the same Discord Developer Portal
   - Go to "General Information" section
   - Copy the Application ID

3. **Guild ID**:
   - Enable Developer Mode in Discord (User Settings > Advanced > Developer Mode)
   - Right-click on your server and select "Copy Server ID"

### 2. Bot Permissions

Make sure your bot has the following permissions:
- `Manage Channels`
- `Send Messages`
- `Use Slash Commands`
- `Connect` (for voice channels)
- `Speak` (for voice channels)

### 3. Configuration

#### Temporary Voice Channel Category ID

You need to update the category ID where temporary voice channels will be created. Open `config/config.go` and change the `CategoryID` constant:

```go
const (
	// CategoryID is the Discord category where temporary voice channels will be created
	CategoryID = "YOUR_CATEGORY_ID_HERE"
)
```

**How to get Category ID:**
1. Enable Developer Mode in Discord
2. Right-click on the category where you want temp channels to appear
3. Select "Copy Category ID"

**Note:** This category ID is used by both the `/sala` command (for creating channels) and the `/voice-limite` command (for identifying temporary channels).

### 4. Installation

1. Install Go (version 1.19 or higher)
2. Clone this repository
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Run the bot:
   ```bash
   go run main.go
   ```

## Commands

### `/sala`
Creates a temporary voice channel with customizable options.

**Options:**
- `name`: Choose from predefined channel names (🏢 Working, 🗣️ Only English, 😊 Chill Talk, etc.)
- `limit`: Set user limit (2-40 users)

**Example:**
```
/sala name:🏢 Working limit:10
```

### `/voice-limite`
Changes the user limit of your current voice channel.

**Options:**
- `limit`: New user limit (0-99, where 0 = unlimited)

**Example:**
```
/voice-limite limit:5
```

## How It Works

1. **Dynamic Voice Creation**: When you use `/sala`, the bot creates a temporary voice channel in the specified category
2. **Auto-cleanup**: The bot monitors temporary channels and automatically deletes them after 10 seconds of being empty
3. **Permission Management**: Users in temporary channels can modify the user limit without needing special permissions

## Troubleshooting

### Common Issues

1. **"Error checking permissions"**: Make sure the bot has the required permissions in your server
2. **Commands not appearing**: Ensure the bot has been properly invited with slash command permissions
3. **Channels not being created**: Verify the category ID is correct and the bot has permission to create channels in that category

### Bot Invite Link

Use this link to invite your bot (replace `YOUR_APP_ID` with your actual Application ID):

```
https://discord.com/api/oauth2/authorize?client_id=YOUR_APP_ID&permissions=8&scope=bot%20applications.commands
```

## Development

The project structure:
```
├── main.go                 # Main bot entry point
├── commands/               # Command handlers
│   ├── dynamic_voice.go    # /sala command handler
│   └── voice_limit.go      # /voice-limite command handler
├── registerCommands/       # Command registration
│   └── register.go         # Command definitions
├── config/                 # Configuration
│   └── config.go           # Global configuration variables
└── .env                    # Environment variables (create this)
```

## License

This project is open source and available under the MIT License.