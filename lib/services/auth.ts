import { sql } from '@/lib/db'
import { hash, compare } from 'bcryptjs'

export interface User {
  id: string
  email: string
  password?: string
  created_at: Date
  updated_at: Date
}

export async function hashPassword(password: string): Promise<string> {
  return hash(password, 10)
}

export async function checkPassword(password: string, hash: string): Promise<boolean> {
  return compare(password, hash)
}

export async function createUser(email: string, password: string): Promise<User> {
  if (!email || !password) {
    throw new Error('Email and password are required')
  }
  
  if (password.length < 8) {
    throw new Error('Password must be at least 8 characters')
  }
  
  const hashedPassword = await hashPassword(password)
  
  try {
    const result = await sql`
      INSERT INTO users (email, password)
      VALUES (${email}, ${hashedPassword})
      RETURNING id, email, created_at, updated_at
    `
    return result.rows[0] as User
  } catch (error: any) {
    if (error.code === '23505') { // Unique violation
      throw new Error('Email already in use')
    }
    throw error
  }
}

export async function getUserByEmail(email: string): Promise<User | null> {
  const result = await sql`
    SELECT id, email, password, created_at, updated_at
    FROM users
    WHERE email = ${email}
    LIMIT 1
  `
  
  if (result.rows.length === 0) {
    return null
  }
  
  return result.rows[0] as User
}

export async function authenticateUser(email: string, password: string): Promise<User> {
  const user = await getUserByEmail(email)
  
  if (!user || !user.password) {
    throw new Error('Invalid credentials')
  }
  
  const valid = await checkPassword(password, user.password)
  if (!valid) {
    throw new Error('Invalid credentials')
  }
  
  // Remove password from return
  const { password: _, ...userWithoutPassword } = user
  return userWithoutPassword as User
}

