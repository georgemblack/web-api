const admin = require('firebase-admin')
const express = require('express')

// Firestore connection
admin.initializeApp({
  credential: admin.credential.applicationDefault()
});
const db = admin.firestore();

// Express
const app = express()
const port = process.env.PORT

app.get('/', (req, res) => res.send('Hello, Cloud Run!'))

app.get('/example', (req, res) => {
  let exampleDocRef = db.collection('personal-web-views').doc('test-doc-id');
  let result = exampleDocRef.set({
    'someKey': 'someValue',
    'someOtherKey': 'someOtherValue'
  })
  res.send('Thanks for visiting :)')
})

app.listen(port, () => console.log(`Listening on port ${port}`))