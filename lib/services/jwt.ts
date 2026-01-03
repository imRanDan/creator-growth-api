import jwt from 'jsonwebtoken'

export interface Claims {
  user_id: string
  email: string
  iat?: number
  exp?: number
}

export function generateStateToken(userID: string, email: string): string {
  const secret = process.env.JWT_SECRET!
  if (!secret) {
    throw new Error('JWT_SECRET not set')
  }
  return jwt.sign(
    { user_id: userID, email },
    secret,
    { expiresIn: '10m' }
  )
}

export function generateToken(userID: string, email: string): string {
  const secret = process.env.JWT_SECRET!
  if (!secret) {
    throw new Error('JWT_SECRET not set')
  }
  return jwt.sign(
    { user_id: userID, email },
    secret,
    { expiresIn: '24h' }
  )
}

export function validateToken(token: string): Claims {
  const secret = process.env.JWT_SECRET!
  if (!secret) {
    throw new Error('JWT_SECRET not set')
  }
  return jwt.verify(token, secret) as Claims
}

