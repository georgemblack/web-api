const express = require('express')
const app = express()
const port = process.env.PORT

app.get('/', (req, res) => res.send('Hello, Cloud Run!'))

app.get('/example', (req, res) => res.send('Example endpoint'))

app.listen(port, () => console.log(`Listening on port ${port}`))