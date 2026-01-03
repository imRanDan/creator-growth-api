export async function sendWelcomeEmail(email: string): Promise<void> {
  const apiKey = process.env.RESEND_API_KEY
  if (!apiKey) {
    console.log('⚠️  RESEND_API_KEY not set - skipping email send')
    return
  }
  
  const fromEmail = process.env.RESEND_FROM_EMAIL || 'onboarding@resend.dev'
  
  const response = await fetch('https://api.resend.com/emails', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${apiKey}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      from: fromEmail,
      to: [email],
      subject: 'Welcome to Creator Growth! 🎉',
      html: `
        <!DOCTYPE html>
        <html>
        <head>
          <meta charset="utf-8">
          <style>
            body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
            .container { max-width: 600px; margin: 0 auto; padding: 20px; }
            .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; }
            .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 10px 10px; }
          </style>
        </head>
        <body>
          <div class="container">
            <div class="header">
              <h1>📈 Welcome to Creator Growth!</h1>
            </div>
            <div class="content">
              <p>Hey there! 👋</p>
              <p>Thanks for joining the waitlist! We're excited to have you on board.</p>
              <p>We'll notify you as soon as we launch. In the meantime, get ready to:</p>
              <ul>
                <li>📊 Track your Instagram engagement in real-time</li>
                <li>🚀 Get smart insights to grow your audience</li>
                <li>⚡ See which posts perform best</li>
              </ul>
              <p>Talk soon!</p>
              <p>— The Creator Growth Team</p>
            </div>
          </div>
        </body>
        </html>
      `,
    }),
  })
  
  if (!response.ok) {
    const text = await response.text()
    throw new Error(`Resend API error: ${response.status} - ${text}`)
  }
}

