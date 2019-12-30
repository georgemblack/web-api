const admin = require('firebase-admin')
const uuid = require('uuid/v4')

const VIEW_COLLECTION_NAME = 'personal-web-views'
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

async function getLinks () {
  const snapshot = await db
    .collection(LINK_COLLECTION_NAME)
    .orderBy('timestamp', 'desc')
    .get()

  const links = snapshot.docs.map(doc => {
    const payload = doc.data()
    return {
      id: doc.id,
      timestamp: payload.timestamp._seconds,
      title: payload.title,
      url: payload.url
    }
  })

  return links
}

module.exports = {
  writeView,
  getLinks
}
