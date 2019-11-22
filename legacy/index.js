const AWS = require('aws-sdk')
const uuid = require('uuid/v4')

const tableName = process.env.TABLE_NAME
const dynamoClient = new AWS.DynamoDB.DocumentClient({ apiVersion: '2012-08-10' })

exports.handler = async (event, context) => {
  // validate request data
  if (
    typeof event.userAgent !== 'string' ||
    event.userAgent === '' ||
    typeof event.hostname !== 'string' ||
    event.hostname === '' ||
    typeof event.pathname !== 'string' ||
    event.pathname === '' ||
    typeof event.referrer !== 'string' ||
    typeof event.windowInnerWidth !== 'number' ||
    !Number.isInteger(event.windowInnerWidth) ||
    typeof event.timezone !== 'string' ||
    event.timezone === ''
  ) return 'Validation error'

  const item = {
    id: uuid(),
    userAgent: event.userAgent,
    hostname: event.hostname,
    pathname: event.pathname,
    windowInnerWidth: event.windowInnerWidth,
    timezone: event.timezone,
    timestamp: new Date().toISOString()
  }

  // only add item if non-empty
  if (event.referrer !== '') item.referrer = event.referrer

  await new Promise((resolve, reject) => {
    dynamoClient.put({
      TableName: tableName,
      Item: item
    }, err => {
      if (err) {
        console.log('Error while writing to DynamoDB')
        console.log(err)
        reject(err)
      } else {
        console.log('Successful write to DynamoDB')
        resolve()
      }
    })
  })

  return 'Success'
}
