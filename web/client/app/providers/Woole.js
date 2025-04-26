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

    parseRequestToCurl(rec) {
        const req = rec.request

        // Construct cURL command
        let curlCommand = `curl -X ${req.method} ${rec.host}${req.url}`

        // Add headers to cURL command
        Object.keys(req.header).forEach(header => {
            curlCommand += ` \\\n -H '${header}: ${req.header[header]}'`
        })

        // Add req body to cURL command
        if (req.body) {
            curlCommand += ` \\\n --data-raw '${req.body}'`
        }

        req.curl = curlCommand
    },

    getHeader(item, header, defaultValue) {
        if (item === void 0 || item.header === void 0) { return defaultValue }
        if (item.header[header] !== void 0) { return item.header[header] }
        if (item.header[header.toLowerCase()] !== void 0) { return item.header[header.toLowerCase()] }
        return defaultValue
    },

    escapeRegex(str) {
        // Didn't escape "/" and "-" because they are not used in regex without other characters
        return str.replace(/[\\^$.*+?()[\]{}|]/g, '\\$&')
    },

    parseAddress(address) {
        let ip, port
    
        if (address.startsWith('[')) {
            // IPv6 with brackets
            const parts = address.split(']:')
            ip = parts[0].slice(1) // Remove the opening '['
            port = parts[1]
        } else if (address.includes(':') && address.includes('.')) {
            // IPv4 with port
            const parts = address.split(':')
            ip = parts[0]
            port = parts[1]
        } else if (address.includes(':')) {
            // IPv6 without brackets
            const parts = address.split(':')
            port = parts.pop() // Last part is the port
            ip = parts.join(':') // Rejoin the rest as the IPv6 address
        } else {
            // No port
            ip = address
            port = null
        }
    
        return { ip, port }
    },

    parseSize(bytes) {
        const sizeUnits = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
        if (!bytes) { return '0 B' }
        let i = 0
        for (; bytes >= 1000 && i < sizeUnits.length; i++) { bytes /= 1024 }
        return bytes.toFixed(2) + ' ' + sizeUnits[i]
    }
})