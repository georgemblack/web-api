const admin = require('firebase-admin')
const uuid = require('uuid/v4')

const VIEW_COLLECTION_NAME = 'personal-web-views-temp'
const LINK_COLLECTION_NAME = 'personal-web-links'

// Firestore connection
admin.initializeApp({
  credential: admin.credential.applicationDefault()
})
const db = admin.firestore()

function writeView (payload) {
  const docRef = db.collection(VIEW_COLLECTION_NAME).doc(uuid())
  docRef.set(payload)
}

module.exports = {
  writeView
}
