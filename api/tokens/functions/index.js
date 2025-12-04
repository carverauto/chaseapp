const functions = require("firebase-functions")
const escapeHtml = require('escape-html')

const api_key = process.env.API_KEY
const api_secret = process.env.API_SECRET
const { connect } = require('getstream')
const client = connect(api_key, api_secret)

exports.getStreamToken = functions.https.onRequest((req, res) => {
    let id
    ({id} = req.body)
    if (id) {
        const userToken = client.createUserToken(escapeHtml(id))
        res.json({
            data: {
                token: userToken
            }
        })
    } else {
        res.json({
            data: {
                message: 'Missing user_id'
            }
        })
    }
})