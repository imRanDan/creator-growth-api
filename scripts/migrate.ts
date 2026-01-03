import { config } from 'dotenv'
import { resolve } from 'path'
import { runMigrations } from '@/lib/db'

// Load .env.local file
config({ path: resolve(process.cwd(), '.env.local') })

async function main() {
  try {
    console.log('Running database migrations...')
    console.log('POSTGRES_URL:', process.env.POSTGRES_URL ? '✅ Found' : '❌ Missing')
    await runMigrations()
    console.log('✅ Migrations completed successfully!')
    process.exit(0)
  } catch (error) {
    console.error('❌ Migration failed:', error)
    process.exit(1)
  }
}

main()
