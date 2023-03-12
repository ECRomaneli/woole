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
            if (queryStr === void 0 || queryStr === null || queryStr === '') { return objList }
    
            queryStr = queryStr.toLowerCase()
            let tokens = queryStr.split(TOKEN_SEPARATOR)

            if (tokens.length === 1) {
                return objList.filter((obj) => findValue(obj, queryStr, '', exclude))
            }

            let query = {
                key: tokens.shift().trim(), 
                value: tokens.join(TOKEN_SEPARATOR).trim().toLowerCase(),
            }

            return objList.filter((obj) => findQuery(obj, query, '', exclude))
        })
    
        function findQuery(obj, query, nestedKeys, excludeKeys) {
            return getObjectKeys(obj).some((key) => {
                let newNestedKeys = nestedKeys + KEY_SEPARATOR + key.toLowerCase()
                if (excludeKeys !== void 0 && exclude(newNestedKeys, excludeKeys)) { return false }
                
                if (newNestedKeys.indexOf(query.key) === UNKNOWN) {
                    return findQuery(obj[key], query, newNestedKeys, excludeKeys)
                } else if (query.value !== EMPTY_STR && !match(query.value, '', obj[key])) {
                    return findValue(obj[key], query.value, newNestedKeys, excludeKeys)
                }
        
                return true
            })
        }
        
        function findValue(obj, value, nestedKeys, excludeKeys) {
            return getObjectKeys(obj).some((key) => {
                let newNestedKeys = nestedKeys + KEY_SEPARATOR + key.toLowerCase()
                if (excludeKeys !== void 0 && exclude(newNestedKeys, excludeKeys)) { return false }
                
                if (match(value, newNestedKeys, obj[key])) { return true }
                return findValue(obj[key], value, newNestedKeys, excludeKeys)
            })
        }

        function match(expectedValue, key, value) {
            if (key.toLowerCase().indexOf(expectedValue) !== UNKNOWN) { return true }
            
            const typeOf = typeof value

            if (typeOf === STRING || typeOf === NUMBER || typeOf === BOOLEAN || value === null || value === void 0) {
                return `${value}`.toLowerCase().indexOf(expectedValue) !== UNKNOWN
            } 

            return false
        }
        
        function exclude(nestedKeys, excludeKeys) {
            return excludeKeys.some((key) => {
                const indexOf = nestedKeys.indexOf(key)
                return indexOf !== UNKNOWN && indexOf === nestedKeys.length - key.length
            })
        }
        
        function getObjectKeys(obj) {
            if (!(obj instanceof Object)) { return EMPTY_ARR }
            return Object.keys(obj)
        }
    }
})