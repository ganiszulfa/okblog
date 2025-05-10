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

1. Generate `WP_posts.json` with this query
```sql
SELECT 
    p.*,
    GROUP_CONCAT(t.name SEPARATOR ', ') AS tags
FROM 
    wp_posts p
LEFT JOIN 
    wp_term_relationships tr ON p.ID = tr.object_id
LEFT JOIN 
    wp_term_taxonomy tt ON tr.term_taxonomy_id = tt.term_taxonomy_id
LEFT JOIN 
    wp_terms t ON tt.term_id = t.term_id
WHERE 
    p.post_type IN ('post', 'page')
    AND tt.taxonomy = 'post_tag'
GROUP BY 
    p.ID
ORDER BY 
    p.post_date DESC;
```

2. Place your WordPress JSON export file named `WP_posts.json` in the same directory as this script
3. Run the script:

```bash
python post_sql_converter.py [--profile-id PROFILE_ID]
```

### Arguments

- `--profile-id`: Optional UUID to use as the profile ID for all posts. If not provided, the script will use a default hardcoded profile ID.

### Example

```bash
python post_sql_converter.py --profile-id "beefbeef-beef-beef-beef-beefbeefbeef"
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