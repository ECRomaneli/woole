app.provide('$woole', {
    parseContentType: (contentType) => {
        let content = {}

        if (contentType === void 0 || contentType === '') { return content }

        // workaround for array headers mixed up with headers separated by semicolon
        let tokens = contentType.toLowerCase().split(";").map(str => str.trim())

        // Parse the xxxx/yyyyy content-type
        let categoryAndType = tokens.shift().split('/')
        content.category = categoryAndType[0]
        content.type = categoryAndType[1]

        // Parse other possible tokens
        for (let token in tokens) {
            token.startsWith("charset=") && (content.charset = token.substring(8))
        }

        return content
    },

    decodeQueryParams(req) {
        const rawQueryArr = req.url.split('?')
        if (rawQueryArr.length < 2) { return }

        rawQueryArr.shift()
        const rawQueryParams = rawQueryArr.join('?').split('&')

        const queryParams = {}
        for (const rawQueryParam of rawQueryParams) {
            const pair = rawQueryParam.split('=')
            queryParams[decodeURIComponent(pair[0])] = pair[1] !== void 0 ? decodeURIComponent(pair[1]) : ''
        }

        req.queryParams = queryParams
    },

    encodeQueryParams(req) {
        if (req.queryParams === void 0) { return }

        const queryParams = Object
            .keys(req.queryParams)
            .filter(key => key !== '')
            .map(key => encodeURIComponent(key) + '=' + encodeURIComponent(req.queryParams[key]))
            .join('&')

        let url = req.url.split('?')[0]
        if (queryParams !== '') { url += '?' + queryParams }

        req.url = url
    },

    decodeBody(item) {
        if (item.body) {
            item.b64Body = item.body
            item.body = atob(item.b64Body)
        }
    },

    encodeBody(item) {
        if (item.b64Body) {
            item.body = btoa(item.b64Body)
            item.b64Body = void 0
        }
    },

    parseRequestToCurl(req) {
        // Construct cURL command
        let curlCommand = `curl -X ${req.method} ${req.host}${req.url}`

        // Add headers to cURL command
        Object.keys(req.header).forEach(header => {
            curlCommand += ` \\\n -H '${header}: ${req.header[header]}'`
        })

        // Add req body to cURL command
        if (req.body) {
            curlCommand += ` \\\n --data-raw '${req.body}'`
        }

        req.curl = curlCommand
    }
});