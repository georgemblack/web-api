const admin = require('firebase-admin')
const express = require('express')
const uuid = require('uuid/v4')

// Firestore connection
admin.initializeApp({
  credential: admin.credential.applicationDefault()
})
const db = admin.firestore()

// Express
const app = express()
app.use(express.json())
const port = process.env.PORT

app.get('/', (req, res) => res.send('How are ya now?'))

app.use((req, res, next) => {
  res.header('Access-Control-Allow-Origin', 'https://georgeblack.me')
  res.header('Access-Control-Allow-Headers', 'Content-Type')
  next()
})

app.options((req, res) => {
  res.send(200)
})

app.post('/views', (req, res) => {
  if (req.hostname !== 'api.georgeblack.me') {
    res.status(403).send('Wrong hostname')
  }

  // validate payload
  if (
    typeof req.body.userAgent !== 'string' ||
    req.body.userAgent === '' ||
    typeof req.body.hostname !== 'string' ||
    req.body.hostname === '' ||
    typeof req.body.pathname !== 'string' ||
    req.body.pathname === '' ||
    typeof req.body.referrer !== 'string' ||
    typeof req.body.windowInnerWidth !== 'number' ||
    !Number.isInteger(req.body.windowInnerWidth) ||
    typeof req.body.timezone !== 'string' ||
    req.body.timezone === ''
  ) res.status(400).send('Validation failed')

  const docPayload = {
    userAgent: req.body.userAgent,
    hostname: req.body.hostname,
    pathname: req.body.pathname,
    ipAddress: req.ip,
    forwardedIpAddress: req.headers['X-Forwarded-For'],
    remoteAddress: req.connection.remoteAddress,
    windowInnerWidth: req.body.windowInnerWidth,
    timezone: req.body.timezone,
    timestamp: new Date().toISOString()
  }

  // append referrer if non-empty
  if (req.body.referrer !== '') docPayload.referrer = req.body.referrer

  // write to firestore
  try {
    const docRef = db.collection('personal-web-views').doc(uuid())
    docRef.set(docPayload)
    res.status(200).send('Thanks for visiting :)')
  } catch (err) {
    res.status(500).send('Internal error')
  }
})

app.listen(port, () => console.log(`Listening on port ${port}`))
