app.use({
    install: (app, _) => {
        const   TOKEN_SEPARATOR = ':', 
                KEY_SEPARATOR = '.', 
                NOT_PREFIX = 'not '
                UNKNOWN = -1, 
                EMPTY_ARR = [], 
                EMPTY_STR = '',
                STRING = 'string', 
                NUMBER = 'number', 
                BOOLEAN = 'boolean',
                BIGINT = 'bigint'
        
        function search(objList, queryStr, exclude) {
            if (objList === void 0 || objList === null) { return [] }
            if (queryStr === void 0 || queryStr === null || queryStr.trim() === EMPTY_STR) { return objList.slice() }
    
            const query = getQuery(queryStr)
            return objList.filter((obj) => query.not !== findQuery(obj, query, EMPTY_STR, exclude))
        }

        function findQuery(obj, query, nestedKeys, excludedKeys, keyFound) {
            return getObjectKeys(obj).some((key) => {
                const newNestedKeys = nestedKeys + KEY_SEPARATOR + key.toLowerCase()

                if (isExcluded(newNestedKeys, excludedKeys)) { return false }
                
                if (keyFound === void 0) {
                    if (newNestedKeys.indexOf(query.key) === UNKNOWN) {
                        return findQuery(obj[key], query, newNestedKeys, excludedKeys) 
                            || (query.value === void 0 && match(query.key, obj[key]))
                    }

                    if (query.value === void 0) { return true }
                }

                return match(query.value, obj[key]) || findQuery(obj[key], query, newNestedKeys, excludedKeys, true)
            })
        }

        function match(expectedValue, value) {
            if (value === null || value === void 0) { return false }
            
            const typeOf = typeof value

            if (typeOf === STRING) {
                return `${value}`.toLowerCase().indexOf(expectedValue) !== UNKNOWN
            }

            if (typeOf === NUMBER || typeOf === BIGINT || typeOf === BOOLEAN)  {
                return `${value}`.indexOf(expectedValue) !== UNKNOWN
            }

            return false
        }

        function getQuery(rawQuery) {
            const query = { not: rawQuery.indexOf(NOT_PREFIX) === 0 }
            if (query.not) { rawQuery = rawQuery.substr(NOT_PREFIX.length) }

            const tokens = rawQuery.trim().toLowerCase().split(TOKEN_SEPARATOR)
            query.key = tokens.shift().trim()
            query.value = tokens.join(TOKEN_SEPARATOR).trim() || void 0

            return query
        }
        
        function isExcluded(nestedKeys, excludedKeys) {
            return excludedKeys !== void 0 && excludedKeys.some((key) => nestedKeys.endsWith(key))
        }
        
        function getObjectKeys(obj) {
            return obj instanceof Object ? Object.keys(obj) : EMPTY_ARR
        }

        app.provide('$search', search)
        document.searchEngine = search
    }
})