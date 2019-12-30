const express = require('express')
const firestore = require('./firestore')

const HOSTNAME = 'api.georgeblack.me'

// Express setup
const app = express()
app.use(express.json())
const port = process.env.PORT || 8080

function validateHostname (req, res, next) {
  if (req.hostname !== HOSTNAME) {
    return res.status(403).send('Wrong hostname')
  }
  next()
}

app.use((req, res, next) => {
  res.header('Access-Control-Allow-Origin', 'https://georgeblack.me')
  res.header('Access-Control-Allow-Methods', 'POST, GET, OPTIONS')
  res.header('Access-Control-Allow-Headers', 'Content-Type')
  next()
})

app.options((req, res) => {
  res.send(200)
})

app.post('/views', validateHostname, (req, res) => {
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
  ) {
    return res.status(400).send('Validation failed')
  }

  const docPayload = {
    userAgent: req.body.userAgent,
    hostname: req.body.hostname,
    pathname: req.body.pathname,
    windowInnerWidth: req.body.windowInnerWidth,
    timezone: req.body.timezone,
    timestamp: new Date()
  }

  // append referrer if non-empty
  if (req.body.referrer !== '') docPayload.referrer = req.body.referrer

  // write to firestore
  try {
    firestore.writeView(docPayload)
  } catch (err) {
    return res.status(500).send('Internal error')
  }
  return res.status(200).send('Thanks for visiting :)')
})

app.listen(port, () => console.log(`Listening on port ${port}`))
