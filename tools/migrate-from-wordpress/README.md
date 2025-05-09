# WordPress to SQL Post Migration Tool

This tool converts WordPress posts exported in JSON format to SQL INSERT statements for importing into a custom blog database.

## Features

- Converts WordPress posts and pages to SQL INSERT statements
- Automatically converts YouTube links to embedded iframes
- Preserves post metadata (title, content, dates, etc.)
- Generates unique UUIDs for posts
- Handles slug conflicts

## Requirements

- Python 3.6+
- WordPress post data in JSON format

## Usage

1. Place your WordPress JSON export file named `WP_posts.json` in the same directory as this script
2. Run the script:

```bash
python post_sql_converter.py [--profile-id PROFILE_ID]
```

### Arguments

- `--profile-id`: Optional UUID to use as the profile ID for all posts. If not provided, the script will use a default hardcoded profile ID.

### Example

```bash
python post_sql_converter.py --profile-id "1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p"
```

## Output

The script will generate a file called `wp_posts_migrated.sql` in the same directory, containing SQL INSERT statements that can be executed against your database.

## SQL Schema Compatibility

The generated SQL is designed to work with a database schema that has a `posts` table with the following columns:
- id (UUID)
- profile_id (UUID)
- type (enum: 'POST', 'PAGE')
- title (string)
- content (text)
- created_at (datetime)
- updated_at (datetime)
- is_published (boolean)
- published_at (datetime)
- slug (string)
- excerpt (text)
- view_count (integer)

The SQL uses the MySQL `UUID_TO_BIN()` function to convert UUIDs to binary format. 