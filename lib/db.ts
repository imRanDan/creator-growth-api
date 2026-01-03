import { Pool } from 'pg'

let pool: Pool | null = null

export function getPool(): Pool {
  if (!pool) {
    const connectionString = process.env.POSTGRES_URL || process.env.DATABASE_URL
    
    if (!connectionString) {
      throw new Error('POSTGRES_URL or DATABASE_URL environment variable is required')
    }
    
    pool = new Pool({
      connectionString,
    })
  }
  
  return pool
}

// Template literal SQL helper
export const sql = new Proxy({} as any, {
  get: () => {
    return async (strings: TemplateStringsArray, ...values: any[]) => {
      const pool = getPool()
      const query = strings.reduce((acc, str, i) => {
        return acc + str + (i < values.length ? `$${i + 1}` : '')
      }, '')
      
      const result = await pool.query(query, values)
      return {
        rows: result.rows,
      }
    }
  }
})

export async function runMigrations() {
  const pool = getPool()
  
  try {
    await pool.query(`
      CREATE EXTENSION IF NOT EXISTS "pgcrypto";
      
      CREATE TABLE IF NOT EXISTS users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        email VARCHAR(255) UNIQUE NOT NULL,
        password VARCHAR(255) NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
      );
    `)
    
    await pool.query(`
      CREATE TABLE IF NOT EXISTS instagram_accounts (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        ig_user_id TEXT UNIQUE NOT NULL,
        username TEXT,
        access_token TEXT NOT NULL,
        token_expires_at TIMESTAMP WITH TIME ZONE,
        followers BIGINT,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
      );
    `)
    
    await pool.query(`
      CREATE TABLE IF NOT EXISTS instagram_posts (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        ig_post_id TEXT NOT NULL,
        account_id UUID NOT NULL REFERENCES instagram_accounts(id) ON DELETE CASCADE,
        caption TEXT,
        media_type TEXT,
        media_url TEXT,
        like_count INT,
        comments_count INT,
        posted_at TIMESTAMP WITH TIME ZONE,
        fetched_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        UNIQUE(account_id, ig_post_id)
      );
    `)
    
    await pool.query(`
      CREATE TABLE IF NOT EXISTS waitlist (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        email VARCHAR(255) UNIQUE NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
      );
    `)
    
    console.log('✅ Migrations completed')
    await pool.end()
  } catch (error) {
    console.error('Migration error:', error)
    throw error
  }
}