/**
 * Formats a raw request body into a document that can be stored in Firestore.
 */
function formatPostPayload(requestBody) {
  const docPayload = {
    published: new Date(req.body.published),
    metadata: req.body.metadata,
    content: req.body.content,
  };

  // If location provided, convert to Firestore geopoint
  if ("location" in docPayload.metadata) {
    const lat = docPayload.metadata.location[0];
    const lon = docPayload.metadata.location[1];
    docPayload.metadata.location = new Firestore.GeoPoint(
      Number(lat),
      Number(lon)
    );
  }

  return docPayload;
}

module.exports = {
  formatPostPayload,
};
