app.provide('$util', {
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
            token.indexOf("charset=") === 0 && (content.charset = token.substring(8))
        }

        return content
    }
});