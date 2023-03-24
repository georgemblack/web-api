/**
 * Formats a raw request body into a document that can be stored in Firestore.
 */
function formatPostPayload(requestBody) {
  const docPayload = {
    title: requestBody.title,
    slug: requestBody.slug,
    published: new Date(requestBody.published),
    content: requestBody.content,
    draft: requestBody.draft,
    tags: requestBody.tags || [],
  };

  return docPayload;
}

/**
 * Formats a raw request body into a document that can be stored in Firestore.
 */
function formatLikePayload(requestBody) {
  return {
    title: requestBody.title,
    url: requestBody.url,
    timestamp: new Date(),
  };
}

export default {
  formatPostPayload,
  formatLikePayload,
};
