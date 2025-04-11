app.use({
    install: (app) => {
        const typesByBeautifier = {
             'html': [ 'html', 'xml' ],
              'css': [ 'css', 'sass', 'scss' ],
               'js': [ 'javascript', 'js' ],
             'json': [ 'json', 'json5' ]
        }

        function getBeautifierByType(type) {
            for (const beautifier in typesByBeautifier) {
                if (typesByBeautifier[beautifier].some(t => type.indexOf(t) !== -1)) {
                    return beautifier
                }
            }
            return void 0
        }

        function supports(type) {
            return getBeautifierByType(type) !== void 0
        }

        function beautify(type, code) {
            const beautifier = getBeautifierByType(type)

            if (beautifier === void 0) {
                console.warn(`[$beautifier] The type ${type} is not supported.`)
                return code
            }

            try {
                switch (beautifier) {
                    case 'html': return html_beautify(code)
                    case 'css':  return css_beautify(code)
                    case 'js':   return js_beautify(code)
                    case 'json': return js_beautify(code) // JSON.stringify(JSON.parse(this.code), void 0, '\t')
                }
            } catch(err) {
                console.error(`[$beautifier] Failed to beautify. Error: ${err}`)
                return code
            }

            
        }

        app.provide('$beautifier', { supports: supports, beautify: beautify })
    }
})