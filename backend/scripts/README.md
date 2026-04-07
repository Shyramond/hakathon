# Kuznica - Database Scripts

**Kuznica** - simplified version with only material upgrades, no subscriptions, no item crafting.

## Structure

```
scripts/
|-- init/
|   |-- db_init_kuznica.go   # Kuznica database initialization (upgrades only)
|-- init_db.go               # Original basic initialization
```

## Usage

### Database Initialization

Initialize kuznica database with material upgrades only:

```bash
cd scripts/init
go run db_init_kuznica.go
```

This script:
- Creates all necessary collections
- Sets up material upgrade recipes (Grey Green Blue Purple Gold)
- Creates database indexes for performance
- Creates admin user with empty inventory
- No item crafting, no subscriptions, no limits

### Alternative Initialization

Run basic database initialization:

```bash
go run init_db.go
```

## Database Collections

### Users
- User authentication
- Basic profile info
- No subscription tracking
- No login streaks

### Inventories
- Player materials (Grey, Green, Blue, Purple, Gold)
- Currency
- Daily rewards (no streak tracking)
- Crafting history

### Crafting Recipes
- Only material upgrade recipes:
  - Grey to Green (100% success, 10 quantity)
  - Green to Blue (100% success, 10 quantity)
  - Blue to Purple (100% success, 10 quantity)
  - Purple to Gold (100% success, 10 quantity)

### Daily Rewards
- Basic daily login rewards
- No streak bonuses
- Simple reward system

## Admin User

Default admin credentials:
- **Username:** admin
- **Password:** admin123
- **Email:** admin@soulforge.com
- **Starting Inventory:** Empty (add materials via API)

## API Endpoints

### Authentication
- `POST /api/v1/register` - Register new user
- `POST /api/v1/login` - User login
- `GET /api/v1/profile` - Get user profile

### Inventory
- `GET /api/v1/inventory` - Get user inventory
- `POST /api/v1/inventory/daily-reward` - Claim daily reward
- `POST /api/v1/inventory/materials` - Add materials (admin)
- `GET /api/v1/inventory/history` - Get crafting history

### Crafting
- `POST /api/v1/crafting/upgrade` - Upgrade materials
- `GET /api/v1/crafting/recipes` - Get upgrade recipes

## Database Configuration

- **Database Name:** `kuznica`
- **Connection:** `mongodb://localhost:27017`
- **Collections:** users, inventories, crafting_recipes, daily_rewards

## Features Removed

- No subscription system
- No item crafting (skins, weapons, etc.)
- No login streaks
- No limited items
- No premium features
- No item recipes
- No backup/restore functionality

## Security Notes

1. Change default admin password after first login
2. Use environment variables for database connection in production
3. Monitor database performance

## Development Workflow

1. **Fresh Install:** Use `make init-kuznica`
2. **Testing:** Use API endpoints directly
3. **Production:** Create database manually if needed

## Troubleshooting

### Connection Issues
- MongoDB running on localhost:27017
- Database `kuznica` exists
- Network connectivity

### Performance Issues
- Check database indexes
- Monitor memory usage
- Optimize queries

## Kuznica Benefits

- Simplified codebase
- No subscription complexity
- Focus on core upgrade mechanics
- Easier to maintain
- Faster development
- Minimal dependencies
