app.use({
    install: (app, _) => {
        const   TOKEN_SEPARATOR = ':', 
                KEY_SEPARATOR = '.', 
                UNKNOWN = -1, 
                EMPTY_ARR = [], 
                EMPTY_STR = '',
                STRING = 'string', 
                NUMBER = 'number', 
                BOOLEAN = 'boolean'

        app.provide('$search', (objList, queryStr, exclude) => {
            if (objList === void 0 || objList === null) { return [] }
            if (queryStr === void 0 || queryStr === null || queryStr === '') { return objList.slice() }
    
            let tokens = queryStr.trim().toLowerCase().split(TOKEN_SEPARATOR)

            if (tokens.length === 1) {
                return objList.filter((obj) => findValue(obj, tokens[0], '', exclude))
            }

            let query = {
                key: tokens.shift().trim(),
                value: tokens.join(TOKEN_SEPARATOR).trim()
            }

            return objList.filter((obj) => findQuery(obj, query, '', exclude))
        })
    
        function findQuery(obj, query, nestedKeys, excludedKeys) {
            return getObjectKeys(obj).some((key) => {
                let newNestedKeys = nestedKeys + KEY_SEPARATOR + key.toLowerCase()
                if (excludedKeys !== void 0 && exclude(newNestedKeys, excludedKeys)) { return false }
                
                if (newNestedKeys.indexOf(query.key) === UNKNOWN) {
                    return findQuery(obj[key], query, newNestedKeys, excludedKeys)
                } else if (query.value !== EMPTY_STR && !match(query.value, '', obj[key])) {
                    return findValue(obj[key], query.value, newNestedKeys, excludedKeys)
                }
        
                return true
            })
        }
        
        function findValue(obj, value, nestedKeys, excludedKeys) {
            return getObjectKeys(obj).some((key) => {
                let newNestedKeys = nestedKeys + KEY_SEPARATOR + key.toLowerCase()
                if (excludedKeys !== void 0 && exclude(newNestedKeys, excludedKeys)) { return false }
                
                if (match(value, newNestedKeys, obj[key])) { return true }
                return findValue(obj[key], value, newNestedKeys, excludedKeys)
            })
        }

        function match(expectedValue, key, value) {
            if (key.indexOf(expectedValue) !== UNKNOWN) { return true }
            
            const typeOf = typeof value

            if (typeOf === STRING || typeOf === NUMBER || typeOf === BOOLEAN || value === null || value === void 0) {
                return `${value}`.toLowerCase().indexOf(expectedValue) !== UNKNOWN
            } 

            return false
        }
        
        function exclude(nestedKeys, excludedKeys) {
            return excludedKeys.some((key) => nestedKeys.endsWith(key))
        }
        
        function getObjectKeys(obj) {
            return obj instanceof Object ? Object.keys(obj) : EMPTY_ARR
        }
    }
})